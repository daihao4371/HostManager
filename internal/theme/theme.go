package theme

import "github.com/nsf/termbox-go"

// 主题配置结构
type Theme struct {
	Background termbox.Attribute `yaml:"background"`
	Foreground termbox.Attribute `yaml:"foreground"`
	Highlight  termbox.Attribute `yaml:"highlight"`
	Border     termbox.Attribute `yaml:"border"`
	Success    termbox.Attribute `yaml:"success"`
	Warning    termbox.Attribute `yaml:"warning"`
	Error      termbox.Attribute `yaml:"error"`
	Info       termbox.Attribute `yaml:"info"`
	
	// 扩展颜色方案
	Accent1    termbox.Attribute `yaml:"accent1"`    // 主要强调色
	Accent2    termbox.Attribute `yaml:"accent2"`    // 次要强调色
	Muted      termbox.Attribute `yaml:"muted"`      // 静音色
	Surface    termbox.Attribute `yaml:"surface"`    // 表面色
	OnSurface  termbox.Attribute `yaml:"on_surface"` // 表面上的文字色
}

// 主题集合
type Themes struct {
	Light Theme `yaml:"light"`
	Dark  Theme `yaml:"dark"`
}

// 获取预定义的浅色主题（现代化设计）
func GetLightTheme() Theme {
	return Theme{
		Background: termbox.ColorWhite,
		Foreground: termbox.ColorBlack,
		Highlight:  termbox.ColorBlue | termbox.AttrBold,
		Border:     termbox.ColorCyan,
		Success:    termbox.ColorGreen | termbox.AttrBold,
		Warning:    termbox.ColorYellow | termbox.AttrBold,
		Error:      termbox.ColorRed | termbox.AttrBold,
		Info:       termbox.ColorMagenta | termbox.AttrBold,
		Accent1:    termbox.ColorBlue,
		Accent2:    termbox.ColorCyan,
		Muted:      termbox.Attribute(8), // 深灰色
		Surface:    termbox.Attribute(15), // 白色
		OnSurface:  termbox.ColorBlack,
	}
}

// 获取预定义的深色主题（现代化设计）
func GetDarkTheme() Theme {
	return Theme{
		Background: termbox.ColorDefault,
		Foreground: termbox.ColorWhite,
		Highlight:  termbox.ColorYellow | termbox.AttrBold,
		Border:     termbox.ColorBlue,
		Success:    termbox.ColorGreen | termbox.AttrBold,
		Warning:    termbox.ColorYellow | termbox.AttrBold,
		Error:      termbox.ColorRed | termbox.AttrBold,
		Info:       termbox.ColorCyan | termbox.AttrBold,
		Accent1:    termbox.ColorMagenta,
		Accent2:    termbox.ColorCyan,
		Muted:      termbox.Attribute(8), // 深灰色
		Surface:    termbox.Attribute(0), // 黑色
		OnSurface:  termbox.ColorWhite,
	}
}

// 获取高对比度主题（无障碍设计）
func GetHighContrastTheme() Theme {
	return Theme{
		Background: termbox.ColorBlack,
		Foreground: termbox.ColorWhite | termbox.AttrBold,
		Highlight:  termbox.ColorYellow | termbox.AttrBold | termbox.AttrReverse,
		Border:     termbox.ColorWhite | termbox.AttrBold,
		Success:    termbox.ColorGreen | termbox.AttrBold | termbox.AttrReverse,
		Warning:    termbox.ColorYellow | termbox.AttrBold | termbox.AttrReverse,
		Error:      termbox.ColorRed | termbox.AttrBold | termbox.AttrReverse,
		Info:       termbox.ColorCyan | termbox.AttrBold | termbox.AttrReverse,
		Accent1:    termbox.ColorMagenta | termbox.AttrBold,
		Accent2:    termbox.ColorCyan | termbox.AttrBold,
		Muted:      termbox.Attribute(7), // 浅灰色
		Surface:    termbox.ColorBlack,
		OnSurface:  termbox.ColorWhite | termbox.AttrBold,
	}
}

// 根据名称获取主题
func (t *Themes) GetTheme(name string) *Theme {
	switch name {
	case "light":
		return &t.Light
	case "dark":
		return &t.Dark
	case "high-contrast":
		theme := GetHighContrastTheme()
		return &theme
	default:
		return &t.Dark
	}
}

// 设置默认主题
func (t *Themes) SetDefaults() {
	if t.Light.Background == 0 {
		t.Light = GetLightTheme()
	}
	if t.Dark.Background == 0 {
		t.Dark = GetDarkTheme()
	}
}