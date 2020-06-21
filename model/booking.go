package model

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type BookStatus int

const (
	BookingStatusPandding BookStatus = 1
	BookingStatusReject   BookStatus = -1
	BookingStatusApprove  BookStatus = 2
)

type BookingType int

const (
	BookingTypeSlotTime    BookingType = 1
	BookingTypeServiceItem BookingType = 2
	BookingTypePackage     BookingType = 3
)

// Booking struct save date time
type Booking struct {
	orm.ModelBase
	BookingType        BookingType        `json:"booking_type"`
	BooQueue           int                `json:"boo_queue" `
	BooLineID          string             `json:"boo_line_id" gorm:"type:varchar(50)"`
	CustomerID         uint               `json:"customer_id"`
	ChatChannelID      uint               `json:"chat_chaneel_id"`
	TransactionID      uint               `json:"transaction_id"`
	AccountID          uint               `json:"account_id"`
	BooStatus          BookStatus         `json:"boo_status"`
	BookedDay          time.Time          `gorm:"column:booked_day" json:"booked_day"`
	BookedDate         time.Time          `gorm:"column:booked_date" json:"booked_date"`
	BookedStart        time.Time          `gorm:"column:booked_start" json:"booked_start"`
	BookedEnd          time.Time          `gorm:"column:booked_end" json:"booked_end"`
	BookingServiceItem BookingServiceItem `json:"booking_service_item" gorm:"BookingID"`
	BookingPackage     *BookingPackage    `json:"booking_package" gorm:"BookingID"`
	Place              *Place             `json:"place" gorm:"ForeignKey:PlaceID"`
	Transaction        *Transaction       `json:"transaction"  gorm:"ForeignKey:TransactionID"`
	Customer           *Customer          `json:"customer" gorm:"ForeignKey:CustomerID"`
	ChatChannel        *ChatChannel       `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
	Account            *Account           `json:"account" gorm:"ForeignKey:AccountID"`
}

type BookingServiceItem struct {
	BookingID     uint        `json:"booking_id"`
	PlaceID       uint        `json:"place_id"`
	TimeSlotID    uint        `json:"time_slot_id"`
	EmployeeID    uint        `json:"employee_id"`
	ServiceItemID uint        `json:"serice_item_id"`
	Place         *Place      `json:"place" gorm:"ForeignKey:PlaceID"`
	Booking       *Booking    `json:"booking" gorm:"ForeignKey:BookingID"`
	TimeSlot      *TimeSlot   `json:"time_slot" gorm:"ForeignKey:TimeSlotID"`
	Employee      Employee    `json:"employee" gorm:"ForeignKey:EmployeeID"`
	ServiceItem   ServiceItem `json:"service_item" gorm:"ForeignKey:ServiceItemID"`
}

type BookingPackage struct {
	BookingID  uint      `json:"booking_id"`
	Booking    Booking   `json:"booking" gorm:"ForeignKey:BookingID"`
	PackageID  uint      `json:"package_id"`
	TimeSlotID uint      `json:"time_slot_id"`
	TimeSlot   *TimeSlot `json:"time_slot" gorm:"ForeignKey:TimeSlotID"`
	Package    Package   `json:"package" gorm:"ForeignKey:PackageID"`
}

type EmploTimeSlot struct {
	EmployeeID string
	Start      time.Time
	End        time.Time
	MBs        []TimeSpent
}

type TimeSpent struct {
	Start time.Time
	End   time.Time
}
type ByTimeSpent []EmploTimeSlot

func (a ByTimeSpent) Len() int {
	return len(a)
}

// BookingStatus is status of booking.
// type BookingStatus struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// BookingStatusPandding is booking status pandding for confirm
// var BookingStatusPandding = BookingStatus{ID: 1, Name: "pandding"}

// BookingStatusReject is booking status after pandding user pick It.
// var BookingStatusReject = BookingStatus{ID: 2, Name: "reject"}

// BookingStatusApprove is status approve.
// var BookingStatusApprove = BookingStatus{ID: 3, Name: "approve"}

// BookingState is state of booking.
type BookingState struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MBStatus int

const (
	PanddingMBStatus MBStatus = 0
	ApproveMBStatus  MBStatus = 1
)

type MasterBooking struct {
	orm.ModelBase

	MBQue      int       `json:"m_b_que"`
	EmployeeID uint      `json:"employee_id"`
	BookingID  uint      `json:"booking_id"`
	PlaceID    uint      `json:"place_id"`
	AccountID  uint      `json:"account_id"`
	MBDay      time.Time `json:"mb_day"`
	MBFrom     time.Time `json:"mb_from"`
	MBTo       time.Time `json:"mb_to"`
	MBStatus   MBStatus  `json:"mb_status" sql:"default:0"`

	Booking  *Booking  `json:"booking" gorm:ForeignKey:BookingID""`
	Employee *Employee `json:"employee" gorm:"ForeignKey:EmployeeID"`
	Place    *Place    `json:"place" gorm:"ForeignKey:PlaceID"`
	Account  *Account  `json:"account" gorm:"ForeignKey:AccountID"`
}

type Pla struct {
	ID     uint
	Amount int
	MBs    []TimeSpent
}

type employee struct {
	id uint
}

// func RemoveIndex(s , index int) []interface{} {
// return append(s[:index], s[index+1:]...)
// }

func (tran *Transaction) MakeMasterBooking(sql *sql.DB) ([]interface{}, string, error) {
	var boos []Booking
	var query string
	var values []interface{}

	rows, err := sql.Query(`
		SELECT 
			bo.id AS booking_id, booked_start, booked_end, booked_day,s.id AS service_id,
			CASE
			WHEN bp.package_id IS NULL 
				THEN 0
				ELSE bp.package_id
			END,
			CASE
			WHEN bsi.employee_id IS NULL 
				THEN 0
				ELSE bsi.employee_id
			END 
		FROM bookings AS bo
		LEFT JOIN booking_service_items AS bsi ON bo.id = bsi.booking_id 
		LEFT JOIN service_items AS si ON  si.id = bsi.service_item_id AND si.deleted_at IS NULL 
		LEFT JOIN services AS s ON s.id = si.service_id AND s.deleted_at IS NULL 
		LEFT JOIN booking_packages AS bp ON bo.id = bp.booking_id 
		LEFT JOIN packages AS pa ON pa.id = bp.package_id AND pa.deleted_at IS NULL 
		WHERE bo.deleted_at IS NULL AND bo.transaction_id = $1`, tran.ID)
	if err != nil {
		return nil, "notFound", err
	}
	var serIID uint
	var booPacID string
	var emploID uint
	for rows.Next() {
		var boo Booking
		rows.Scan(&boo.ID, &boo.BookedStart, &boo.BookedEnd, &boo.BookedDay, &serIID, &booPacID, &emploID)
		boos = append(boos, boo)
	}

	for _, boo := range boos {
		if serIID != 0 {
			values, query, err = SetMasterBookingSer(sql, emploID, serIID, boo)
			if err != nil {
				return nil, query, err
			}
		} else {
			// _, serIDs, err := GetPack(sql, booPacID)
			// if err != nil {
			// 	fmt.Println("====1")
			// 	return nil, "", err
			// }
			// values, query, err = setMasterPack(sql, serIDs, boo)
			// if err != nil {
			// 	return nil, "", err
			// }
		}

	}
	return values, query, nil
}

func setMasterPack(sql *sql.DB, serIDs []string, boo Booking) ([]interface{}, string, error) {
	values := []interface{}{}
	var query string

	serIDsStr := strings.Join(serIDs, ",")
	qePlas := fmt.Sprintf(`
		SELECT 
			pl.id, pl.plac_amount, service_id
		FROM place_service AS ps
		INNER JOIN places AS pl ON pl.id = ps.place_id AND plac_active = true
		WHERE service_id IN (%s) AND pl.deleted_at IS NULL
		GROUP BY service_id, pl.id
		ORDER BY service_id, pl.id;
	`, serIDsStr)

	PlaRows, err := sql.Query(qePlas)
	if err != nil {
		return nil, "", err
	}
	for PlaRows.Next() {
		PlaRows.Scan()
	}

	qeEmpo := fmt.Sprintf(`
		SELECT 
			e.id 
		FROM employees AS e
		INNER JOIN time_slots AS ts ON ts.employee_id = e.id AND ts.deleted_at IS NULL AND time_day = $1 AND time_active = true
		INNER JOIN employee_service AS es ON e.id = es.employee_id AND service_id IN (%s)
		WHERE e.deleted_at IS NULL AND e.empo_is_active = true
		ORDER BY e.id;
	`, serIDsStr)
	EmpoRows, err := sql.Query(qeEmpo)
	if err != nil {
		return nil, "", err
	}
	for EmpoRows.Next() {
		EmpoRows.Scan()
	}
	return values, query, nil
}

func (boo Booking) MasterBookingSer(plaMDs []MasterBooking) ([]interface{}, string, error) {
	values := []interface{}{}
	var query string
	var plaMD MasterBooking
	diff := boo.BookedEnd.Sub(boo.BookedStart) / RowDur
	for i := 0; i < int(diff); i++ {
		var from time.Time
		from = boo.BookedStart.Add(RowDur * time.Duration(i))
		to := boo.BookedStart.Add(RowDur * time.Duration(i+1))
		plaMD.EmployeeID = boo.BookingServiceItem.EmployeeID
		plaMD.MBDay = boo.BookedDay
		plaMD.MBFrom = from
		plaMD.MBTo = to
		plaMD.PlaceID = boo.BookingServiceItem.PlaceID
		if len(plaMDs) == 0 {
			plaMD.MBQue = 1
		} else {
			for _, pla := range plaMDs {
				if plaMD.MBFrom == pla.MBFrom && pla.MBTo == plaMD.MBTo {
					plaMD.MBQue = pla.MBQue + 1
				}
			}
		}

		// fmt.Println(plaMD.EmployeeID, emploID)
		values = append(values,
			plaMD.MBQue, plaMD.EmployeeID, boo.ID, plaMD.PlaceID, plaMD.MBDay, plaMD.MBFrom, plaMD.MBTo, boo.AccountID)
		numFields := 8
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query + `now(), now()),`
	}
	return values, query[:len(query)-1], nil
}

func SetMasterBookingSer(sql *sql.DB, emploID, serIID uint, boo Booking) ([]interface{}, string, error) {
	var empIDs []string
	var plas []Pla
	var plaIDs []string
	values := []interface{}{}
	var query string

	rows, err := sql.Query(`
		SELECT 
			pl.id, pl.plac_amount
		FROM place_service AS ps
		INNER JOIN places AS pl ON pl.id = ps.place_id AND plac_active = true
		WHERE service_id = $1 AND pl.deleted_at IS NULL
		ORDER BY pl.id;`, serIID)
	if err != nil {
		return nil, "notPlace", err
	}
	for rows.Next() {
		var pla Pla
		rows.Scan(&pla.ID, &pla.Amount)
		plaIDs = append(plaIDs, fmt.Sprintf("%d", pla.ID))
		plas = append(plas, pla)
	}
	if emploID == 0 {
		rows, err := sql.Query(`
				SELECT 
					e.id 
				FROM employees AS e
				INNER JOIN time_slots AS ts ON ts.employee_id = e.id AND ts.deleted_at IS NULL AND time_day = $1 AND time_active = true
				INNER JOIN employee_service AS es ON e.id = es.employee_id AND service_id = $2
				WHERE e.deleted_at IS NULL AND e.empo_is_active = true
				ORDER BY e.id;`, int(boo.BookedDay.Add(7*time.Hour).Weekday()), serIID)
		if err != nil {
			return nil, "notEmployee", err
		}
		for rows.Next() {
			var empID uint
			rows.Scan(&empID)
			empIDs = append(empIDs, fmt.Sprintf("%d", empID))
		}
	} else {
		empIDs = append(empIDs, fmt.Sprintf("%d", emploID))
	}
	rows, err = sql.Query(`
			SELECT
				employee_id
			FROM master_bookings AS mb
			WHERE
				deleted_at IS NULL AND employee_id IN ($1) AND account_id = $2 AND mb_day = $3 AND mb_from BETWEEN $4 AND $5 OR mb_to BETWEEN $6 AND $7
			GROUP BY employee_id
			ORDER BY employee_id`, strings.Join(empIDs, ","), boo.AccountID, boo.BookedDay,
		boo.BookedStart, boo.BookedEnd,
		boo.BookedStart, boo.BookedEnd)

	fmt.Println(rows, err, "====4")

	var emplID uint
	var isReady bool = true

	for rows.Next() {
		rows.Scan(&emplID)
		for _, empID := range empIDs {
			iduint, _ := strconv.ParseUint(empID, 10, 64)
			if emplID == uint(iduint) {
				isReady = false
				break
			}
		}
		if isReady {
			break
		}
	}
	if emplID == 0 {
		iduint, _ := strconv.ParseUint(empIDs[0], 10, 64)
		emplID = uint(iduint)
	}

	if !isReady {
		return nil, "notEmployeeReady", errors.New("employee not ready.")
	}

	rows, err = sql.Query(`
		SELECT
			place_id, MAX(mb_que), mb_from, mb_to
		FROM master_bookings AS mb
		WHERE
			deleted_at IS NULL AND place_id IN ($1) AND account_id = $2 AND mb_day = $3 AND mb_from BETWEEN $4 AND $5 OR mb_to BETWEEN $6 AND $7
		GROUP BY place_id, mb_from, mb_to
		ORDER BY place_id, mb_from, mb_to`, strings.Join(plaIDs, ","), boo.AccountID, boo.BookedDay,
		boo.BookedStart, boo.BookedEnd,
		boo.BookedStart, boo.BookedEnd)
	var plaMD MasterBooking
	var plaMDs []MasterBooking
	var isPlaReady bool = true
	for _, plaI := range plas {
		for rows.Next() {
			// var pla MasterBooking
			rows.Scan(&plaMD.PlaceID, &plaMD.MBQue, &plaMD.MBFrom, &plaMD.MBTo)
			if plaI.ID == plaMD.PlaceID {
				fmt.Println(plaI.Amount, plaMD.MBQue, "plaI.Amount >= plaMD.MBQue")
				if plaI.Amount <= plaMD.MBQue {
					plaMDs = make([]MasterBooking, 0)
					isPlaReady = false
				} else {
					plaMDs = append(plaMDs, plaMD)
				}
				// if len(plas) == 1 {

				// } else {
				// 	plas = plas[1:]
				// }
			}
		}

		if isPlaReady {
			break
		}
	}
	fmt.Println("isPlaReady", isPlaReady)
	if !isPlaReady {
		return nil, "notPlaceReady", errors.New("place not ready.")
	}
	if plaMD.ID == 0 {
		plaMD.PlaceID = plas[0].ID
		plaMD.MBQue = 1
	}

	diff := boo.BookedEnd.Sub(boo.BookedStart) / RowDur
	for i := 0; i < int(diff); i++ {
		var from time.Time
		from = boo.BookedStart.Add(RowDur * time.Duration(i))
		to := boo.BookedStart.Add(RowDur * time.Duration(i+1))
		plaMD.EmployeeID = emplID
		plaMD.MBDay = boo.BookedDay
		plaMD.MBFrom = from
		plaMD.MBTo = to
		for _, pla := range plaMDs {
			if plaMD.MBFrom == pla.MBFrom && pla.MBTo == plaMD.MBTo {
				plaMD.MBQue = pla.MBQue + 1
			}
		}
		values = append(values,
			plaMD.MBQue, plaMD.EmployeeID, boo.ID, plaMD.PlaceID, plaMD.MBDay, plaMD.MBFrom, plaMD.MBTo, boo.AccountID)
		numFields := 8
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	return values, query[:len(query)-1], nil
}

func (b *Booking) CreateSql(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT INTO bookings (
		booking_type ,boo_queue ,boo_line_id ,customer_id ,chat_channel_id ,transaction_id ,account_id ,boo_status ,booked_day ,booked_date ,booked_start ,booked_end ,created_at ,updated_at 
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, now(), now()) RETURNING id`)
	if err != nil {
		return err
	}
	result := stmt.QueryRow(b.BookingType, b.BooQueue, b.BooLineID, b.CustomerID, b.ChatChannelID, b.TransactionID, b.AccountID, b.BooStatus, b.BookedDay, b.BookedDate, b.BookedStart, b.BookedEnd)
	if err != nil {
		return err
	}
	if err := result.Scan(&b.ID); err != nil {
		return err
	}
	return nil
}

