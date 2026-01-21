# DcMiniGames

Discord 小游戏机器人项目，使用 Go 语言开发，采用领域驱动设计 (DDD) 架构。

## 项目概述

这是一个 Discord 机器人，提供多种小游戏功能：

- **UNO 卡牌游戏** - 经典多人对战卡牌游戏
- **宝可梦对战** - 基于 PokeAPI 数据的 1v1 宝可梦对战系统

## 技术栈

- **语言**: Go 1.25.4
- **Discord SDK**: github.com/bwmarrin/discordgo v0.29.0
- **配置管理**: gopkg.in/yaml.v3
- **UUID 生成**: github.com/google/uuid v1.6.0
- **数据来源**: PokeAPI GitHub CSV (宝可梦数据)

## 项目结构

```
DcMiniGames/
├── cmd/
│   └── bot/
│       └── main.go                    # 程序入口
├── internal/
│   ├── application/
│   │   ├── pokemon/
│   │   │   └── handler.go             # 宝可梦对战应用层处理器
│   │   └── uno/
│   │       └── handler.go             # UNO 应用层处理器
│   ├── domain/
│   │   ├── pokemon/
│   │   │   ├── entity/
│   │   │   │   ├── battle.go          # 对战实体
│   │   │   │   ├── battler.go         # 对战中的宝可梦
│   │   │   │   └── pokemon.go         # 宝可梦实体与技能
│   │   │   └── valueobject/
│   │   │       ├── ability.go         # 特性值对象
│   │   │       ├── battlemode.go      # 对战模式
│   │   │       ├── event.go           # 领域事件
│   │   │       ├── item.go            # 道具值对象
│   │   │       ├── nature.go          # 性格值对象
│   │   │       ├── poketype.go        # 属性类型与克制表
│   │   │       └── weather.go         # 天气系统
│   │   ├── shared/
│   │   │   └── llm/
│   │   │       └── client.go          # LLM 客户端接口（预留）
│   │   └── uno/
│   │       ├── entity/
│   │       │   ├── card.go            # 卡牌实体
│   │       │   ├── game.go            # 游戏实体
│   │       │   └── player.go          # 玩家实体
│   │       └── valueobject/
│   │           ├── cardtype.go        # 卡牌类型值对象
│   │           └── color.go           # 颜色值对象
│   ├── infrastructure/
│   │   ├── discord/
│   │   │   └── bot.go                 # Discord Bot 封装
│   │   ├── imaging/
│   │   │   └── card_renderer.go       # 卡牌图片渲染
│   │   ├── persistence/
│   │   │   └── memory/
│   │   │       ├── battle_repo.go     # 宝可梦对战仓储
│   │   │       └── game_repo.go       # UNO 游戏仓储
│   │   └── pokeapi/
│   │       └── client.go              # PokeAPI CSV 数据客户端
│   └── interfaces/
│       └── discord/
│           ├── commands/
│           │   ├── pokemon_commands.go # 宝可梦对战命令处理
│           │   └── uno_commands.go     # UNO 斜杠命令处理
│           └── components/
├── pkg/
│   └── config/
│       └── config.go                  # 配置加载
├── assets/
│   ├── pokemon/
│   │   └── abilities.json             # 特性数据
│   └── uno/                           # UNO 卡牌图片资源
├── config.yaml                        # 配置文件
├── go.mod
└── go.sum
```

## 架构设计

采用 **领域驱动设计 (DDD)** 分层架构：

### 领域层 (Domain Layer)
- `entity/`: 核心业务实体
  - UNO: Card, Game, Player
  - Pokemon: Pokemon, PokemonBuild, Battle, Battler, Move
- `valueobject/`: 值对象
  - UNO: Color, CardType
  - Pokemon: PokeType, Nature, Ability, Item, Weather
- `repository/`: 仓储接口定义（待实现）
- `service/`: 领域服务（待实现）
- `event/`: 领域事件（待实现）

### 应用层 (Application Layer)
- `application/uno/handler.go`: UNO 游戏用例逻辑
- `application/pokemon/handler.go`: 宝可梦对战用例逻辑，包含配置管理和预设系统

### 基础设施层 (Infrastructure Layer)
- `discord/bot.go`: Discord API 封装
- `imaging/card_renderer.go`: 图片渲染服务
- `persistence/memory/`: 内存存储实现
- `pokeapi/client.go`: PokeAPI 数据获取客户端

### 接口层 (Interfaces Layer)
- `discord/commands/`: Discord 斜杠命令处理器

## 配置

配置文件 `config.yaml`:

```yaml
discord:
  token: "YOUR_BOT_TOKEN"
  guild_id: ""  # 留空则全局注册命令

uno:
  assets_path: "./assets/uno"

llm:
  provider: "openai"
  api_key: ""
  base_url: ""
  model: "gpt-4"
```

## 命令

### 构建与运行

```bash
# 构建
go build -o bot ./cmd/bot

# 运行
./bot -config config.yaml

# 或直接运行
go run ./cmd/bot -config config.yaml
```

### 依赖管理

```bash
# 下载依赖
go mod download

# 整理依赖
go mod tidy
```

---

## UNO 游戏功能

### Discord 斜杠命令

- `/uno` - 打开 UNO 游戏面板

### 游戏流程

1. **创建游戏**: 点击「创建游戏」按钮
2. **加入游戏**: 其他玩家点击「加入游戏」
3. **开始游戏**: 房主点击「开始游戏」（至少2人）
4. **进行游戏**: 
   - 查看手牌
   - 点击可出的牌打出
   - 摸牌或跳过回合
