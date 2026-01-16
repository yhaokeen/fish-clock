package app

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// CreateMainWindow 创建主窗口
func (a *App) CreateMainWindow() {
	a.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "FishClock",
		Width:          380,
		Height:         200,
		Frameless:      true,
		DisableResize:  true,
		BackgroundType: application.BackgroundTypeTranslucent,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			BackdropType: application.Acrylic,
		},
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
	})
}

// OpenSettingsWindow 打开设置窗口
func (a *App) OpenSettingsWindow() {
	// 如果设置窗口已存在且有效，显示它
	if a.settingsWindow != nil {
		a.settingsWindow.Show()
		a.settingsWindow.Focus()
		return
	}

	// 创建新的设置窗口
	a.settingsWindow = a.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "设置",
		Width:            500,
		Height:           450,
		DisableResize:    false,
		BackgroundType:   application.BackgroundTypeSolid,
		BackgroundColour: application.NewRGB(240, 240, 240),
		URL:              "/#/settings",
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropNormal,
		},
	})

	// 监听窗口关闭事件，清除引用
	a.settingsWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		a.settingsWindow = nil
	})
}
