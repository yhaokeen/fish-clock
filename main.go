package main

import (
	"embed"
	"fish-clock/app"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 创建应用
	application := app.New(assets)

	// 创建主窗口
	application.CreateMainWindow()

	// 设置系统托盘
	application.SetupSystemTray()

	// 运行应用
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
