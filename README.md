# DcMiniGames

Discord 小游戏机器人项目，使用 Go 语言开发，采用领域驱动设计 (DDD) 架构。

## 项目概述

这是一个 Discord 机器人，提供多种小游戏功能：

- **UNO 卡牌游戏** - 经典多人对战卡牌游戏
- **宝可梦对战** - 基于 PokeAPI 数据的宝可梦对战系统，支持 1v1/3v3/6v6 多种模式

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
│   │   │   ├── ability/               # 特性效果系统 (新增)
│   │   │   │   ├── effect.go          # 特性效果接口定义
│   │   │   │   ├── effects_calc.go    # 计算修正类特性
│   │   │   │   ├── effects_entry.go   # 出场触发类特性
│   │   │   │   ├── effects_formchange.go # 形态变化类特性
│   │   │   │   ├── effects_hit.go     # 受击触发类特性
│   │   │   │   ├── effects_status.go  # 状态免疫类特性
│   │   │   │   ├── effects_turnend.go # 回合结束类特性
│   │   │   │   ├── registry.go        # 特性效果注册表
│   │   │   │   └── service.go         # 特性效果服务
│   │   │   ├── entity/
│   │   │   │   ├── battle.go          # 对战实体 (支持多模式)
│   │   │   │   ├── battler.go         # 对战中的宝可梦
│   │   │   │   ├── battler_adapter.go # Battler 接口适配器
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
��       └── discord/
│           ├── commands/
│           │   ├── pokemon_commands.go # 宝可梦对战命令处理
│           │   └── uno_commands.go     # UNO 斜杠命令处理
│           └── components/
├── pkg/
│   └── config/
│       └── config.go                  # 配置加载
├── assets/
│   ├── pokemon/
│   │   ├── abilities.json             # 特性数据
│   │   └── pending_abilities.md       # 待实现特性列表
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
  - Pokemon: Pokemon, PokemonBuild, Battle, Battler, BattlePlayer, Move
- `valueobject/`: 值对象
  - UNO: Color, CardType
  - Pokemon: PokeType, Nature, Ability, Item, Weather
- `ability/`: 特性效果系统（新增）
  - Effect 接口与 BaseEffect 基础实现
  - 按触发时机分类的特性效果实现
  - Registry 注册表与 Service 服务层

### 应用层 (Application Layer)
- `application/uno/handler.go`: UNO 游戏用例逻辑
- `application/pokemon/handler.go`: 宝可梦对战用例逻辑，包含配置管理、预设系统和 AI 对战

### 基础设施层 (Infrastructure Layer)
- `discord/bot.go`: Discord API 封装
- `imaging/card_renderer.go`: 图片渲染服务
- `persistence/memory/`: 内存存储实现
- `pokeapi/client.go`: PokeAPI 数据获取客户端

### 接口层 (Interfaces Layer)
- `discord/commands/`: Discord 斜杠命令处理器

## 配置

首次运行时会自动创建配置文件 `config.yaml`。配置文件示例：

```yaml
discord:
  token: "YOUR_BOT_TOKEN"  # 必填：Discord Bot Token
  guild_id: ""              # 可选：留空则全局注册命令，填写则仅在指定服务器注册

uno:
  assets_path: "./assets/uno"  # UNO 卡牌图片资源路径

llm:  # LLM 集成（预留，暂未使用）
  provider: "openai"
  api_key: ""
  base_url: ""
  model: "gpt-4"
```

**重要提示**: 请勿将包含真实 Token 的 `config.yaml` 提交到版本控制系统。

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

### 对战模式

支持多种对战模式：

| 模式 | 说明 |
|------|------|
| **单挑 (1v1)** | 经典 1v1 对战，双方各选 1 只宝可梦 |
| **3v3 单打** | 双方各选 3 只宝可梦，可换人 |
| **6v6 单打** | 完整队伍对战，双方各选 6 只宝可梦 |
| **人机对战** | 与 AI 训练师对战（支持 1v1/3v3/6v6） |

### 游戏流程

1. **创建对战**: 选择对战模式（PVP 或人机）
2. **加入对战**: 其他玩家使用 `/pokemon` 加入（PVP 模式）
3. **选择宝可梦**: 
   - 搜索（名称或图鉴编号）
   - 浏览图鉴
   - 快捷选择热门宝可梦
   - 3v3/6v6 模式需选择对应数量的宝可梦
4. **配置宝可梦**:
   - 选择性格（影响能力成长 ±10%）
   - 选择特性（普通特性/隐藏特性）
   - 选择 4 个技能
   - 可保存/加载预设
