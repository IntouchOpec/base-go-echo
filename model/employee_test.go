package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
)

func TestGetEmployeeReady(t *testing.T) {
	emplos := model.Emplos{
		model.Employee{
			EmpoAoumt: 2,
		},
		model.Employee{
			EmpoAoumt: 3,
		},
	}
	ids := []uint{1, 2}
	for index, id := range ids {
		emplos[index].ID = id
	}
	memplos := model.MEmplos{
		model.MasterEmployee{
			MEmpFrom:   time.Now(),
			MEmpTo:     time.Now().Add(1 * time.Hour),
			MEmpQue:    1,
			EmployeeID: 1,
		},
		model.MasterEmployee{
			MEmpFrom:   time.Now(),
			MEmpTo:     time.Now().Add(1 * time.Hour),
			MEmpQue:    2,
			EmployeeID: 1,
		}}
	emplo := memplos.GetEmptyEmployee(emplos)
	fmt.Println("emploEmPla", emplo)
	if emplo == nil {
		t.Errorf("Must have emploce empty")
	}
	if emplo.ID != 2 {
		t.Errorf("Place ready must be emploce id 2")
	}
}

func TestMakeMasterEmployee(t *testing.T) {
	start := time.Now()
	end := start.Add(1 * time.Hour)
	emplo := model.Employee{EmpoAoumt: 2}
	memplos, _ := model.MakeMasterEmployees(start, end, 1, 1, time.Now(), emplo)
	if len(memplos) != 4 {
		t.Errorf("Get %s To %s memplos may have %d", start.Format("15:04"), end.Format("15:04"), len(memplos))
	}
	for index, memplo := range memplos {
		if index != 0 {
			if memplo.MEmpTo.Sub(memplos[index].MEmpFrom) == 0 {
				t.Errorf("%s To %s memplos may have ", memplo.MEmpTo.Format("15:04"), memplos[index].MEmpFrom.Format("15:04"))
			}
			if int(memplo.MEmpTo.Sub(memplo.MEmpFrom)) != int(model.RowDur) {
				t.Errorf("%s To %s may %d diff %d", memplo.MEmpTo.Format("15:04"), memplo.MEmpFrom.Format("15:04"), int(memplo.MEmpTo.Sub(memplo.MEmpFrom)), model.RowDur)
			}
		}
	}
	if memplos[0].MEmpFrom != start {
		t.Errorf("first master emploce may be equal start but %s == %s", start.Format("15:04:01"), memplos[1].MEmpFrom.Format("15:04:01"))
	}
	if memplos[len(memplos)-1].MEmpTo != end {
		t.Errorf("last master must be equal end but %s == %s", memplos[len(memplos)-1].MEmpTo.Format("15:04:01"), end.Format("15:04:01"))
	}
}
