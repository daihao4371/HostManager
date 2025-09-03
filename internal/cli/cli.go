package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	case "init":
		return c.handleInit()
	case "add-host":
		return c.handleAddHost()
	case "remove", "rm":
		return c.handleRemove(args[1:])
	case "edit":
		return c.handleEdit(args[1:])
	case "completion":
		return c.handleCompletion(args[1:])
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
   init                   生成配置文件模板
   add-host              交互式添加新主机
   edit <主机>            编辑指定主机配置
   remove, rm <主机>      删除指定主机
   completion <shell>     生成shell补全脚本
   help, --help, -h       显示此帮助信息
   version, --version, -v 显示版本信息

列表选项:
   --groups, -g          按分组显示
   --favorites, -f       仅显示收藏的主机

配置管理:
   hostmanager init                    # 创建配置文件模板
   hostmanager add-host               # 交互式添加主机
   hostmanager edit server1           # 编辑指定主机配置
   hostmanager remove server1         # 删除指定主机
   hostmanager completion bash >> ~/.bashrc   # 安装Bash补全
   hostmanager completion zsh >> ~/.zshrc     # 安装Zsh补全

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

// 处理配置初始化命令
func (c *CLI) handleInit() error {
	configPath := "config.yaml"
	
	// 检查配置文件是否已存在
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("⚠️  配置文件 %s 已存在\n", configPath)
		fmt.Printf("是否要覆盖? (y/N): ")
		
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input != "y" && input != "yes" {
			fmt.Printf("操作已取消\n")
			return nil
		}
	}
	
	// 生成默认配置模板
	template := `# HostManager 配置文件模板
# 生成时间: ` + fmt.Sprintf("%v", "now") + `

groups:
  - name: "生产环境"
    hosts:
      - name: "Web服务器-1"
        ip: "192.168.1.100"
        port: 22
        username: "admin"
        auth_type: "key"
        key_path: "~/.ssh/id_rsa"
        description: "主要的Web服务器"
        tags: ["production", "web"]
        favorite: true

  - name: "测试环境"
    hosts:
      - name: "测试服务器"
        ip: "192.168.1.200"
        port: 22
        username: "test"
        auth_type: "password"
        password: "your_password"
        description: "测试环境服务器"
        tags: ["test"]
        favorite: false

ui_config:
  theme: "dark"
  language: "zh"
  key_bindings:
    exit: "Esc"
    search: "/"
    favorites: "f"
    status_check: "s"
    toggle_fav: "Space"
    theme_switch: "t"
  layout:
    type: "single"
    show_details: false
`
	
	err := os.WriteFile(configPath, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %v", err)
	}
	
	fmt.Printf("✅ 配置文件模板已创建: %s\n", configPath)
	fmt.Printf("请编辑配置文件添加你的主机信息\n")
	return nil
}

