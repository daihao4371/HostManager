package main

import (
	"log"
	"os"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/cli"
	"github.com/daihao4371/hostmanager/internal/config"
	"github.com/daihao4371/hostmanager/internal/ui"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 检查命令行参数
	args := os.Args[1:] // 去掉程序名

	if len(args) > 0 {
		// CLI模式：有命令行参数时使用命令行接口
		cliHandler := cli.NewCLI(cfg)
		err := cliHandler.HandleCommand(args)
		if err != nil {
			log.Printf("❌ 错误: %v", err)
			os.Exit(1)
		}
	} else {
		// UI模式：无参数时启动交互式全屏界面
		err = termbox.Init()
		if err != nil {
			log.Fatalf("无法初始化termbox: %v", err)
		}
		defer termbox.Close()

		// 创建并运行菜单
		menu := ui.NewMenu(cfg)
		menu.Run()
	}
}