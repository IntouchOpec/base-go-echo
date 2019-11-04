package model

import (
	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// service souce service and service.
type Service struct {
	orm.ModelBase
	Name         string         `form:"name" json:"name" gorm:"type:varchar(25)"`
	Detail       string         `form:"detail" json:"detail" gorm:"type:varchar(25)"`
	Price        float32        `form:"price" json:"price"`
	Active       bool           `form:"active" json:"active"`
	AccountID    uint           `form:"account_id" json:"account_id" gorm:"not null;"`
	Account      Account        `gorm:"ForeignKey:id"`
	Image        string         `form:"image" json:"image" gorm:"type:varchar(255)"`
	ChatChannels []*ChatChannel `json:"chat_channels" gorm:"many2many:service_chat_channel"`
	ServiceSlots []*ServiceSlot `gorm:"ForeignKey:ServiceID;" json:"sub_services"`
	Promotions   []*Promotion   `json:"promotions" gorm:"many2many:service_promotion;"`
}

// ServiceSlot service set.
type ServiceSlot struct {
	orm.ModelBase
	Start     string     `json:"start"`
	End       string     `json:"end"`
	Day       int        `json:"day"`
	Amount    int        `json:"amount"`
	Active    bool       `json:"active"`
	ServiceID uint       `json:"service_id"`
	Bookings  []*Booking `json:"bookings"`
	Service   Service    `json:"service" gorm:"ForeignKey:ServiceID;"`
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

func (subservice *ServiceSlot) CreateServiceSlot() *ServiceSlot {
	if err := DB().Create(&subservice).Error; err != nil {
		return nil
	}
	return subservice
}

func (subservice *ServiceSlot) UpdateServiceSlot(id int) *ServiceSlot {
	if err := DB().Find(&subservice, id).Error; err != nil {
		return nil
	}

	if err := DB().Save(&subservice).Error; err != nil {
		return nil
	}

	return subservice
}

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

func DeleteServiceSlot(id string) *ServiceSlot {
	subservice := ServiceSlot{}
	if err := DB().Find(&subservice, id).Error; err != nil {
		return nil
	}
	if err := DB().Delete(&subservice).Error; err != nil {
		return nil
	}
	return &subservice
}
