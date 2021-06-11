package lb

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"sync"
)

type PoolSync interface {
	Sync(pool TrinoPool) error
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
	ToRemove []CoordinatorRef
}

func (p *PoolStateSync) Sync(pool TrinoPool) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ctx := context.Background()

	actualCoordinators := pool.Fetch(FetchRequest{})
	expectedCoordinators, err := p.storage.All(ctx)

	if err != nil {
		return err
	}

	syncAction := getSyncAction(actualCoordinators, expectedCoordinators)

	if len(syncAction.ToAdd) != 0 || len(syncAction.ToRemove) != 0 {
		p.logger.Info("new pool coordinators retrieved: Add %d, Remove: %d", len(syncAction.ToAdd), len(syncAction.ToRemove))

		for _, removed := range syncAction.ToRemove {
			if err := pool.Remove(removed.ID); err != nil {
				return err
			}
		}

		for _, added := range syncAction.ToAdd {
			if err := pool.Add(added); err != nil {
				return err
			}
		}
	}

	// After pool items sync we sync items coordinators
	for _, currItem := range actualCoordinators {
		for _, stateItem := range expectedCoordinators {
			if currItem.Name == stateItem.Name {
				if err := pool.Update(currItem.ID, stateItem); err != nil {
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

func getSyncAction(current []CoordinatorRef, state []models.Coordinator) syncAction {
	toAdd := make([]models.Coordinator, 0)
	toRemove := make([]CoordinatorRef, 0)

	for _, curr := range current {
		if !containsCoord(state, curr.Coordinator) {
			toRemove = append(toRemove, curr)
		}
	}

	for _, stat := range state {
		if !containsTarget(current, stat) {
			toAdd = append(toAdd, stat)
		}
	}

	return syncAction{
		ToAdd:    toAdd,
		ToRemove: toRemove,
	}
}

func containsTarget(src []CoordinatorRef, target models.Coordinator) bool {
	for _, c := range src {
		if target.Name == c.Name {
			return true
		}
	}
	return false
}

func containsCoord(src []models.Coordinator, target models.Coordinator) bool {
	for _, c := range src {
		if target.Name == c.Name {
			return true
		}
	}
	return false
}

type NoOpSync struct{}

func (n NoOpSync) Sync(pool TrinoPool) error {
	return nil
}
