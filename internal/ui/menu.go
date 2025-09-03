package ui

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/config"
	"github.com/daihao4371/hostmanager/internal/i18n"
	"github.com/daihao4371/hostmanager/internal/models"
	"github.com/daihao4371/hostmanager/internal/ssh"
	"github.com/daihao4371/hostmanager/internal/theme"
)

// 菜单管理器（增强版）
type Menu struct {
	groups            []models.Group
	filteredGroups    []models.Group
	currentGroup      int
	currentHost       int
	inGroup           bool
	searchMode        bool
	searchQuery       string
	connectionHistory []models.Host
	showFavorites     bool
	statusCheckMode   bool
	config            *config.Config
	currentTheme      *theme.Theme
	texts             i18n.Texts

	// 高级UI功能
	renderEngine     *RenderEngine     // 渲染引擎
	toastManager     *ToastManager     // 通知管理器
	animationManager *AnimationManager // 动画管理器
	loadingStates    map[string]bool   // 加载状态
	lastFrame        time.Time         // 上一帧时间
	needsRedraw      bool              // 是否需要重绘
}

// Toast通知管理器
type ToastManager struct {
	toasts    []Toast
	maxToasts int
}

// Toast通知
type Toast struct {
	ID        string
	Message   string
	Type      string // "success", "error", "warning", "info"
	StartTime time.Time
	Duration  time.Duration
	Progress  float32
}

// 动画管理器
type AnimationManager struct {
	animations map[string]*Animation
	ticker     *time.Ticker
	running    bool
}

// 创建新的菜单实例（增强版）
func NewMenu(cfg *config.Config) *Menu {
	menu := &Menu{
		groups:            cfg.Groups,
		filteredGroups:    cfg.Groups,
		currentGroup:      0,
		currentHost:       0,
		inGroup:           false,
		searchMode:        false,
		searchQuery:       "",
		connectionHistory: []models.Host{},
		showFavorites:     false,
		statusCheckMode:   false,
		config:            cfg,
		loadingStates:     make(map[string]bool),
		needsRedraw:       true,
	}

	// 初始化主题和国际化
	menu.currentTheme = cfg.UIConfig.Themes.GetTheme(cfg.UIConfig.Theme)
	menu.texts = i18n.GetTexts(cfg.UIConfig.Language)

	// 初始化高级UI组件
	menu.renderEngine = NewRenderEngine(menu.currentTheme)
	menu.toastManager = &ToastManager{
		toasts:    []Toast{},
		maxToasts: 3,
	}
	menu.animationManager = &AnimationManager{
		animations: make(map[string]*Animation),
		running:    false,
	}

	return menu
}

// 运行主循环（增强版）
func (m *Menu) Run() {
	// 启动动画管理器
	m.startAnimationManager()
	defer m.stopAnimationManager()

	for {
		currentTime := time.Now()

		// 更新动画和Toast
		m.updateAnimations()
		m.updateToasts(currentTime)

		// 检查是否需要重绘
		if m.needsRedraw || m.renderEngine.HasActiveAnimations() || len(m.toastManager.toasts) > 0 {
			m.draw()
			m.needsRedraw = false
		}

		// 处理输入
		if !m.handleInput() {
			break
		}

		m.lastFrame = currentTime

		// 控制帧率，减少CPU使用
		time.Sleep(16 * time.Millisecond) // ~60 FPS
	}
}

// 显示Toast通知
func (m *Menu) showToast(message, toastType string, duration time.Duration) {
	if len(m.toastManager.toasts) >= m.toastManager.maxToasts {
		// 移除最旧的Toast
		m.toastManager.toasts = m.toastManager.toasts[1:]
	}

	toast := Toast{
		ID:        fmt.Sprintf("toast_%d", time.Now().UnixNano()),
		Message:   message,
		Type:      toastType,
		StartTime: time.Now(),
		Duration:  duration,
		Progress:  0.0,
	}

	m.toastManager.toasts = append(m.toastManager.toasts, toast)
	m.needsRedraw = true
}

// 更新Toast状态
func (m *Menu) updateToasts(currentTime time.Time) {
	activeToasts := []Toast{}

	for _, toast := range m.toastManager.toasts {
		elapsed := currentTime.Sub(toast.StartTime)
		if elapsed < toast.Duration {
			toast.Progress = float32(elapsed) / float32(toast.Duration)
			activeToasts = append(activeToasts, toast)
		}
	}

	if len(activeToasts) != len(m.toastManager.toasts) {
		m.needsRedraw = true
	}

	m.toastManager.toasts = activeToasts
}

// 启动动画管理器
func (m *Menu) startAnimationManager() {
	m.animationManager.running = true
	m.animationManager.ticker = time.NewTicker(16 * time.Millisecond)

	go func() {
		for range m.animationManager.ticker.C {
			if !m.animationManager.running {
				break
			}
			m.renderEngine.UpdateAnimations()
		}
	}()
}

// 停止动画管理器
func (m *Menu) stopAnimationManager() {
	m.animationManager.running = false
	if m.animationManager.ticker != nil {
		m.animationManager.ticker.Stop()
	}
}

// 更新动画
func (m *Menu) updateAnimations() {
	// 这里可以添加特定的动画逻辑
	m.renderEngine.UpdateAnimations()
}

// 过滤主机基于搜索查询
func (m *Menu) filterHosts() {
	if m.searchQuery == "" {
		m.filteredGroups = m.groups
		return
	}

	m.filteredGroups = []models.Group{}
	for _, group := range m.groups {
		filteredGroup := models.Group{Name: group.Name, Hosts: []models.Host{}}
		for _, host := range group.Hosts {
			if containsIgnoreCase(host.Name, m.searchQuery) ||
				containsIgnoreCase(host.IP, m.searchQuery) ||
				containsIgnoreCase(host.Username, m.searchQuery) {
				filteredGroup.Hosts = append(filteredGroup.Hosts, host)
			}
		}
		if len(filteredGroup.Hosts) > 0 {
			m.filteredGroups = append(m.filteredGroups, filteredGroup)
		}
	}
}

