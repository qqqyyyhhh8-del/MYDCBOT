package activity

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/user/dcminigames/pkg/config"
)

type Server struct {
	config     config.ActivityConfig
	httpServer *http.Server
	viteCmd    *exec.Cmd
	viteReady  bool
	viteMu     sync.RWMutex
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type TokenRequest struct {
	Code string `json:"code"`
}

func NewServer(cfg config.ActivityConfig) *Server {
	return &Server{config: cfg}
}

func (s *Server) Start() error {
	// 如果是开发模式，先启动 Vite
	if s.config.DevMode {
		if err := s.startVite(); err != nil {
			return fmt.Errorf("启动 Vite 失败: %w", err)
		}
	}

	// 从 public_url 提取路径前缀
	basePath := s.getBasePath()
	if basePath != "" {
		log.Printf("路径前缀: %s", basePath)
	}

	mux := http.NewServeMux()

	// OAuth2 token 交换端点
	mux.HandleFunc("/api/token", s.handleTokenExchange)
	if basePath != "" {
		mux.HandleFunc(basePath+"/api/token", s.handleTokenExchange)
	}

	// 健康检查
	mux.HandleFunc("/api/health", s.handleHealth)
	if basePath != "" {
		mux.HandleFunc(basePath+"/api/health", s.handleHealth)
	}

	// 无名杀文件系统 API - 同时注册有前缀和无前缀的路由
	fsRoutes := []struct {
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"/checkFile", s.handleCheckFile},
		{"/checkDir", s.handleCheckDir},
		{"/readFile", s.handleReadFile},
		{"/readFileAsText", s.handleReadFileAsText},
		{"/writeFile", s.handleWriteFile},
		{"/getFileList", s.handleGetFileList},
		{"/createDir", s.handleCreateDir},
		{"/removeDir", s.handleRemoveDir},
		{"/removeFile", s.handleRemoveFile},
	}
	for _, route := range fsRoutes {
		mux.HandleFunc(route.path, route.handler)
		if basePath != "" {
			mux.HandleFunc(basePath+route.path, route.handler)
		}
	}

	if s.config.DevMode {
		// 开发模式：代理到 Vite 开发服务器
		viteURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", s.config.VitePort))
		proxy := httputil.NewSingleHostReverseProxy(viteURL)

		// 自定义 Director 来处理 WebSocket 升级和路径前缀
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			// 去掉路径前缀
			if basePath != "" && strings.HasPrefix(req.URL.Path, basePath) {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, basePath)
				if req.URL.Path == "" {
					req.URL.Path = "/"
				}
			}
			originalDirector(req)
			req.Host = viteURL.Host
		}

		// 处理 WebSocket (Vite HMR)
		proxy.ModifyResponse = func(resp *http.Response) error {
			return nil
		}

		proxyHandler := s.corsMiddleware(s.injectDiscordSDKProxy(proxy, basePath))
		mux.Handle("/", proxyHandler)
		if basePath != "" {
			mux.Handle(basePath+"/", proxyHandler)
		}
	} else {
		// 生产模式：直接提供静态文件
		gameDir := s.config.GamePath
		if !filepath.IsAbs(gameDir) {
			wd, _ := os.Getwd()
			gameDir = filepath.Join(wd, gameDir)
		}

		if _, err := os.Stat(gameDir); os.IsNotExist(err) {
			return fmt.Errorf("游戏目录不存在: %s", gameDir)
		}

		fileServer := http.FileServer(http.Dir(gameDir))
		staticHandler := s.corsMiddleware(s.injectDiscordSDK(http.StripPrefix(basePath, fileServer), gameDir, basePath))
		mux.Handle("/", staticHandler)
		if basePath != "" {
			mux.Handle(basePath+"/", staticHandler)
		}
	}

	addr := fmt.Sprintf(":%d", s.config.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Activity 服务器启动在 %s", addr)
	if s.config.DevMode {
		log.Printf("开发模式: 代理到 Vite (端口 %d)", s.config.VitePort)
	} else {
		gameDir := s.config.GamePath
		if !filepath.IsAbs(gameDir) {
			wd, _ := os.Getwd()
			gameDir = filepath.Join(wd, gameDir)
		}
		log.Printf("游戏目录: %s", gameDir)
	}
	if s.config.PublicURL != "" {
		log.Printf("公网地址: %s", s.config.PublicURL)
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Activity 服务器错误: %v", err)
		}
	}()

	return nil
}

