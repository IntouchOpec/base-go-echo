package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type ProviderService struct {
	orm.ModelBase
	PSPrice    float64     `form:"price" json:"ps_price"`
	ProviderID uint        `json:"provider_id"`
	ServiceID  uint        `form:"service_id" json:"service_id"`
	Provider   Provider    `json:"provider" gorm:"ForeignKey:ProviderID"`
	Service    Service     `json:"service" gorm:"ForeignKey:ServiceID"`
	TimeSlots  []*TimeSlot `json:"time_slots"`
	Bookings   []*Booking  `json:"bookings"`
}

func GetProviderServiceDetail(id string) (*ProviderService, error) {
	providerService := ProviderService{}
	if err := DB().Find(&providerService).Error; err != nil {
		return nil, err
	}
	return &providerService, nil
}

func GetProviderServiceDetailByFild(Params ...interface{}) (*ProviderService, error) {
	providerService := ProviderService{}
	if err := DB().Find(&providerService).Error; err != nil {
		return nil, err
	}
	return &providerService, nil
}

func (prov *ProviderService) CreateProviderService() error {
	if err := DB().Create(&prov).Error; err != nil {
		return err
	}
	return nil
}

func (prov *ProviderService) UpdateProviderService() error {
	if err := DB().Save(&prov).Error; err != nil {
		return err
	}
	return nil
}

func RemoveProviderService(id string) (*ProviderService, error) {
	prov := ProviderService{}
	if err := DB().Find(&prov, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&prov).Error; err != nil {
		return nil, err
	}
	return &prov, nil
}
