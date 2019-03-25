//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/JustonDavies/go_serverless_api/pkg/services/task"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------
func insertTask(test *testing.T, input *task.Task) {
	var ctx context.Context
	var store task.Store
	var service task.Service

	ctx = context.Background()

	store = task.NewPostgresStore()
	if err := store.Open(os.Getenv(`DATABASE_CONNECTION_PARAMETERS`)); err != nil {
		test.Fatalf(`an unexpected error occured while opening the database: %s`, err)
	}

	service = task.NewService(nil, store)
	defer func() {
		if err := service.Shutdown(); err != nil {
			test.Fatalf(`an unrecoverable error has occured while tring to shutdown the service: %s`, err)
		}
	}()

	if err := service.Create(ctx, input); err != nil {
		test.Fatalf(`an unexpected error occured while opening the database: %s`, err)
	}
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestDeleteTask(test *testing.T) {
	//-- Shared Variables ----------
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	var subject task.Task

	//-- Test Parameters ----------
	var name = `Test API delete task`
	var details = `Testing deleting a task through the event handler API`
	var resolvedAt = time.Now()

	//-- Pre-conditions ----------
	subject = task.Task{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
	}
	insertTask(test, &subject)

	ctx = context.Background()

	request = events.APIGatewayProxyRequest{PathParameters: map[string]string{`id`: fmt.Sprintf(`%d`, subject.ID)}, Resource: `fake test resource`}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, subject.ID, output.ID)
		assert.Equal(test, subject.Name, output.Name)
		assert.Equal(test, *subject.Details, *output.Details)
		assert.Equal(test, subject.ResolvedAt.Unix(), output.ResolvedAt.Unix())
		assert.Equal(test, subject.CreatedAt.Unix(), output.CreatedAt.Unix())
		assert.Equal(test, subject.UpdatedAt, output.UpdatedAt) //It hasn't been updated yet so these should be nil
	}
}

func TestDeleteTaskNotFound(test *testing.T) {
	//-- Shared Variables ----------
	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------

	//-- Pre-conditions ----------
	ctx = context.Background()

	request = events.APIGatewayProxyRequest{PathParameters: map[string]string{`id`: fmt.Sprintf(`%d`, 0)}, Resource: `fake test resource`}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusNotFound, response.StatusCode)
}
