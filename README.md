# HostManager - macOS iTerm2 SSH会话管理工具

🍎 **专为macOS iTerm2设计**的SSH会话管理工具，让你在Mac上更高效地管理和连接SSH服务器。

## ✨ 为什么选择HostManager？

### 🍎 macOS原生体验
- **iTerm2深度集成**: 专为iTerm2优化的SSH会话管理
- **原生快捷键**: 支持macOS键盘快捷键习惯
- **Spotlight集成**: 可通过Spotlight快速启动
- **多窗口支持**: 完美适配iTerm2的多标签和分屏

### 🚀 SSH会话管理
- **一键连接**: 告别复杂的SSH命令记忆
- **会话分组**: 按项目、环境智能分类SSH连接
- **连接历史**: 自动记录最近使用的SSH会话
- **状态监控**: 实时检查服务器连接状态
- **收藏夹**: 快速访问常用服务器
- **智能搜索**: 按服务器名称、IP快速定位
- **双界面**: 图形化菜单 + 命令行，适合不同使用场景
- **📁 Zmodem 支持**: 内置 sz/rz 文件传输功能

### 📁 文件传输功能
- **sz/rz 命令**: 连接后直接使用 sz 发送文件，rz 接收文件
- **自动检测**: 智能检测系统是否安装 lrzsz 工具
- **无需配置**: 默认为所有主机启用 Zmodem 支持
- **跨平台**: 支持 macOS、Linux、Windows 客户端

## 🚀 快速开始

### 系统要求
- **macOS 10.15+** (专为macOS设计)
- **iTerm2 3.0+** (必需，最佳体验)
- **Go 1.24.4+** (仅构建时需要)

### 安装

#### 方法 1: 从 GitHub Releases 下载（推荐）