func CreateMasterBooking(vStr string, sql *sql.Tx, values []interface{}) error {
	stmt, err := sql.Prepare(fmt.Sprintf("INSERT INTO master_bookings (mb_que,employee_id,booking_id,place_id,mb_day,mb_from,mb_to,account_id, created_at, updated_at) VALUES %s ", vStr))
	if err != nil {
		return err
	}
	// fmt.Println(values)
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}
	// count, _ := result.RowsAffected()
	// fmt.Println(count)
	return nil
}

// SaveBooking is function create chat answer.
// func (booking *Booking) SaveBooking() (*Booking, error) {
// 	db := DB()
// 	booked := Booking{}
// 	db.Preload("ServiceSlot").Where("Booked_Date = ? and Sub_service_ID = ?", booking.BookedDate, booking.TimeSlotID).Last(&booked)
// 	if booked.TimeSlot.TimeAmount == 0 {
// 		booking.BooQueue = 1
// 	} else if booked.TimeSlot.TimeAmount > booked.BooQueue {
// 		booking.BooQueue = booked.BooQueue + 1
// 	} else {
// 		return nil, errors.New("can't insert booking case queue full")
// 	}

// 	if err := db.Create(&booking).Error; err != nil {
// 		return nil, err
// 	}
// 	return booking, nil
// }

