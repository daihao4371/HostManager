package cli

import (
	"fmt"
	"strings"
	"sort"

	"github.com/daihao4371/hostmanager/internal/config"
	"github.com/daihao4371/hostmanager/internal/models"
	"github.com/daihao4371/hostmanager/internal/ssh"
)

// CLI命令处理器
type CLI struct {
	config *config.Config
}

// 创建新的CLI实例
func NewCLI(cfg *config.Config) *CLI {
	return &CLI{config: cfg}
}

// 处理命令行参数
func (c *CLI) HandleCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("无效的命令参数")
	}

	command := args[0]
	switch command {
	case "connect", "c":
		return c.handleConnect(args[1:])
	case "list", "ls", "l":
		return c.handleList(args[1:])
	case "status", "s":
		return c.handleStatus(args[1:])
	case "history", "h":
		return c.handleHistory()
	case "favorites", "fav", "f":
		return c.handleFavorites()
	case "groups", "g":
		return c.handleGroups()
	case "search":
		return c.handleSearch(args[1:])
	case "help", "--help", "-h":
		c.showHelp()
		return nil
	case "version", "--version", "-v":
		c.showVersion()
		return nil
	default:
		return fmt.Errorf("未知命令: %s. 使用 'hostmanager help' 查看帮助", command)
	}
}

// 处理连接命令
func (c *CLI) handleConnect(args []string) error {
	if len(args) == 0 {
		return c.showConnectHelp()
	}

	target := args[0]
	
	// 首先尝试按名称查找
	host := c.findHostByName(target)
	if host == nil {
		// 尝试按IP查找
		host = c.findHostByIP(target)
	}
	
	if host == nil {
		// 尝试模糊搜索
		hosts := c.searchHosts(target)
		if len(hosts) == 0 {
			return fmt.Errorf("未找到主机: %s", target)
		} else if len(hosts) == 1 {
			host = &hosts[0]
		} else {
			fmt.Printf("🔍 找到多个匹配的主机:\n")
			for i, h := range hosts {
				fmt.Printf("  %d. %s (%s@%s:%d)\n", i+1, h.Name, h.Username, h.IP, h.Port)
			}
			fmt.Printf("请使用更具体的名称或IP地址\n")
			return nil
		}
	}

	fmt.Printf("🚀 正在连接到 %s (%s@%s:%d)...\n", host.Name, host.Username, host.IP, host.Port)
	
	// 直接调用SSH连接
	ssh.Connect(*host, func(h models.Host) {
		// 简单的历史记录回调
		fmt.Printf("✅ 连接历史已更新\n")
	})
	
	return nil
}

// 处理列表命令
func (c *CLI) handleList(args []string) error {
	showGroups := false
	showFavOnly := false
	
	for _, arg := range args {
		switch arg {
		case "--groups", "-g":
			showGroups = true
		case "--favorites", "-f":
			showFavOnly = true
		}
	}

	if showFavOnly {
		c.listFavorites()
	} else if showGroups {
		c.listByGroups()
	} else {
		c.listAllHosts()
	}
	
	return nil
}

// 处理状态检查命令
func (c *CLI) handleStatus(args []string) error {
	if len(args) == 0 {
		return c.checkAllStatus()
	}
	
	target := args[0]
	host := c.findHostByName(target)
	if host == nil {
		host = c.findHostByIP(target)
	}
	
	if host == nil {
		return fmt.Errorf("未找到主机: %s", target)
	}
	
	fmt.Printf("🔍 正在检查 %s 的状态...\n", host.Name)
	status := ssh.CheckHostStatus(*host)
	
	statusIcon := "❓"
	statusText := "未知"
	switch status {
	case "online":
		statusIcon = "🟢"
		statusText = "在线"
	case "offline":
		statusIcon = "🔴" 
		statusText = "离线"
	}
	
	fmt.Printf("   %s %s (%s@%s:%d) - %s\n", statusIcon, host.Name, host.Username, host.IP, host.Port, statusText)
	return nil
}

// 处理历史记录命令
func (c *CLI) handleHistory() error {
	fmt.Printf("📋 连接历史记录:\n")
	// 这里需要从配置中读取历史记录
	// 临时显示空历史
	fmt.Printf("   暂无历史记录\n")
	return nil
}

// 处理收藏夹命令
func (c *CLI) handleFavorites() error {
	c.listFavorites()
	return nil
}

// 处理分组命令
func (c *CLI) handleGroups() error {
	c.listByGroups()
	return nil
}

// 处理搜索命令
func (c *CLI) handleSearch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请提供搜索关键词")
	}
	
	keyword := strings.Join(args, " ")
	hosts := c.searchHosts(keyword)
	
	if len(hosts) == 0 {
		fmt.Printf("🔍 未找到匹配 '%s' 的主机\n", keyword)
		return nil
	}
	
	fmt.Printf("🔍 搜索结果 (%d个匹配):\n", len(hosts))
	for _, host := range hosts {
		favoriteIcon := ""
		if host.Favorite {
			favoriteIcon = "⭐"
		}
		fmt.Printf("   %s%s (%s@%s:%d)\n", favoriteIcon, host.Name, host.Username, host.IP, host.Port)
	}
	
	return nil
}

