package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/models"
)

// 绘制收藏夹
func (m *Menu) drawFavorites(y int) {
	m.printThemedString(0, y, "⭐ "+m.texts.Favorites, m.currentTheme.Info)
	y++
	favorites := m.getFavoriteHosts()
	if len(favorites) == 0 {
		m.printThemedString(0, y, "   "+m.texts.NoFavorites, m.currentTheme.Border)
	} else {
		for i, host := range favorites {
			color := m.currentTheme.Foreground
			prefix := "   "
			if i == m.currentHost {
				color = m.currentTheme.Highlight
				prefix = "🔹 "
			}

			statusIcon := m.getStatusIcon(host.Status)
			authIcon := m.getAuthIcon(host)
			hostInfo := fmt.Sprintf("%s%s%s %s (%s@%s:%d)", prefix, statusIcon, authIcon, host.Name, host.Username, host.IP, host.Port)

			if host.Description != "" {
				hostInfo += fmt.Sprintf(" - %s", host.Description)
			}
			m.printThemedString(0, y, hostInfo, color)
			y++
		}
	}
}

// 绘制分组列表
func (m *Menu) drawGroups(y int, groups []models.Group) {
	if !m.searchMode {
		m.printThemedString(0, y, "📁 "+m.texts.ServerGroups, m.currentTheme.Success)
		y++
	}

	width, _ := termbox.Size()
	
	for i, group := range groups {
		color := m.currentTheme.Foreground
		prefix := "   "
		if i == m.currentGroup {
			color = m.currentTheme.Highlight
			prefix = "🔸 "
		}

		hostCount := len(group.Hosts)
		groupDisplay := fmt.Sprintf("%s%s (%d台主机)", prefix, group.Name, hostCount)
		
		// 确保分组名称不被截断 - 使用简化的截断逻辑
		maxDisplayWidth := width - 2
		if len([]rune(groupDisplay)) > maxDisplayWidth {
			// 计算可显示的分组名称长度
			prefixLen := len([]rune(prefix))
			suffixLen := len([]rune(fmt.Sprintf(" (%d台主机)", hostCount)))
			maxNameLen := maxDisplayWidth - prefixLen - suffixLen - 3 // 为省略号留空间
			
			if maxNameLen > 0 {
				nameRunes := []rune(group.Name)
				if len(nameRunes) > maxNameLen {
					truncatedName := string(nameRunes[:maxNameLen]) + "..."
					groupDisplay = fmt.Sprintf("%s%s (%d台主机)", prefix, truncatedName, hostCount)
				}
			}
		}
		
		m.printThemedString(0, y, groupDisplay, color)
		y++
	}
}

// 绘制主机列表
func (m *Menu) drawHosts(y int, groups []models.Group) {
	if m.currentGroup >= len(groups) {
		m.currentGroup = 0
	}
	group := groups[m.currentGroup]

	if !m.searchMode {
		groupHeader := fmt.Sprintf("📂 分组: %s (按ESC返回，空格收藏)", group.Name)
		m.printThemedString(0, y, groupHeader, m.currentTheme.Info)
		y++
		m.printThemedString(0, y, "───────────────────────────────────────────────────────────", m.currentTheme.Border)
		y++
	}

	for i, host := range group.Hosts {
		color := m.currentTheme.Foreground
		prefix := "   "
		if i == m.currentHost {
			color = m.currentTheme.Highlight
			prefix = "🔹 "
		}

		statusIcon := m.getStatusIcon(host.Status)
		authIcon := m.getAuthIcon(host)
		favoriteIcon := ""
		if host.Favorite {
			favoriteIcon = "⭐"
		}

		hostInfo := fmt.Sprintf("%s%s%s%s %s (%s@%s:%d)", prefix, statusIcon, authIcon, favoriteIcon, host.Name, host.Username, host.IP, host.Port)
		if host.Description != "" {
			hostInfo += fmt.Sprintf(" - %s", host.Description)
		}
		m.printThemedString(0, y, hostInfo, color)
		y++
	}
}

// 左栏绘制（分栏布局）
func (m *Menu) drawLeftColumn(x, y, width, height int) {
	currentY := y

	// 标题
	title := "🖥️  " + m.texts.Title
	m.printThemedStringInBounds(x, currentY, title, m.currentTheme.Info, width)
	currentY += 2

	// 快速连接
	if len(m.connectionHistory) > 0 {
		m.printThemedStringInBounds(x, currentY, "⚡ "+m.texts.QuickConnect, m.currentTheme.Success, width)
		currentY++
		for i, host := range m.connectionHistory {
			if i >= 3 || currentY >= height-2 {
				break
			}
			statusIcon := m.getStatusIcon(host.Status)
			historyInfo := fmt.Sprintf("  %d. %s %s", i+1, statusIcon, host.Name)
			m.printThemedStringInBounds(x, currentY, historyInfo, m.currentTheme.Border, width)
			currentY++
		}
		currentY++
	}

	// 分组列表
	m.printThemedStringInBounds(x, currentY, "📁 "+m.texts.ServerGroups, m.currentTheme.Success, width)
	currentY++
	for i, group := range m.filteredGroups {
		if currentY >= height-1 {
			break
		}
		color := m.currentTheme.Foreground
		prefix := "  "
		if i == m.currentGroup {
			color = m.currentTheme.Highlight
			prefix = "▶ "
		}

		// 智能截断分组名称以适应左栏宽度
		hostCount := len(group.Hosts)
		suffix := fmt.Sprintf(" (%d台)", hostCount)
		maxNameWidth := width - len([]rune(prefix)) - len([]rune(suffix)) - 4 // 留出边距
		
		displayName := group.Name
		if len([]rune(displayName)) > maxNameWidth && maxNameWidth > 3 {
			nameRunes := []rune(displayName)
			displayName = string(nameRunes[:maxNameWidth-3]) + "..."
		}
		
		groupDisplay := fmt.Sprintf("%s%s%s", prefix, displayName, suffix)
		m.printThemedStringInBounds(x, currentY, groupDisplay, color, width)
		currentY++
	}
}

