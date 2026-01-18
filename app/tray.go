package app

import "github.com/wailsapp/wails/v3/pkg/application"

// SetupSystemTray 设置系统托盘
func (a *App) SetupSystemTray() {
	tray := a.SystemTray.New()
	trayMenu := a.NewMenu()

	trayMenu.Add("显示").OnClick(func(ctx *application.Context) {
		// TODO: 实现显示窗口逻辑
	})

	trayMenu.Add("设置").OnClick(func(ctx *application.Context) {
		a.OpenConfigWindow()
	})

	trayMenu.AddSeparator()

	trayMenu.Add("退出").OnClick(func(ctx *application.Context) {
		a.Quit()
	})

	tray.SetMenu(trayMenu)
}
