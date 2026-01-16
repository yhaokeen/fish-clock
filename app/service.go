package app

import (
	"fish-clock/internal/festival"
)

// FestivalService 节日服务
type FestivalService struct{}

// FestivalInfo 节日信息结构（返回给前端）
type FestivalInfo struct {
	Name string `json:"name"`
	Days int    `json:"days"`
	Type string `json:"type"`
}

// GetNextFestival 获取下一个节日
func (s *FestivalService) GetNextFestival() *FestivalInfo {
	// 从今天开始查找60天内最近的节日
	f := festival.Today().GetNearestFestival(60)
	if f == nil {
		return &FestivalInfo{
			Name: "无",
			Days: 0,
			Type: "",
		}
	}

	// 计算距离天数
	today := festival.Today()
	days := int(f.SolarDay.GetJulianDay().Subtract(today.GetJulianDay()))

	return &FestivalInfo{
		Name: f.Name,
		Days: days,
		Type: f.Type.String(),
	}
}
