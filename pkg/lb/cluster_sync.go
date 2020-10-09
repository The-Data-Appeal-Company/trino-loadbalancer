package lb

import (
	"context"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"sync"
)

type PoolSync interface {
	Sync(pool HttpPool) error
}

type PoolStateSync struct {
	storage discovery.Storage
	logger  logging.Logger
	mutex   *sync.Mutex
}

func NewPoolStateSync(storage discovery.Storage, logger logging.Logger) *PoolStateSync {
	return &PoolStateSync{
		storage: storage,
		logger:  logger,
		mutex:   &sync.Mutex{},
	}
}

type syncAction struct {
	ToAdd    []models.Coordinator
	ToRemove []models.Coordinator
}

func (p *PoolStateSync) Sync(pool HttpPool) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ctx := context.Background()

	current := coordinators(pool.AllBackends())
	state, err := p.storage.All(ctx)

	if err != nil {
		return err
	}

	syncAction := getSyncAction(current, state)

	if len(syncAction.ToAdd) != 0 || len(syncAction.ToRemove) != 0 {
		p.logger.Info("new pool state retrieved: Add %d, Remove: %d", len(syncAction.ToAdd), len(syncAction.ToRemove))

		for _, removed := range syncAction.ToRemove {
			if err := pool.Remove(removed.Name); err != nil {
				return err
			}
		}

		for _, added := range syncAction.ToAdd {
			if err := pool.Add(added); err != nil {
				return err
			}
		}
	}

	// After pool items sync we sync items state
	for _, currItem := range current {
		for _, stateItem := range state {
			if currItem.Name == stateItem.Name {
				if err := pool.Update(currItem.Name, stateItem); err != nil {
					return err
				}

				if currItem.Enabled != stateItem.Enabled {
					p.logger.Info("cluster %s status: %b", currItem.Name, currItem.Enabled)
				}
			}
		}
	}

	return nil
}

func getSyncAction(current []models.Coordinator, state []models.Coordinator) syncAction {
	toAdd := make([]models.Coordinator, 0)
	toRemove := make([]models.Coordinator, 0)

	for _, curr := range current {
		if !contains(state, curr) {
			toRemove = append(toRemove, curr)
		}
	}

	for _, stat := range state {
		if !contains(current, stat) {
			toAdd = append(toAdd, stat)
		}
	}

	return syncAction{
		ToAdd:    toAdd,
		ToRemove: toRemove,
	}
}

func coordinators(backends []*CoordinatorConnection) []models.Coordinator {
	coordinators := make([]models.Coordinator, len(backends))
	for i := range backends {
		coordinators[i] = backends[i].Backend
	}
	return coordinators
}

func contains(src []models.Coordinator, target models.Coordinator) bool {
	for _, c := range src {
		if target.Name == c.Name {
			return true
		}
	}
	return false
}

type NoOpSync struct{}

func (n NoOpSync) Sync(pool HttpPool) error {
	return nil
}
