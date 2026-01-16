package app

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	WorkStartHour int     `json:"workStartHour"` // 上班时间（小时）
	WorkEndHour   int     `json:"workEndHour"`   // 下班时间（小时）
	MonthlySalary int     `json:"monthlySalary"` // 月薪
	Payday        int     `json:"payday"`        // 发薪日
	Opacity       float64 `json:"opacity"`       // 窗口透明度 0-1
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		WorkStartHour: 9,
		WorkEndHour:   18,
		MonthlySalary: 15000,
		Payday:        25,
		Opacity:       0.9,
	}
}

// ConfigService 配置服务
type ConfigService struct {
	config     *Config
	configPath string
}

// NewConfigService 创建配置服务
func NewConfigService() *ConfigService {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".fish-clock", "config.json")

	service := &ConfigService{
		config:     DefaultConfig(),
		configPath: configPath,
	}

	service.Load()
	return service
}

// Load 加载配置
func (s *ConfigService) Load() error {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		// 配置文件不存在，使用默认配置并保存
		log.Printf("配置文件不存在，使用默认配置: %s", s.configPath)
		return s.Save()
	}

	if err := json.Unmarshal(data, s.config); err != nil {
		log.Printf("配置文件解析失败，使用默认配置: %v", err)
		s.config = DefaultConfig()
		return s.Save()
	}

	return nil
}

// Save 保存配置
func (s *ConfigService) Save() error {
	// 确保目录存在
	dir := filepath.Dir(s.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.configPath, data, 0644)
}

// GetConfig 获取配置
func (s *ConfigService) GetConfig() *Config {
	return s.config
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(config *Config) error {
	s.config = config
	if err := s.Save(); err != nil {
		log.Printf("保存配置失败: %v", err)
		return err
	}
	log.Printf("配置已保存: %+v", config)
	return nil
}
