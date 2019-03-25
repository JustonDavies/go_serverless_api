//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------
func defaultConnectionString() string {
	if str := os.Getenv(`DATABASE_CONNECTION_PARAMETERS`); len(str) > 0 {
		return str
	}
	return `postgres://task_service_user:task_service_password@127.0.0.1/task_service_database?sslmode=disable&timezone=UTC`
}

func defaultMigrationPath() string {
	if str := os.Getenv(`DATABASE_MIGRATION_PATH`); len(str) > 0 {
		return str
	}
	return `file://migrations`
}

func openStore(test *testing.T) Store {
	var store = NewPostgresStore()

	if store == nil {
		test.Fatal(`an unexpected error occurred while initializing the store`)
	} else if err := store.Open(defaultConnectionString()); err != nil {
		test.Fatalf(`an unexpected error occured while connecting to the database: %s`, err.Error())
	} else if store.(*postgresStore).database == nil {
		test.Fatal(`an unrecoverable condition occurred where the database pointer was empty`)
	}

	return store
}

func closeStore(test *testing.T, store Store) {
	if err := store.(*postgresStore).Close(); err != nil {
		test.Fatalf(`an unexpected error occurred while closing the database: %s`, err.Error())
	}
}

func zeroStore(test *testing.T, store Store) {
	if err := store.(*postgresStore).drop(defaultMigrationPath()); err != nil && err != database.ErrLocked {
		test.Fatalf(`an unexpected error occurred while dropping the database: %s`, err.Error())
	} else if err == database.ErrLocked {
		//TODO: Investigate in more depth and log a ticket with the right library provider if required, I think it has to do with not being able to operate on sequences
		//test.Logf(`A mostly benin error occurred where the database driver reported an inability to obtain a lock but the operation was successful: %s`, err.Error())
	}
}

func resetStore(test *testing.T, store Store) {
	zeroStore(test, store)

	if err := store.(*postgresStore).up(defaultMigrationPath()); err != nil {
		test.Fatalf(`an unexpected error occurred while migrating the database: %s`, err.Error())
	}
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestStoreNewStore(test *testing.T) {
	//-- Shared Variables ----------
	var store Store

	//-- Test Parameters ----------

	//-- Pre-conditions ----------

	//-- Action ----------
	store = NewPostgresStore()

	//-- Post-conditions ----------
	assert.NotNil(test, store)
	assert.Nil(test, store.(*postgresStore).database)
}

func TestStoreConnect(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var connectErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = NewPostgresStore()

	//-- Action ----------
	connectErr = store.Open(defaultConnectionString())
	defer closeStore(test, store)

	//-- Post-conditions ----------
	assert.NotNil(test, store)
	assert.NotNil(test, store.(*postgresStore).database)
	assert.Nil(test, connectErr)

}

func TestStoreConnectBadConnectionParameters(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var connectErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = NewPostgresStore()

	//-- Action ----------
	connectErr = store.Open(`postgres://bad_user:bad_password@256.256.256.256/bad_database`)

	//-- Post-conditions ----------
	assert.NotNil(test, store)
	assert.Nil(test, store.(*postgresStore).database)
	assert.NotNil(test, connectErr)
}

func TestStoreDisconnect(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var disconnectErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)

	//-- Action ----------
	disconnectErr = store.Close()

	//-- Post-conditions ----------
	assert.Nil(test, disconnectErr)
}

func TestStoreDisconnectNotConnected(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var disconnectErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = NewPostgresStore()

	//-- Action ----------
	disconnectErr = store.Close()

	//-- Post-conditions ----------
	assert.NotNil(test, disconnectErr)
}

func TestStorePrepareUp(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	zeroStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`up`, defaultMigrationPath())

	//-- Post-conditions ----------
	if prepareError != nil && prepareError != migrate.ErrNoChange {
		test.Fatal(`an unexpected error occurred while migrating (up) the database: `, prepareError.Error())
	}
}

func TestStorePrepareUpBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	zeroStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`up`, `file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, prepareError)
}

func TestStorePrepareDown(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`down`, defaultMigrationPath())

	//-- Post-conditions ----------
	if prepareError != nil && prepareError != migrate.ErrNoChange {
		test.Fatal(`an unexpected error occurred while migrating (down) the database: `, prepareError.Error())
	}
}

func TestStorePrepareDownBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`down`, `file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, prepareError)
}

func TestStorePrepareDrop(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`drop`, defaultMigrationPath())

	//-- Post-conditions ----------
	if prepareError != nil && prepareError != database.ErrLocked {
		test.Fatal(`an unexpected error occurred while dropping the database: `, prepareError.Error())
	}
}

func TestStorePrepareDropBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`drop`, `file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, prepareError)
}

func TestStorePrepareInvalid(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var prepareError error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)

	//-- Action ----------
	prepareError = store.Prepare(`invalid_option`, defaultMigrationPath())

	//-- Post-conditions ----------
	assert.NotNil(test, prepareError)
}

func TestStoreUp(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	zeroStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).up(defaultMigrationPath())

	//-- Post-conditions ----------
	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		test.Fatal(`an unexpected error occurred while migrating (up) the database: `, migrationErr.Error())
	}
}

func TestStoreUpBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	zeroStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).up(`file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, migrationErr)
}

func TestStoreDown(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).down(defaultMigrationPath())

	//-- Post-conditions ----------
	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		test.Fatal(`an unexpected error occurred while migrating (down) the database: `, migrationErr.Error())
	}
}

func TestStoreDownBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).down(`file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, migrationErr)
}

func TestStoreDrop(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).drop(defaultMigrationPath())

	//-- Post-conditions ----------
	if migrationErr != nil && migrationErr != database.ErrLocked {
		test.Fatal(`an unexpected error occurred while dropping the database: `, migrationErr.Error())
	}
}

func TestStoreDropBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	migrationErr = store.(*postgresStore).drop(`file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, migrationErr)
}

func TestStoreVersion(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	_, _, migrationErr = store.(*postgresStore).version(defaultMigrationPath())

	//-- Post-conditions ----------
	if migrationErr != nil && migrationErr != database.ErrLocked {
		test.Fatal(`an unexpected error occurred while reading the version of the database: `, migrationErr.Error())
	}
}

func TestStoreVersionBadPath(test *testing.T) {
	//-- Shared Variables ----------
	var store Store
	var migrationErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	_, _, migrationErr = store.(*postgresStore).version(`file://not_real`)

	//-- Post-conditions ----------
	assert.NotNil(test, migrationErr)
}

func TestStoreInsert(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var insertErr error

	//-- Test Parameters ----------
	var name = `Testing valid model store insert`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	insertErr = store.(*postgresStore).insert(ctx, model)

	//-- Post-conditions ----------
	assert.Nil(test, insertErr)
	assert.NotEqual(test, uint(0), model.ID)
	assert.NotEqual(test, time.Time{}, model.CreatedAt.Unix())
	assert.Nil(test, model.UpdatedAt)
}

func TestStoreInsertInvalid(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var insertErr error

	//-- Test Parameters ----------
	var name = `Testing invalid model store insert ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	insertErr = store.(*postgresStore).insert(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, insertErr)
	assert.Equal(test, uint(0), model.ID)
	assert.Equal(test, time.Time{}, model.CreatedAt)
}

func TestStoreInsertIllAdvised(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var insertErr error

	//-- Test Parameters ----------
	var id uint = 1
	var name = `Testing zero-value ID model store insert`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.ID = id
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	insertErr = store.(*postgresStore).insert(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, insertErr)
	assert.Equal(test, time.Time{}, model.CreatedAt)
}

func TestStoreInsertDuplicate(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, duplicateTask *Task
	var store Store
	var insertErr error

	//-- Test Parameters ----------
	var name = `Testing inadvisable model store insert`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	if err := store.(*postgresStore).insert(ctx, model); err != nil {
		test.Fatalf(`unexpected error when inserting record: %s`, err)
	}

	//-- Action ----------
	duplicateTask = new(Task)
	*duplicateTask = *model

	insertErr = store.(*postgresStore).insert(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, insertErr)
	assert.True(test, model.compare(*duplicateTask))
}

func TestStoreUpdate(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, udpatedTask *Task
	var store Store
	var updateErr error

	//-- Test Parameters ----------
	var name = `Testing valid update`
	var updatedName = `Testing valid name update`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	if err := store.(*postgresStore).insert(ctx, model); err != nil {
		test.Fatalf(`unexpected error when inserting record: %s`, err)
	}

	//-- Action ----------
	udpatedTask = new(Task)
	*udpatedTask = *model
	udpatedTask.Name = updatedName

	updateErr = store.(*postgresStore).update(ctx, udpatedTask)

	//-- Post-conditions ----------
	assert.Nil(test, updateErr)
	assert.Equal(test, model.ID, udpatedTask.ID)
	assert.NotEqual(test, model.Name, udpatedTask.Name)
	assert.NotNil(test, udpatedTask.UpdatedAt)
	assert.Nil(test, model.UpdatedAt)
}

func TestStoreUpdateInvalid(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var updateErr error

	//-- Test Parameters ----------
	var name = `Testing valid update`
	var updatedName = `Testing valid name update ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	if err := store.(*postgresStore).insert(ctx, model); err != nil {
		test.Fatalf(`unexpected error when inserting record: %s`, err)
	}

	//-- Action ----------
	model.Name = updatedName

	updateErr = store.(*postgresStore).update(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, updateErr)
	assert.Nil(test, model.UpdatedAt)
	assert.Nil(test, model.UpdatedAt)
}

