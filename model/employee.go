package model

import (
	"errors"
	"time"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

type Employee struct {
	orm.ModelBase
	EmpoAoumt     int         `json:"empo_amount" form:"empo_amount"`
	EmpoName      string      `form:"empo_name" json:"empo_name" gorm:"type:varchar(25)"`
	EmpoDetail    string      `form:"empo_detail" json:"empo_detail"`
	EmpoLineID    string      `form:"empo_line_id" json:"empo_line_id" gorm:"type:varchar(50)"`
	EmpoImage     string      `form:"image" json:"empo_image" gorm:"type:varchar(255)"`
	EmpoIsActive  bool        `json:"empo_is_active" sql:"default:true"`
	AccountID     uint        `form:"account_id" json:"account_id"`
	ChatChannelID uint        `form:"chat_channel_id" json:"chat_channel_id"`
	Services      []*Service  `json:"services" gorm:"many2many:employee_service"`
	TimeSlots     []*TimeSlot `josn:"time_slots" gorm:"foreignkey:EmployeeID;association_foreignkey:EmployeeID"`
	ChatChannel   ChatChannel `json:"chat_channel" gorm:"ForeignKey:ChatChannelID"`
	Account       Account     `json:"account" gorm:"ForeignKey:AccountID"`
}

type MasterEmployee struct {
	orm.ModelBase

	MEmpQue int `json:"m_emp_que"`
	// MEmpStatus MEmpStatus `json:"m_emp_status"`
	MEmpDay    time.Time `json:"m_emp_day"`
	MEmpFrom   time.Time `json:"m_emp_from"`
	MEmpTo     time.Time `json:"m_emp_to"`
	EmployeeID uint      `json:"employee_id"`
	Employee   *Employee `json:"employee" gorm:"ForeignKey:EmployeeID"`
	AccountID  uint      `json:"account_id"`
	Account    Account   `json:"account" gorm:"ForeignKey:AccountID"`
}

func (prov *Employee) CreateEmployee() error {
	if err := DB().Create(&prov).Error; err != nil {
		return err
	}

	return nil
}

func (prov *Employee) UpdateEmployee() error {
	if err := DB().Save(&prov).Error; err != nil {
		return err
	}

	return nil
}

func GetEmployeeList(accID uint) ([]*Employee, error) {
	provs := []*Employee{}
	if err := DB().Where("account_id = ?", accID).Find(&provs).Error; err != nil {
		return nil, err
	}

	return provs, nil
}

func GetEmployeeDetail(id string, accID uint) (*Employee, error) {
	prov := Employee{}
	if err := DB().Preload("Services").Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func GetEmployeeServiceTimeSlotList(id string, accID uint) (*Employee, error) {
	prov := Employee{}
	if err := DB().Preload("EmployeeServices", func(db *gorm.DB) *gorm.DB {
		return db.Preload("TimeSlots").Preload("Service")
	}).Where("account_id = ?", accID).Find(&prov, id).Error; err != nil {
		return nil, err
	}

	return &prov, nil
}

func RemoveEmployee(id string) (*Employee, error) {
	prov := Employee{}
	db := DB()
	if err := db.Find(&prov, id).Error; err != nil {
		return nil, err
	}

	if err := db.Delete(&prov).Error; err != nil {
		return nil, err
	}
	return &prov, nil
}

func (pro *Employee) RemoveImage() error {
	pro.EmpoImage = ""

	if err := DB().Save(&pro).Error; err != nil {
		return err
	}

	return nil
}

func (mEmp MasterEmployee) GetMEmplo(epos Emplos, db *gorm.DB) (MEmplos, error) {
	var memps MEmplos
	var empIDs []uint
	for _, epo := range epos {
		empIDs = append(empIDs, epo.ID)
	}
	if err := db.Where("employee_id in (?) and account_id = ? and m_emp_day = ? and m_emp_from BETWEEN ? and ? or m_emp_to BETWEEN ? and ?",
		empIDs,
		mEmp.AccountID,
		mEmp.MEmpDay,
		mEmp.MEmpFrom,
		mEmp.MEmpTo,
		mEmp.MEmpFrom,
		mEmp.MEmpTo).Find(&memps).Error; err != nil {
		return nil, err
	}
	return memps, nil
}

func SetTimeSlotMixEmployee(timeSs []TimeSlot, empls []*Employee) []uint {
	var empIDs []uint
	for _, timeS := range timeSs {
		for _, emp := range empls {
			if emp.ID == timeS.EmployeeID {
				empIDs = append(empIDs, emp.ID)
			}
		}
	}
	return empIDs
}

func (mEmp MasterEmployee) IsReadyEmployee(epos Emplos, db *gorm.DB) (bool, error) {
	var memps MEmplos
	var empIDs []uint
	for _, epo := range epos {
		empIDs = append(empIDs, epo.ID)
	}
	if err := db.Where("employee_id in (?) and account_id = ? and m_emp_day = ?",
		empIDs,
		mEmp.AccountID,
		mEmp.MEmpDay).Find(&memps).Error; err != nil {
		return false, err
	}
	if len(memps) == 0 {
		return true, nil
	}

	return false, nil
}

func (eMps Emplos) GetEmployeeReady(mEmp MasterEmployee, db *gorm.DB) (*Employee, error) {
	memp, err := mEmp.GetMEmplo(eMps, db)
	if err != nil {
		return nil, err
	}
	emp := memp.GetEmptyEmployee(eMps)
	if err != nil {
		return nil, errors.New("")
	}
	return emp, nil
}

type MEmplos []MasterEmployee
type Emplos []Employee

func (eMps Emplos) GetEmployeeEmpty(mEmp MasterEmployee, db *gorm.DB) (*Employee, error) {
	mEmplos, err := mEmp.GetMEmplo(eMps, db)
	if err != nil {
		return nil, err
	}
	emp := mEmplos.GetEmptyEmployee(eMps)
	if err != nil {
		return nil, errors.New("")
	}
	return emp, nil
}

func (memp MEmplos) Create(db *gorm.DB) error {
	if err := db.Create(&memp).Error; err != nil {
		return err
	}
	return nil
}

func (memp MEmplos) GetEmptyEmployee(emps []Employee) *Employee {
	var isEmpty bool
	for _, emp := range emps {
		isEmpty = true
		for _, memp := range memp {
			if memp.EmployeeID != emp.ID {
				continue
			}
			if memp.MEmpQue >= 1 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			return &emp
		}
	}
	return nil
}

func MakeMasterEmployees(start, end time.Time, accID uint, day time.Time, emp Employee) ([]*MasterEmployee, error) {
	var memps []*MasterEmployee
	diff := end.Sub(start) / RowDur

	for i := 0; i < int(diff); i++ {
		var from time.Time
		from = start.Add(RowDur * time.Duration(i))
		to := start.Add(RowDur * time.Duration(i+1))
		mEmp := MasterEmployee{
			// MEmpQue:    quo,
			EmployeeID: emp.ID,
			MEmpDay:    day,
			AccountID:  accID,
			MEmpFrom:   from,
			MEmpTo:     to,
		}
		memps = append(memps, &mEmp)
	}

	return memps, nil
}
