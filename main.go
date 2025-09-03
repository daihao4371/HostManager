package main

import (
	"log"

	"github.com/nsf/termbox-go"

	"github.com/daihao4371/hostmanager/internal/config"
	"github.com/daihao4371/hostmanager/internal/ui"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 初始化终端界面
	err = termbox.Init()
	if err != nil {
		log.Fatalf("无法初始化termbox: %v", err)
	}
	defer termbox.Close()

	// 创建并运行菜单
	menu := ui.NewMenu(cfg)
	menu.Run()
}