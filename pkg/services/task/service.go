//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"log"
)

//-- Constants ---------------------------------------------------------------------------------------------------------

//-- Structs -----------------------------------------------------------------------------------------------------------
type taskService struct {
	store  Store
	logger log.Logger
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func NewService(middlewares []Middleware, store Store) Service {
	var service Service

	service = taskService{store: store}
	for _, middleware := range middlewares {
		service = middleware(service)
	}

	return service
}

func (service taskService) Create(ctx context.Context, task *Task) error {
	if err := service.store.insert(ctx, task); err != nil {
		return err
	} else {
		return nil
	}
}

func (service taskService) Update(ctx context.Context, task *Task) error {
	if err := service.store.update(ctx, task); err != nil {
		return err
	} else {
		return nil
	}
}

func (service taskService) Read(ctx context.Context, id uint) (*Task, error) {
	if task, err := service.store.read(ctx, id); err != nil {
		return nil, err
	} else {
		return task, nil
	}
}

func (service taskService) Delete(ctx context.Context, id uint) (*Task, error) {
	if task, err := service.store.delete(ctx, id); err != nil {
		return nil, err
	} else {
		return task, nil
	}
}

func (service taskService) List(ctx context.Context, limit uint, offset uint) ([]Task, error) {
	if tasks, err := service.store.list(ctx, limit, offset); err != nil {
		return nil, err
	} else {
		return tasks, nil
	}
}

func (service taskService) Shutdown() error {
	if err := service.store.Close(); err != nil {
		return err
	} else {
		return nil
	}
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
