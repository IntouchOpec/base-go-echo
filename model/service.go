package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// service souce service and service.
type Service struct {
	orm.ModelBase
	SerName      string         `form:"name" json:"ser_name" gorm:"type:varchar(25)"`
	SerDetail    string         `form:"detail" json:"ser_detail" gorm:"type:varchar(25)"`
	SerPrice     float32        `form:"price" json:"ser_price"`
	SerActive    bool           `form:"active" json:"ser_active"`
	SerAccountID uint           `form:"account_id" json:"ser_account_id" gorm:"not null;"`
	SerImage     string         `form:"image" json:"ser_image" gorm:"type:varchar(255)"`
	Account      Account        `json:"account" gorm:"ForeignKey:id"`
	ChatChannels []*ChatChannel `json:"chat_channels" gorm:"many2many:service_chat_channel"`
}

// ServiceSlot service set.
// type ServiceSlot struct {
// 	orm.ModelBase
// 	Start     string     `json:"start"`
// 	End       string     `json:"end"`
// 	Day       int        `json:"day"`
// 	Amount    int        `json:"amount"`
// 	Active    bool       `json:"active"`
// 	ServiceID uint       `json:"service_id"`
// 	Bookings  []*Booking `json:"bookings"`
// 	Service   Service    `json:"service" gorm:"ForeignKey:ServiceID;"`
// }

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

func Getservice(chatchannelID int) *[]Service {
	services := []Service{}
	if err := DB().Where("").Find(&services).Error; err != nil {
		return nil
	}
	return &services
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