func (s *Server) startVite() error {
	gameDir := s.config.GamePath
	if !filepath.IsAbs(gameDir) {
		wd, _ := os.Getwd()
		gameDir = filepath.Join(wd, gameDir)
	}

	if _, err := os.Stat(gameDir); os.IsNotExist(err) {
		return fmt.Errorf("游戏目录不存在: %s", gameDir)
	}

	// 检查 node_modules 是否存在
	nodeModules := filepath.Join(gameDir, "node_modules")
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		log.Println("正在安装依赖 (pnpm install)...")
		installCmd := exec.Command("pnpm", "install")
		installCmd.Dir = gameDir
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("安装依赖失败: %w", err)
		}
	}

	// 启动 Vite 开发服务器
	log.Printf("正在启动 Vite 开发服务器 (端口 %d)...", s.config.VitePort)

	// 使用 npx vite，添加环境变量禁用文件监视以避免 ENOSPC 错误
	s.viteCmd = exec.Command("npx", "vite", "--port", fmt.Sprintf("%d", s.config.VitePort), "--host", "127.0.0.1")
	s.viteCmd.Dir = gameDir
	// 设置环境变量禁用 chokidar 的文件监视
	s.viteCmd.Env = append(os.Environ(), "CHOKIDAR_USEPOLLING=false", "VITE_DISABLE_WATCH=true")

	// 获取 stdout 管道来检测 Vite 是否就绪
	stdout, err := s.viteCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("获取 stdout 失败: %w", err)
	}

	stderr, err := s.viteCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("获取 stderr 失败: %w", err)
	}

	if err := s.viteCmd.Start(); err != nil {
		return fmt.Errorf("启动 Vite 失败: %w", err)
	}

	// 等待 Vite 就绪
	readyChan := make(chan bool, 1)

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[Vite] %s", line)

			// 检测 Vite 是否就绪
			if strings.Contains(line, "Local:") || strings.Contains(line, "ready in") {
				s.viteMu.Lock()
				s.viteReady = true
				s.viteMu.Unlock()
				select {
				case readyChan <- true:
				default:
				}
			}
		}
	}()

	// 等待 Vite 就绪或超时
	select {
	case <-readyChan:
		log.Println("Vite 开发服务器已就绪")
	case <-time.After(60 * time.Second):
		return fmt.Errorf("等待 Vite 启动超时")
	}

	return nil
}

func (s *Server) Stop() error {
	// 先停止 Vite
	if s.viteCmd != nil && s.viteCmd.Process != nil {
		log.Println("正在停止 Vite 开发服务器...")
		if err := s.viteCmd.Process.Kill(); err != nil {
			log.Printf("停止 Vite 失败: %v", err)
		}
		s.viteCmd.Wait()
	}

	// 再停止 HTTP 服务器
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleTokenExchange(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	// 向 Discord 交换 token
	tokenResp, err := s.exchangeCodeForToken(req.Code)
	if err != nil {
		log.Printf("Token 交换失败: %v", err)
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}

func (s *Server) exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Discord API 错误 (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	status := map[string]interface{}{"status": "ok"}
	if s.config.DevMode {
		s.viteMu.RLock()
		status["vite_ready"] = s.viteReady
		s.viteMu.RUnlock()
	}

	json.NewEncoder(w).Encode(status)
}

func (s *Server) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// injectDiscordSDKProxy 为代理模式注入 Discord SDK（通过修改响应）
func (s *Server) injectDiscordSDKProxy(proxy *httputil.ReverseProxy, basePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 对于 WebSocket 请求，直接代理
		if r.Header.Get("Upgrade") == "websocket" {
			proxy.ServeHTTP(w, r)
			return
		}

		// 规范化路径，去掉前缀
		reqPath := r.URL.Path
		if basePath != "" && strings.HasPrefix(reqPath, basePath) {
			reqPath = strings.TrimPrefix(reqPath, basePath)
			if reqPath == "" {
				reqPath = "/"
			}
		}

		// 只对 index.html 或根路径进行注入
		if reqPath == "/" || reqPath == "/index.html" {
			// 先从 Vite 获取响应
			viteURL := fmt.Sprintf("http://127.0.0.1:%d%s", s.config.VitePort, reqPath)
			resp, err := http.Get(viteURL)
			if err != nil {
				http.Error(w, "Vite 服务器不可用", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "读取响应失败", http.StatusInternalServerError)
				return
			}

			// 注入 Discord SDK
			modifiedContent := strings.Replace(
				string(body),
				"</head>",
				s.getDiscordSDKScript(basePath)+"</head>",
				1,
			)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(modifiedContent))
			return
		}

		// 其他请求直接代理
		proxy.ServeHTTP(w, r)
	})
}

