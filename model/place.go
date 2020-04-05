package model

import (
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type PlaceType string

const (
	PlaceRoom PlaceType = "room"
	// PlaceRoom PlaceType = "room"
)

const formatDate string = "2006-01-02 15:04"

type Place struct {
	orm.ModelBase

	PlacName     string         `json:"plac_name" form:"name" gorm:"type:varchar(50)"`
	PlacDetail   string         `json:"plac_detail" form:"detail"`
	PlacActive   bool           `json:"plac_active" form:"active" sql:"default:true"`
	PlacType     string         `form:"type" json:"plac_type" gorm:"type:varchar(50)"`
	PlacAmount   int            `form:"amount" json:"plac_amount"`
	PlacImage    string         `form:"image" json:"plac_image" gorm:"type:varchar(255)"`
	ChatChannels []*ChatChannel `json:"chat_channels" gorm:"many2many:place_chat_channel"`
	AccountID    uint           `json:"account_id"`
	Services     []*Service     `json:"services" gorm:"many2many:place_service"`
	Account      Account        `json:"account" gorm:"ForeignKey:AccountID"`
}

type PlaceService struct {
	Place     []Place
	ServiceID uint
}

func (pla *Place) CreatePlace() error {
	if err := DB().Create(&pla).Error; err != nil {
		return err
	}
	return nil
}

func (pla *Place) Update() error {
	if err := DB().Save(&pla).Error; err != nil {
		return err
	}
	return nil
}

func GetPlaceList(AccountID uint) ([]*Place, error) {
	places := []*Place{}
	if err := DB().Where("AccountID = ?", AccountID).Find(&places).Error; err != nil {
		return nil, err
	}
	return places, nil
}

func GetPlaceDetail(id string, AccountID uint) (*Place, error) {
	place := Place{}
	if err := DB().Where("account_id = ?", AccountID).Find(&place, id).Error; err != nil {
		return nil, err
	}
	return &place, nil
}

func DeletePlaceByID(id string) (*Place, error) {
	place := Place{}

	if err := DB().Find(&place, id).Error; err != nil {
		return nil, err
	}

	if err := DB().Delete(&place).Error; err != nil {
		return nil, err
	}

	return &place, nil
}

func (pla *Place) RemoveImage() error {
	pla.PlacImage = ""
	if err := DB().Save(&pla).Error; err != nil {
		return err
	}
	return nil
}

const RowDur time.Duration = 15 * time.Minute

func MakeHour(d time.Time) (time.Time, error) {
	ho := d.Hour() + 1
	d, err := time.Parse("2006-01-02", "0001-01-01")
	if err != nil {
		return time.Now(), err
	}
	return d.Add(time.Duration(ho) * time.Hour), nil
}

func MakeTime(d time.Time) (time.Time, error) {
	ho := d.Hour()
	mi := d.Minute()
	if mi > 0 && mi <= 15 {
		mi = 45
	} else if mi > 15 && mi <= 30 {
		mi = 00
		ho++
	} else if mi > 30 && mi <= 45 {
		ho++
		mi = 30
	} else {
		ho++
		mi = 45
	}
	d, err := time.Parse("2006-01-02", "0001-01-01")
	if err != nil {
		return time.Now(), err
	}
	return d.Add((time.Duration(mi) * time.Minute) + (time.Duration(ho) * time.Hour)), nil
}

func MakeTimeStartAndTimeEnd(d time.Time, timeUsed time.Duration) (time.Time, time.Time, error) {
	d, err := MakeTime(d)
	if err != nil {
		return d, d, err
	}
	return d, d.Add(timeUsed), nil
}

func SetPlaces(places []*Place) []uint {
	var plaIDs []uint
	for _, pla := range places {
		plaIDs = append(plaIDs, pla.ID)
	}
	return plaIDs
}
