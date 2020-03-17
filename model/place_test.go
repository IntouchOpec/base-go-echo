package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/IntouchOpec/base-go-echo/model"
)

func TestGetEmptyPlace(t *testing.T) {
	plas := model.Plas{
		model.Place{
			PlacAmount: 2,
		},
		model.Place{
			PlacAmount: 3,
		},
	}
	ids := []uint{1, 2}
	for index, id := range ids {
		plas[index].ID = id
	}
	mplas := model.MPlas{
		model.MasterPlace{
			MPlaFrom: time.Now(),
			MPlaTo:   time.Now().Add(1 * time.Hour),
			MPlaQue:  1,
			PlaceID:  1,
		},
		model.MasterPlace{
			MPlaFrom: time.Now(),
			MPlaTo:   time.Now().Add(1 * time.Hour),
			MPlaQue:  2,
			PlaceID:  1,
		}}
	plaEmPla := mplas.GetEmptyPlace(plas)
	fmt.Println("plaEmPla", plaEmPla)
	if plaEmPla == nil {
		t.Errorf("Must have place empty")
	}
	if plaEmPla.ID != 2 {
		t.Errorf("Place ready must be place id 2")
	}
}

func TestMakeMasterPlaces(t *testing.T) {
	start := time.Now()
	end := start.Add(1 * time.Hour)
	pla := model.Place{PlacAmount: 2}
	mplas, _ := model.MakeMasterPlaces(start, end, 1, 1, time.Now(), pla)
	if len(mplas) != 4 {
		t.Errorf("Get %s To %s mplas may have %d", start.Format("15:04"), end.Format("15:04"), len(mplas))
	}
	for index, mpla := range mplas {
		if index != 0 {
			if mpla.MPlaTo.Sub(mplas[index].MPlaFrom) == 0 {
				t.Errorf("%s To %s mplas may have ", mpla.MPlaTo.Format("15:04"), mplas[index].MPlaFrom.Format("15:04"))
			}
			if int(mpla.MPlaTo.Sub(mpla.MPlaFrom)) != int(model.RowDur) {
				t.Errorf("%s To %s may %d diff %d", mpla.MPlaTo.Format("15:04"), mpla.MPlaFrom.Format("15:04"), int(mpla.MPlaTo.Sub(mpla.MPlaFrom)), model.RowDur)
			}
		}
	}
	if mplas[0].MPlaFrom != start {
		t.Errorf("first master emplo may be equal start but %s == %s", start.Format("15:04:01"), mplas[1].MPlaFrom.Format("15:04:01"))
	}
	if mplas[len(mplas)-1].MPlaTo != end {
		t.Errorf("last master must be equal end but %s == %s", mplas[len(mplas)-1].MPlaTo.Format("15:04:01"), end.Format("15:04:01"))
	}
}

func TestMakeTime(t *testing.T) {
	d, _ := time.Parse("2006-01-02 15:04:01", "2006-01-02 15:04:01")
	d, _ = model.MakeTime(d)
	if d.Minute() != 45 {
		t.Errorf("give %s want %s but %d ", "15:04", "15:45", d.Minute())
	}
}
