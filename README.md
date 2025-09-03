# HostManager - SSH连接管理工具

🖥️ 为iTerm2和终端用户优化的SSH连接管理工具，提供直观的UI界面和强大的命令行功能。

## ✨ 主要特性

- 🎯 **双模式支持**: 全屏UI界面 + 命令行工具
- 📁 **分组管理**: 按环境、项目等维度组织主机
- ⭐ **收藏功能**: 快速访问常用主机
- 🔍 **智能搜索**: 按名称、IP、用户名搜索
- 📊 **状态监控**: 批量检查主机在线状态
- 🚀 **快速连接**: 一键连接常用主机
- 📝 **连接历史**: 自动记录连接记录
- 🎨 **主题切换**: 支持明亮/暗色主题
- 🔧 **iTerm2集成**: 完美配合iTerm2使用

## 🚀 快速开始

### 安装

```bash
# 克隆项目
git clone https://github.com/daihao4371/hostmanager.git
cd hostmanager

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

## 🛠️ 开发

### 项目结构

```
hostmanager/
├── main.go                 # 主入口（支持CLI/UI双模式）
├── config.yaml            # 配置文件示例
├── install.sh             # iTerm2集成安装脚本
├── internal/
│   ├── cli/               # 命令行接口
│   │   └── cli.go         # CLI命令处理
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── ssh/              # SSH连接
│   ├── theme/            # 主题管理
│   ├── i18n/             # 国际化
│   └── ui/               # UI界面
│       ├── menu.go       # 菜单管理（增强版）
│       ├── render.go     # 渲染引擎（高级）
│       ├── components.go # UI组件（增强）
│       ├── input.go      # 输入处理（修复）
│       ├── interaction.go # 交互管理
│       └── draw.go       # 绘制功能
└── README.md             # 项目文档
```

### 构建

```bash
# 开发构建
go build .

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o hostmanager-linux .
GOOS=darwin GOARCH=amd64 go build -o hostmanager-macos .
GOOS=windows GOARCH=amd64 go build -o hostmanager-windows.exe .
```

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License