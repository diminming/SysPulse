package housekeeper

import (
	"syspulse/model"
	"time"
)

func NewHouseKeeper() *Housekeeper {
	housekeeper := new(Housekeeper)
	housekeeper.topoCleaningInterval = 3600 * time.Second
	housekeeper.topoTimeout = 12 // unit: Hour
	return housekeeper
}

type Housekeeper struct {
	topoCleaningInterval time.Duration
	topoTimeout          int64
}

func (housekeeper *Housekeeper) Run() {
	housekeeper.clearTopo()
}

func (housekeeper *Housekeeper) clearTopo() {

	for {
		ticker := time.NewTicker(housekeeper.topoCleaningInterval)
		<-ticker.C
		timestamp := time.Now().UnixMilli() - housekeeper.topoTimeout*60*60*1000
		model.DeleteTimeoutTopo(timestamp)
	}
}
