# HostManager - 企业级SSH连接管理工具

🖥️ **简洁、强大、可维护** 的SSH连接管理工具，遵循世界级代码架构标准，专为iTerm2和终端用户设计。

## ✨ 核心特性

### 🏗️ 架构设计
- **简洁原则**: 遵循KISS和YAGNI原则，避免过度工程化
- **可读性优先**: 代码即文档，直观易懂的逻辑结构
- **模块化设计**: 清晰的责任分离和低耦合架构

### 🎯 功能特性
- **双模式支持**: 全屏UI界面 + 命令行工具
- **分组管理**: 按环境、项目等维度组织主机
- **收藏功能**: 快速访问常用主机
- **智能搜索**: 按名称、IP、用户名搜索
- **状态监控**: 批量检查主机在线状态
- **快速连接**: 一键连接常用主机
- **连接历史**: 自动记录连接记录
- **主题切换**: 支持明亮/暗色主题
- **iTerm2集成**: 完美配合iTerm2使用

## 🚀 快速开始

### 系统要求
- Go 1.24.4+
- macOS/Linux/Windows
- 终端支持（推荐iTerm2）

### 安装

```bash
# 克隆项目
git clone https://github.com/daihao4371/hostmanager.git
cd hostmanager

# 安装依赖（自动处理）
go mod download

# 编译
go build .

# 安装到系统（推荐）
chmod +x install.sh
./install.sh
```

### 基本使用

```bash
# 启动全屏UI界面
./hostmanager

# 显示帮助信息
./hostmanager help

# 列出所有主机
./hostmanager list

# 按分组显示
./hostmanager list --groups

# 连接到主机
./hostmanager connect server1

# 搜索主机
./hostmanager search web

# 检查主机状态
./hostmanager status server1
```

## 📋 命令行接口

### 主要命令

| 命令 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `connect` | `c` | 连接到指定主机 | `hostmanager connect server1` |
| `list` | `ls`, `l` | 显示主机列表 | `hostmanager list --groups` |
| `status` | `s` | 检查主机状态 | `hostmanager status server1` |
| `search` | - | 搜索主机 | `hostmanager search web` |
| `favorites` | `fav`, `f` | 显示收藏夹 | `hostmanager favorites` |
| `groups` | `g` | 按分组显示 | `hostmanager groups` |
| `history` | `h` | 显示连接历史 | `hostmanager history` |
| `help` | `--help`, `-h` | 显示帮助 | `hostmanager help` |
| `version` | `--version`, `-v` | 显示版本 | `hostmanager version` |

### 便捷别名（安装后可用）

```bash
hm                    # = hostmanager
hml                   # = hostmanager list
hmg                   # = hostmanager list --groups
hmf                   # = hostmanager favorites
hms server1           # = hostmanager status server1
hmc server1           # = hostmanager connect server1
hm-connect server1    # 智能连接函数
hm-search web         # 智能搜索函数
```

## 🎮 全屏UI界面

无参数运行时进入全屏交互界面：

```bash
./hostmanager
```

### UI快捷键

- `↑↓` : 导航
- `Enter` : 选择/连接
- `Esc` : 返回/退出
- `Space` : 切换收藏
- `f` : 收藏夹模式
- `s` : 状态检查
- `t` : 切换主题
- `l` : 切换布局
- `/` : 搜索模式
- `1-5` : 快速连接历史记录
- `q` : 退出

## 🔧 iTerm2 集成

### 推荐设置

1. **创建快捷键**：
   - `Cmd+Shift+H`: 运行 `hostmanager list`
   - `Cmd+Shift+C`: 运行 `hostmanager connect`

2. **Profile设置**：
   - 为不同环境设置不同的Profile
   - 配置不同的颜色主题区分环境

3. **Split Pane使用**：
   ```bash
   # 在不同面板中同时连接多台服务器
   hostmanager connect server1  # 左侧面板
   hostmanager connect server2  # 右侧面板
   ```

### 高级集成