// 字符串包含检查（忽略大小写）
func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) &&
		len(substr) > 0 &&
		findSubstring(strings.ToLower(str), strings.ToLower(substr))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// 获取收藏的主机列表
func (m *Menu) getFavoriteHosts() []models.Host {
	var favorites []models.Host
	for _, group := range m.groups {
		for _, host := range group.Hosts {
			if host.Favorite {
				favorites = append(favorites, host)
			}
		}
	}
	return favorites
}

// 切换主机收藏状态
func (m *Menu) toggleFavorite() {
	if m.inGroup && m.currentGroup < len(m.filteredGroups) && m.currentHost < len(m.filteredGroups[m.currentGroup].Hosts) {
		targetHost := m.filteredGroups[m.currentGroup].Hosts[m.currentHost]
		for i := range m.groups {
			for j := range m.groups[i].Hosts {
				if m.groups[i].Hosts[j].IP == targetHost.IP &&
					m.groups[i].Hosts[j].Port == targetHost.Port &&
					m.groups[i].Hosts[j].Username == targetHost.Username {
					m.groups[i].Hosts[j].Favorite = !m.groups[i].Hosts[j].Favorite
					m.saveConfig()
					m.filterHosts()
					return
				}
			}
		}
	}
}

// 保存配置到文件
func (m *Menu) saveConfig() {
	m.config.Groups = m.groups
	m.config.Save("config.yaml")
}

// 批量检查所有主机状态
func (m *Menu) checkAllHostsStatus() {
	for i := range m.groups {
		for j := range m.groups[i].Hosts {
			go func(groupIndex, hostIndex int) {
				status := ssh.CheckHostStatus(m.groups[groupIndex].Hosts[hostIndex])
				m.groups[groupIndex].Hosts[hostIndex].Status = status
				// 同步更新历史记录中的状态
				for k := range m.connectionHistory {
					if m.connectionHistory[k].IP == m.groups[groupIndex].Hosts[hostIndex].IP &&
						m.connectionHistory[k].Port == m.groups[groupIndex].Hosts[hostIndex].Port {
						m.connectionHistory[k].Status = status
					}
				}
			}(i, j)
		}
	}
}

// 添加主机到连接历史
func (m *Menu) addToHistory(host models.Host) {
	// 检查是否已在历史中
	for i, h := range m.connectionHistory {
		if h.IP == host.IP && h.Port == host.Port && h.Username == host.Username {
			// 移到最前面
			m.connectionHistory = append([]models.Host{host}, append(m.connectionHistory[:i], m.connectionHistory[i+1:]...)...)
			return
		}
	}

	// 添加到历史前面，保持最多5个
	m.connectionHistory = append([]models.Host{host}, m.connectionHistory...)
	if len(m.connectionHistory) > 5 {
		m.connectionHistory = m.connectionHistory[:5]
	}
}

// 切换主题
func (m *Menu) toggleTheme() {
	if m.config.UIConfig.Theme == "light" {
		m.config.UIConfig.Theme = "dark"
	} else {
		m.config.UIConfig.Theme = "light"
	}
	m.currentTheme = m.config.UIConfig.Themes.GetTheme(m.config.UIConfig.Theme)
	m.saveConfig()
}

// 切换布局
func (m *Menu) toggleLayout() {
	if m.config.UIConfig.Layout.Type == "single" {
		m.config.UIConfig.Layout.Type = "columns"
	} else {
		m.config.UIConfig.Layout.Type = "single"
	}
	m.saveConfig()
}

// 重新加载配置
func (m *Menu) reloadConfig() {
	newConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		return
	}
	m.config = newConfig
	m.groups = newConfig.Groups
	m.currentTheme = m.config.UIConfig.Themes.GetTheme(m.config.UIConfig.Theme)
	m.texts = i18n.GetTexts(m.config.UIConfig.Language)
	m.filterHosts()
	m.currentGroup = 0
	m.currentHost = 0
	m.inGroup = false
}

// 连接SSH（包装函数）
func (m *Menu) connectSSH(host models.Host) {
	ssh.Connect(host, m.addToHistory)

	// 连接断开后的恢复处理
	m.recoverFromSSHDisconnect()
}

// SSH连接断开后的恢复处理
func (m *Menu) recoverFromSSHDisconnect() {
	// 显示提示信息
	fmt.Printf("按任意键返回主菜单...\n")

	// 使用简单的字符读取，避免需要回车
	var b = make([]byte, 1)
	os.Stdin.Read(b)

	// 完全重新初始化termbox
	termbox.Close()                    // 确保完全关闭
	time.Sleep(100 * time.Millisecond) // 短暂等待

	err := termbox.Init()
	if err != nil {
		log.Printf("重新初始化termbox失败: %v", err)
		return
	}

	// 设置termbox模式
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.Output256)

	// 强制清屏和重置
	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_ = termbox.Flush()

	// 重新初始化UI状态
	m.needsRedraw = true
	m.renderEngine = NewRenderEngine(m.currentTheme)

	// 重置菜单状态
	m.inGroup = false
	m.currentHost = 0
	m.searchMode = false
	m.searchQuery = ""
	m.filterHosts() // 重置过滤状态

	// 显示恢复Toast
	m.showToast("已返回主菜单", "success", 2*time.Second)
}
