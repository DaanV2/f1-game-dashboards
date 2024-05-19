package sessions

import (
	"fmt"
	"maps"
	"sync"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/hooks"
)

type (
	// ChairManager is a struct that manages all the chairs
	ChairManager struct {
		chairs_lock sync.RWMutex
		chairs      map[string]Chair

		OnChairAdded   hooks.Hook[Chair]
		OnChairUpdated hooks.Hook[Chair]
		OnChairRemoved hooks.Hook[Chair]
	}

	// Chair is a readonly struct that represents a chair
	Chair struct {
		Active bool   `json:"is_active"` // readonly, If the chair is active
		Name   string `json:"name"`      // readonly, The name of the chair
		Port   int    `json:"port"`      // readonly, The upd port of the chair
	}
)

// NewChair creates a new chair
func NewChair(name string, port int, active bool) Chair {
	return Chair{
		active,
		name,
		port,
	}
}

// Id returns the id of the chair
func (c *Chair) Id() string {
	return fmt.Sprint(c.Port)
}

// NewChairManager creates a new chair manager
func NewChairManager() *ChairManager {
	return &ChairManager{
		chairs_lock: sync.RWMutex{},
		chairs:      make(map[string]Chair),
	}
}

// Add adds a chair to the chair manager
func (cm *ChairManager) Add(chair Chair) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	cm.chairs[chair.Id()] = chair
	cm.OnChairAdded.Call(chair)
}

// Update updates a chair in the chair manager
func (cm *ChairManager) Update(chair Chair) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	cm.chairs[chair.Id()] = chair
	cm.OnChairUpdated.Call(chair)
}

// Get gets a chair from the chair manager
func (cm *ChairManager) Get(id string) (Chair, bool) {
	cm.chairs_lock.RLock()
	defer cm.chairs_lock.RUnlock()

	c, ok := cm.chairs[id]
	return c, ok
}

// Remove removes a chair from the chair manager
func (cm *ChairManager) Remove(id string) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	ch, ok := cm.chairs[id]
	if !ok {
		return
	}

	delete(cm.chairs, id)
	cm.OnChairRemoved.Call(ch)
}

// All returns all the chairs in the chair manager
func (cm *ChairManager) All() map[string]Chair {
	cm.chairs_lock.RLock()
	defer cm.chairs_lock.RUnlock()

	return maps.Clone(cm.chairs)
}