func (booking *Booking) UpdateBooking(id string) *Booking {
	db := DB()
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}

	if err := db.Save(&booking).Error; err != nil {
		return nil
	}
	return booking
}

func GetBookingList(chatChannelID string) *[]Booking {
	bookings := []Booking{}
	db := DB()
	if err := db.Where("chat_channel_id = ?", chatChannelID).Find(&bookings).Error; err != nil {
		return nil
	}
	return &bookings
}

func GetBooking(id string) *Booking {
	db := DB()
	booking := Booking{}
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}
	return &booking
}

func (booking *Booking) DeleteBooking(id string) *Booking {
	db := DB()
	if err := db.Find(&booking, id).Error; err != nil {
		return nil
	}
	if err := db.Delete(&booking, id).Error; err != nil {
		return nil
	}
	return booking
}

func (book *Booking) BookingAcjectStatus(status string) (*Booking, error) {

	return book, nil
}

func (booking *Booking) CrateBookingOnTran(tx *gorm.DB) error {
	if err := tx.Create(&booking).Error; err != nil {
		return err
	}
	return nil
}

func (b *Booking) BookServiceItem(db *gorm.DB) error {
	// start, end, err := MakeTimeStartAndTimeEnd(d, serI.SSTime)

	return nil
}

// func (booking *Booking) BookPackage(db *gorm.DB) error {
// 	return nil
// }

