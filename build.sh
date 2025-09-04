#!/bin/bash

# HostManager 跨平台构建脚本
# 用法: ./build.sh [版本号]

set -e

VERSION=${1:-"dev"}
OUTPUT_DIR="dist"

echo "🚀 开始构建 HostManager ${VERSION}..."

# 清理输出目录
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"

# 定义构建目标
declare -a targets=(
    "linux/amd64"
    "linux/arm64" 
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# 构建每个目标平台
for target in "${targets[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "${target}"
    echo "📦 构建 ${GOOS}/${GOARCH}..."
    
    # 设置二进制文件名
    BINARY_NAME="hostmanager"
    if [ "${GOOS}" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    # 构建
    env GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 \
        go build -ldflags="-s -w -X main.Version=${VERSION}" -o "${BINARY_NAME}" .
    
    # 创建发布包目录
    PACKAGE_NAME="hostmanager-${VERSION}-${GOOS}-${GOARCH}"
    PACKAGE_DIR="${OUTPUT_DIR}/${PACKAGE_NAME}"
    mkdir -p "${PACKAGE_DIR}"
    
    # 复制文件
    cp "${BINARY_NAME}" "${PACKAGE_DIR}/"
    cp config.example.yaml "${PACKAGE_DIR}/config.yaml"
    cp install.sh "${PACKAGE_DIR}/"
    cp install-global.sh "${PACKAGE_DIR}/"
    cp completion.bash "${PACKAGE_DIR}/"
    cp completion.zsh "${PACKAGE_DIR}/"
    cp README.md "${PACKAGE_DIR}/"
    
    # 创建安装说明
    cat > "${PACKAGE_DIR}/INSTALL.txt" << EOF
HostManager ${VERSION} 安装说明

1. 解压文件到任意目录
2. 编辑 config.yaml 配置文件，添加你的服务器信息
3. 运行以下命令进行全局安装:
   ./install-global.sh

4. 安装完成后，在任意目录使用 'hm' 命令

使用方法:
- hm              启动交互式界面
- hm list         列出所有主机  
- hm connect 主机名 连接指定主机

Zmodem 文件传输:
连接后可使用 sz/rz 命令传输文件
需要先安装 lrzsz: 
- macOS: brew install lrzsz
- Ubuntu: apt install lrzsz
- CentOS: yum install lrzsz
EOF
    
    # 创建压缩包
    cd "${OUTPUT_DIR}"
    if [ "${GOOS}" = "windows" ]; then
        zip -r "${PACKAGE_NAME}.zip" "${PACKAGE_NAME}"
        echo "✅ 已创建: ${PACKAGE_NAME}.zip"
    else
        tar -czf "${PACKAGE_NAME}.tar.gz" "${PACKAGE_NAME}"
        echo "✅ 已创建: ${PACKAGE_NAME}.tar.gz"
    fi
    cd ..
    
    # 清理二进制文件
    rm -f "${BINARY_NAME}"
    rm -rf "${PACKAGE_DIR}"
done

echo ""
echo "🎉 构建完成！输出目录: ${OUTPUT_DIR}"
echo ""
echo "发布包:"
ls -la "${OUTPUT_DIR}"/*.{tar.gz,zip} 2>/dev/null || true
echo ""
echo "使用方法："
echo "1. 上传发布包到 GitHub Releases"
echo "2. 或者推送 tag 触发自动发布: git tag v${VERSION} && git push origin v${VERSION}"