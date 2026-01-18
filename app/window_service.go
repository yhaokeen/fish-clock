package app

// WindowService 窗口服务
type WindowService struct {
	app *App
}

// NewWindowService 创建窗口服务
func NewWindowService(app *App) *WindowService {
	return &WindowService{app: app}
}

// CloseCurrentWindow 关闭当前窗口
func (s *WindowService) CloseCurrentWindow() {
	if s.app.configWindow != nil {
		s.app.configWindow.Close()
	}
}