func (msEmplos *MEmplos) Get(db *gorm.DB, where string, values ...interface{}) error {
	// "account_id = ? and m_emplo_day = ? and employee_id in (?) and m_empo_from >= ?"
	if err := db.Order("employee_id, m_empo_from").Where(where, values).Find(&msEmplos).Error; err != nil {
		return err
	}
	return nil
}

type serEmp struct {
	mts []EmploTimeSlot
	id  uint
}

type sersPla struct {
	plas []Pla
	id   uint
}

// type pack struct {
// 	id       uint
// 	serEmps  []serEmp
// 	sersPlas []sersPla
// }

func (b *Booking) PackAppoint(db *sql.DB, pack Pack, serIDs []string, d time.Time) error {
	var strSerIDs string
	for i := 3; i < len(serIDs)+3; i++ {
		strSerIDs += fmt.Sprintf("%s,", serIDs[i-3])
	}
	strSerIDs = fmt.Sprintf("(%s)", strSerIDs[:len(strSerIDs)-1])
	qe := fmt.Sprintf(`
	SELECT 
		es.service_id, es.employee_id 
	FROM employee_service AS es
	INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id 
		AND ts.deleted_at IS NULL 
		AND time_day = $1
		AND time_start < $2 
	INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL AND e.empo_is_active = true
	WHERE es.service_id IN %s
	ORDER BY es.service_id, es.employee_id`, strSerIDs)
	rows, err := db.Query(qe, int(d.Weekday()), d)
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan()
	}
	// var serEmps []serEmp
	// var serID string
	// var empID uint
	// var serEmpSt serEmp
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

	// 			serEmpSt = serEmp{}
	// 		}
	// 	}
	// 	serID = idStr
	// }
	// fmt.Println(empID)
	// if len(serEmps) == 0 {
	// 	if empID != 0 {
	// 		uintID, _ := strconv.ParseUint(serID, 10, 64)
	// 		serEmpSt.id = uint(uintID)
	// 		serEmpSt.emps = append(serEmpSt.emps, empID)
	// 		serEmps = append(serEmps, serEmpSt)
	// 	} else {
	// 		return errors.New("not found employee")
	// 	}
	// }
	rows, err = db.Query(fmt.Sprintf(`
			SELECT 
				ps.service_id, pl.id, pl.plac_amount
			FROM places AS pl
			INNER JOIN place_service AS ps ON ps.place_id = pl.id AND service_id in %s
			WHERE pl.deleted_at IS NULL 
			GROUP BY ps.service_id, pl.id, pl.plac_amount
			ORDER BY pl.id`, strSerIDs))
	if err != nil {
		return err
	}
	// var pla sersPla
	// var p Pla
	// var sersPlas []sersPla
	// serID = ""
	// for rows.Next() {
	// 	var amount int
	// 	var plaID uint
	// 	var id uint
	// 	rows.Scan(&id, &plaID, &amount)
	// 	idStr := fmt.Sprintf("%d", plaID)
	// 	if serID == idStr {
	// 		p.ID = id
	// 		p.Amount = amount
	// 		pla.plas = append(pla.plas, p)
	// 	} else {
	// 		if serID != "" {
	// 			uintID, _ := strconv.ParseUint(serID, 10, 64)
	// 			pla.id = uint(uintID)
	// 			sersPlas = append(sersPlas, pla)
	// 			p = Pla{}
	// 		}
	// 	}
	// }
	// if len(sersPlas) == 0 {
	// 	uintID, _ := strconv.ParseUint(serID, 10, 64)
	// 	pla.id = uint(uintID)
	// 	sersPlas = append(sersPlas, pla)
	// }
	return nil
}

