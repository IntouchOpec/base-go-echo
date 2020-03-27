package model

import (
	"fmt"
	"time"
)

type ms struct {
	id uint 
	start time.Time
	end time.Time
	ms struct {
		start time.Time
		end time.Time
	}
}
func TestMakeTimeSlotBookingNow(t *testing.T) {
	ms1 = ms{
		id:1, 
		start: time.Parse("2006-01-02 15:02:01", "0001-01-01 08:00:00"),
		end: time.Parse("2006-01-02 15:02:01", "0001-01-01 16:00:00") ,
	ms: [struct{start: time.Parse("2006-01-02 15:02:00", "0001-01-01 16:15:00") ,end : time.Parse("2006-01-02 15:02:01","0001-01-01 17:15:00") }]} 
	ms2 = ms{id:3, start: time.Parse("2006-01-02 15:02:01", "0000-12-31 22:00:00"),end: time.Parse("2006-01-02 15:02:01", "0001-01-01 16:00:00") ,
	ms: [struct{start: time.Parse("2006-01-02 15:02:00", "0001-01-01 16:15:00") ,end : time.Parse("2006-01-02 15:02:01","0001-01-01 17:30:00") }]} 
	ms3 = ms{id:4, start: time.Parse("2006-01-02 15:02:01", "0001-01-01 10:00:00"),end: time.Parse("2006-01-02 15:02:01", "0001-01-01 16:00:00") ,
	ms: [struct{start: time.Parse("2006-01-02 15:02:00", "0001-01-01 16:15:00") ,end : time.Parse("2006-01-02 15:02:01","0001-01-01 17:30:00") }]}
}
