package sessions

import (
	"maps"
	"sync"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/hooks"
)

type ChairManager struct {
	lock sync.RWMutex

	chairs map[int]*Chair

	OnChairAdded   hooks.Hook[*Chair]
	OnChairRemoved hooks.Hook[*Chair]
}

type Chair struct {
	Active bool   // If the chair is active
	Name   string // The name of the chair
	Port   int    // The upd port of the chair
}

func NewChair(name string, port int) *Chair {
	return &Chair{
		Active: true,
		Name:   name,
		Port:   port,
	}
}

func (c *Chair) Id() int {
	return c.Port
}

func NewChairManager() *ChairManager {
	return &ChairManager{
		lock:   sync.RWMutex{},
		chairs: make(map[int]*Chair),
	}
}

func (cm *ChairManager) AddChair(chair *Chair) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.chairs[chair.Id()] = chair
	cm.OnChairAdded.Call(chair)
}

func (cm *ChairManager) GetChair(id int) *Chair {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	return cm.chairs[id]
}

func (cm *ChairManager) RemoveChair(id int) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.chairs, id)
	cm.OnChairRemoved.Call(cm.GetChair(id))
}

func (cm *ChairManager) Chairs() map[int]*Chair {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	return maps.Clone(cm.chairs)
}