```bash
# 在iTerm2中创建自定义函数（需要安装fzf）
function quick-ssh() {
    local host=$(hostmanager list | fzf --header="选择要连接的主机" | awk '{print $2}')
    if [[ -n "$host" ]]; then
        hostmanager connect "$host"
    fi
}

# 快速状态检查
function check-servers() {
    hostmanager status | grep -E "(🔴|❓)" && echo "发现离线或未知状态的服务器！"
}
```

## ⚙️ 配置文件

配置文件位置：`config.yaml`

```yaml
groups:
  - name: "生产环境"
    hosts:
      - name: "Web服务器-1"
        ip: "192.168.1.100"
        port: 22
        username: "admin"
        password: "your_password"
        favorite: true
        
  - name: "测试环境"
    hosts:
      - name: "测试服务器"
        ip: "192.168.1.200"
        port: 22
        username: "test"
        auth_type: "key"
        key_path: "~/.ssh/id_rsa"

ui_config:
  theme: "dark"  # "light" or "dark"
  language: "zh-CN"
  layout:
    type: "single"  # "single" or "columns"
```

## 🛠️ 开发指南

### 架构原则

本项目严格遵循以下架构原则：

1. **简洁性**: 每个模块只负责单一职责
2. **可读性**: 代码自文档化，逻辑清晰直观
3. **可维护性**: 低耦合高内聚，易于扩展

### 项目结构

```
hostmanager/
├── main.go                 # 应用入口点（CLI/UI路由）
├── config.yaml            # 配置文件示例
├── config.example.yaml    # 配置模板文件
├── install.sh             # iTerm2集成安装脚本
├── go.mod                 # Go模块依赖管理
├── internal/              # 内部包（不对外暴露）
│   ├── cli/               # 命令行接口层
│   │   └── cli.go         # CLI命令处理和路由
│   ├── config/            # 配置管理模块
│   │   └── config.go      # 配置文件解析和验证
│   ├── models/            # 数据模型层
│   │   └── host.go        # 主机数据结构定义
│   ├── ssh/               # SSH连接核心逻辑
│   │   └── connection.go  # SSH连接实现
│   ├── theme/             # 主题管理系统
│   │   └── theme.go       # 主题配置和切换逻辑
│   ├── i18n/              # 国际化支持
│   │   └── texts.go       # 多语言文本管理
│   └── ui/                # 用户界面层
│       ├── menu.go        # 菜单状态管理
│       ├── render.go      # 渲染引擎核心
│       ├── components.go  # 可复用UI组件
│       ├── input.go       # 用户输入处理
│       ├── interaction.go # 用户交互逻辑
│       ├── layout.go      # 布局管理系统
│       └── draw.go        # 底层绘制功能
└── README.md              # 项目文档
```

### 代码质量标准

- **函数复杂度**: 每个函数不超过30行
- **类大小限制**: 每个结构体不超过300行
- **嵌套控制**: 逻辑嵌套深度不超过3层
- **参数限制**: 函数参数不超过4个
- **错误处理**: 所有可能失败的操作都有明确的错误处理

### 构建和测试

```bash
# 开发环境构建
go build -v .

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...

# 交叉编译发布版本
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/hostmanager-linux .
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/hostmanager-macos .
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/hostmanager-windows.exe .
```

### 开发规范

遵循项目的 `CLAUDE.md` 中定义的架构师级别代码标准：

- 使用 Guard Clauses 减少嵌套
- 函数职责单一，易于测试
- 错误处理显式且完整
- 变量命名清晰描述用途
- 代码注释解释"为什么"而非"做什么"

## 🤝 贡献指南

### 代码质量要求

在提交代码前，请确保：

1. **遵循架构原则**: 查看 `CLAUDE.md` 了解详细规范
2. **通过所有检查**: 运行 `go test ./... && go vet ./... && go fmt ./...`
3. **保持简洁性**: 函数不超过30行，嵌套不超过3层
4. **错误处理**: 所有可能失败的操作都要有明确的错误处理

### 提交流程

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 遵循代码规范提交代码
4. 提交Pull Request

## 📄 许可证

MIT License - 详见 LICENSE 文件