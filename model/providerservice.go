package model

type ProviderService struct {
	ProviderID uint        `json:"provider_id"`
	ServiceID  uint        `json:"service_id"`
	Providers  []*Provider `json:"providers" gorm:"ForeignKey:ProviderID"`
	Services   []*Service  `json:"dervice" gorm:"ForeignKey:ProviderID"`
	PSPrice    float64     `json:"ps_price"`
}