5. **对战阶段**:
   - 每回合选择技能
   - 3v3/6v6 模式支持换人
   - 支持认输
   - 实时战斗日志
6. **对战结束**: 击倒对方所有宝可梦的玩家获胜

### 对战系统特性

#### 完整伤害计算公式
- 基础伤害公式: `((2*Level/5+2)*Power*Atk/Def)/50 + 2`
- 属性克制 (18 种属性完整克制表)
- 同属性加成 (STAB): 1.5x
- 会心一击: 1.5x
- 随机因子: 85%-100%

#### 特性效果系统

项目实现了完整的特性效果系统（已实现 74 个特性），按触发时机分类：

**出场触发类 (12个)**
- 威吓 (Intimidate) - 降低对手攻击
- 降雨 (Drizzle) - 召唤雨天
- 日照 (Drought) - 召唤晴天
- 扬沙 (Sand Stream) - 召唤沙暴
- 降雪 (Snow Warning) - 召唤冰雹
- 压迫感 (Pressure) - 消耗对手 PP
- 紧张感 (Unnerve) - 阻止对手吃树果
- 下载 (Download) - 根据对手防御/特防提升攻击/特攻
- 察觉 (Frisk) - 识破对手道具
- 不挠之剑 (Intrepid Sword) - 出场时攻击+1
- 不屈之盾 (Dauntless Shield) - 出场时防御+1
- 复制 (Trace) - 复制对手特性

**计算修正类-攻击方 (19个)**
- 大力士/瑜伽之力 (Huge Power/Pure Power) - 物攻×2
- 技术高手 (Technician) - 威力≤60 技能×1.5
- 硬爪 (Tough Claws) - 接触技能×1.3
- 强壮之颚 (Strong Jaw) - 咬类技能×1.5
- 适应力 (Adaptability) - STAB 加成×2
- 强行 (Sheer Force) - 附加效果技能×1.3
- 茂盛/猛火/激流/虫之预感 (Overgrow/Blaze/Torrent/Swarm) - HP≤1/3 时属性技能×1.5
- 铁拳 (Iron Fist) - 拳类技能×1.2
- 狙击手 (Sniper) - 会心伤害×1.5
- 沙之力 (Sand Force) - 沙暴中岩/地/钢技能×1.3
- 有色眼镜 (Tinted Lens) - 不很有效技能×2
- 脑核之力 (Neuroforce) - 效果绝佳技能×1.25
- 舍身 (Reckless) - 反伤技能×1.2
- 超级发射器 (Mega Launcher) - 波动/脉冲技能×1.5
- 钢能力者 (Steelworker) - 钢属性技能×1.5

**计算修正类-防御方 (17个)**
- 厚脂肪 (Thick Fat) - 火/冰伤害减半
- 漂浮 (Levitate) - 免疫地面
- 神奇守护 (Wonder Guard) - 只受弱点伤害
- 多重鳞片 (Multiscale) - 满血时伤害减半
- 避雷针 (Lightning Rod) - 吸引电系技能
- 蓄电 (Volt Absorb) - 电系技能回复HP
- 储水 (Water Absorb) - 水系技能回复HP
- 引火 (Flash Fire) - 吸收火系技能提升火威力
- 毛皮大衣 (Fur Coat) - 物理伤害减半
- 耐热 (Heatproof) - 火伤害减半
- 干燥皮肤 (Dry Skin) - 水回复/火伤害×1.25
- 引水 (Storm Drain) - 吸引水系技能
- 食草 (Sap Sipper) - 草系技能回复HP并提升攻击
- 电气引擎 (Motor Drive) - 电系技能提升速度
- 坚硬岩石 (Solid Rock) - 效果绝佳伤害×0.75
- 过滤 (Filter) - 效果绝佳伤害×0.75
- 棱镜装甲 (Prism Armor) - 效果绝佳伤害×0.75

**受击触发类 (12个)**
- 静电 (Static) - 30% 麻痹接触者
- 诅咒之躯 (Cursed Body) - 30% 封印对手技能
- 恶臭 (Stench) - 10% 畏缩
- 毒刺 (Poison Point) - 30% 中毒
- 火焰之躯 (Flame Body) - 30% 灼伤
- 粗糙皮肤 (Rough Skin) - 反伤1/8 HP
- 孢子 (Effect Spore) - 接触时可能中毒/麻痹/睡眠
- 铁刺 (Iron Barbs) - 反伤1/8 HP
- 迷人之躯 (Cute Charm) - 30% 着迷
- 黏滑 (Gooey) - 降低接触者速度
- 卷发 (Tangling Hair) - 降低接触者速度
- 木乃伊 (Mummy) - 接触者特性变为木乃伊

