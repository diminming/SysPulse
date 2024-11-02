package housekeeper

import (
	"time"

	"github.com/syspulse/model"
)

func NewHouseKeeper() *Housekeeper {
	housekeeper := new(Housekeeper)
	housekeeper.topoCleaningInterval = 10 * time.Second
	housekeeper.topoTimeout = 1 * time.Hour
	housekeeper.markOverdueInterval = 10 * time.Second
	housekeeper.makrOverdueTimeout = 1 * time.Hour
	return housekeeper
}

type Housekeeper struct {
	topoCleaningInterval time.Duration
	topoTimeout          time.Duration
	markOverdueInterval  time.Duration
	makrOverdueTimeout   time.Duration
}

func (housekeeper *Housekeeper) Run() {
	housekeeper.clearTopo()
}

func (housekeeper *Housekeeper) clearTopo() {
	for {
		timestamp := time.Now().Add(-1 * housekeeper.topoTimeout)
		model.DeleteTimeoutTopo(timestamp.UnixMilli())
		ticker := time.NewTicker(housekeeper.topoCleaningInterval)
		<-ticker.C
	}
}

func (housekeeper *Housekeeper) MarkOverdue() {
	for {
		timestamp := time.Now().Add(-1 * housekeeper.makrOverdueTimeout)
		model.JobMarkOverdue(timestamp.UnixMilli())
		ticker := time.NewTicker(housekeeper.markOverdueInterval)
		<-ticker.C
	}
}