func (b *Booking) MakeMBPacks(db *sql.DB, pack *Pack, serIDs []string) error {
	var strSerIDs string

	for i := 0; i < len(serIDs); i++ {
		strSerIDs += serIDs[i] + ","
	}
	strSerIDs = fmt.Sprintf("(%s)", strSerIDs[:len(strSerIDs)-1])
	qe := fmt.Sprintf(`
		SELECT 
			es.employee_id ,
			ts.time_end, ts.time_start, ts.id,
			e.empo_image
		FROM employee_service AS es
		INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL
		INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id AND ts.deleted_at IS NULL
			AND ts.deleted_at IS NULL 
			AND time_day = $1 
			AND time_start < $2
			AND time_end < $3
			AND ts.time_active = true
		WHERE es.service_id IN %s
		ORDER BY es.employee_id`, serIDs)
	rows, err := db.Query(qe, int(b.BookedDay.Weekday()), b.BookedStart, b.BookedEnd)
	if err != nil {
		return err
	}
	for rows.Next() {
		rows.Scan()
	}
	// var serEmps []*serEmp
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
	// 		return errors.New("not found employee")
	// 	}
	// }

	return nil
}

func (b *Booking) PackNow(db *sql.DB, pack PackSerI) error {
	d := time.Now()
	start, end, err := MakeTimeStartAndTimeEnd(d, pack.TimeUse)
	if err != nil {
		return err
	}
	day, _ := time.Parse("2006-01-02", d.Format("2006-01-02"))
	var strSerIDs string
	for index, _ := range pack.PSerIs {
		strSerIDs += fmt.Sprintf("%d,", pack.PSerIs[index].ServiceID)
	}
	strSerIDs = fmt.Sprintf("(%s)", strSerIDs[:len(strSerIDs)-1])

	qe := fmt.Sprintf(`
	SELECT 
		es.service_id, es.employee_id, time_start, time_end
	FROM employee_service AS es
	INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id 
		AND ts.deleted_at IS NULL 
		AND time_day = $1
		AND time_start < $2 
	INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL AND e.empo_is_active = true
	WHERE es.service_id IN %s
	ORDER BY es.service_id, es.employee_id, time_start, time_end`, strSerIDs)

	rows, err := db.Query(qe, int(d.Weekday()), start.Add(-7*time.Hour))
	if err != nil {
		return err
	}
	var serID uint
	var emploIDs string
	for rows.Next() {
		var emt EmploTimeSlot
		rows.Scan(&serID, &emt.EmployeeID, &emt.Start, &emt.End)
		emploIDs += fmt.Sprintf("%s,", emt.EmployeeID)
		for index, si := range pack.PSerIs {
			if si.ID == serID {
				pack.PSerIs[index].EmploTimeSlots = append(pack.PSerIs[index].EmploTimeSlots, emt)
			}
		}
	}

	if serID == 0 || err != nil {
		return errors.New("not found employee")
	}

	mss, err := b.GetMasterBooking(db, day, start.Add(-7*time.Hour), emploIDs[:len(emploIDs)-1])
	if err != nil {
		return err
	}
	var emploID uint
	var MSStart time.Time
	var MSEnd time.Time
	var emploMs []EmploTimeSlot
	var emploM EmploTimeSlot
	for index, ms := range mss {
		if ms.EmployeeID == emploID {
			if !ms.MBTo.Equal(MSEnd) {
				emploM.MBs = append(emploM.MBs, TimeSpent{Start: MSStart, End: MSEnd})
				MSStart = ms.MBFrom
			}
		} else if index != 0 {
			emploM.EmployeeID = fmt.Sprintf("%d,", ms.EmployeeID)
			emploM.MBs = append(emploM.MBs, TimeSpent{Start: MSStart, End: ms.MBTo})
			emploMs = append(emploMs, emploM)
		}
		emploID = ms.EmployeeID
		MSEnd = ms.MBTo
		if len(mss)-1 == index {
			emploM.EmployeeID = fmt.Sprintf("%d", emploID)
			emploM.MBs = append(emploM.MBs, TimeSpent{Start: MSStart, End: MSEnd})
			emploMs = append(emploMs, emploM)
		}
	}
	serStart := start
	for index, serI := range pack.PSerIs {
		serEnd := start.Add(serI.UseTime)
		for _, emp := range serI.EmploTimeSlots {
			for _, ems := range emploMs {
				if emp.EmployeeID == ems.EmployeeID {
					if inTimeSpan(emp.Start, emp.End, serStart) && inTimeSpan(emp.Start, emp.End, serEnd) {

						for _, mb := range ems.MBs {
							if !inTimeSpan(mb.Start, mb.End, serStart) && !inTimeSpan(mb.Start, mb.End, serEnd) {
								u64, _ := strconv.ParseUint(ems.EmployeeID, 10, 32)
								pack.PSerIs[index].EmployeeID = uint(u64)
							}
							break
						}

					}
					break
				}
			}
		}
	}

	rows, err = db.Query(fmt.Sprintf(`
			SELECT 
				ps.service_id, pl.id, pl.plac_amount
			FROM places AS pl
			INNER JOIN place_service AS ps ON ps.place_id = pl.id AND service_id in %s
			WHERE pl.deleted_at IS NULL 
			GROUP BY ps.service_id, pl.id, pl.plac_amount
			ORDER BY ps.service_id, pl.id`, strSerIDs))
	if err != nil {
		return err
	}

	var serviceID uint
	var plaIDs string
	for rows.Next() {
		var pla Pla
		rows.Scan(&serviceID, &pla.ID, &pla.Amount)
		plaIDs += fmt.Sprintf("%d,", pla.ID)
		for index, si := range pack.PSerIs {
			if si.ID == serviceID {
				pack.PSerIs[index].Plas = append(pack.PSerIs[index].Plas, pla)
			}
		}
	}
	qe = fmt.Sprintf(`
		SELECT 
			place_id, MAX(mb_que), mb_from, mb_to
		FROM master_bookings AS mb 
		WHERE 
			deleted_at IS NULL AND place_id IN (%s) 
			AND account_id = $1
			AND mb_day = $2
			AND mb_from > $3
		GROUP BY place_id, mb_from, mb_to
		ORDER BY place_id, mb_from, mb_to`, plaIDs[:len(plaIDs)-1])
	rows, err = db.Query(qe, b.AccountID, d, start)

	if err != nil {
		return err
	}
	var plaMB []MasterBooking
	var plaID uint
	var plaStart time.Time
	var plaEnd time.Time
	var i int = 0
	for rows.Next() {
		var mb MasterBooking
		rows.Scan(&mb.PlaceID, &mb.MBQue, &mb.MBFrom, &mb.MBTo)
		if plaID == mb.PlaceID {
			if plaEnd.Equal(mb.MBFrom) {
			} else {
				plaStart = mb.MBTo
				plaMB = append(plaMB, MasterBooking{PlaceID: plaID, MBFrom: plaStart, MBTo: plaEnd})
			}
		} else if i != 0 {
			plaStart = mb.MBTo
			plaMB = append(plaMB, MasterBooking{PlaceID: plaID, MBFrom: plaStart, MBTo: plaEnd})
		} else {
			plaStart = mb.MBTo
		}
		plaEnd = mb.MBTo
		plaID = mb.PlaceID
		i++
	}
	if i >= 1 {
		plaMB = append(plaMB, MasterBooking{PlaceID: plaID, MBFrom: plaStart, MBTo: plaEnd})
	}

	for x, serI := range pack.PSerIs {
		for y, pla := range serI.Plas {
			for _, plaM := range plaMB {
				if pla.ID == plaM.PlaceID {
					pack.PSerIs[x].Plas[y].MBs = append(pack.PSerIs[x].Plas[y].MBs, TimeSpent{Start: plaM.MBFrom, End: plaM.MBTo})
				}
			}
		}
	}
	for _, serI := range pack.PSerIs {
		for _, pla := range serI.Plas {
			fmt.Println(pla)
			// for _, ems := range emploMs {
			// 	if emp.EmployeeID == ems.EmployeeID {
			// 		if inTimeSpan(emp.Start, emp.End, serStart) && inTimeSpan(emp.Start, emp.End, serEnd) {

			// 			for _, mb := range ems.MBs {
			// 				if !inTimeSpan(mb.Start, mb.End, serStart) && !inTimeSpan(mb.Start, mb.End, serEnd) {
			// 					u64, _ := strconv.ParseUint(ems.EmployeeID, 10, 32)
			// 					pack.PSerIs[index].EmployeeID = uint(u64)
			// 				}
			// 				break
			// 			}

			// 		}
			// 		break
			// 	}
			// }
		}
	}
	d, _ = time.Parse("2006-01-02", d.Format("2006-01-02"))
	b.BookedDay = d
	b.BookedStart = start
	b.BookedEnd = end

	return nil
}

