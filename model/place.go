package model

import (
	"errors"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
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
	MasterPlaces []*MasterPlace `json:"master_places"`
}

type MasterPlace struct {
	orm.ModelBase

	PlaceID   uint      `json:"place_id"`
	MPlaQue   int       `json:"m_pla_que"`
	AccountID uint      `json:"account_id"`
	MPlaDay   time.Time `json:"m_pla_day"`
	MPlaFrom  time.Time `json:"m_pla_from"`
	MPlaTo    time.Time `json:"m_pla_to"`
	Place     *Place    `json:"place" gorm:"ForeignKey:PlaceID"`
	Account   Account   `json:"account" gorm:"ForeignKey:AccountID"`
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

func (mpla *MasterPlace) CreateMasterPlace() (*MasterPlace, error) {
	if err := DB().Create(&mpla).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func (mpla *MasterPlace) UpdateMasterPlace() error {
	if err := DB().Save(&mpla).Error; err != nil {
		return err
	}
	return nil
}

func GetMasterPlace(from, to, day time.Time) (*MasterPlace, error) {
	startDate, _ := time.Parse(formatDate, "2019-07-09")
	endDate, _ := time.Parse(formatDate, "2019-07-15")

	mpla := MasterPlace{}
	if err := DB().Where("mpla_day = ? and mpla_to >= ? and mpla_to <= ? and mpla_from >= ? and mpla_from <= ?",
		day, startDate, startDate, endDate, endDate).Find(&mpla).Error; err != nil {
		return nil, err
	}
	return &mpla, nil
}

func (pla *Place) RemoveImage() error {
	pla.PlacImage = ""
	if err := DB().Save(&pla).Error; err != nil {
		return err
	}
	return nil
}

type Plas []*Place

func (plas Plas) GetPlaceEmpty(mpla MasterPlace, db *gorm.DB) (*Place, error) {
	mplas, err := mpla.GetMPlas(plas, db)
	if err != nil {
		return nil, err
	}
	pla := mplas.GetEmptyPlace(plas)
	if err != nil {
		return nil, errors.New("")
	}
	return pla, nil
}

type MPlas []MasterPlace

func (mpla MasterPlace) GetMPlas(plas []*Place, db *gorm.DB) (MPlas, error) {
	var mplas MPlas
	var placeIDs []uint
	for _, pla := range plas {
		placeIDs = append(placeIDs, pla.ID)
	}
	if err := db.Where("place_id in (?) and account_id = ? and m_pla_day = ? and m_pla_from BETWEEN ? and ? or m_pla_to BETWEEN ? and ?",
		placeIDs,
		mpla.AccountID,
		mpla.MPlaDay,
		mpla.MPlaFrom,
		mpla.MPlaTo,
		mpla.MPlaFrom,
		mpla.MPlaTo).Find(&mplas).Error; err != nil {
		return nil, err
	}
	return mplas, nil
}

func (mplas MPlas) Create(db *gorm.DB) error {
	if err := db.Create(&mplas).Error; err != nil {
		return err
	}
	return nil
}

func (mplas MPlas) GetEmptyPlace(plas []*Place) *Place {
	var isEmpty bool
	for _, pla := range plas {
		isEmpty = true
		for _, mpla := range mplas {
			if mpla.PlaceID != pla.ID {
				continue
			}
			if mpla.MPlaQue >= pla.PlacAmount {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			return pla
		}
	}
	return nil
}

const RowDur time.Duration = 15 * time.Minute

func MakeMasterPlaces(start, end time.Time, accID uint, day time.Time, pla *Place) ([]*MasterPlace, error) {
	var mplas []*MasterPlace
	diff := end.Sub(start) / RowDur
	for i := 0; i < int(diff); i++ {
		var from time.Time
		from = start.Add(RowDur * time.Duration(i))
		to := start.Add(RowDur * time.Duration(i+1))
		pla := MasterPlace{
			// MPlaQue:   quo,
			PlaceID:   pla.ID,
			MPlaDay:   day,
			AccountID: accID,
			MPlaFrom:  from,
			MPlaTo:    to,
		}
		mplas = append(mplas, &pla)
	}

	return mplas, nil
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
