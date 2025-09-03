package i18n

// 国际化文本配置
type Texts struct {
	Title             string `yaml:"title"`
	SearchMode        string `yaml:"search_mode"`
	SearchPlaceholder string `yaml:"search_placeholder"`
	NoMatches         string `yaml:"no_matches"`
	FoundHosts        string `yaml:"found_hosts"`
	QuickConnect      string `yaml:"quick_connect"`
	ServerGroups      string `yaml:"server_groups"`
	Operations        string `yaml:"operations"`
	Favorites         string `yaml:"favorites"`
	NoFavorites       string `yaml:"no_favorites"`
	Connecting        string `yaml:"connecting"`
	ConnectionClosed  string `yaml:"connection_closed"`
	PressAnyKey       string `yaml:"press_any_key"`
	Online            string `yaml:"online"`
	Offline           string `yaml:"offline"`
	Unknown           string `yaml:"unknown"`
}

// 获取中文文本
func GetChineseTexts() Texts {
	return Texts{
		Title:             "SSH 连接管理器",
		SearchMode:        "搜索模式",
		SearchPlaceholder: "搜索关键词: ",
		NoMatches:         "未找到匹配的主机",
		FoundHosts:        "找到 %d 个匹配的主机",
		QuickConnect:      "快速连接 (按数字键1-5直接连接):",
		ServerGroups:      "服务器分组:",
		Operations:        "操作: ↑↓选择 | 回车连接 | /搜索 | f收藏夹 | s状态检查 | r重载 | t主题 | l布局 | ESC退出",
		Favorites:         "收藏的主机 (按f退出收藏模式):",
		NoFavorites:       "暂无收藏的主机，在主机列表中按空格键添加收藏",
		Connecting:        "正在连接到 %s (%s@%s:%d)...",
		ConnectionClosed:  "与 %s 的连接已断开",
		PressAnyKey:       "按任意键返回主菜单...",
		Online:            "在线",
		Offline:           "离线",
		Unknown:           "未知",
	}
}

// 获取英文文本
func GetEnglishTexts() Texts {
	return Texts{
		Title:             "SSH Connection Manager",
		SearchMode:        "Search Mode",
		SearchPlaceholder: "Search Keywords: ",
		NoMatches:         "No matching hosts found",
		FoundHosts:        "Found %d matching hosts",
		QuickConnect:      "Quick Connect (Press number key 1-5):",
		ServerGroups:      "Server Groups:",
		Operations:        "Operations: ↑↓Select | Enter Connect | /Search | f Favorites | s Status | r Reload | t Theme | l Layout | ESC Exit",
		Favorites:         "Favorite Hosts (Press f to exit favorites mode):",
		NoFavorites:       "No favorite hosts. Press Space in host list to add favorites",
		Connecting:        "Connecting to %s (%s@%s:%d)...",
		ConnectionClosed:  "Connection to %s closed",
		PressAnyKey:       "Press any key to return to main menu...",
		Online:            "Online",
		Offline:           "Offline",
		Unknown:           "Unknown",
	}
}

// 根据语言代码获取文本
func GetTexts(language string) Texts {
	switch language {
	case "en":
		return GetEnglishTexts()
	case "zh":
		return GetChineseTexts()
	default:
		return GetChineseTexts()
	}
}