**状态免疫类 (9个)**
- 免疫 (Immunity) - 免疫中毒
- 精神力 (Inner Focus) - 免疫畏缩
- 柔软 (Limber) - 免疫麻痹
- 不眠 (Insomnia) - 免疫睡眠
- 干劲 (Vital Spirit) - 免疫睡眠
- 熔岩铠甲 (Magma Armor) - 免疫冰冻
- 水幕 (Water Veil) - 免疫灼伤
- 我行我素 (Own Tempo) - 免疫混乱
- 迟钝 (Oblivious) - 免疫着迷/挑衅

**回合结束类 (6个)**
- 加速 (Speed Boost) - 每回合速度+1
- 雨盘 (Rain Dish) - 雨天回复HP
- 冰冻之躯 (Ice Body) - 冰雹回复HP
- 蜕皮 (Shed Skin) - 30% 治愈异常状态
- 毒疗 (Poison Heal) - 中毒时回复HP
- 太阳之力 (Solar Power) - 晴天特攻×1.5但损失HP

**速度修正类 (4个)**
- 悠游自如 (Swift Swim) - 雨天速度×2
- 叶绿素 (Chlorophyll) - 晴天速度×2
- 拨沙 (Sand Rush) - 沙暴速度×2
- 拨雪 (Slush Rush) - 冰雹速度×2

**优先度修正类 (2个)**
- 恶作剧之心 (Prankster) - 变化技能优先度+1
- 疾风之翼 (Gale Wings) - 满血时飞行技能优先度+1

**击倒触发类 (3个)**
- 自信过剩 (Moxie) - 击倒对手后攻击+1
- 异兽提升 (Beast Boost) - 击倒对手后最高能力+1
- 魂心 (Soul-Heart) - 场上宝可梦濒死时特攻+1

**形态变化类 (4个)**
- 羁绊变身 (Battle Bond) - 甲贺忍蛙击倒对手后变为小智版
- 达摩模式 (Zen Mode) - 达摩狒狒 HP≤50% 时变为达摩模式
- 群聚变形 (Power Construct) - 基格尔德 HP≤50% 时变为完全体
- 战斗切换 (Stance Change) - 坚盾剑怪根据技能类型切换剑/盾形态

#### 支持的机制
- **性格系统**: 25 种性格，影响能力值 ±10%
- **特性系统**: 完整特性效果实现，支持多种触发时机
- **能力等级**: -6 到 +6 阶段变化
- **异常状态**: 中毒、剧毒、灼伤、麻痹、睡眠、冰冻
- **临时状态**: 混乱、着迷、挑衅、定身法、寄生种子、替身等
- **道具效果**: 讲究头带/眼镜/围巾、生命宝珠、气势披带等
- **技能优先度**: -7 到 +5
- **充能技能**: 破坏光线、终极冲击等
- **队伍系统**: 3v3/6v6 模式支持换人

#### AI 对战系统
- 支持人机对战模式
- AI 自动从热门宝可梦中选择队伍
- AI 自动选择技能进行对战

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
- 特性注册表使用 `sync.Once` 单例模式

### 接口设计

- 使用接口避免循环依赖（如 `ability.Battler` 接口）
- 适配器模式连接不同层（如 `MoveAdapter`）
- 注册表模式管理特性效果

---

## 扩展计划

### 项目进度

**已完成**
- ✅ UNO 卡牌游戏完整实现
- ✅ 宝可梦对战核心系统（1v1/3v3/6v6）
- ✅ 完整伤害计算公式与属性克制系统
- ✅ 74 个特性效果实现（包括形态变化）
- ✅ 性格系统、能力等级、异常状态
- ✅ 道具效果、天气系统基础设施
- ✅ AI 对战系统
- ✅ 预设系统

**进行中**
- 🔄 更多特性效果实现（74/~270 已完成）
- 🔄 形态变化特性扩展

### LLM 集成（预留）

`internal/domain/shared/llm/client.go` 定义了 LLM 客户端接口：

```go
type Client interface {
    Chat(ctx context.Context, messages []Message) (string, error)
    ChatStream(ctx context.Context, messages []Message) (<-chan string, error)
    GetModel() string
}
```

可用于未来添加 AI 对战解说、智能策略建议等功能。

### 待实现功能

**核心功能**
- [ ] 持久化存储（Redis/数据库）
- [ ] 游戏统计与排行榜
- [ ] 更多小游戏
- [ ] 智能 AI 玩家（基于 LLM）
- [ ] 游戏房间管理

