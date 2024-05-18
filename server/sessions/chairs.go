package sessions

import (
	"maps"
	"sync"

	"github.com/DaanV2/f1-game-dashboards/server/pkg/hooks"
)

type (
	// ChairManager is a struct that manages all the chairs
	ChairManager struct {
		chairs_lock sync.RWMutex
		chairs      map[int]Chair

		OnChairAdded   hooks.Hook[Chair]
		OnChairUpdated hooks.Hook[Chair]
		OnChairRemoved hooks.Hook[Chair]
	}

	// Chair is a struct that represents a chair
	Chair struct {
		active bool   // If the chair is active
		name   string // The name of the chair
		port   int    // The upd port of the chair
	}
)

// NewChair creates a new chair
func NewChair(name string, port int, active bool) *Chair {
	return &Chair{
		active,
		name,
		port,
	}
}

// Id returns the id of the chair
func (c *Chair) Id() int {
	return c.port
}

// Port returns the upd port of the chair
func (c *Chair) Port() int {
	return c.port
}

// IsActive returns if the chair is active
func (c *Chair) IsActive() bool {
	return c.active
}

// Name returns the name of the chair
func (c *Chair) Name() string {
	return c.name
}

// NewChairManager creates a new chair manager
func NewChairManager() *ChairManager {
	return &ChairManager{
		chairs_lock: sync.RWMutex{},
		chairs:      make(map[int]Chair),
	}
}

// AddChair adds a chair to the chair manager
func (cm *ChairManager) AddChair(chair Chair) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	cm.chairs[chair.Id()] = chair
	cm.OnChairAdded.Call(chair)
}

// UpdateChair updates a chair in the chair manager
func (cm *ChairManager) UpdateChair(chair Chair) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	cm.chairs[chair.Id()] = chair
	cm.OnChairUpdated.Call(chair)
}

// GetChair gets a chair from the chair manager
func (cm *ChairManager) GetChair(id int) Chair {
	cm.chairs_lock.RLock()
	defer cm.chairs_lock.RUnlock()

	return cm.chairs[id]
}

// RemoveChair removes a chair from the chair manager
func (cm *ChairManager) RemoveChair(id int) {
	cm.chairs_lock.Lock()
	defer cm.chairs_lock.Unlock()

	delete(cm.chairs, id)
	cm.OnChairRemoved.Call(cm.GetChair(id))
}

// Chairs returns all the chairs in the chair manager
func (cm *ChairManager) Chairs() map[int]Chair {
	cm.chairs_lock.RLock()
	defer cm.chairs_lock.RUnlock()

	return maps.Clone(cm.chairs)
}
