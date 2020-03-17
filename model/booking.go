package model

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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

	Booking  *Booking  `json:"booking" gorm:ForeignKey:BookingID""`
	Employee *Employee `json:"employee" gorm:"ForeignKey:EmployeeID"`
	Place    *Place    `json:"place" gorm:"ForeignKey:PlaceID"`
	Account  *Account  `json:"account" gorm:"ForeignKey:AccountID"`
}

type place struct {
	id     uint
	amount int
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
	values := []interface{}{}

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
		fmt.Println("====1 ", err)
		return nil, "notFound", err
	}
	var serIID uint
	var booPacID uint
	var emploID uint
	for rows.Next() {
		var boo Booking
		err := rows.Scan(&boo.ID, &boo.BookedStart, &boo.BookedEnd, &boo.BookedDay, &serIID, &booPacID, &emploID)
		fmt.Println(err)
		boos = append(boos, boo)
	}
	fmt.Println(emploID, tran.ID)
	for _, boo := range boos {
		var empIDs []string
		var plas []place
		var plaIDs []string

		if serIID != 0 {
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
				var pla place
				rows.Scan(&pla.id, &pla.amount)
				plaIDs = append(plaIDs, fmt.Sprintf("%d", pla.id))
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
					fmt.Println(empID)
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
			fmt.Println(err, "====5")
			var plaMD MasterBooking
			var plaMDs []MasterBooking
			var isPlaReady bool = true
			for _, plaI := range plas {
				for rows.Next() {
					// var pla MasterBooking
					rows.Scan(&plaMD.PlaceID, &plaMD.MBQue, &plaMD.MBFrom, &plaMD.MBTo)
					if plaI.id == plaMD.PlaceID {
						if plaI.amount >= plaMD.MBQue {
							plaMDs = make([]MasterBooking, 0)
							isPlaReady = false
						} else {
							plaMDs = append(plaMDs, plaMD)
						}
						plas = plas[1:]
					}
				}

				if isPlaReady {
					break
				}
			}
			if !isPlaReady {
				return nil, "notPlaceReady", errors.New("place not ready.")
			}
			if plaMD.ID == 0 {
				plaMD.PlaceID = plas[0].id
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
				fmt.Println(plaMD.EmployeeID, emplID)
				values = append(values,
					plaMD.MBQue, plaMD.EmployeeID, boo.ID, plaMD.PlaceID, plaMD.MBDay, plaMD.MBFrom, plaMD.MBTo)
				numFields := 7
				n := i * numFields

				query += `(`
				for j := 0; j < numFields; j++ {
					query += `$` + strconv.Itoa(n+j+1) + `,`
				}
				query = query[:len(query)-1] + `),`
			}
		} else {

		}

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
	stmt, err := sql.Prepare(fmt.Sprintf("INSERT INTO master_bookings (mb_que,employee_id,booking_id,place_id,mb_day,mb_from,mb_to) VALUES %s ", vStr))
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

func MakeTimeSlotBookingNow(d time.Time, serI ServiceItem, db *gorm.DB) (time.Time, time.Time, error) {
	var msPlas []MasterPlace
	var msEmplos MEmplos
	var timeSs []TimeSlot
	if len(serI.Service.Places) == 0 {
		return d, d, errors.New("")
	}
	plaIDs := SetPlaces(serI.Service.Places)
	start, end, err := MakeTimeStartAndTimeEnd(d, serI.SSTime)
	if err != nil {

	}
	db.Order("employee_id").Where("account_id = ? and time_day = ? and time_start >= ?",
		serI.AccountID, int(d.Weekday()), start).Find(&timeSs)

	if len(timeSs) == 0 {
		return d, d, errors.New("")
	}

	empIDs := SetTimeSlotMixEmployee(timeSs, serI.Service.Employees)

	db.Order("m_empo_from").Where("account_id = ? and m_emplo_day = ? and employee_id in (?) and m_empo_from >= ?",
		serI.AccountID, d, empIDs, start).Find(&msEmplos)
	db.Order("m_empo_from").Where("account_id = ? and m_pla_day = ? and place_id in (?) and m_emp_from >= ?",
		serI.AccountID, d, plaIDs, start).Find(&msPlas)
	isMsPlas := len(msPlas) == 0
	isMsEmplos := len(msEmplos) == 0
	if isMsPlas && isMsEmplos {
		return start, end, nil
	}
	if !isMsPlas {
		// empl := msEmplos.GetEmptyEmployee(serI.Service.Employees)
	}
	if !isMsEmplos {

	}
	return start, end, nil
}

func (msEmplos *MEmplos) Get(db *gorm.DB, where string, values ...interface{}) error {
	// "account_id = ? and m_emplo_day = ? and employee_id in (?) and m_empo_from >= ?"
	if err := db.Order("employee_id, m_empo_from").Where(where, values).Find(&msEmplos).Error; err != nil {
		return err
	}
	return nil
}

func (msPlas *MPlas) Get(db *gorm.DB, where string, values ...interface{}) error {
	// "account_id = ? and m_pla_day = ? and place_id in (?) and m_emp_from >= ?"
	if err := db.Order("place_id, m_empo_from").Where(where, values).Find(&msPlas).Error; err != nil {
		return err
	}
	return nil
}

type serEmp struct {
	emps []uint
	id   uint
}

type sersPla struct {
	plas []place
	id   uint
}

type pack struct {
	id       uint
	serEmps  []serEmp
	sersPlas []sersPla
}

func clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func (b *Booking) PackNow(db *sql.DB, pack Pack, serIDs []string) error {
	d := time.Now()
	var sersPlas []sersPla
	start, end, err := MakeTimeStartAndTimeEnd(d, pack.TimeUse)
	if err != nil {
		return err
	}
	var strSerIDs string
	for i := 4; i < len(serIDs)+4; i++ {
		strSerIDs += fmt.Sprintf("%s,", serIDs[i-4])
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
		AND time_end > $3
	INNER JOIN employees AS e ON e.id = es.employee_id AND e.deleted_at IS NULL AND e.empo_is_active = true
	WHERE es.service_id IN %s
	ORDER BY es.service_id, es.employee_id`, strSerIDs)
	rows, err := db.Query(qe, int(d.Weekday()), start, end)
	if err != nil {
		return err
	}
	var serEmps []serEmp
	var serID string
	var empID uint
	var serEmp serEmp
	for rows.Next() {
		rows.Scan(&serID, &empID)
		idStr := fmt.Sprintf("%d", serEmp.id)
		if serID == idStr {
			serEmp.emps = append(serEmp.emps, empID)
		} else {
			if serID != "" {
				uintID, _ := strconv.ParseUint(serID, 10, 64)
				serEmp.id = uint(uintID)
				serEmps = append(serEmps, serEmp)
				serEmp.emps = append(serEmp.emps, empID)
				clear(serEmp)
			}
		}
		serID = idStr
	}
	if len(serEmps) == 0 {
		uintID, _ := strconv.ParseUint(serID, 10, 64)
		serEmp.id = uint(uintID)
		serEmp.emps = append(serEmp.emps, empID)
		serEmps = append(serEmps, serEmp)
	}
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
	var pla sersPla
	var p place
	serID = ""
	for rows.Next() {
		var amount int
		var plaID uint
		var id uint
		rows.Scan(&id, &plaID, &amount)
		idStr := fmt.Sprintf("%d", plaID)
		if serID == idStr {
			p.id = id
			p.amount = amount
			pla.plas = append(pla.plas, p)
		} else {
			if serID != "" {
				uintID, _ := strconv.ParseUint(serID, 10, 64)
				pla.id = uint(uintID)

				sersPlas = append(sersPlas, pla)
				clear(pla)
			}
		}
	}
	if len(sersPlas) == 0 {

	}
	d, _ = time.Parse("2006-01-02", d.Format("2006-01-02"))
	fmt.Println(d, start, end)
	b.BookedDay = d
	b.BookedStart = start
	b.BookedEnd = end

	return nil
}

func (b *Booking) Now(db *sql.DB, serI ServiceItem) error {
	var emploIDs []string
	var plaIDs []string
	var plas []place
	d := time.Now()
	start, end, err := MakeTimeStartAndTimeEnd(d, serI.SSTime)
	if err != nil {
		return err
	}
	start, end = start.Add(-(7 * time.Hour)), end.Add(-(7 * time.Hour))
	rows, err := db.Query(`
			SELECT 
				es.employee_id 
			FROM employee_service AS es
			INNER JOIN time_slots AS ts ON ts.employee_id = es.employee_id 
				AND ts.deleted_at IS NULL 
				AND time_day = $1
				AND time_start < $2 
				AND time_end > $3
			WHERE es.service_id = $4 AND deleted_at IS NULL
			ORDER BY es.employee_id`, int(d.Weekday()), start, end, serI.Service.ID)
	if err == nil {
		for rows.Next() {
			var id string
			rows.Scan(&id)
			// fmt.Println(id, "=id")
			emploIDs = append(emploIDs, id)
		}
	}
	if len(emploIDs) == 0 {
		return errors.New("not employee")
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
		return err
	}
	for rows.Next() {
		var pla place
		rows.Scan(&pla.id, pla.amount)
		plaIDs = append(plaIDs, fmt.Sprintf("%d", pla.id))
		plas = append(plas, pla)
	}
	if len(plaIDs) == 0 {
		return errors.New("")
	}
	rows, err = db.Query(`
		SELECT 
			employee_id
		FROM master_bookings AS mb 
		WHERE 
			deleted_at IS NULL AND employee_id IN ($1) AND account_id = $2 AND mb_day = $3 AND mb_from BETWEEN $4 AND $5 OR mb_to BETWEEN $6 AND $7
		GROUP BY employee_id
		ORDER BY employee_id`, strings.Join(emploIDs, ","), b.AccountID, d,
		start, end, start, end)
	var emplID uint
	var isReady bool = true
	if err == nil {
		for rows.Next() {
			rows.Scan(&emplID)
			for _, empID := range emploIDs {
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
	} else {
		emploID, _ := strconv.ParseUint(emploIDs[0], 10, 64)
		emplID = uint(emploID)
	}

	if !isReady {
		return errors.New("place not ready.")
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
	fmt.Println("error", err)
	var plaMD MasterBooking
	var plaMDs []MasterBooking
	var isPlaReady bool = true
	for _, plaI := range plas {
		for rows.Next() {
			rows.Scan(&plaMD.PlaceID, &plaMD.MBQue, &plaMD.MBFrom, &plaMD.MBTo)
			if plaI.id == plaMD.PlaceID {
				if plaI.amount >= plaMD.MBQue {
					plaMDs = make([]MasterBooking, 0)
					isPlaReady = false
				} else {
					plaMDs = append(plaMDs, plaMD)
				}
				if len(plas) == 1 {
					plas = make([]place, 0)
					break
				}
				plas = plas[1:]
			}
		}

		if isPlaReady {
			break
		}
	}
	if !isPlaReady {
		return errors.New("place not ready.")
	}

	d, _ = time.Parse("2006-01-02", d.Format("2006-01-02"))
	b.BookedDay = d
	b.BookedStart = start
	b.BookedEnd = end

	return nil
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
