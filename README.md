# HostManager - SSH连接管理器

## 项目结构

```
HostManager/
├── main.go                   # 程序入口 (28行)
├── config.yaml              # 配置文件
├── internal/                # 内部包
│   ├── models/              # 数据模型 (21行)
│   │   └── host.go         # 主机和分组结构定义
│   ├── config/              # 配置管理 (151行)
│   │   └── config.go       # 配置加载、保存和默认值设置
│   ├── theme/               # 主题管理 (70行)
│   │   └── theme.go        # 主题定义和管理
│   ├── i18n/                # 国际化 (76行)
│   │   └── texts.go        # 多语言文本资源
│   ├── ssh/                 # SSH连接 (139行)
│   │   └── connection.go   # SSH连接逻辑和工具
│   └── ui/                  # 用户界面 (671行)
│       ├── menu.go         # 菜单管理和业务逻辑 (216行)
│       ├── draw.go         # 基础绘制功能 (174行)
│       ├── layout.go       # 布局和显示组件 (247行)
│       └── input.go        # 输入处理和快捷键 (234行)
└── README.md                # 项目文档
```

## 架构设计优势

### 模块化设计
- **单一职责**: 每个包专注于特定功能领域
- **低耦合**: 包之间通过接口和数据结构交互
- **高内聚**: 相关功能集中在同一包内

### 代码可维护性
- **main.go**: 从1250+行精简到28行，只负责程序启动
- **功能分离**: UI、配置、SSH、主题等功能完全分离
- **易于测试**: 每个包可以独立测试

### 扩展性
- **新功能**: 可以轻松添加新的包和功能模块
- **新主题**: 在theme包中添加新的主题方案
- **新语言**: 在i18n包中添加新的语言支持

## 用户体验功能

### 1. 主题切换
- **快捷键**: `t` - 在深色/浅色主题之间切换
- **配置**: 在 `config.yaml` 中的 `ui_config.theme` 设置默认主题
- **支持主题**: `dark`（深色）, `light`（浅色）

### 2. 快捷键自定义
所有快捷键都可以在 `config.yaml` 的 `ui_config.key_bindings` 中自定义：

```yaml
key_bindings:
  exit: "Esc"         # 退出/返回
  search: "/"         # 搜索
  favorites: "f"      # 收藏夹
  status_check: "s"   # 状态检查
  reload: "r"         # 重载配置
  toggle_fav: "Space" # 切换收藏
  theme_switch: "t"   # 主题切换
  layout_switch: "l"  # 布局切换
```

### 3. 布局优化
- **单栏布局**: 传统的垂直列表显示（默认）
- **分栏布局**: 左右分栏显示，左栏显示分组，右栏显示主机详情
- **快捷键**: `l` - 在单栏/分栏布局之间切换
- **详细信息**: 在分栏模式下可显示主机的详细信息

### 4. 国际化支持
- **支持语言**: 中文（zh）, 英文（en）
- **配置**: 在 `config.yaml` 中的 `ui_config.language` 设置
- **动态切换**: 修改配置文件后按 `r` 重载即可切换语言

## 编译和运行

```bash
# 编译项目
go build -o hostmanager

# 运行程序
./hostmanager
```

## 开发指南

### 添加新功能
1. 在 `internal/` 下创建新的包
2. 在相应包中实现功能
3. 在 `ui` 包中集成用户界面
4. 更新配置结构（如需要）

### 添加新主题
1. 在 `internal/theme/theme.go` 中添加新主题函数
2. 在 `config.yaml` 的 `themes` 节中配置颜色
3. 更新主题获取逻辑

### 添加新语言
1. 在 `internal/i18n/texts.go` 中添加新语言的文本
2. 更新 `GetTexts()` 函数支持新语言代码

## 配置文件结构
完整的配置文件现在支持以下结构：
- `groups`: 主机分组配置
- `ui_config`: 用户界面配置
  - `theme`: 主题选择
  - `language`: 语言选择
  - `key_bindings`: 快捷键配置
  - `layout`: 布局配置
  - `themes`: 主题颜色配置