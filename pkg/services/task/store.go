//-- Package Declaration -----------------------------------------------------------------------------------------------
package task

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//-- Constants ---------------------------------------------------------------------------------------------------------
var (
	queryMap = map[string]string{
		`insertTask`: `INSERT INTO tasks(name, details, resolved_at, created_at) VALUES($1, $2, $3, $4) RETURNING id`,
		`updateTask`: `UPDATE tasks SET name = $2, details = $3, resolved_at = $4, updated_at = $5 WHERE id = $1 RETURNING id`,
		`readTask`:   `SELECT * FROM tasks WHERE id = $1 LIMIT 1`,
		`deleteTask`: `DELETE FROM tasks WHERE id = $1 RETURNING *`,
		`listTasks`:  `SELECT * FROM tasks ORDER BY id LIMIT $1 OFFSET $2 ROWS`,
	}

	ErrIllAdvisedInsert = errors.New(`inserting a Task with non-zero ID in inadvisable; either pass a clean struct or do an update if this is an existing record`)
)

//-- Structs -----------------------------------------------------------------------------------------------------------
type postgresStore struct {
	database *sql.DB
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func NewPostgresStore() Store {
	return new(postgresStore)
}

func (store *postgresStore) Open(options string) error {
	//-- Open to database ----------
	{
		if db, err := sql.Open(`postgres`, options); err != nil {
			return err
		} else if err := db.Ping(); err != nil {
			return err
		} else {
			store.database = db
		}
	}

	//-- Return ----------
	return nil
}

func (store *postgresStore) Close() error {
	if store.database == nil {
		return errors.New(`database was not initialized or in a connected state`)
	} else if err := store.database.Close(); err != nil {
		return err
	}
	return nil
}

func (store *postgresStore) Prepare(option string, parameter string) error {
	switch option {
	case `up`:
		return store.up(parameter)
	case `down`:
		return store.down(parameter)
	case `drop`:
		return store.drop(parameter)
	default:
		return errors.New(`no valid option provided, unable to prepare`)
	}
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func (store *postgresStore) handleTransactionError(transaction *sql.Tx, original error) error {
	if err := transaction.Rollback(); err != nil {
		return errors.New(fmt.Sprintf(`an unrecoverable exception has occured rolling back the transaction (%s) - > (%s)`, original, err))
	}

	return original
}

func (store *postgresStore) up(migrationPath string) error {
	if driver, err := postgres.WithInstance(store.database, &postgres.Config{}); err != nil {
		return err
	} else if migration, err := migrate.NewWithDatabaseInstance(migrationPath, `postgres`, driver); err != nil {
		return err
	} else if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (store *postgresStore) down(migrationPath string) error {
	if driver, err := postgres.WithInstance(store.database, &postgres.Config{}); err != nil {
		return err
	} else if migration, err := migrate.NewWithDatabaseInstance(migrationPath, `postgres`, driver); err != nil {
		return err
	} else if err := migration.Down(); err != nil {
		return err
	}
	return nil
}

func (store *postgresStore) drop(migrationPath string) error {
	if driver, err := postgres.WithInstance(store.database, &postgres.Config{}); err != nil {
		return err
	} else if migration, err := migrate.NewWithDatabaseInstance(migrationPath, `postgres`, driver); err != nil {
		return err
	} else if err := migration.Drop(); err != nil {
		return err
	}
	return nil
}

func (store *postgresStore) version(migrationPath string) (uint, bool, error) {
	if driver, err := postgres.WithInstance(store.database, &postgres.Config{}); err != nil {
		return 0, false, err
	} else if migration, err := migrate.NewWithDatabaseInstance(migrationPath, `postgres`, driver); err != nil {
		return 0, false, err
	} else if ver, dirty, err := migration.Version(); err != nil && err != migrate.ErrNoChange {
		return 0, false, err
	} else {
		return ver, dirty, nil
	}
}

func (store *postgresStore) insert(ctx context.Context, task *Task) error {
	//-- Common variables ----------
	var id int
	var timestamp = time.Now().UTC()
	var query = queryMap[`insertTask`]

	//-- Parameter checking ----------
	if task.ID != 0 {
		return ErrIllAdvisedInsert
	}

	//-- Sanitize & validate ---------
	if err := task.sanitize(); err != nil {
		return err
	} else if err := task.validate(); err != nil {
		return err
	}

	//-- Insert Transaction ----------
	{
		if transaction, err := store.database.BeginTx(ctx, nil); err != nil {
			return err
		} else if err := transaction.QueryRow(query, task.Name, task.Details, task.ResolvedAt, timestamp).Scan(&id); err != nil {
			return store.handleTransactionError(transaction, err)
		} else if err := transaction.Commit(); err != nil {
			return err
		} else {
			task.ID = uint(id)
			task.CreatedAt = timestamp
			task.UpdatedAt = nil
			return nil
		}
	}
}

func (store *postgresStore) update(ctx context.Context, task *Task) error {
	//-- Common variables ----------
	var id int
	var timestamp = time.Now().UTC()
	var query = queryMap[`updateTask`]

	//-- Sanitize & validate ---------
	if err := task.sanitize(); err != nil {
		return err
	} else if err := task.validate(); err != nil {
		return err
	}

	//-- Insert Transaction ----------
	{
		if transaction, err := store.database.BeginTx(ctx, nil); err != nil {
			return err
		} else if err := transaction.QueryRow(query, task.ID, task.Name, task.Details, task.ResolvedAt, timestamp).Scan(&id); err != nil {
			return store.handleTransactionError(transaction, err)
		} else if err := transaction.Commit(); err != nil {
			return err
		} else {
			task.UpdatedAt = &timestamp
			return nil
		}
	}
}

func (store *postgresStore) read(ctx context.Context, id uint) (*Task, error) {
	//-- Common variables ----------
	var task = new(Task)
	var query = queryMap[`readTask`]

	//-- Insert Transaction ----------
	{
		if transaction, err := store.database.BeginTx(ctx, nil); err != nil {
			return nil, err
		} else if err := transaction.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Details, &task.ResolvedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, store.handleTransactionError(transaction, err)
		} else if err := transaction.Commit(); err != nil {
			return nil, err
		} else {
			return task, nil
		}
	}
}

func (store *postgresStore) delete(ctx context.Context, id uint) (*Task, error) {
	//-- Common variables ----------
	var task = new(Task)
	var query = queryMap[`deleteTask`]

	//-- Insert Transaction ----------
	{
		if transaction, err := store.database.BeginTx(ctx, nil); err != nil {
			return nil, err
		} else if err := transaction.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Details, &task.ResolvedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, store.handleTransactionError(transaction, err)
		} else if err := transaction.Commit(); err != nil {
			return nil, err
		} else {
			return task, nil
		}
	}
}

func (store *postgresStore) list(ctx context.Context, limit uint, offset uint) ([]Task, error) {
	//-- Common variables ----------
	var tasks = make([]Task, 0)
	var query = queryMap[`listTasks`]

	//-- Insert Transaction ----------
	{
		var transaction *sql.Tx
		if t, err := store.database.BeginTx(ctx, nil); err != nil {
			return nil, err
		} else {
			transaction = t
		}

		var results, err = transaction.Query(query, limit, offset)
		if err != nil {
			return tasks, err
		}

		var resultsScanError error
		for results.Next() {
			var task = new(Task)
			if err := results.Scan(&task.ID, &task.Name, &task.Details, &task.ResolvedAt, &task.CreatedAt, &task.UpdatedAt); err != nil {
				resultsScanError = err
				break
			}
			tasks = append(tasks, *task)
		}

		if err := results.Close(); err != nil {
			return nil, err
		} else if err := transaction.Commit(); err != nil {
			return nil, err
		} else if resultsScanError != nil {
			return nil, store.handleTransactionError(transaction, resultsScanError)
		}

		return tasks, nil
	}
}
