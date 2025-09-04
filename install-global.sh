#!/bin/bash

# HostManager 全局安装脚本
# 用法: ./install-global.sh

set -e

# 定义路径
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/hostmanager"
CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🔧 开始安装 HostManager..."

# 创建配置目录
echo "📁 创建配置目录..."
mkdir -p "$CONFIG_DIR"

# 复制配置文件
echo "📋 复制配置文件..."
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    cp "$CURRENT_DIR/config.yaml" "$CONFIG_DIR/config.yaml"
    echo "✅ 配置文件已复制到 $CONFIG_DIR/config.yaml"
else
    echo "ℹ️  配置文件已存在，跳过复制"
fi

# 复制可执行文件到系统路径
echo "📦 安装可执行文件..."
sudo cp "$CURRENT_DIR/hostmanager" "$INSTALL_DIR/hm"
sudo chmod +x "$INSTALL_DIR/hm"

# 更新配置文件查找路径
echo "🔄 更新配置文件路径..."
cat > /tmp/hm_wrapper.sh << 'EOF'
#!/bin/bash
# HostManager 启动脚本

CONFIG_PATHS=(
    "$HOME/.config/hostmanager/config.yaml"
    "$HOME/config.yaml" 
    "./config.yaml"
)

# 查找配置文件
CONFIG_FILE=""
for path in "${CONFIG_PATHS[@]}"; do
    if [ -f "$path" ]; then
        CONFIG_FILE="$path"
        break
    fi
done

if [ -z "$CONFIG_FILE" ]; then
    echo "❌ 错误: 找不到配置文件"
    echo "请确保以下位置之一存在 config.yaml:"
    printf "   %s\n" "${CONFIG_PATHS[@]}"
    exit 1
fi

# 切换到配置文件目录执行程序
cd "$(dirname "$CONFIG_FILE")"
exec /usr/local/bin/hostmanager "$@"
EOF

sudo mv /tmp/hm_wrapper.sh "$INSTALL_DIR/hm"
sudo chmod +x "$INSTALL_DIR/hm"

echo ""
echo "🎉 安装完成！"
echo ""
echo "现在你可以在任何位置使用 'hm' 命令"
echo "配置文件位置: $CONFIG_DIR/config.yaml"
echo ""
echo "使用方法:"
echo "  hm              - 启动交互式界面"
echo "  hm list         - 列出所有主机"
echo "  hm connect 主机名 - 连接指定主机"
echo ""