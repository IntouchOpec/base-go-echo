package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// service souce service and service.
type Service struct {
	orm.ModelBase
	SerName          string             `form:"name" json:"ser_name" gorm:"type:varchar(25)"`
	SerDetail        string             `form:"detail" json:"ser_detail" gorm:"type:varchar(25)"`
	SerPrice         float32            `form:"price" json:"ser_price"`
	SerActive        bool               `form:"active" json:"ser_active"`
	AccountID        uint               `form:"account_id" json:"account_id" gorm:"not null;"`
	SerTime          string             `form:"time" json:"ser_time" gorm:"type:varchar(10)"`
	SerImage         string             `form:"image" json:"ser_image" gorm:"type:varchar(255)"`
	Account          Account            `json:"account" gorm:"ForeignKey:AccountID"`
	ChatChannels     []*ChatChannel     `json:"chat_channels" gorm:"many2many:service_chat_channel"`
	ProviderServices []*ProviderService `json:"provider_services" gorm:"many2many:ProviderService"`
}

// Saveservice is function create service.
func (service *Service) Saveservice() *Service {
	if err := DB().Create(&service).Error; err != nil {
		return nil
	}
	return service
}

func (service *Service) Updateservice(id int) *Service {

	if err := DB().Find(&service, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&service).Error; err != nil {
		return nil
	}

	return service
}

// func (subservice *ServiceSlot) CreateServiceSlot() *ServiceSlot {
// 	if err := DB().Create(&subservice).Error; err != nil {
// 		return nil
// 	}
// 	return subservice
// }

// func (subservice *ServiceSlot) UpdateServiceSlot(id int) *ServiceSlot {
// 	if err := DB().Find(&subservice, id).Error; err != nil {
// 		return nil
// 	}

// 	if err := DB().Save(&subservice).Error; err != nil {
// 		return nil
// 	}

// 	return subservice
// }

func GetServiceList(accID uint) (*[]Service, error) {
	services := []Service{}
	if err := DB().Where("ser_account_id = ?", accID).Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}

func GetserviceByID(chatchannelID, id int) *Service {
	service := Service{}
	if err := DB().Where("ChatChannelID = ?", chatchannelID).Find(&service, id).Error; err != nil {
		return nil
	}
	return &service
}

func DeleteserviceByID(id string) *Service {
	service := Service{}
	if err := DB().Find(&service, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&service).Error; err != nil {
		return nil
	}
	return &service
}

// func DeleteServiceSlot(id string) *ServiceSlot {
// 	subservice := ServiceSlot{}
// 	if err := DB().Find(&subservice, id).Error; err != nil {
// 		return nil
// 	}
// 	if err := DB().Delete(&subservice).Error; err != nil {
// 		return nil
// 	}
// 	return &subservice
// }
