package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/daihao4371/hostmanager/internal/models"
	"github.com/daihao4371/hostmanager/internal/theme"
)

// 快捷键配置
type KeyBindings struct {
	Exit         string `yaml:"exit"`          // 退出/返回
	Search       string `yaml:"search"`        // 搜索
	Favorites    string `yaml:"favorites"`     // 收藏夹
	StatusCheck  string `yaml:"status_check"`  // 状态检查
	Reload       string `yaml:"reload"`        // 重载配置
	ToggleFav    string `yaml:"toggle_fav"`    // 切换收藏
	ThemeSwitch  string `yaml:"theme_switch"`  // 主题切换
	LayoutSwitch string `yaml:"layout_switch"` // 布局切换
}

// 布局配置
type Layout struct {
	Type        string `yaml:"type"`         // "single" 或 "columns"
	ShowDetails bool   `yaml:"show_details"` // 是否显示详细信息
	ColumnWidth int    `yaml:"column_width"` // 列宽度
}

// 用户界面配置
type UIConfig struct {
	Theme       string          `yaml:"theme"`        // "light" 或 "dark"
	Language    string          `yaml:"language"`     // "zh" 或 "en"
	KeyBindings KeyBindings     `yaml:"key_bindings"`
	Layout      Layout          `yaml:"layout"`
	Themes      theme.Themes    `yaml:"themes"`
}

// 主配置结构
type Config struct {
	Groups   []models.Group `yaml:"groups"`
	UIConfig UIConfig       `yaml:"ui_config"`
}

// 查找配置文件的可能位置
func findConfigFile(filename string) (string, error) {
	// 配置文件查找优先级
	searchPaths := []string{
		filename,                                    // 当前目录
		filepath.Join(".", filename),               // 当前目录（显式）
		filepath.Join(os.Getenv("HOME"), filename), // 用户家目录
		"/etc/hostmanager/" + filename,             // 系统配置目录
	}

	// 如果是相对路径，尝试找到程序所在目录
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		searchPaths = append(searchPaths, filepath.Join(execDir, filename))
	}

	// 按优先级查找配置文件
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", os.ErrNotExist
}

// 加载配置文件
func LoadConfig(filePath string) (*Config, error) {
	// 查找配置文件的实际位置
	actualPath, err := findConfigFile(filePath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(actualPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// 设置主机默认值
	for i := range config.Groups {
		for j := range config.Groups[i].Hosts {
			setHostDefaults(&config.Groups[i].Hosts[j])
		}
	}

	// 设置UI配置默认值
	setUIDefaults(&config.UIConfig)

	return &config, nil
}

// 保存配置到文件（全局函数）
func SaveConfig(filePath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// 保存配置文件
func (c *Config) Save(filePath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// 设置主机默认值
func setHostDefaults(host *models.Host) {
	if host.Port == 0 {
		host.Port = 22
	}
	if host.Username == "" {
		host.Username = "app"
	}
	if host.AuthType == "" {
		host.AuthType = "key"
	}
	if host.KeyPath == "" && host.AuthType == "key" {
		host.KeyPath = "/Users/daihao/.ssh/id_rsa"
	}
}

// 设置UI配置默认值
func setUIDefaults(ui *UIConfig) {
	if ui.Theme == "" {
		ui.Theme = "dark"
	}
	if ui.Language == "" {
		ui.Language = "zh"
	}

	// 设置默认快捷键
	setKeyBindingDefaults(&ui.KeyBindings)

	// 设置默认布局
	setLayoutDefaults(&ui.Layout)

	// 设置默认主题
	ui.Themes.SetDefaults()
}

// 设置快捷键默认值
func setKeyBindingDefaults(kb *KeyBindings) {
	if kb.Exit == "" {
		kb.Exit = "Esc"
	}
	if kb.Search == "" {
		kb.Search = "/"
	}
	if kb.Favorites == "" {
		kb.Favorites = "f"
	}
	if kb.StatusCheck == "" {
		kb.StatusCheck = "s"
	}
	if kb.Reload == "" {
		kb.Reload = "r"
	}
	if kb.ToggleFav == "" {
		kb.ToggleFav = "Space"
	}
	if kb.ThemeSwitch == "" {
		kb.ThemeSwitch = "t"
	}
	if kb.LayoutSwitch == "" {
		kb.LayoutSwitch = "l"
	}
}

// 设置布局默认值
func setLayoutDefaults(layout *Layout) {
	if layout.Type == "" {
		layout.Type = "single"
	}
	if layout.ColumnWidth == 0 {
		layout.ColumnWidth = 80
	}
}