package app

import "github.com/wailsapp/wails/v3/pkg/application"

// SetupSystemTray 设置系统托盘
func (a *App) SetupSystemTray() {
	tray := a.SystemTray.New()
	trayMenu := a.NewMenu()
	onTop := trayMenu.AddCheckbox("置顶", false)
	onTop.OnClick(func(ctx *application.Context) {
		if a.mainWindow != nil {
			// 切换置顶状态
			isChecked := onTop.Checked()
			a.mainWindow.SetAlwaysOnTop(isChecked)
		}
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
