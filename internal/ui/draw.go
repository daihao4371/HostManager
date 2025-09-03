package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/models"
)

// 主绘制函数
func (m *Menu) draw() {
	termbox.Clear(m.currentTheme.Background, m.currentTheme.Background)

	if m.config.UIConfig.Layout.Type == "columns" {
		m.drawColumnsLayout()
	} else {
		m.drawSingleLayout()
	}

	termbox.Flush()
}

// 单栏布局绘制
func (m *Menu) drawSingleLayout() {
	y := 0

	if m.searchMode {
		y = m.drawSearchMode(y)
	} else {
		y = m.drawHeader(y)
		y = m.drawQuickConnect(y)
	}

	m.drawMainContent(y)
}

// 分栏布局绘制
func (m *Menu) drawColumnsLayout() {
	width, height := termbox.Size()
	leftWidth := width / 2

	// 左栏：分组和快速连接
	m.drawLeftColumn(0, 0, leftWidth, height)

	// 分隔线
	for i := 0; i < height; i++ {
		termbox.SetCell(leftWidth, i, '│', m.currentTheme.Border, m.currentTheme.Background)
	}

	// 右栏：主机列表或详细信息
	m.drawRightColumn(leftWidth+1, 0, width-leftWidth-1, height)
}

// 绘制搜索模式
func (m *Menu) drawSearchMode(y int) int {
	m.printThemedString(0, y, "═══════════════════════════════════════════════════════════", m.currentTheme.Border)
	y++
	m.printThemedString(0, y, "🔍 "+m.texts.SearchMode, m.currentTheme.Info)
	m.printThemedString(15, y, "按ESC退出搜索，回车连接第一个匹配项", m.currentTheme.Foreground)
	y++
	m.printThemedString(0, y, "═══════════════════════════════════════════════════════════", m.currentTheme.Border)
	y++

	// 搜索输入框
	searchDisplay := fmt.Sprintf("%s%s█", m.texts.SearchPlaceholder, m.searchQuery)
	m.printThemedString(0, y, searchDisplay, m.currentTheme.Info)
	y += 2

	// 搜索结果提示
	if len(m.filteredGroups) == 0 {
		m.printThemedString(0, y, "❌ "+m.texts.NoMatches, m.currentTheme.Error)
	} else {
		totalHosts := 0
		for _, group := range m.filteredGroups {
			totalHosts += len(group.Hosts)
		}
		resultInfo := fmt.Sprintf("✅ "+m.texts.FoundHosts, totalHosts)
		m.printThemedString(0, y, resultInfo, m.currentTheme.Success)
	}
	y++
	m.printThemedString(0, y, "───────────────────────────────────────────────────────────", m.currentTheme.Border)
	return y + 1
}

// 绘制标题区域
func (m *Menu) drawHeader(y int) int {
	m.printThemedString(0, y, "═══════════════════════════════════════════════════════════", m.currentTheme.Border)
	y++
	m.printThemedString(0, y, "🖥️  "+m.texts.Title, m.currentTheme.Info)
	y++
	m.printThemedString(0, y, "═══════════════════════════════════════════════════════════", m.currentTheme.Border)
	y++

	// 显示当前主题和布局信息
	themeInfo := fmt.Sprintf("主题: %s | 布局: %s", m.config.UIConfig.Theme, m.config.UIConfig.Layout.Type)
	m.printThemedString(0, y, themeInfo, m.currentTheme.Border)
	y++

	// 操作提示
	m.printThemedString(0, y, m.texts.Operations, m.currentTheme.Foreground)
	y++
	m.printThemedString(0, y, "───────────────────────────────────────────────────────────", m.currentTheme.Border)
	return y + 1
}

// 绘制快速连接区域
func (m *Menu) drawQuickConnect(y int) int {
	if len(m.connectionHistory) > 0 {
		m.printThemedString(0, y, "⚡ "+m.texts.QuickConnect, m.currentTheme.Success)
		y++
		for i, host := range m.connectionHistory {
			if i >= 5 {
				break
			}
			statusIcon := m.getStatusIcon(host.Status)
			favoriteIcon := ""
			if host.Favorite {
				favoriteIcon = "⭐"
			}

			historyInfo := fmt.Sprintf("   %d. %s%s %s (%s@%s:%d)", i+1, statusIcon, favoriteIcon, host.Name, host.Username, host.IP, host.Port)
			m.printThemedString(0, y, historyInfo, m.currentTheme.Border)
			y++
		}
		m.printThemedString(0, y, "───────────────────────────────────────────────────────────", m.currentTheme.Border)
		y++
	}
	return y
}

// 绘制主要内容区域
func (m *Menu) drawMainContent(y int) {
	groups := m.filteredGroups
	if len(groups) == 0 {
		return
	}

	if !m.inGroup {
		if m.showFavorites {
			m.drawFavorites(y)
		} else {
			m.drawGroups(y, groups)
		}
	} else {
		m.drawHosts(y, groups)
	}
}

// 获取状态图标
func (m *Menu) getStatusIcon(status string) string {
	switch status {
	case "online":
		return "🟢"
	case "offline":
		return "🔴"
	default:
		return "❓"
	}
}

// 获取认证类型图标
func (m *Menu) getAuthIcon(host models.Host) string {
	if host.AuthType == "password" || (host.AuthType == "key" && host.KeyPath == "" && host.Password != "") {
		return "🔐" // 密码认证
	}
	return "🔑" // 密钥认证
}

// 带主题的字符串打印
func (m *Menu) printThemedString(x, y int, str string, color termbox.Attribute) {
	for i, r := range str {
		termbox.SetCell(x+i, y, r, color, m.currentTheme.Background)
	}
}