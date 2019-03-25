//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
)

//-- Constants ---------------------------------------------------------------------------------------------------------
type Model interface {
	compare() bool
	sanitize() error
	validate() error
}

type Service interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error

	Read(ctx context.Context, id uint) (*Task, error)
	Delete(ctx context.Context, id uint) (*Task, error)

	List(ctx context.Context, limit uint, offset uint) ([]Task, error)

	Shutdown() error
}

type Middleware func(Service) Service

type Store interface {
	Open(options string) error
	Close() error

	Prepare(option string, parameter string) error

	insert(ctx context.Context, task *Task) error
	update(ctx context.Context, task *Task) error

	read(ctx context.Context, id uint) (*Task, error)
	delete(ctx context.Context, id uint) (*Task, error)

	list(ctx context.Context, limit uint, offset uint) ([]Task, error)
}

//-- Structs -----------------------------------------------------------------------------------------------------------
type ConnectionParameters struct {
	Driver  string
	Options string
}

//-- Exported Functions ------------------------------------------------------------------------------------------------

//-- Internal Functions ------------------------------------------------------------------------------------------------
