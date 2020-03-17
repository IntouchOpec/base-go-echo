package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
)

type Package struct {
	orm.ModelBase

	PacName      string         `form:"pac_name" json:"pac_name" gorm:"type:varchar(50)"`
	PacDetail    string         `form:"pac_detail" json:"pac_detail" `
	PacOrder     int            `form:"pac_order" json:"pac_order"`
	PacPrice     float64        `form:"pac_price" json:"pac_price"`
	PacType      string         `form:"pac_type" json:"pac_type" gorm:"type:varchar(50)"`
	PacImage     string         `form:"pac_image" json:"pac_image" gorm:"type:varchar(255)"`
	PacIsActive  bool           `form:"pac_is_active" json:"pac_is_active" sql:"default:true" gorm:"default:true"`
	PacTime      time.Duration  `json:"pac_time"`
	ServiceItems []*ServiceItem `json:"service_items" gorm:"many2many:package_service_item"`
	AccountID    uint           `json:"account_id"`
	Account      *Account       `json:"account" gorm:"ForeignKey:AccountID"`
}

type Pack struct {
	Name    string
	TimeUse time.Duration
	Price   float64
	ID      uint
}

type SerIs struct {
	id       uint
	emploIDs []uint
	placeIDs []uint
}

func GetPack(sqlDb *sql.DB, id string) (*Pack, []string, error) {
	var p *Pack
	var serIDs []string
	var sers []SerIs
	rows, err := sqlDb.Query(`
		SELECT 
			pa.id, pac_name, si.service_id , pa.pac_price, pa.pac_time
		FROM packages AS pa
		INNER JOIN package_service_item AS psi ON psi.package_id = pa.id
		INNER JOIN service_items AS si ON si.id = psi.service_item_id AND si.deleted_at IS NULL
		INNER JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL
		WHERE pa.deleted_at IS NULL AND pa.id = $1 AND pac_is_active = true`, id)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var ser SerIs
		rows.Scan(&p.ID, &p.Name, &ser.id, &p.Price, &p.TimeUse)
		sers = append(sers, ser)
		serIDs = append(serIDs, fmt.Sprintf("%d", ser.id))
	}
	return p, serIDs, nil
}