1. 从 [Releases 页面](https://github.com/daihao4371/hostmanager/releases) 下载适合你系统的预编译版本
2. 解压到任意目录
3. 运行安装脚本进行全局安装：

```bash
# 解压下载的文件
tar -xzf hostmanager-v1.0.2-darwin-amd64.tar.gz  # macOS Intel
# 或
tar -xzf hostmanager-v1.0.2-darwin-arm64.tar.gz  # macOS Apple Silicon

cd hostmanager-v1.0.2-*
chmod +x install-global.sh
./install-global.sh
```

安装完成后可在任意位置使用 `hm` 命令。

#### 方法 2: 使用 Go 安装

```bash
go install github.com/daihao4371/hostmanager@latest
```

#### 方法 3: 源码编译

```bash
# 克隆项目
git clone https://github.com/daihao4371/hostmanager.git
cd hostmanager

# 编译
go build -o hostmanager .

# 全局安装
./install-global.sh
```

### 在iTerm2中使用

```bash
# 启动全屏SSH会话选择界面
./hostmanager

# 快速连接（支持tab补全）
./hostmanager connect server1

# 查看所有SSH会话
./hostmanager list

# 按环境分组查看
./hostmanager list --groups

# 搜索特定服务器
./hostmanager search web

# 检查服务器在线状态  
./hostmanager status server1
```

### 📁 文件传输使用

连接到服务器后，可直接使用 sz/rz 命令：

```bash
# 发送文件到本地
sz filename.txt
sz *.log        # 发送多个文件
sz /path/to/dir/* # 发送目录下所有文件

# 从本地接收文件到远程
rz              # 会弹出文件选择对话框
```

**前提条件**: 需要安装 lrzsz 工具包
```bash
# macOS
brew install lrzsz

# Ubuntu/Debian
sudo apt install lrzsz

# CentOS/RHEL
sudo yum install lrzsz
```

## 📋 SSH会话管理命令

### 核心命令

| 命令 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `connect` | `c` | 在iTerm2中连接SSH会话 | `hostmanager connect server1` |
| `list` | `ls`, `l` | 显示SSH会话列表 | `hostmanager list --groups` |
| `status` | `s` | 检查服务器连接状态 | `hostmanager status server1` |
| `search` | - | 搜索SSH会话 | `hostmanager search web` |
| `favorites` | `fav`, `f` | 显示收藏的会话 | `hostmanager favorites` |
| `groups` | `g` | 按项目环境分组显示 | `hostmanager groups` |
| `history` | `h` | 显示SSH连接历史 | `hostmanager history` |
| `add-host` | - | 添加新的SSH会话 | `hostmanager add-host` |
| `edit` | - | 编辑SSH会话配置 | `hostmanager edit server1` |
| `remove` | `rm` | 删除SSH会话 | `hostmanager remove server1` |
| `init` | - | 初始化配置文件 | `hostmanager init` |
| `help` | `--help`, `-h` | 显示帮助 | `hostmanager help` |
| `version` | `--version`, `-v` | 显示版本 | `hostmanager version` |

### iTerm2快捷别名

```bash
hm                    # = hostmanager
hml                   # = hostmanager list  
hmg                   # = hostmanager list --groups
hmf                   # = hostmanager favorites
hms server1           # = hostmanager status server1  
hmc server1           # = hostmanager connect server1
hm-connect server1    # iTerm2智能连接函数
hm-search web         # iTerm2智能搜索函数
```

## 🎮 iTerm2全屏会话管理界面

在iTerm2中无参数运行，启动专用的SSH会话管理界面：

```bash
./hostmanager
```

### 专为iTerm2优化的快捷键

- `↑↓` : 导航SSH会话列表
- `Enter` : 在iTerm2新标签中连接SSH会话
- `Esc` : 返回/退出会话管理界面  
- `Space` : 切换SSH会话收藏状态
- `f` : 显示收藏的SSH会话
- `s` : 批量检查服务器状态
- `t` : 切换iTerm2主题（明亮/暗色）
- `l` : 切换显示布局
- `/` : 搜索SSH会话
- `1-5` : 快速连接最近5个SSH会话
- `q` : 退出HostManager

## 🔧 iTerm2深度集成配置

### macOS系统集成

1. **iTerm2全局快捷键**：
   - `Cmd+Shift+H`: 启动HostManager会话管理
   - `Cmd+Shift+S`: 快速SSH连接菜单

2. **iTerm2 Profile配置**：
   - 为生产环境设置红色背景提醒
   - 为测试环境设置绿色背景区分
   - 配置不同的字体和透明度

3. **多窗口SSH会话管理**：
   ```bash
   # 在不同iTerm2面板中同时管理多个SSH会话
   hostmanager connect prod-server1  # 左侧面板
   hostmanager connect test-server2   # 右侧面板
   hostmanager connect db-server3     # 底部面板
   ```

### 高级工作流集成

```bash
# 在iTerm2中创建智能SSH连接函数（需要安装fzf）
function quick-ssh() {
    local host=$(hostmanager list | fzf --header="选择SSH会话连接" | awk '{print $2}')
    if [[ -n "$host" ]]; then
        hostmanager connect "$host"
    fi
}

# 批量服务器状态检查
function check-all-servers() {
    hostmanager status | grep -E "(🔴|❓)" && echo "⚠️ 发现离线或异常的服务器！"
}

# 在iTerm2中创建SSH会话组快速启动
function start-dev-env() {
    osascript -e 'tell app "iTerm2" to create window with default profile'
    hostmanager connect dev-web &
    sleep 1
    hostmanager connect dev-db &  
    sleep 1
    hostmanager connect dev-cache &
}
```

## ⚙️ SSH会话配置文件

专为macOS用户设计的SSH会话配置：`config.yaml`

```yaml
groups:
  - name: "生产环境 🔴"
    hosts:
      - name: "Web服务器-1"
        ip: "192.168.1.100" 
        port: 22
        username: "admin"
        auth_type: "key"
        key_path: "~/.ssh/id_rsa"
        description: "主要的Web服务器"
        tags: ["production", "web"]
        favorite: true
        
  - name: "开发环境 🟢" 
    hosts:
      - name: "开发服务器"
        ip: "192.168.1.200"
        port: 22
        username: "dev"
        auth_type: "key"  
        key_path: "~/.ssh/id_rsa"
        description: "开发测试服务器"
        tags: ["development"]

# iTerm2专用界面配置
ui_config:
  theme: "dark"  # 适配iTerm2暗色主题
  language: "zh-CN"
  key_bindings:
    exit: "Esc"
    search: "/"
    favorites: "f"
    status_check: "s"
    toggle_fav: "Space"
    theme_switch: "t"
  layout:
    type: "single"
    show_details: true  # 在iTerm2中显示详细信息
```

## 🛠️ macOS开发者指南

### 项目目标
专注于为macOS开发者提供最佳的iTerm2 SSH会话管理体验。

### 设计理念

macOS原生应用的设计原则：

1. **简洁优雅**: 符合Apple Human Interface Guidelines
2. **直觉操作**: 遵循macOS用户使用习惯
3. **深度集成**: 与iTerm2无缝配合
4. **高性能**: 针对macOS系统优化

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