**宝可梦对战增强**
- [ ] 双打对战模式
- [ ] 超级进化 / 极巨化 / 太晶化
- [ ] 更多特性效果实现（约 200 个待实现，详见 `assets/pokemon/pending_abilities.md`）
- [ ] 会心相关特性（战斗盔甲、硬壳盔甲、超幸运等）
- [ ] 命中/闪避相关特性（复眼、沙隐、雪隐等）
- [ ] 追加效果相关特性（天恩、鳞粉等）
- [ ] 更多形态变化特性（预知梦、花之礼等）

---

## 资源文件

### UNO 卡牌图片
位于 `assets/uno/` 目录，命名规则：
- 数字牌: `{Color}{Number}.jpg` (如 `Red5.jpg`)
- 功能牌: `{Color}{Type}.jpg` (如 `BlueSkip.jpg`)
- 万能牌: `Wild.jpg`, `WildDraw.jpg`

### 宝可梦数据
- `assets/pokemon/abilities.json`: 特性数据（辅助）
- `assets/pokemon/pending_abilities.md`: 待实现特性列表（约 200 个）
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
7. **特性系统**: 新增特性效果需在 `registry.go` 的 `registerAllEffects` 中注册
8. **接口适配**: 添加新的 Battler/Move 方法时需同步更新 `battler_adapter.go`
9. **形态变化**: 形态变化特性需实现 `OnFormChange` 方法，返回 `FormChangeResult` 结构体
10. **特性分类**: 特性按文件分类存放（`effects_calc.go`、`effects_entry.go`、`effects_formchange.go` 等），便于维护

---

## 🎴 三国杀 / 无名杀 (Discord Activity)

> ⚠️ **重要提示：此功能目前无法正常运行！**

### 功能介绍

本项目计划集成 [无名杀](https://github.com/libccy/noname) 作为 Discord Activity，让用户可以直接在 Discord 语音频道中游玩三国杀。

无名杀是一款开源的三国杀网页版游戏，支持：
- 🎮 多种游戏模式（身份局、国战、欢乐成双等）
- 🃏 丰富的武将和卡牌扩展包
- 👥 多人在线对战
- 🎨 精美的界面和动画效果

### 配置说明

在 `config.yaml` 中配置 Activity：

```yaml
activity:
  enabled: false  # 是否启用 Activity 功能
  client_id: "your-application-id"  # Discord Application ID
  client_secret: "your-client-secret"  # Discord Application Client Secret
  port: 8080  # Activity Web 服务器端口
  public_url: ""  # 公网访问地址 (如: https://your-domain.com)
  game_path: "./noname"  # 无名杀游戏文件路径
  dev_mode: true  # 开发模式：自动启动 Vite 开发服务器并代理请求
  vite_port: 5173  # Vite 开发服务器端口
```

### 相关文件

- `noname/` - 无名杀游戏源代码目录
- `internal/infrastructure/activity/` - Activity 基础设施层
- `internal/interfaces/discord/commands/activity_commands.go` - Activity 命令处理

---

## ⚠️ 需要帮助：三国杀功能修复

<table>
<tr>
<td>

### 🚨 当前状态：无法正常运行

三国杀 Discord Activity 功能目前**因不明原因无法正常工作**。

我们已经尝试了多种方法，但仍未能成功让无名杀在 Discord Activity 环境中正常运行。

### 可能的问题方向

- Discord Embedded App SDK 集成问题
- 无名杀在 iframe 环境中的兼容性
- WebSocket 连接或网络代理配置
- Activity 服务器与游戏客户端的通信
- 跨域资源共享 (CORS) 配置
- 其他未知的技术限制

### 🙏 寻求社区帮助

如果您有以下经验，我们非常期待您的帮助：

- Discord Activity / Embedded App SDK 开发经验
- 无名杀源码修改或部署经验
- Web 游戏嵌入式环境适配经验
- 类似项目的成功案例

### 如何贡献

1. Fork 本仓库
2. 尝试修复 Activity 功能
3. 提交 Pull Request
4. 或在 Issues 中分享您的发现和建议

**任何形式的帮助都非常感谢！** 🙏

</td>
</tr>
</table>

---

## 许可证

本项目仅供学习交流使用。

- UNO 是 Mattel 公司的注册商标
- 宝可梦是 Nintendo/Creatures Inc./GAME FREAK inc. 的注册商标
- 三国杀是游卡桌游的注册商标
- 无名杀遵循其原项目的开源协议
