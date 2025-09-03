package ui

import (
	"fmt"
	"strings"

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

	// 绘制Toast通知（在最上层）
	m.drawToasts()

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
	width, _ := termbox.Size()

	// 绘制顶部边框，根据屏幕宽度调整
	borderLine := strings.Repeat("═", width-1)
	m.printThemedString(0, y, borderLine, m.currentTheme.Border)
	y++

	// 标题居中显示（使用正确的显示宽度计算）
	title := "🖥️  " + m.texts.Title
	titleDisplayWidth := getDisplayWidth(title)
	if titleDisplayWidth < width {
		padding := (width - titleDisplayWidth) / 2
		if padding > 0 {
			m.printThemedString(padding, y, title, m.currentTheme.Info)
		} else {
			m.printThemedString(0, y, title, m.currentTheme.Info)
		}
	} else {
		m.printThemedString(0, y, title, m.currentTheme.Info)
	}
	y++

	m.printThemedString(0, y, borderLine, m.currentTheme.Border)
	y++

	// 显示当前主题和布局信息
	themeInfo := fmt.Sprintf("主题: %s | 布局: %s", m.config.UIConfig.Theme, m.config.UIConfig.Layout.Type)
	m.printThemedString(0, y, themeInfo, m.currentTheme.Border)
	y++

	// 操作提示 - 根据屏幕宽度智能分行显示
	operations := m.texts.Operations
	maxOperationWidth := width - 4 // 留出边距

	// 针对不同宽度采用不同的显示策略
	if maxOperationWidth < 50 {
		// 极窄屏幕：显示最基础的操作提示
		m.printThemedString(0, y, "操作: ↑↓选择 | 回车连接 | /搜索 | ESC退出", m.currentTheme.Foreground)
		y++
	} else if maxOperationWidth < 80 {
		// 中等宽度：分两行显示
		m.printThemedString(0, y, "操作: ↑↓选择 | 回车连接 | /搜索 | f收藏夹", m.currentTheme.Foreground)
		y++
		m.printThemedString(0, y, "s状态检查 | r重载 | t主题 | l布局 | ESC退出", m.currentTheme.Foreground)
		y++
	} else if getDisplayWidth(operations) > maxOperationWidth {
		// 宽度充足但操作文本太长：使用智能分割
		operationLines := splitStringToFitWidth(operations, maxOperationWidth)
		for _, line := range operationLines {
			m.printThemedString(0, y, line, m.currentTheme.Foreground)
			y++
		}
	} else {
		// 宽度充足：直接显示完整操作提示
		m.printThemedString(0, y, operations, m.currentTheme.Foreground)
		y++
	}

	// 底部分隔线
	separatorLine := strings.Repeat("─", width-1)
	m.printThemedString(0, y, separatorLine, m.currentTheme.Border)
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

// 带主题的字符串打印（正确处理宽字符）
func (m *Menu) printThemedString(x, y int, str string, color termbox.Attribute) {
	width, height := termbox.Size()
	
	// 边界检查
	if y >= height || x >= width {
		return
	}
	
	screenX := x
	for _, r := range str {
		// 检查是否超出屏幕右边界
		if screenX >= width {
			break
		}
		
		// 设置主字符
		termbox.SetCell(screenX, y, r, color, m.currentTheme.Background)
		screenX++
		
		// 如果是宽字符（中文、emoji等），需要设置占位符
		if r >= 0x4E00 && r <= 0x9FFF || // 中文字符
		   r >= 0x3400 && r <= 0x4DBF || // 中文扩展A
		   r >= 0x20000 && r <= 0x2A6DF || // 中文扩展B  
		   r >= 0x1F300 && r <= 0x1F9FF { // Emoji
			if screenX < width {
				// 为宽字符的第二个位置设置空格占位符
				termbox.SetCell(screenX, y, ' ', color, m.currentTheme.Background)
				screenX++
			}
		}
	}
}

// 绘制Toast通知
func (m *Menu) drawToasts() {
	width, height := termbox.Size()
	startY := 1 // 从顶部第二行开始，避免覆盖标题

	for i, toast := range m.toastManager.toasts {
		if i >= 3 { // 最多显示3个通知
			break
		}

		y := startY + i
		if y >= height-1 {
			break
		}

		// Toast背景和边框
		toastWidth := minInt(len(toast.Message)+4, width-4)
		startX := width - toastWidth - 2

		// 根据类型选择颜色
		var bgColor, textColor termbox.Attribute
		var icon string

		switch toast.Type {
		case "success":
			bgColor = m.currentTheme.Success
			textColor = m.currentTheme.Background
			icon = "✅"
		case "error":
			bgColor = m.currentTheme.Error
			textColor = m.currentTheme.Background
			icon = "❌"
		case "warning":
			bgColor = m.currentTheme.Warning
			textColor = m.currentTheme.Background
			icon = "⚠️"
		default:
			bgColor = m.currentTheme.Info
			textColor = m.currentTheme.Background
			icon = "ℹ️"
		}

		// 绘制Toast背景
		for x := startX; x < startX+toastWidth; x++ {
			termbox.SetCell(x, y, ' ', textColor, bgColor)
		}

		// 绘制Toast文本
		toastText := fmt.Sprintf("%s %s", icon, toast.Message)
		runeIndex := 0
		for _, r := range toastText {
			if runeIndex >= toastWidth-2 {
				break
			}
			termbox.SetCell(startX+1+runeIndex, y, r, textColor, bgColor)
			runeIndex++
		}
	}
}

// 辅助函数：获取两个整数的较小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 计算字符串的显示宽度（精确处理emoji和中文字符）
func getDisplayWidth(s string) int {
	width := 0
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		// ASCII字符占1个位置
		if r < 128 {
			width++
		} else if r >= 0x1F300 && r <= 0x1F9FF { // Emoji范围
			width += 2
		} else if r >= 0x4E00 && r <= 0x9FFF { // 中文字符范围
			width += 2
		} else if r >= 0x3400 && r <= 0x4DBF { // 中文扩展A区
			width += 2
		} else if r >= 0x20000 && r <= 0x2A6DF { // 中文扩展B区
			width += 2
		} else {
			width += 1 // 其他字符默认占1个位置
		}
	}
	return width
}