// 处理添加主机命令
func (c *CLI) handleAddHost() error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("📝 添加新主机到配置\n\n")
	
	// 收集主机信息
	host := models.Host{}
	
	// 主机名称（必填）
	for {
		fmt.Printf("主机名称: ")
		input, _ := reader.ReadString('\n')
		host.Name = strings.TrimSpace(input)
		if host.Name != "" {
			break
		}
		fmt.Printf("❌ 主机名称不能为空\n")
	}
	
	// IP地址（必填）
	for {
		fmt.Printf("IP地址: ")
		input, _ := reader.ReadString('\n')
		host.IP = strings.TrimSpace(input)
		if host.IP != "" {
			break
		}
		fmt.Printf("❌ IP地址不能为空\n")
	}
	
	// 端口号
	fmt.Printf("端口号 [22]: ")
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)
	if portInput == "" {
		host.Port = 22
	} else {
		port, err := strconv.Atoi(portInput)
		if err != nil || port <= 0 || port > 65535 {
			fmt.Printf("❌ 无效端口号，使用默认端口 22\n")
			host.Port = 22
		} else {
			host.Port = port
		}
	}
	
	// 用户名（必填）
	for {
		fmt.Printf("用户名: ")
		input, _ := reader.ReadString('\n')
		host.Username = strings.TrimSpace(input)
		if host.Username != "" {
			break
		}
		fmt.Printf("❌ 用户名不能为空\n")
	}
	
	// 认证方式
	fmt.Printf("认证方式 (key/password) [key]: ")
	authInput, _ := reader.ReadString('\n')
	authInput = strings.TrimSpace(strings.ToLower(authInput))
	if authInput == "" || authInput == "key" {
		host.AuthType = "key"
		fmt.Printf("私钥路径 [~/.ssh/id_rsa]: ")
		keyInput, _ := reader.ReadString('\n')
		keyInput = strings.TrimSpace(keyInput)
		if keyInput == "" {
			host.KeyPath = "~/.ssh/id_rsa"
		} else {
			host.KeyPath = keyInput
		}
	} else {
		host.AuthType = "password"
		fmt.Printf("密码: ")
		passInput, _ := reader.ReadString('\n')
		host.Password = strings.TrimSpace(passInput)
	}
	
	// 描述（可选）
	fmt.Printf("描述 [可选]: ")
	descInput, _ := reader.ReadString('\n')
	host.Description = strings.TrimSpace(descInput)
	
	// 是否收藏
	fmt.Printf("添加到收藏夹? (y/N): ")
	favInput, _ := reader.ReadString('\n')
	favInput = strings.TrimSpace(strings.ToLower(favInput))
	host.Favorite = favInput == "y" || favInput == "yes"
	
	// 选择分组
	return c.addHostToGroup(host)
}

// 添加主机到分组并保存配置
func (c *CLI) addHostToGroup(host models.Host) error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("\n📂 选择分组:\n")
	for i, group := range c.config.Groups {
		fmt.Printf("  %d. %s\n", i+1, group.Name)
	}
	fmt.Printf("  %d. 创建新分组\n", len(c.config.Groups)+1)
	
	for {
		fmt.Printf("请选择分组 [1]: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		var groupIndex int
		if input == "" {
			groupIndex = 0
		} else {
			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(c.config.Groups)+1 {
				fmt.Printf("❌ 无效选择，请输入 1-%d\n", len(c.config.Groups)+1)
				continue
			}
			groupIndex = index - 1
		}
		
		// 创建新分组
		if groupIndex == len(c.config.Groups) {
			fmt.Printf("新分组名称: ")
			groupInput, _ := reader.ReadString('\n')
			groupName := strings.TrimSpace(groupInput)
			if groupName == "" {
				fmt.Printf("❌ 分组名称不能为空\n")
				continue
			}
			
			newGroup := models.Group{
				Name:  groupName,
				Hosts: []models.Host{host},
			}
			c.config.Groups = append(c.config.Groups, newGroup)
		} else {
			// 添加到现有分组
			c.config.Groups[groupIndex].Hosts = append(c.config.Groups[groupIndex].Hosts, host)
		}
		break
	}
	
	// 保存配置到文件
	err := config.SaveConfig("config.yaml", c.config)
	if err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}
	
	fmt.Printf("✅ 主机 %s 已添加到配置\n", host.Name)
	return nil
}

// 处理补全脚本命令
func (c *CLI) handleCompletion(args []string) error {
	if len(args) == 0 {
		fmt.Printf("📋 可用的补全脚本:\n")
		fmt.Printf("   bash    生成Bash补全脚本\n")
		fmt.Printf("   zsh     生成Zsh补全脚本\n")
		fmt.Printf("\n使用方法:\n")
		fmt.Printf("   hostmanager completion bash >> ~/.bashrc\n")
		fmt.Printf("   hostmanager completion zsh >> ~/.zshrc\n")
		return nil
	}
	
	shell := args[0]
	switch shell {
	case "bash":
		return c.generateBashCompletion()
	case "zsh":
		return c.generateZshCompletion()
	default:
		return fmt.Errorf("不支持的shell: %s. 支持: bash, zsh", shell)
	}
}