// injectDiscordSDK 为静态文件模式注入 Discord SDK
func (s *Server) injectDiscordSDK(next http.Handler, gameDir string, basePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 规范化路径，去掉前缀
		reqPath := r.URL.Path
		if basePath != "" && strings.HasPrefix(reqPath, basePath) {
			reqPath = strings.TrimPrefix(reqPath, basePath)
			if reqPath == "" {
				reqPath = "/"
			}
		}

		// 只对 index.html 或根路径进行注入
		if reqPath == "/" || reqPath == "/index.html" {
			indexPath := filepath.Join(gameDir, "index.html")
			content, err := os.ReadFile(indexPath)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// 注入 Discord SDK 初始化脚本
			modifiedContent := strings.Replace(
				string(content),
				"</head>",
				s.getDiscordSDKScript(basePath)+"</head>",
				1,
			)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(modifiedContent))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getDiscordSDKScript(basePath string) string {
	apiPath := basePath + "/api/token"
	return fmt.Sprintf(`
<script type="module">
// Discord Activity SDK 初始化脚本
(async function() {
    // 检测是否在 Discord Activity 环境中
    const urlParams = new URLSearchParams(window.location.search);
    const isDiscordEnvironment = urlParams.has('frame_id') || urlParams.has('instance_id') || window.self !== window.top;
    
    console.log('[Discord] 环境检测:', {
        frame_id: urlParams.get('frame_id'),
        instance_id: urlParams.get('instance_id'),
        isInIframe: window.self !== window.top,
        isDiscordEnvironment: isDiscordEnvironment
    });

    // 初始化默认的辅助函数（无论是否在 Discord 环境中）
    window.discordReady = false;
    window.discordSdk = null;
    window.discordUser = null;
    window.discordError = null;
    window.isDiscordReady = function() { return window.discordReady; };
    window.getDiscordUser = function() { return window.discordUser; };
    window.getDiscordParticipants = function() { return window.discordParticipants || []; };
    window.inviteToActivity = function() { console.warn('[Discord] 不在 Discord 环境中或 SDK 未就绪'); };

    if (!isDiscordEnvironment) {
        console.log('[Discord] 非 Discord 环境，跳过 SDK 初始化');
        // 触发事件通知游戏可以继续加载（无 Discord 功能）
        window.dispatchEvent(new CustomEvent('discordSkipped', { detail: { reason: 'not_in_discord' } }));
        return;
    }

    try {
        console.log('[Discord] 正在加载 SDK...');
        
        // 带超时的 SDK 加载
        const loadSDKWithTimeout = async (timeout = 15000) => {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), timeout);
            
            try {
                const module = await import('https://esm.sh/@discord/embedded-app-sdk');
                clearTimeout(timeoutId);
                return module;
            } catch (err) {
                clearTimeout(timeoutId);
                throw err;
            }
        };
        
        const { DiscordSDK } = await loadSDKWithTimeout();
        console.log('[Discord] SDK 加载成功');
        
        window.discordSdk = new DiscordSDK('%s');
        
        // 等待 SDK 就绪
        await window.discordSdk.ready();
        console.log('[Discord] SDK 已就绪');
        console.log('[Discord] 频道ID:', window.discordSdk.channelId);
        console.log('[Discord] 服务器ID:', window.discordSdk.guildId);
        
        // 获取授权码
        const { code } = await window.discordSdk.commands.authorize({
            client_id: '%s',
            response_type: 'code',
            state: '',
            prompt: 'none',
            scope: ['identify', 'guilds'],
        });
        console.log('[Discord] 授权成功');
        
        // 交换 token
        const response = await fetch('%s', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code }),
        });
        
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error('Token 交换失败: ' + response.status + ' - ' + errorText);
        }
        
        const { access_token } = await response.json();
        console.log('[Discord] Token 交换成功');
        
        // 使用 token 进行认证
        const auth = await window.discordSdk.commands.authenticate({ access_token });
        console.log('[Discord] 认证成功');
        
        // 存储认证信息供游戏使用
        window.discordAuth = auth;
        window.discordUser = auth.user;
        window.discordChannelId = window.discordSdk.channelId;
        window.discordGuildId = window.discordSdk.guildId;
        window.discordInstanceId = window.discordSdk.instanceId;
        window.discordReady = true;
        
        console.log('[Discord] 用户:', auth.user.username);
        console.log('[Discord] 用户ID:', auth.user.id);
        
        // 更新邀请函数
        window.inviteToActivity = async function() {
            if (!window.discordSdk) return;
            try {
                await window.discordSdk.commands.openInviteDialog();
            } catch (err) {
                console.error('[Discord] 打开邀请对话框失败:', err);
            }
        };
        
        // 触发自定义事件，通知游戏 Discord 已就绪
        window.dispatchEvent(new CustomEvent('discordReady', { 
            detail: { 
                user: auth.user,
                channelId: window.discordSdk.channelId,
                guildId: window.discordSdk.guildId,
                instanceId: window.discordSdk.instanceId
            } 
        }));
        
        // 订阅 Activity 参与者变化
        window.discordSdk.subscribe('ACTIVITY_INSTANCE_PARTICIPANTS_UPDATE', (data) => {
            console.log('[Discord] 参与者更新:', data.participants);
            window.discordParticipants = data.participants;
            window.dispatchEvent(new CustomEvent('discordParticipantsUpdate', { 
                detail: { participants: data.participants } 
            }));
        });
        
    } catch (err) {
        console.error('[Discord] SDK 初始化失败:', err);
        window.discordError = err;
        window.dispatchEvent(new CustomEvent('discordError', { detail: { error: err } }));
        // 即使 Discord 初始化失败，也不应该阻止游戏加载
    }
})();
</script>
`, s.config.ClientID, s.config.ClientID, apiPath)
}

func (s *Server) GetPublicURL() string {
	if s.config.PublicURL != "" {
		return s.config.PublicURL
	}
	return fmt.Sprintf("http://localhost:%d", s.config.Port)
}

// getBasePath 从 public_url 提取路径前缀
func (s *Server) getBasePath() string {
	if s.config.PublicURL == "" {
		return ""
	}
	parsed, err := url.Parse(s.config.PublicURL)
	if err != nil {
		return ""
	}
	path := strings.TrimSuffix(parsed.Path, "/")
	return path
}

// ==================== 无名杀文件系统 API ====================

type FSResponse struct {
	Success  bool        `json:"success"`
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (s *Server) getGameDir() string {
	gameDir := s.config.GamePath
	if !filepath.IsAbs(gameDir) {
		wd, _ := os.Getwd()
		gameDir = filepath.Join(wd, gameDir)
	}
	return gameDir
}

func (s *Server) isInProject(reqPath string) bool {
	gameDir := s.getGameDir()
	fullPath := filepath.Join(gameDir, reqPath)
	normalized := filepath.Clean(fullPath)
	return strings.HasPrefix(normalized, gameDir)
}

func (s *Server) sendFSResponse(w http.ResponseWriter, resp FSResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleCheckFile(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "fileName is required"})
		return
	}
	if !s.isInProject(fileName) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), fileName)
	info, err := os.Stat(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件不存在或无法访问"})
		return
	}
	if info.IsDir() {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "不是一个文件"})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200})
}

