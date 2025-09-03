#!/bin/bash
# HostManager iTerm2集成安装脚本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$SCRIPT_DIR/hostmanager"
INSTALL_DIR="/usr/local/bin"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}🖥️  HostManager iTerm2 集成安装器${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

install_binary() {
    print_info "正在安装HostManager到 $INSTALL_DIR ..."
    
    if [ ! -f "$BINARY_PATH" ]; then
        print_error "未找到可执行文件: $BINARY_PATH"
        print_info "请先运行 'go build .' 编译程序"
        return 1
    fi
    
    # 检查是否有写权限
    if [ ! -w "$INSTALL_DIR" ]; then
        print_info "需要管理员权限来安装到 $INSTALL_DIR"
        sudo cp "$BINARY_PATH" "$INSTALL_DIR/hostmanager"
        sudo chmod +x "$INSTALL_DIR/hostmanager"
    else
        cp "$BINARY_PATH" "$INSTALL_DIR/hostmanager"
        chmod +x "$INSTALL_DIR/hostmanager"
    fi
    
    if [ $? -eq 0 ]; then
        print_success "HostManager已安装到 $INSTALL_DIR/hostmanager"
        return 0
    else
        print_error "安装失败"
        return 1
    fi
}

setup_aliases() {
    print_info "设置便捷别名..."
    
    SHELL_RC=""
    if [ "$SHELL" = "/bin/zsh" ] || [ "$SHELL" = "/usr/bin/zsh" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ "$SHELL" = "/bin/bash" ] || [ "$SHELL" = "/usr/bin/bash" ]; then
        SHELL_RC="$HOME/.bashrc"
    else
        print_warning "未识别的shell: $SHELL，请手动设置别名"
        return 1
    fi
    
    # 备份原配置文件
    if [ -f "$SHELL_RC" ]; then
        cp "$SHELL_RC" "$SHELL_RC.backup.$(date +%Y%m%d_%H%M%S)"
    fi
    
    # 添加别名到配置文件
    cat >> "$SHELL_RC" << 'EOF'

# HostManager 别名配置
alias hm='hostmanager'
alias hml='hostmanager list'
alias hmg='hostmanager list --groups'
alias hmf='hostmanager favorites'
alias hms='hostmanager status'
alias hmc='hostmanager connect'
alias hmh='hostmanager history'

# 快捷函数
hm-connect() {
    if [ $# -eq 0 ]; then
        hostmanager list
        echo "用法: hm-connect <主机名>"
    else
        hostmanager connect "$1"
    fi
}

hm-search() {
    if [ $# -eq 0 ]; then
        echo "用法: hm-search <关键词>"
    else
        hostmanager search "$@"
    fi
}

# 自动补全支持（基础版）
_hostmanager_complete() {
    local cur commands
    cur="${COMP_WORDS[COMP_CWORD]}"
    commands="connect list status history favorites groups search help version"
    
    if [ $COMP_CWORD -eq 1 ]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
    fi
}

complete -F _hostmanager_complete hostmanager hm

EOF

    print_success "别名已添加到 $SHELL_RC"
    print_info "运行 'source $SHELL_RC' 或重启终端来启用别名"
}

show_usage_examples() {
    echo ""
    print_info "🚀 使用示例:"
    echo ""
    echo -e "${GREEN}# 基本命令${NC}"
    echo "  hostmanager                  # 启动全屏UI界面"
    echo "  hostmanager help             # 显示帮助信息"
    echo "  hostmanager list             # 显示所有主机"
    echo "  hostmanager list --groups    # 按分组显示主机"
    echo ""
    echo -e "${GREEN}# 连接管理${NC}"
    echo "  hostmanager connect server1  # 连接到server1"
    echo "  hostmanager status server1   # 检查server1状态"
    echo "  hostmanager search web       # 搜索包含'web'的主机"
    echo ""
    echo -e "${GREEN}# 便捷别名（安装后可用）${NC}"
    echo "  hm                          # = hostmanager"
    echo "  hml                         # = hostmanager list"
    echo "  hmg                         # = hostmanager list --groups"
    echo "  hmc server1                 # = hostmanager connect server1"
    echo "  hm-connect server1          # 智能连接函数"
    echo "  hm-search web               # 智能搜索函数"
    echo ""
    echo -e "${GREEN}# iTerm2 集成技巧${NC}"
    echo "  - 在iTerm2中设置快捷键调用 'hostmanager list'"
    echo "  - 使用iTerm2的Split Pane功能同时管理多个SSH会话"
    echo "  - 结合iTerm2的Profile功能为不同环境设置不同的主题"
}

main() {
    print_header
    
    echo -e "这个脚本将会:"
    echo -e "  1. 安装 hostmanager 到系统路径"
    echo -e "  2. 设置便捷的别名和函数"
    echo -e "  3. 配置自动补全"
    echo ""
    
    read -p "继续安装吗？(y/N): " -r
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "安装已取消"
        exit 0
    fi
    
    # 安装二进制文件
    if ! install_binary; then
        exit 1
    fi
    
    # 设置别名
    setup_aliases
    
    print_success "🎉 HostManager iTerm2集成安装完成！"
    show_usage_examples
    
    echo ""
    print_warning "请运行以下命令来启用别名："
    if [ "$SHELL" = "/bin/zsh" ] || [ "$SHELL" = "/usr/bin/zsh" ]; then
        echo -e "${YELLOW}  source ~/.zshrc${NC}"
    else
        echo -e "${YELLOW}  source ~/.bashrc${NC}"
    fi
}

# 运行主函数
main "$@"