// 右栏绘制（分栏布局）
func (m *Menu) drawRightColumn(x, y, width, height int) {
	if m.showFavorites {
		m.drawFavoritesInColumn(x, y, width, height)
	} else if m.inGroup && m.currentGroup < len(m.filteredGroups) {
		m.drawHostsInColumn(x, y, width, height)
	} else {
		// 显示操作说明
		m.printThemedStringInBounds(x, y, "操作说明:", m.currentTheme.Info, width)
		y++
		m.printThemedStringInBounds(x, y, "↑↓ 选择项目", m.currentTheme.Foreground, width)
		y++
		m.printThemedStringInBounds(x, y, "回车 进入分组/连接主机", m.currentTheme.Foreground, width)
		y++
		m.printThemedStringInBounds(x, y, "t 切换主题", m.currentTheme.Foreground, width)
		y++
		m.printThemedStringInBounds(x, y, "l 切换布局", m.currentTheme.Foreground, width)
	}
}

// 在指定范围内绘制字符串（正确处理宽字符）
func (m *Menu) printThemedStringInBounds(x, y int, str string, color termbox.Attribute, maxWidth int) {
	if maxWidth <= 0 {
		return
	}

	runes := []rune(str)
	currentDisplayWidth := 0
	screenX := x
	
	for _, r := range runes {
		// 计算字符显示宽度
		charWidth := 1
		if r >= 0x1F300 && r <= 0x1F9FF { // Emoji范围
			charWidth = 2
		} else if r >= 0x4E00 && r <= 0x9FFF { // 中文字符范围
			charWidth = 2
		} else if r >= 0x3400 && r <= 0x4DBF { // 中文扩展A区
			charWidth = 2
		} else if r >= 0x20000 && r <= 0x2A6DF { // 中文扩展B区
			charWidth = 2
		} else if r >= 128 {
			charWidth = 2
		}
		
		// 检查是否还有空间显示这个字符
		if currentDisplayWidth + charWidth > maxWidth || screenX >= x + maxWidth {
			break
		}
		
		// 设置主字符
		termbox.SetCell(screenX, y, r, color, m.currentTheme.Background)
		screenX++
		currentDisplayWidth++
		
		// 如果是宽字符，设置第二个位置的占位符
		if charWidth == 2 && screenX < x + maxWidth {
			termbox.SetCell(screenX, y, ' ', color, m.currentTheme.Background)
			screenX++
			currentDisplayWidth++
		}
	}
}

// 在右栏绘制收藏夹
func (m *Menu) drawFavoritesInColumn(x, y, width, height int) {
	m.printThemedStringInBounds(x, y, "⭐ "+m.texts.Favorites, m.currentTheme.Info, width)
	y++
	favorites := m.getFavoriteHosts()
	if len(favorites) == 0 {
		m.printThemedStringInBounds(x, y, m.texts.NoFavorites, m.currentTheme.Border, width)
	} else {
		for i, host := range favorites {
			if y >= height-1 {
				break
			}
			color := m.currentTheme.Foreground
			prefix := "  "
			if i == m.currentHost {
				color = m.currentTheme.Highlight
				prefix = "▶ "
			}

			statusIcon := m.getStatusIcon(host.Status)
			hostInfo := fmt.Sprintf("%s%s %s", prefix, statusIcon, host.Name)
			m.printThemedStringInBounds(x, y, hostInfo, color, width)
			y++
		}
	}
}

// 在右栏绘制主机列表
func (m *Menu) drawHostsInColumn(x, y, width, height int) {
	group := m.filteredGroups[m.currentGroup]

	groupHeader := fmt.Sprintf("📂 %s", group.Name)
	m.printThemedStringInBounds(x, y, groupHeader, m.currentTheme.Info, width)
	y++

	for i, host := range group.Hosts {
		if y >= height-1 {
			break
		}
		color := m.currentTheme.Foreground
		prefix := "  "
		if i == m.currentHost {
			color = m.currentTheme.Highlight
			prefix = "▶ "
		}

		statusIcon := m.getStatusIcon(host.Status)
		authIcon := m.getAuthIcon(host)
		favoriteIcon := ""
		if host.Favorite {
			favoriteIcon = "⭐"
		}

		hostInfo := fmt.Sprintf("%s%s%s%s %s", prefix, statusIcon, authIcon, favoriteIcon, host.Name)
		m.printThemedStringInBounds(x, y, hostInfo, color, width)
		y++

		// 在分栏模式下显示更多详细信息
		if m.config.UIConfig.Layout.ShowDetails && prefix == "▶ " {
			detailInfo := fmt.Sprintf("    %s@%s:%d", host.Username, host.IP, host.Port)
			m.printThemedStringInBounds(x, y, detailInfo, m.currentTheme.Border, width)
			y++
			if host.Description != "" {
				descInfo := fmt.Sprintf("    %s", host.Description)
				m.printThemedStringInBounds(x, y, descInfo, m.currentTheme.Border, width)
				y++
			}
		}
	}
}