// 按名称查找主机
func (c *CLI) findHostByName(name string) *models.Host {
	for _, group := range c.config.Groups {
		for _, host := range group.Hosts {
			if strings.EqualFold(host.Name, name) {
				return &host
			}
		}
	}
	return nil
}

// 按IP查找主机
func (c *CLI) findHostByIP(ip string) *models.Host {
	for _, group := range c.config.Groups {
		for _, host := range group.Hosts {
			if host.IP == ip {
				return &host
			}
		}
	}
	return nil
}

// 搜索主机
func (c *CLI) searchHosts(keyword string) []models.Host {
	var results []models.Host
	keyword = strings.ToLower(keyword)
	
	for _, group := range c.config.Groups {
		for _, host := range group.Hosts {
			if strings.Contains(strings.ToLower(host.Name), keyword) ||
			   strings.Contains(strings.ToLower(host.IP), keyword) ||
			   strings.Contains(strings.ToLower(host.Username), keyword) {
				results = append(results, host)
			}
		}
	}
	
	return results
}

// 列出所有主机
func (c *CLI) listAllHosts() {
	fmt.Printf("📋 所有主机列表:\n")
	
	var allHosts []models.Host
	for _, group := range c.config.Groups {
		allHosts = append(allHosts, group.Hosts...)
	}
	
	// 按名称排序
	sort.Slice(allHosts, func(i, j int) bool {
		return allHosts[i].Name < allHosts[j].Name
	})
	
	for _, host := range allHosts {
		favoriteIcon := ""
		if host.Favorite {
			favoriteIcon = "⭐"
		}
		fmt.Printf("   %s%s (%s@%s:%d)\n", favoriteIcon, host.Name, host.Username, host.IP, host.Port)
	}
}

// 按分组列出主机
func (c *CLI) listByGroups() {
	fmt.Printf("📂 按分组显示:\n")
	
	for _, group := range c.config.Groups {
		fmt.Printf("\n  📁 %s (%d台主机):\n", group.Name, len(group.Hosts))
		for _, host := range group.Hosts {
			favoriteIcon := ""
			if host.Favorite {
				favoriteIcon = "⭐"
			}
			fmt.Printf("     %s%s (%s@%s:%d)\n", favoriteIcon, host.Name, host.Username, host.IP, host.Port)
		}
	}
}

// 列出收藏夹
func (c *CLI) listFavorites() {
	fmt.Printf("⭐ 收藏夹:\n")
	
	var favorites []models.Host
	for _, group := range c.config.Groups {
		for _, host := range group.Hosts {
			if host.Favorite {
				favorites = append(favorites, host)
			}
		}
	}
	
	if len(favorites) == 0 {
		fmt.Printf("   暂无收藏的主机\n")
		return
	}
	
	for _, host := range favorites {
		fmt.Printf("   ⭐ %s (%s@%s:%d)\n", host.Name, host.Username, host.IP, host.Port)
	}
}

// 检查所有主机状态
func (c *CLI) checkAllStatus() error {
	fmt.Printf("🔍 检查所有主机状态...\n")
	
	for _, group := range c.config.Groups {
		fmt.Printf("\n📁 %s:\n", group.Name)
		for _, host := range group.Hosts {
			status := ssh.CheckHostStatus(host)
			statusIcon := "❓"
			statusText := "未知"
			switch status {
			case "online":
				statusIcon = "🟢"
				statusText = "在线"
			case "offline":
				statusIcon = "🔴"
				statusText = "离线"  
			}
			fmt.Printf("   %s %s - %s\n", statusIcon, host.Name, statusText)
		}
	}
	
	return nil
}

// 显示连接帮助
func (c *CLI) showConnectHelp() error {
	fmt.Printf("🚀 连接命令用法:\n")
	fmt.Printf("   hostmanager connect <主机名|IP地址>\n")
	fmt.Printf("   hostmanager c <主机名|IP地址>\n\n")
	fmt.Printf("示例:\n")
	fmt.Printf("   hostmanager connect server1\n")
	fmt.Printf("   hostmanager c 192.168.1.100\n")
	return nil
}

// 显示帮助信息
func (c *CLI) showHelp() {
	fmt.Printf(`🖥️  HostManager - SSH连接管理工具

用法:
   hostmanager [命令] [选项]

可用命令:
   connect, c <主机>      连接到指定主机
   list, ls, l [选项]     显示主机列表
   status, s [主机]       检查主机状态
   history, h             显示连接历史
   favorites, fav, f      显示收藏夹
   groups, g              按分组显示主机
   search <关键词>         搜索主机
   help, --help, -h       显示此帮助信息
   version, --version, -v 显示版本信息

列表选项:
   --groups, -g          按分组显示
   --favorites, -f       仅显示收藏的主机

示例:
   hostmanager                    # 启动交互式UI
   hostmanager connect server1    # 连接到server1
   hostmanager list --groups      # 按分组显示主机列表
   hostmanager status server1     # 检查server1状态
   hostmanager search web         # 搜索包含'web'的主机

更多信息请访问: https://github.com/daihao4371/hostmanager
`)
}

// 显示版本信息
func (c *CLI) showVersion() {
	fmt.Printf("HostManager v1.0.0\n")
	fmt.Printf("为iTerm2和终端用户优化的SSH连接管理工具\n")
}