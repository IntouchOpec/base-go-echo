package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

// service souce service and service.
type Service struct {
	orm.ModelBase
	SerName          string             `form:"name" json:"ser_name" gorm:"type:varchar(25)"`
	SerDetail        string             `form:"detail" json:"ser_detail" gorm:"type:varchar(25)"`
	SerPrice         float64            `form:"price" json:"ser_price"`
	SerTime          string             `form:"time" json:"ser_time" gorm:"type:varchar(10)"`
	SerActive        bool               `form:"active" json:"ser_active"`
	SerImage         string             `form:"image" json:"ser_image" gorm:"type:varchar(255)"`
	AccountID        uint               `form:"account_id" json:"account_id" gorm:"not null;"`
	ServiceItems     []*ServiceItem     `json:"service_items"`
	Places           []*Place           `json:"places" gorm:"many2many:place_service"`
	Account          *Account           `json:"account" gorm:"ForeignKey:AccountID"`
	ChatChannels     []*ChatChannel     `json:"chat_channels" gorm:"many2many:service_chat_channel"`
	EmployeeServices []*EmployeeService `json:"employee_services" gorm:"many2many:EmployeeService"`
}

type ServiceItem struct {
	orm.ModelBase
	SSTime     time.Time  `form:"s_s_time" json:"s_s_time"`
	SSPrice    float64    `form:"s_s_price" json:"s_s_price"`
	SSName     string     `form:"s_s_name" json:"s_s_name" gorm:"type:varchar(10)"`
	SSIsActive bool       `json:"s_s_is_active" form:"s_s_is_active" sql:"default:false"`
	ServiceID  uint       `json:"service_id"`
	Packages   []*Package `json:"packages" gorm:"many2many:package_service_item"`
	Service    Service    `json:"service" gorm:"ForeignKey:ServiceID"`
	AccountID  uint       `json:"account_id" gorm:"not null;"`
	Account    Account    `json:"account" gorm:"ForeignKey:AccountID"`
}

// Saveservice is function create service.
func (service *Service) SaveService() error {
	if err := DB().Create(&service).Error; err != nil {
		return err
	}
	return nil
}

func (service *Service) UpdateService(id string) error {

	if err := DB().Save(&service).Error; err != nil {
		return err
	}

	return nil
}

func (service *Service) RemoveImage(id string) error {
	db := DB()
	service.SerImage = ""

	if err := db.Save(&service).Error; err != nil {
		return err
	}

	return nil
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
	if err := DB().Where("account_id = ?", accID).Find(&services).Error; err != nil {
		return nil, err
	}
	return &services, nil
}

func GetserviceByID(chatchannelID, id int) *Service {
	service := Service{}
	if err := DB().Where("chat_channel_id = ?", chatchannelID).Find(&service, id).Error; err != nil {
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