// 智能分割字符串以适应指定宽度（增强版）
func splitStringToFitWidth(s string, maxWidth int) []string {
	if getDisplayWidth(s) <= maxWidth {
		return []string{s}
	}

	var lines []string
	
	// 首先尝试按 " | " 分割
	words := strings.Split(s, " | ")
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " | "
		}
		testLine += word

		if getDisplayWidth(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			// 如果当前行不为空，保存它
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// 单个词太长，需要强制截断
				if getDisplayWidth(word) > maxWidth {
					lines = append(lines, truncateStringToWidth(word, maxWidth))
				} else {
					currentLine = word
				}
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// 将字符串截断到指定宽度，在末尾添加省略号
func truncateStringToWidth(s string, maxWidth int) string {
	if getDisplayWidth(s) <= maxWidth {
		return s
	}
	
	if maxWidth < 3 {
		return strings.Repeat(".", maxWidth) // 太窄时只显示点
	}
	
	runes := []rune(s)
	currentWidth := 0
	cutPoint := 0
	
	for i, r := range runes {
		charWidth := 1
		if r >= 0x1F300 && r <= 0x1F9FF { // Emoji
			charWidth = 2
		} else if r >= 0x4E00 && r <= 0x9FFF { // 中文
			charWidth = 2
		} else if r >= 0x3400 && r <= 0x4DBF { // 中文扩展A
			charWidth = 2
		} else if r >= 0x20000 && r <= 0x2A6DF { // 中文扩展B
			charWidth = 2
		} else if r >= 128 { // 其他非ASCII字符
			charWidth = 2
		}
		
		if currentWidth + charWidth + 3 > maxWidth { // 为省略号留出空间
			cutPoint = i
			break
		}
		currentWidth += charWidth
	}
	
	if cutPoint == 0 {
		return "..."
	}
	
	return string(runes[:cutPoint]) + "..."
}