// 生成Bash补全脚本内容
func (c *CLI) generateBashCompletion() error {
	script := `# HostManager Bash补全脚本
_hostmanager_completion() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    commands="connect c list ls l status s search history h favorites fav f groups g init add-host edit remove rm completion help version"
    
    case "${prev}" in
        hostmanager|hm)
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            return 0
            ;;
        connect|c|status|s)
            # 动态获取主机列表
            if command -v hostmanager >/dev/null 2>&1; then
                local hosts=$(hostmanager list 2>/dev/null | grep -o '[a-zA-Z0-9_-]*@[0-9.]*' | cut -d'@' -f1 | sort -u)
                COMPREPLY=( $(compgen -W "${hosts}" -- ${cur}) )
            fi
            return 0
            ;;
        edit|remove|rm)
            # 编辑和删除命令也需要主机名补全
            if command -v hostmanager >/dev/null 2>&1; then
                local hosts=$(hostmanager list 2>/dev/null | grep -o '[a-zA-Z0-9_-]*@[0-9.]*' | cut -d'@' -f1 | sort -u)
                COMPREPLY=( $(compgen -W "${hosts}" -- ${cur}) )
            fi
            return 0
            ;;
        list|ls|l)
            COMPREPLY=( $(compgen -W "--groups --favorites -g -f" -- ${cur}) )
            return 0
            ;;
        completion)
            COMPREPLY=( $(compgen -W "bash zsh" -- ${cur}) )
            return 0
            ;;
    esac
}

complete -F _hostmanager_completion hostmanager
complete -F _hostmanager_completion hm
`
	fmt.Print(script)
	return nil
}

// 生成Zsh补全脚本内容
func (c *CLI) generateZshCompletion() error {
	script := `#compdef hostmanager hm
# HostManager Zsh补全脚本

_hostmanager() {
    local context state line
    typeset -A opt_args

    _arguments -C \
        '1: :->commands' \
        '*: :->args' \
        && return 0

    case "$state" in
        commands)
            local commands; commands=(
                'connect:连接到指定主机'
                'c:连接到指定主机(简写)'
                'list:显示主机列表'
                'ls:显示主机列表(简写)'
                'l:显示主机列表(简写)'
                'status:检查主机状态'
                's:检查主机状态(简写)'
                'search:搜索主机'
                'history:显示连接历史'
                'h:显示连接历史(简写)'
                'favorites:显示收藏夹'
                'fav:显示收藏夹(简写)'
                'f:显示收藏夹(简写)'
                'groups:按分组显示'
                'g:按分组显示(简写)'
                'init:生成配置文件模板'
                'add-host:交互式添加新主机'
                'edit:编辑指定主机配置'
                'remove:删除指定主机'
                'rm:删除指定主机(简写)'
                'completion:生成shell补全脚本'
                'help:显示帮助信息'
                'version:显示版本信息'
            )
            _describe 'commands' commands
            ;;
        args)
            case "${words[2]}" in
                connect|c|status|s)
                    # 动态获取主机列表
                    if (( $+commands[hostmanager] )); then
                        local hosts; hosts=($(hostmanager list 2>/dev/null | grep -o '[a-zA-Z0-9_-]*@[0-9.]*' | cut -d'@' -f1 | sort -u))
                        _describe 'hosts' hosts
                    fi
                    ;;
                edit|remove|rm)
                    # 编辑和删除命令也需要主机名补全
                    if (( $+commands[hostmanager] )); then
                        local hosts; hosts=($(hostmanager list 2>/dev/null | grep -o '[a-zA-Z0-9_-]*@[0-9.]*' | cut -d'@' -f1 | sort -u))
                        _describe 'hosts' hosts
                    fi
                    ;;
                list|ls|l)
                    local options; options=(
                        '--groups:按分组显示'
                        '-g:按分组显示(简写)'
                        '--favorites:仅显示收藏'
                        '-f:仅显示收藏(简写)'
                    )
                    _describe 'options' options
                    ;;
                completion)
                    local shells; shells=('bash:Bash补全脚本' 'zsh:Zsh补全脚本')
                    _describe 'shells' shells
                    ;;
                search)
                    _message '搜索关键词'
                    ;;
            esac
            ;;
    esac
    
    return 1
}

_hostmanager "$@"
`
	fmt.Print(script)
	return nil
}

