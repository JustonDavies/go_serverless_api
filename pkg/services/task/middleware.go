//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

//-- Constants ---------------------------------------------------------------------------------------------------------
var (
	logFormat = "\n\tRequest ID -> %s \n\tFunction   -> %s \n\tParameters -> %s \n\tResponse   -> %v \n\tError      -> %v"
)

//-- Structs -----------------------------------------------------------------------------------------------------------
type logMiddleware struct {
	next   Service
	logger log.Logger
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func NewLogMiddleare(logger log.Logger) Middleware {
	return func(next Service) Service {
		return logMiddleware{next, logger}
	}
}

func (middleware logMiddleware) Create(ctx context.Context, task *Task) error {
	var err error
	var parameterCapture string

	parameterCapture = fmt.Sprintf(`%v`, task)
	err = middleware.next.Create(ctx, task)

	middleware.logger.Printf(logFormat, uuid.New().String(), `task create`, parameterCapture, task, err)
	return err
}

func (middleware logMiddleware) Update(ctx context.Context, task *Task) error {
	var err error
	var parameterCapture string

	parameterCapture = fmt.Sprintf(`%v`, task)
	err = middleware.next.Update(ctx, task)

	middleware.logger.Printf(logFormat, uuid.New().String(), `task update`, parameterCapture, task, err)
	return err
}

func (middleware logMiddleware) Read(ctx context.Context, id uint) (*Task, error) {
	var err error
	var result *Task
	var parameterCapture string

	parameterCapture = fmt.Sprintf(`%d`, id)
	result, err = middleware.next.Read(ctx, id)

	middleware.logger.Printf(logFormat, uuid.New().String(), `task read`, parameterCapture, result, err)
	return result, err
}

func (middleware logMiddleware) Delete(ctx context.Context, id uint) (*Task, error) {
	var err error
	var result *Task
	var parameterCapture string

	parameterCapture = fmt.Sprintf(`%d`, id)
	result, err = middleware.next.Delete(ctx, id)

	middleware.logger.Printf(logFormat, uuid.New().String(), `task delete`, parameterCapture, result, err)
	return result, err
}

func (middleware logMiddleware) List(ctx context.Context, limit uint, offset uint) ([]Task, error) {
	var err error
	var result []Task
	var parameterCapture string

	parameterCapture = fmt.Sprintf(`{Limit: %d, Offset: %d}`, limit, offset)
	result, err = middleware.next.List(ctx, limit, offset)

	middleware.logger.Printf(logFormat, uuid.New().String(), `task list`, parameterCapture, result, err)
	return result, err
}

func (middleware logMiddleware) Shutdown() error {
	var err error

	err = middleware.next.Shutdown()

	middleware.logger.Printf(logFormat, uuid.New().String(), `task shutdown`, ``, ``, err)
	return err
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
