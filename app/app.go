package app

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App 应用实例包装
type App struct {
	*application.App
	ConfigService *ConfigService
	configWindow  *application.WebviewWindow
}

// New 创建新的应用实例
func New(assets embed.FS) *App {
	configService := NewConfigService()

	appInstance := &App{
		ConfigService: configService,
	}

	windowService := NewWindowService(appInstance)

	app := application.New(application.Options{
		Name: "fish-clock",
		Services: []application.Service{
			application.NewService(&FestivalService{}),
			application.NewService(configService),
			application.NewService(windowService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	appInstance.App = app
	return appInstance
}
