package models

// 主机配置结构
type Host struct {
	Name         string   `yaml:"name"`
	IP           string   `yaml:"ip"`
	Port         int      `yaml:"port"`
	Username     string   `yaml:"username"`
	AuthType     string   `yaml:"auth_type"`
	KeyPath      string   `yaml:"key_path,omitempty"`
	Password     string   `yaml:"password,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Favorite     bool     `yaml:"favorite,omitempty"`
	ZmodemEnable *bool    `yaml:"zmodem_enable,omitempty"` // 启用 Zmodem 文件传输支持，默认 true
	Status       string   `yaml:"-"`                       // 运行时状态，不保存到配置文件
}

// 分组配置结构
type Group struct {
	Name  string `yaml:"name"`
	Hosts []Host `yaml:"hosts"`
}

// 获取 Zmodem 启用状态，默认为 true
func (h *Host) IsZmodemEnabled() bool {
	if h.ZmodemEnable == nil {
		return true // 默认启用
	}
	return *h.ZmodemEnable
}