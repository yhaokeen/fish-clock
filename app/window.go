package app

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// CreateMainWindow 创建主窗口
func (a *App) CreateMainWindow() {
	a.mainWindow = a.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "FishClock",
		Width:          380,
		Height:         200,
		Frameless:      true,
		DisableResize:  true,
		BackgroundType: application.BackgroundTypeTransparent,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			// BackdropType: application.Acrylic,
		},
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
	})
}

// OpenSettingsWindow 打开设置窗口
func (a *App) OpenConfigWindow() {
	// 如果设置窗口已存在且有效，显示它
	if a.configWindow != nil {
		a.configWindow.Show()
		a.configWindow.Focus()
		return
	}

	a.createConfigWindow()

	// 监听窗口关闭事件，清除引用
	a.configWindow.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		a.configWindow = nil
	})
}

func (a *App) createConfigWindow() {
	a.configWindow = a.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "设置",
		Width:            500,
		Height:           450,
		Frameless:        true,
		DisableResize:    true,
		BackgroundType:   application.BackgroundTypeTranslucent,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/#/settings",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{
			BackdropType: application.Acrylic,
		},
	})
}
