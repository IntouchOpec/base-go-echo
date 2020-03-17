package main

import "github.com/IntouchOpec/base-go-echo/model"

func CheckTime(t *testing.T) {
	// model.MasterPlace
	mplas :=  []*model.MasterPlace{}
	start := time.Now()
	end := start.Add(1 *time.Hour)
	model.IsEmptyMasterPlace(start,end, mplas)
}