func TestStoreUpdateNotFound(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model *Task
	var store Store
	var updateErr error

	//-- Test Parameters ----------
	var name = `Testing not-found update`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	updateErr = store.(*postgresStore).update(ctx, model)

	//-- Post-conditions ----------
	assert.NotNil(test, updateErr)
	assert.Nil(test, model.UpdatedAt)
}

func TestStoreRead(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, readTask *Task
	var store Store
	var readErr error

	//-- Test Parameters ----------
	var name = `Testing valid read`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	if err := store.(*postgresStore).insert(ctx, model); err != nil {
		test.Fatalf(`unexpected error when inserting record: %s`, err)
	}

	//-- Action ----------
	readTask, readErr = store.(*postgresStore).read(ctx, model.ID)

	//-- Post-conditions ----------
	assert.Nil(test, readErr)
	assert.True(test, model.compare(*readTask))
}

func TestStoreReadNotFound(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, readTask *Task
	var store Store
	var readErr error

	//-- Test Parameters ----------
	var name = `Testing not-found read`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	readTask, readErr = store.(*postgresStore).read(ctx, uint(0))

	//-- Post-conditions ----------
	assert.NotNil(test, readErr)
	assert.Nil(test, readTask)
}

func TestStoreDelete(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, deleteTask *Task
	var store Store
	var deleteErr error

	//-- Test Parameters ----------
	var name = `Testing valid delete`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	if err := store.(*postgresStore).insert(ctx, model); err != nil {
		test.Fatalf(`unexpected error when inserting record: %s`, err)
	}

	//-- Action ----------
	deleteTask, deleteErr = store.(*postgresStore).delete(ctx, model.ID)

	//-- Post-conditions ----------
	assert.Nil(test, deleteErr)
	assert.True(test, model.compare(*deleteTask))
}

func TestStoreDeleteNotFound(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var model, deleteTask *Task
	var store Store
	var deleteErr error

	//-- Test Parameters ----------
	var name = `Testing no-found delete`

	//-- Pre-conditions ----------
	ctx = context.Background()

	model = newValidTask()
	model.Name = name

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	//-- Action ----------
	deleteTask, deleteErr = store.(*postgresStore).delete(ctx, model.ID)

	//-- Post-conditions ----------
	assert.NotNil(test, deleteErr)
	assert.Nil(test, deleteTask)
}

func TestStoreList(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing valid list`
	var quantity = 10

	var limit uint = 25
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, quantity, len(listTasks))
}

func TestStoreListLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing limit list`
	var quantity = 10

	var limit uint = 5
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 5, len(listTasks))
}

func TestStoreListZeroLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing zero limit list`
	var quantity = 10

	var limit uint = 0
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 0, len(listTasks))
}

func TestStoreListOverLimit(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing over limit list`
	var quantity = 10

	var limit uint = 200
	var offset uint = 0

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 10, len(listTasks))
}

func TestStoreListOffset(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing offset list`
	var quantity = 10

	var limit uint = 10
	var offset uint = 5

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 5, len(listTasks))
}

func TestStoreListOverOffset(test *testing.T) {
	//-- Shared Variables ----------
	var ctx context.Context
	var store Store
	var listTasks []Task
	var listErr error

	//-- Test Parameters ----------
	var name = `Testing over offset list`
	var quantity = 10

	var limit uint = 10
	var offset uint = 10

	//-- Pre-conditions ----------
	ctx = context.Background()

	store = openStore(test)
	defer closeStore(test, store)
	resetStore(test, store)

	for i := 0; i < quantity; i++ {
		var model = newValidTask()
		model.Name = fmt.Sprintf(`%s %d`, name, i)

		if err := store.(*postgresStore).insert(ctx, model); err != nil {
			test.Fatalf(`unexpected error when inserting record: %s`, err)
		}
	}

	//-- Action ----------
	listTasks, listErr = store.(*postgresStore).list(ctx, limit, offset)

	//-- Post-conditions ----------
	assert.Nil(test, listErr)
	assert.Equal(test, 0, len(listTasks))
}
