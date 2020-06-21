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
	Image   string
	// SerIs   []SerI
}

type SerI struct {
	ID       uint
	emploIDs []uint
	placeIDs []uint
	TimeUse  time.Time
}

type PackSerI struct {
	ID      uint
	PSerIs  []PSerI
	Name    string
	TimeUse time.Duration
	Price   float64
	Image   string
}

type PSerI struct {
	ID             uint
	EmployeeID     uint
	PlaceID        uint
	Plas           []Pla
	ServiceID      uint
	UseTime        time.Duration
	EmploTimeSlots []EmploTimeSlot
	// 	emploIDs []uint
	// 	placeIDs []uint
	// 	TimeUse  time.Time
}

// type Pla struct {
// 	ID         uint
// 	TimeSpents []TimeSpent
// }

func GetPack(sqlDb *sql.DB, id string) (*Pack, []string, error) {
	var p *Pack
	var serIDs []string
	var sers []SerI
	rows, err := sqlDb.Query(`
		SELECT 
			pa.id, pac_name, si.service_id , pa.pac_price, pa.pac_time,
			si.ss_time
		FROM packages AS pa
		INNER JOIN package_service_item AS psi ON psi.package_id = pa.id
		INNER JOIN service_items AS si ON si.id = psi.service_item_id AND si.deleted_at IS NULL
		INNER JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL
		WHERE pa.deleted_at IS NULL AND pa.id = $1 AND pac_is_active = true`, id)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var ser SerI
		rows.Scan(&p.ID, &p.Name, &ser.ID, &p.Price, &p.TimeUse, &ser.TimeUse)
		sers = append(sers, ser)
		serIDs = append(serIDs, fmt.Sprintf("%d", ser.ID))
	}
	return p, serIDs, nil
}

func GetPackEmployees(sqlDb *sql.DB, d time.Time, serIDs string) ([]*serEmp, error) {
	// qe := fmt.Sprintf(`SELECT
	// 		es.employee_id ,
	// 		ts.time_end, ts.time_start, ts.id,
	// 		e.empo_image
	// 	FROM employee_service AS es
	// 	INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL
	// 	INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id AND ts.deleted_at IS NULL
	// 		AND ts.deleted_at IS NULL
	// 		AND time_day = $1
	// 		AND time_start < $2
	// 		AND ts.time_active = true
	// 	WHERE es.service_id IN %s
	// 	ORDER BY es.employee_id`, serIDs)
	// rows, err := sqlDb.Query(qe, int(d.Weekday()), d)
	// if err != nil {
	// 	return nil, err
	// }
	var serEmps []*serEmp
	// var serID string
	// var empID uint
	// var serEmpSt *serEmp
	// for rows.Next() {
	// 	rows.Scan(&serID, &empID)
	// 	idStr := fmt.Sprintf("%d", serEmpSt.id)
	// 	if serID == idStr {
	// 		serEmpSt.emps = append(serEmpSt.emps, empID)
	// 	} else {
	// 		if serID != "" {
	// 			uintID, _ := strconv.ParseUint(serID, 10, 64)
	// 			serEmpSt.id = uint(uintID)
	// 			serEmps = append(serEmps, serEmpSt)
	// 			serEmpSt.emps = append(serEmpSt.emps, empID)
	// 			serEmpSt = &serEmp{}
	// 		}
	// 	}
	// 	serID = idStr
	// }
	// if len(serEmps) == 0 {
	// 	if empID != 0 {
	// 		uintID, _ := strconv.ParseUint(serID, 10, 64)
	// 		serEmpSt.id = uint(uintID)
	// 		serEmpSt.emps = append(serEmpSt.emps, empID)
	// 		serEmps = append(serEmps, serEmpSt)
	// 	} else {
	// 		return nil, errors.New("not found employee")
	// 	}
	// }
	return serEmps, nil
}
