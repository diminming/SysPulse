package housekeeper

import (
	"time"

	"github.com/syspulse/model"
)

func NewHouseKeeper() *Housekeeper {
	housekeeper := new(Housekeeper)
	housekeeper.topoCleaningInterval = 10 * time.Second
	housekeeper.topoTimeout = 1 // unit: Hour
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
		timestamp := time.Now().UnixMilli() - housekeeper.topoTimeout*60*60*1000
		model.DeleteTimeoutTopo(timestamp)
		ticker := time.NewTicker(housekeeper.topoCleaningInterval)
		<-ticker.C
	}
}
