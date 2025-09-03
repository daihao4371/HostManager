package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

// 主输入处理函数（调试版本）
func (m *Menu) handleInput() bool {
	ev := termbox.PollEvent()

	// 调试信息
	log.Printf("事件类型: %d, 按键: %d, 字符: %c", ev.Type, ev.Key, ev.Ch)

	switch ev.Type {
	case termbox.EventKey:
		m.needsRedraw = true
		if m.searchMode {
			return m.handleSearchInput(ev)
		} else {
			return m.handleNormalInput(ev)
		}
	case termbox.EventResize:
		m.renderEngine = NewRenderEngine(m.currentTheme)
		m.needsRedraw = true
	case termbox.EventError:
		log.Printf("Termbox事件错误: %v", ev.Err)
		return false
	}
	return true
}

// 处理搜索模式的输入
func (m *Menu) handleSearchInput(ev termbox.Event) bool {
	switch ev.Key {
	case termbox.KeyEsc:
		m.searchMode = false
		m.searchQuery = ""
		m.filterHosts()
		m.currentGroup = 0
		m.currentHost = 0
		m.inGroup = false
		m.needsRedraw = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			m.filterHosts()
			m.needsRedraw = true
		}
	case termbox.KeyEnter:
		if len(m.filteredGroups) > 0 && len(m.filteredGroups[0].Hosts) > 0 {
			m.searchMode = false
			termbox.Close()
			host := m.filteredGroups[0].Hosts[0]

			fmt.Printf("\n🎯 搜索选择: %s\n", host.Name)
			m.connectSSH(host)

			err := termbox.Init()
			if err != nil {
				log.Printf("重新初始化termbox失败: %v", err)
				return false
			}

			// 强制重绘屏幕
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			m.needsRedraw = true

			m.searchQuery = ""
			m.filterHosts()
		}
	default:
		if ev.Ch != 0 {
			m.searchQuery += string(ev.Ch)
			m.filterHosts()
			m.needsRedraw = true
		}
	}
	return true
}

// 处理正常模式的输入
func (m *Menu) handleNormalInput(ev termbox.Event) bool {
	switch ev.Key {
	case termbox.KeyArrowUp:
		m.moveUp()
	case termbox.KeyArrowDown:
		m.moveDown()
	case termbox.KeyEnter:
		return m.handleEnter()
	case termbox.KeyEsc:
		if m.inGroup {
			m.inGroup = false
			m.currentHost = 0
		} else if m.showFavorites {
			m.showFavorites = false
		} else {
			return false // 退出程序
		}
	case termbox.KeySpace:
		if m.inGroup {
			m.toggleFavorite()
		}
	default:
		// 处理字符输入
		switch ev.Ch {
		case 'q', 'Q':
			return false // 退出
		case 'r', 'R':
			m.reloadConfig()
			m.showToast("配置已重新加载", "success", 2*time.Second)
		case 't', 'T':
			m.toggleTheme()
			m.showToast("主题已切换", "info", 2*time.Second)
		case 'l', 'L':
			m.toggleLayout()
			m.showToast("布局已切换", "info", 2*time.Second)
		case 'f', 'F':
			m.showFavorites = !m.showFavorites
			m.currentHost = 0
		case 's', 'S':
			m.checkAllHostsStatus()
			m.showToast("正在检查主机状态...", "info", 3*time.Second)
		case '/':
			m.searchMode = true
			m.searchQuery = ""
		case '1', '2', '3', '4', '5':
			return m.handleQuickConnect(ev.Ch)
		}
	}
	return true
}

// 处理向上移动
func (m *Menu) moveUp() {
	if m.showFavorites {
		if m.currentHost > 0 {
			m.currentHost--
		}
	} else if m.inGroup {
		if m.currentHost > 0 {
			m.currentHost--
		}
	} else {
		if m.currentGroup > 0 {
			m.currentGroup--
		}
	}
}

// 处理向下移动
func (m *Menu) moveDown() {
	if m.showFavorites {
		favorites := m.getFavoriteHosts()
		if m.currentHost < len(favorites)-1 {
			m.currentHost++
		}
	} else if m.inGroup {
		if m.currentGroup < len(m.filteredGroups) &&
			m.currentHost < len(m.filteredGroups[m.currentGroup].Hosts)-1 {
			m.currentHost++
		}
	} else {
		if m.currentGroup < len(m.filteredGroups)-1 {
			m.currentGroup++
		}
	}
}

// 处理回车键
func (m *Menu) handleEnter() bool {
	if m.showFavorites {
		favorites := m.getFavoriteHosts()
		if m.currentHost < len(favorites) {
			termbox.Close()
			host := favorites[m.currentHost]
			m.connectSSH(host)

			err := termbox.Init()
			if err != nil {
				log.Printf("重新初始化termbox失败: %v", err)
				return false
			}

			// 强制重绘屏幕和重置状态
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			m.needsRedraw = true
			m.renderEngine = NewRenderEngine(m.currentTheme) // 重新初始化渲染引擎
		}
	} else if m.inGroup {
		termbox.Close()
		host := m.filteredGroups[m.currentGroup].Hosts[m.currentHost]
		m.connectSSH(host)

		err := termbox.Init()
		if err != nil {
			log.Printf("重新初始化termbox失败: %v", err)
			return false
		}

		// 强制重绘屏幕和重置状态
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		m.needsRedraw = true
		m.renderEngine = NewRenderEngine(m.currentTheme) // 重新初始化渲染引擎
	} else {
		m.inGroup = true
		m.currentHost = 0
	}
	return true
}

// 处理快速连接
func (m *Menu) handleQuickConnect(ch rune) bool {
	if !m.searchMode && !m.showFavorites && !m.inGroup {
		index := int(ch - '1')
		if index < len(m.connectionHistory) {
			termbox.Close()
			host := m.connectionHistory[index]
			fmt.Printf("\n⚡ 快速连接: %s\n", host.Name)
			m.connectSSH(host)

			err := termbox.Init()
			if err != nil {
				log.Printf("重新初始化termbox失败: %v", err)
				return false
			}

			// 强制重绘屏幕和重置状态
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			m.needsRedraw = true
			m.renderEngine = NewRenderEngine(m.currentTheme) // 重新初始化渲染引擎
		}
	}
	return true
}