func (s *Server) handleCheckDir(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "dir is required"})
		return
	}
	if !s.isInProject(dir) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), dir)
	info, err := os.Stat(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件夹不存在或无法访问"})
		return
	}
	if !info.IsDir() {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "不是一个文件夹"})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200})
}

func (s *Server) handleReadFile(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "fileName is required"})
		return
	}
	if !s.isInProject(fileName) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件不存在"})
		return
	}
	byteArray := make([]int, len(data))
	for i, b := range data {
		byteArray[i] = int(b)
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: byteArray})
}

func (s *Server) handleReadFileAsText(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "fileName is required"})
		return
	}
	if !s.isInProject(fileName) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), fileName)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件不存在"})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: string(data)})
}

type WriteFileRequest struct {
	Path string `json:"path"`
	Data []int  `json:"data"`
}

func (s *Server) handleWriteFile(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 405, ErrorMsg: "Method not allowed"})
		return
	}
	var req WriteFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "Invalid request body"})
		return
	}
	if !s.isInProject(req.Path) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), req.Path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	data := make([]byte, len(req.Data))
	for i, b := range req.Data {
		data[i] = byte(b)
	}
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: true})
}

func (s *Server) handleGetFileList(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "dir is required"})
		return
	}
	if !s.isInProject(dir) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), dir)
	info, err := os.Stat(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件夹不存在"})
		return
	}
	if !info.IsDir() {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "getFileList只适用于文件夹"})
		return
	}
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	var files, folders []string
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > 0 && name[0] != '.' && name[0] != '_' {
			if entry.IsDir() {
				folders = append(folders, name)
			} else {
				files = append(files, name)
			}
		}
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: map[string][]string{"files": files, "folders": folders}})
}

func (s *Server) handleCreateDir(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "dir is required"})
		return
	}
	if !s.isInProject(dir) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), dir)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: true})
}

func (s *Server) handleRemoveDir(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "dir is required"})
		return
	}
	if !s.isInProject(dir) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), dir)
	if err := os.RemoveAll(fullPath); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: true})
}

func (s *Server) handleRemoveFile(w http.ResponseWriter, r *http.Request) {
	s.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "fileName is required"})
		return
	}
	if !s.isInProject(fileName) {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 403, ErrorMsg: "Access denied"})
		return
	}
	fullPath := filepath.Join(s.getGameDir(), fileName)
	info, err := os.Stat(fullPath)
	if err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 404, ErrorMsg: "文件不存在"})
		return
	}
	if info.IsDir() {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 400, ErrorMsg: "不能删除文件夹"})
		return
	}
	if err := os.Remove(fullPath); err != nil {
		s.sendFSResponse(w, FSResponse{Success: false, Code: 500, ErrorMsg: err.Error()})
		return
	}
	s.sendFSResponse(w, FSResponse{Success: true, Code: 200, Data: true})
}