func (a ByTimeSpent) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimeSpent) Less(i, j int) bool {
	switch {
	case len(a[j].MBs) == 0:
		return true
	case a[i].MBs[0].End.Before(a[j].MBs[0].End):
		return true
	default:
		return false
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func (ets EmploTimeSlot) EmploYeeReady(start, end time.Time) bool {
	if !inTimeSpan(ets.Start, ets.End, start) && !inTimeSpan(ets.Start, ets.End, end) {
		return false
	}

	for _, ms := range ets.MBs {
		if inTimeSpan(ms.Start, ms.End, start) && inTimeSpan(ms.Start, ms.End, end) {
			return false
		}
	}
	return true
}

func (b *Booking) ServiceItemNow(db *sql.DB, serI ServiceItem) ([]MasterBooking, error) {
	var emplos []EmploTimeSlot
	var emploIDs string
	var plaIDs []string
	var plas []Pla
	d := time.Now()
	day, _ := time.Parse("2006-01-02", d.Format("2006-01-02"))

	start, end, err := MakeTimeStartAndTimeEnd(d, serI.SSTime)
	if err != nil {
		// fmt.(err)
		fmt.Println("err")
		return nil, err
	}
	start, end = start.Add(-(7 * time.Hour)), end.Add(-(7 * time.Hour))
	rows, err := db.Query(`
			SELECT 
				es.employee_id, time_start, time_end
			FROM employee_service AS es
			INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id 
				AND ts.deleted_at IS NULL 
				AND time_day = $1
				AND time_start < $2 
			WHERE es.service_id = $3 AND deleted_at IS NULL
			ORDER BY es.employee_id`, int(d.Weekday()), start, serI.Service.ID)
	// fmt.Println(int(d.Weekday()), start, serI.Service.ID)
	// AND time_end > $3
	fmt.Println(err)
	if err == nil {
		for rows.Next() {
			var emplo EmploTimeSlot
			rows.Scan(&emplo.EmployeeID, &emplo.Start, &emplo.End)
			// fmt.Println(emplos, "=id")
			emploIDs += emplo.EmployeeID + ","
			emplos = append(emplos, emplo)
		}
	}
	if len(emplos) == 0 {
		return nil, errors.New("not employee")
	}
	rows, err = db.Query(`
			SELECT 
				pl.id, pl.plac_amount
			FROM places AS pl
			INNER JOIN place_service AS ps ON ps.place_id = pl.id AND service_id = $1
			WHERE pl.deleted_at IS NULL 
			GROUP BY pl.id, pl.plac_amount
			ORDER BY pl.id`, serI.Service.ID)
	if err != nil {
		fmt.Println(err, "err1")
		return nil, err
	}
	for rows.Next() {
		var pla Pla
		rows.Scan(&pla.ID, pla.Amount)
		plaIDs = append(plaIDs, fmt.Sprintf("%d", pla.ID))
		plas = append(plas, pla)
	}
	if len(plaIDs) == 0 {
		fmt.Println(err, "err4")
		return nil, errors.New("")
	}
	emploIDs = emploIDs[:len(emploIDs)-1]
	Mbs, err := b.GetMasterBooking(db, day, start, emploIDs)
	if err != nil {
		return nil, errors.New("")
	}
	var isReady bool = true
	var emplID uint
	for rows.Next() {
		var Mb MasterBooking
		rows.Scan(&Mb.EmployeeID, &Mb.MBFrom, &Mb.MBTo)
		Mbs = append(Mbs, Mb)
	}
	if len(Mbs) == 0 {
		emploID, _ := strconv.ParseUint(emplos[0].EmployeeID, 10, 64)
		emplID = uint(emploID)
	} else {

		for index, emp := range emplos {
			iduint, _ := strconv.ParseUint(emp.EmployeeID, 10, 64)
			var fristTime time.Time
			var lastTime time.Time
			isReady = true
			for _, ms := range Mbs {
				if ms.EmployeeID == uint(iduint) {
					isReady = false
					if fristTime.Hour() == 0 {
						fristTime = ms.MBFrom
					}
					lastTime = ms.MBTo
					if !lastTime.Equal(ms.MBFrom) {
						if ms.MBFrom.Sub(lastTime) > serI.SSTime {
							emplos[index].MBs = append(emp.MBs, TimeSpent{Start: fristTime, End: lastTime})
							fristTime = time.Time{}
							break
						} else {
							emplos[index].MBs = []TimeSpent{TimeSpent{Start: fristTime, End: lastTime}}
						}
					}
					if len(Mbs) == 0 {
						Mbs = make([]MasterBooking, 0)
					} else {
						Mbs = Mbs[1:]
					}
				} else {
					// fmt.Println("====")
					break
				}
			}
			if isReady {
				emploID, _ := strconv.ParseUint(emp.EmployeeID, 10, 64)
				emplID = uint(emploID)
				break
			}
		}
	}
	if !isReady {
		var findSpendTime bool
		var count int
		sort.Sort(ByTimeSpent(emplos))
		for !findSpendTime {
			var em EmploTimeSlot
			notFoundTimeSlot := true
			for index, emplo := range emplos {
				if len(emplos[index].MBs) == 0 {
					emplos[index].MBs = make([]TimeSpent, 0)
					emploID, _ := strconv.ParseUint(emplo.EmployeeID, 10, 64)
					emplID = uint(emploID)
					findSpendTime = true
					notFoundTimeSlot = false
					em = emplo
					break
				}
				if !inTimeSpan(emplo.Start, emplo.End, start) && !inTimeSpan(emplo.Start, emplo.End, end) {
					continue
				} else {
					notFoundTimeSlot = false
				}
				if len(emplo.MBs)-1 < count {
					emploID, _ := strconv.ParseUint(emplo.EmployeeID, 10, 64)
					emplID = uint(emploID)
					findSpendTime = true
					notFoundTimeSlot = false
					em = emplo
					break
				}
				if inTimeSpan(emplo.MBs[count].Start, emplo.MBs[count].End, start) && inTimeSpan(emplo.MBs[count].Start, emplo.MBs[count].End, end) {
					continue
				}

				if emplos[index].MBs[count].End.After(start) || start.Equal(emplos[index].MBs[count].End) && count > 0 {
					start = emplos[index].MBs[count].End
					end = start.Add(serI.SSTime)
				}
			}
			if em.EmployeeID != "" {
				for _, ms := range em.MBs {
					if inTimeSpan(ms.Start, ms.End, start) && inTimeSpan(ms.Start, ms.End, end) {
						findSpendTime = false
						break
					}
				}
			}
			if notFoundTimeSlot {
				return nil, errors.New("employee not ready")
			}
			count++
		}

	}

	rows, err = db.Query(`
			SELECT 
				place_id, MAX(mb_que), mb_from, mb_to
			FROM master_bookings AS mb 
			WHERE 
				deleted_at IS NULL AND place_id IN ($1) 
				AND account_id = $2 
				AND mb_day = $3 
				AND mb_from BETWEEN $4 AND $5
				OR mb_to BETWEEN $6 AND $7
			GROUP BY place_id, mb_from, mb_to
			ORDER BY place_id, mb_from, mb_to`, strings.Join(plaIDs, ","), b.AccountID, d,
		start, end, start, end)

	var plaMD MasterBooking
	var plaMDs []MasterBooking
	var isPlaReady bool = true
	for _, plaI := range plas {
		for rows.Next() {
			rows.Scan(&plaMD.PlaceID, &plaMD.MBQue, &plaMD.MBFrom, &plaMD.MBTo)
			if plaI.ID == plaMD.PlaceID {
				if plaI.Amount >= plaMD.MBQue {
					plaMDs = make([]MasterBooking, 0)
					isPlaReady = false
				} else {
					plaMDs = append(plaMDs, plaMD)
				}
				if len(plas) == 1 {
					plas = make([]Pla, 0)
					break
				}
				plas = plas[1:]
			} else {
				break
			}
		}

		if isPlaReady {
			break
		}
	}
	if !isPlaReady {
		return nil, errors.New("place not ready.")
	}
	d, _ = time.Parse("2006-01-02", d.Format("2006-01-02"))
	b.BookedDay = d
	b.BookedStart = start
	b.BookedEnd = end
	b.BookingServiceItem.EmployeeID = emplID
	return plaMDs, nil
}

func (b Booking) GetMasterBooking(db *sql.DB, day, start time.Time, emploIDs string) ([]MasterBooking, error) {
	qe := fmt.Sprintf(`
		SELECT 
			employee_id, mb_from, mb_to
		FROM master_bookings AS mb 
		WHERE 
			deleted_at IS NULL AND employee_id IN (%s) AND account_id = $1 AND mb_day = $2 AND mb_from > $3
		GROUP BY employee_id, mb_from, mb_to
		ORDER BY employee_id, mb_from, mb_to`, emploIDs)
	rows, err := db.Query(qe, b.AccountID, day, start)
	if err != nil {
		return nil, err
	}
	fmt.Println(" b.AccountID, day, start", b.AccountID, day, start, emploIDs)

	var Mbs []MasterBooking
	for rows.Next() {
		var Mb MasterBooking
		rows.Scan(&Mb.EmployeeID, &Mb.MBFrom, &Mb.MBTo)
		Mbs = append(Mbs, Mb)
	}
	return Mbs, nil
}

// var timeSs []TimeSlot
// var msEmplos MEmplos
// var msPlas MPlas
// if err := db.Order("employee_id").Where("account_id = ? and time_day = ? and time_start >= ?",
// 	serI.AccountID, int(d.Weekday()), start).Find(&timeSs).Error; err != nil {
// 	fmt.Println(err, "===")
// 	return err
// }
// empIDs := SetTimeSlotMixEmployee(timeSs, serI.Service.Employees)
// if len(empIDs) == 0 {
// 	fmt.Println(err)
// 	return errors.New("")
// }
// plaIDs := SetPlaces(serI.Service.Places)
// if len(plaIDs) == 0 {
// 	return errors.New("")
// }

// msEmplos.Get(db,
// 	"account_id = ? and m_emplo_day = ? and employee_id in (?) and m_empo_from >= ?",
// 	serI.AccountID, d, empIDs, start)

// msPlas.Get(db,
// 	"account_id = ? and m_pla_day = ? and place_id in (?) and m_emp_from >= ?",
// 	serI.AccountID, d, plaIDs, start)

// isMsPlas := len(msPlas) == 0
// isMsEmplos := len(msEmplos) == 0
// if isMsPlas && isMsEmplos {
// 	b.BookedStart = start
// 	b.BookedEnd = end
// 	b.BookedDay = d
// 	return nil
// }
// if !isMsPlas {
// 	for _, pla := range serI.Service.Places {
// 		for _, msPla := range msPlas {
// 			if pla.ID != msPla.PlaceID {
// 				continue
// 			}
// 			if pla.PlacAmount >= msPla.MPlaQue {
// 				if end.Hour() > msPla.MPlaTo.Hour() && start.Hour() < msPla.MPlaTo.Hour() {
// 					start, end, _ = MakeTimeStartAndTimeEnd(end, serI.SSTime)
// 				}
// 			}
// 		}
// 	}
// }

// if !isMsEmplos {
// 	for _, empID := range empIDs {
// 		for _, msEmplo := range msEmplos {
// 			if empID != msEmplo.EmployeeID {
// 				continue
// 			}
// 			if end.Hour() > msEmplo.MEmpTo.Hour() && start.Hour() < msEmplo.MEmpTo.Hour() {
// 				start, end, _ = MakeTimeStartAndTimeEnd(end, serI.SSTime)
// 			}
// 		}
// 	}
// }
// func (b *Booking) Appointment(db *gorm.DB, serI ServiceItem) error {
// 	var timeSs []TimeSlot
// 	var msEmplos MEmplos
// 	var msPlas MPlas
// }

// func (mbs []*MasterBooking) Create(tx *sql.Tx) error {
// 	stmt, err := tx.Prepare("")
// 	if err != nil {
// 		return err
// 	}
// }