5. **游戏结束**: 首位打完手牌的玩家获胜

### 卡牌类型

| 类型 | 说明 |
|------|------|
| 数字牌 (0-9) | 红、蓝、绿、黄四色 |
| Skip | 跳过下一位玩家 |
| Reverse | 反转出牌方向 |
| Draw Two | 下一位玩家摸2张并跳过 |
| Wild | 万能牌，可选择颜色 |
| Wild Draw Four | 万能牌，下一位摸4张 |

### 游戏规则

- 每位玩家初始7张牌
- 必须出与当前牌颜色、数字或类型相同的牌
- 万能牌可随时打出并选择颜色
- 2人游戏时 Reverse 等同于 Skip
- 最多支持10位玩家

---

## 宝可梦对战功能

### Discord 斜杠命令

- `/pokemon` - 打开宝可梦对战面板

### 游戏流程

1. **创建对战**: 点击「创建对战」按钮
2. **加入对战**: 其他玩家使用 `/pokemon` 加入
3. **选择宝可梦**: 
   - 搜索（名称或图鉴编号）
   - 浏览图鉴
   - 快捷选择热门宝可梦
4. **配置宝可梦**:
   - 选择性格（影响能力成长 ±10%）
   - 选择特性
   - 选择 4 个技能
   - 可保存/加载预设
5. **对战阶段**:
   - 每回合选择技能
   - 支持认输
   - 实时战斗日志
6. **对战结束**: 首个击倒对方宝可梦的玩家获胜

### 对战系统特性

#### 完整伤害计算公式
- 基础伤害公式: `((2*Level/5+2)*Power*Atk/Def)/50 + 2`
- 属性克制 (18 种属性完整克制表)
- 同属性加成 (STAB): 1.5x
- 会心一击: 1.5x
- 随机因子: 85%-100%

#### 支持的机制
- **性格系统**: 25 种性格，影响能力值 ±10%
- **特性系统**: 常用特性定义（茂盛、猛火、激流等）
- **能力等级**: -6 到 +6 阶段变化
- **异常状态**: 中毒、灼伤、麻痹、睡眠、冰冻
- **道具效果**: 讲究头带/眼镜/围巾、生命宝珠、气势披带等
- **技能优先度**: -7 到 +5
- **充能技能**: 破坏光线、终极冲击等

#### 预设系统
- 保存宝可梦配置（性格、特性、技能）
- 每用户最多 10 个预设
- 快速加载/删除预设

### 数据来源

使用 PokeAPI GitHub CSV 数据：
- 支持第1-9世代 (1025 只宝可梦)
- 简体中文名称/技能/特性
- 完整种族值、属性、可学技能
- 精灵图 (Showdown 风格动图)

```go
// 数据获取接口
pokeapi.GetPredefinedPokemon(id int) *Pokemon
pokeapi.SearchPredefinedPokemon(keyword string) []*Pokemon
pokeapi.GetAllPredefinedPokemon() []*Pokemon
```

---

## 代码规范

### 命名约定

- 包名：小写单词
- 导出函数/类型：大驼峰 (PascalCase)
- 私有函数/变量：小驼峰 (camelCase)
- 常量：大驼峰

### 错误处理

- 使用 `fmt.Errorf` 包装错误
- 中文错误消息面向用户
- 日志使用 `log` 标准库

### 并发安全

- 仓储使用 `sync.RWMutex` 保护共享状态
- 游戏/对战状态修改通过仓储统一管理
- 配置缓存使用读写锁保护

---

## 扩展计划

### LLM 集成（预留）

`internal/domain/shared/llm/client.go` 定义了 LLM 客户端接口：

```go
type Client interface {
    Chat(ctx context.Context, messages []Message) (string, error)
    ChatStream(ctx context.Context, messages []Message) (<-chan string, error)
    GetModel() string
}
```

可用于未来添加 AI 对战、游戏解说等功能。

### 待实现功能

- [ ] 持久化存储（Redis/数据库）
- [ ] 游戏统计与排行榜
- [ ] 更多小游戏
- [ ] AI 玩家
- [ ] 宝可梦队伍对战（多只宝可梦）
- [ ] 天气系统实装
- [ ] 超级进化 / 极巨化 / 太晶化
- [ ] 游戏房间管理

---

## 资源文件

### UNO 卡牌图片
位于 `assets/uno/` 目录，命名规则：
- 数字牌: `{Color}{Number}.jpg` (如 `Red5.jpg`)
- 功能牌: `{Color}{Type}.jpg` (如 `BlueSkip.jpg`)
- 万能牌: `Wild.jpg`, `WildDraw.jpg`

### 宝可梦数据
- `assets/pokemon/abilities.json`: 特性数据（辅助）
- 主要数据通过 PokeAPI GitHub CSV 在线获取

### 精灵图 URL
```
https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/showdown/{id}.gif
```

---

## 开发注意事项

1. **配置安全**: 不要将包含真实 Token 的 config.yaml 提交到版本控制
2. **图片资源**: 确保 assets 目录包含所有必需的 UNO 卡牌图片
3. **命令同步**: Bot 启动时会自动清理并重新注册斜杠命令
4. **内存存储**: 当前使用内存存储，重启后游戏数据丢失
5. **网络依赖**: 宝可梦数据首次加载需要网络连接（从 GitHub 获取 CSV）
6. **数据缓存**: PokeAPI 数据加载后会缓存在内存中，避免重复请求