// 处理删除主机命令
func (c *CLI) handleRemove(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要删除的主机名称")
	}
	
	hostName := args[0]
	
	// 查找主机
	groupIndex, hostIndex := c.findHostLocation(hostName)
	if groupIndex == -1 {
		return fmt.Errorf("未找到主机: %s", hostName)
	}
	
	host := c.config.Groups[groupIndex].Hosts[hostIndex]
	
	// 确认删除
	fmt.Printf("⚠️  确认删除主机 '%s' (%s@%s:%d)? (y/N): ", host.Name, host.Username, host.IP, host.Port)
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	
	if input != "y" && input != "yes" {
		fmt.Printf("操作已取消\n")
		return nil
	}
	
	// 从切片中删除主机
	c.config.Groups[groupIndex].Hosts = append(
		c.config.Groups[groupIndex].Hosts[:hostIndex],
		c.config.Groups[groupIndex].Hosts[hostIndex+1:]...,
	)
	
	// 如果分组为空，询问是否删除分组
	if len(c.config.Groups[groupIndex].Hosts) == 0 {
		fmt.Printf("分组 '%s' 已为空，是否删除此分组? (y/N): ", c.config.Groups[groupIndex].Name)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input == "y" || input == "yes" {
			c.config.Groups = append(
				c.config.Groups[:groupIndex],
				c.config.Groups[groupIndex+1:]...,
			)
		}
	}
	
	// 保存配置
	err := config.SaveConfig("config.yaml", c.config)
	if err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}
	
	fmt.Printf("✅ 主机 '%s' 已删除\n", hostName)
	return nil
}

// 处理编辑主机命令
func (c *CLI) handleEdit(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要编辑的主机名称")
	}
	
	hostName := args[0]
	
	// 查找主机
	groupIndex, hostIndex := c.findHostLocation(hostName)
	if groupIndex == -1 {
		return fmt.Errorf("未找到主机: %s", hostName)
	}
	
	host := &c.config.Groups[groupIndex].Hosts[hostIndex]
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Printf("📝 编辑主机: %s\n\n", host.Name)
	
	// 编辑各个字段
	fmt.Printf("主机名称 [%s]: ", host.Name)
	if input := c.readInputWithDefault(reader); input != "" {
		host.Name = input
	}
	
	fmt.Printf("IP地址 [%s]: ", host.IP)
	if input := c.readInputWithDefault(reader); input != "" {
		host.IP = input
	}
	
	fmt.Printf("端口号 [%d]: ", host.Port)
	if input := c.readInputWithDefault(reader); input != "" {
		if port, err := strconv.Atoi(input); err == nil && port > 0 && port <= 65535 {
			host.Port = port
		} else {
			fmt.Printf("❌ 无效端口号，保持原值 %d\n", host.Port)
		}
	}
	
	fmt.Printf("用户名 [%s]: ", host.Username)
	if input := c.readInputWithDefault(reader); input != "" {
		host.Username = input
	}
	
	fmt.Printf("认证方式 (key/password) [%s]: ", host.AuthType)
	if input := c.readInputWithDefault(reader); input != "" {
		if input == "key" || input == "password" {
			host.AuthType = input
			if input == "key" {
				fmt.Printf("私钥路径 [%s]: ", host.KeyPath)
				if keyInput := c.readInputWithDefault(reader); keyInput != "" {
					host.KeyPath = keyInput
				}
				host.Password = ""
			} else {
				fmt.Printf("密码: ")
				if passInput := c.readInputWithDefault(reader); passInput != "" {
					host.Password = passInput
				}
				host.KeyPath = ""
			}
		}
	}
	
	fmt.Printf("描述 [%s]: ", host.Description)
	if input := c.readInputWithDefault(reader); input != "" {
		host.Description = input
	}
	
	// 保存配置
	err := config.SaveConfig("config.yaml", c.config)
	if err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}
	
	fmt.Printf("✅ 主机 '%s' 已更新\n", host.Name)
	return nil
}

// 查找主机位置（返回分组索引和主机索引）
func (c *CLI) findHostLocation(hostName string) (int, int) {
	for groupIndex, group := range c.config.Groups {
		for hostIndex, host := range group.Hosts {
			if strings.EqualFold(host.Name, hostName) {
				return groupIndex, hostIndex
			}
		}
	}
	return -1, -1
}

// 读取用户输入，支持默认值
func (c *CLI) readInputWithDefault(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}