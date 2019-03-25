//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------
func dummyLogger() log.Logger {
	var buffer bytes.Buffer
	var writer *bufio.Writer
	var logger log.Logger

	writer = bufio.NewWriter(&buffer)

	logger.SetOutput(writer)

	return logger
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestMiddlewareLoggerNewService(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var logger Middleware
	var service Service

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	resetStore(test, store)
	logger = NewLogMiddleare(dummyLogger())

	//-- Action ----------
	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	//-- Post-conditions ----------
	assert.NotNil(test, service)
}

func TestMiddlewareLoggerCreate(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var logger Middleware
	var service Service
	var createErr error

	//-- Test Parameters ----------
	var name = `Test create`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	//-- Action ----------
	createErr = service.Create(ctx, model)

	//-- Post-conditions ----------
	assert.Nil(test, createErr)
	assert.NotEqual(test, uint(0), model.ID)
	assert.NotEqual(test, time.Time{}, model.CreatedAt)
	assert.Nil(test, model.UpdatedAt)
}

func TestMiddlewareLoggerCreateErr(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var logger Middleware
	var service Service
	var createErr error

	//-- Test Parameters ----------
	var name = `Test invalid create ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	ctx = context.Background()
	store = openStore(test)
	resetStore(test, store)
	logger = NewLogMiddleare(dummyLogger())
	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	//-- Action ----------
	createErr = service.Create(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, createErr)
	assert.Equal(test, uint(0), model.ID)
	assert.Equal(test, time.Time{}, model.CreatedAt)
	assert.Nil(test, model.UpdatedAt)
}

func TestMiddlewareLoggerUpdate(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, updatedModel *Task
	var store Store
	var logger Middleware
	var service Service
	var updatedErr error

	//-- Test Parameters ----------
	var name = `Test update`

	var updatedName = `Test updated name`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	if err := service.Create(ctx, model); err != nil {
		test.Fatalf(`unexpected error when creating record: %s`, err)
	}

	//-- Action ----------
	updatedModel = new(Task)
	*updatedModel = *model
	updatedModel.Name = updatedName

	updatedErr = service.Update(ctx, updatedModel)

	//-- Post-conditions ----------
	assert.Nil(test, updatedErr)
	assert.Equal(test, model.ID, updatedModel.ID)
	assert.NotEqual(test, model.Name, updatedModel.Name)
	assert.NotNil(test, updatedModel.UpdatedAt)
	assert.Nil(test, model.UpdatedAt)
}

func TestMiddlewareLoggerUpdateErr(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var logger Middleware
	var service Service
	var updatedErr error

	//-- Test Parameters ----------
	var name = `Test update`

	var updatedName = `Test updated name ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	if err := service.Create(ctx, model); err != nil {
		test.Fatalf(`unexpected error when creating record: %s`, err)
	}

	//-- Action ----------
	model.Name = updatedName

	updatedErr = service.Update(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, updatedErr)
	assert.Nil(test, model.UpdatedAt)
}

func TestMiddlewareLoggerRead(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, readModel *Task
	var store Store
	var logger Middleware
	var service Service
	var readErr error

	//-- Test Parameters ----------
	var name = `Test read`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	if err := service.Create(ctx, model); err != nil {
		test.Fatalf(`unexpected error when creating record: %s`, err)
	}

	//-- Action ----------
	readModel, readErr = service.Read(ctx, model.ID)

	//-- Post-conditions ----------
	assert.Nil(test, readErr)
	assert.True(test, model.compare(*readModel))
}

func TestMiddlewareLoggerReadErr(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, readModel *Task
	var store Store
	var logger Middleware
	var service Service
	var readErr error

	//-- Test Parameters ----------
	var name = `Test invalid read`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	//-- Action ----------
	readModel, readErr = service.Read(ctx, model.ID)

	//-- Post-conditions ----------
	assert.NotNil(test, readErr)
	assert.Nil(test, readModel)
}

func TestMiddlewareLoggerDelete(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, deleteModel *Task
	var store Store
	var logger Middleware
	var service Service
	var deleteErr error

	//-- Test Parameters ----------
	var name = `Test delete`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())
	service = NewService([]Middleware{logger}, store)

	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	if err := service.Create(ctx, model); err != nil {
		test.Fatalf(`unexpected error when creating record: %s`, err)
	}

	//-- Action ----------
	deleteModel, deleteErr = service.Delete(ctx, model.ID)

	//-- Post-conditions ----------
	assert.Nil(test, deleteErr)
	assert.True(test, model.compare(*deleteModel))
}

func TestMiddlewareLoggerDeleteErr(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, deleteModel *Task
	var store Store
	var logger Middleware
	var service Service
	var deleteErr error

	//-- Test Parameters ----------
	var name = `Test invalid delete`

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	model = newValidTask()
	model.Name = name

	//-- Action ----------
	deleteModel, deleteErr = service.Delete(ctx, model.ID)

	//-- Post-conditions ----------
	assert.NotNil(test, deleteErr)
	assert.Nil(test, deleteModel)
}

func TestMiddlewareLoggerList(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing valid list`
	var quantity = 10

	var limit uint = 25
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, quantity, len(listModels))
}

func TestMiddlewareLoggerListLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing limit list`
	var quantity = 10

	var limit uint = 5
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 5, len(listModels))
}

func TestMiddlewareLoggerListZeroLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing zero limit list`
	var quantity = 10

	var limit uint = 0
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 0, len(listModels))
}

func TestMiddlewareLoggerListOverLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing over limit list`
	var quantity = 10

	var limit uint = 200
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 10, len(listModels))
}

func TestMiddlewareLoggerListOffset(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing offset list`
	var quantity = 10

	var limit uint = 10
	var offset uint = 5

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 5, len(listModels))
}

func TestMiddlewareLoggerListOverOffset(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var logger Middleware
	var service Service
	var listErr error
	var listModels []Task

	//-- Test Parameters ----------
	var name = `Testing over offset list`
	var quantity = 10

	var limit uint = 10
	var offset uint = 10

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	resetStore(test, store)

	logger = NewLogMiddleare(dummyLogger())

	service = NewService([]Middleware{logger}, store)
	defer shutdownService(test, service)

	for i := 0; i < quantity; i++ {
		var task = newValidTask()
		task.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := service.Create(ctx, task); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listModels, listErr = service.List(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 0, len(listModels))
}
