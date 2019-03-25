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
func TestUpdateTask(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	var subject task.Task

	//-- Test Parameters ----------
	var name = `Test API update task`
	var details = `Testing updating a task through the event handler API`
	var resolvedAt = time.Now()

	var updatedName = `Test API updated task name`

	//-- Pre-conditions ----------
	subject = task.Task{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
	}
	insertTask(test, &subject)

	ctx = context.Background()

	input = Request{
		ID:         subject.ID,
		Name:       updatedName,
		Details:    subject.Details,
		ResolvedAt: subject.ResolvedAt,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{PathParameters: map[string]string{`id`: fmt.Sprintf(`%d`, subject.ID)}, Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, subject.ID, output.ID)
		assert.Equal(test, updatedName, output.Name)
		assert.NotNil(test, output.UpdatedAt)
	}
}

func TestUpdateTaskInvalid(test *testing.T) {
	//-- Shared Variables ----------
	var input Request

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	var subject task.Task

	//-- Test Parameters ----------
	var name = `Test API update task`
	var details = `Testing updating a task through the event handler API`
	var resolvedAt = time.Now()

	var updatedName = `Test API invalid updated task name ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`

	//-- Pre-conditions ----------
	subject = task.Task{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
	}
	insertTask(test, &subject)

	ctx = context.Background()

	input = Request{
		ID:         subject.ID,
		Name:       updatedName,
		Details:    subject.Details,
		ResolvedAt: subject.ResolvedAt,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{PathParameters: map[string]string{`id`: fmt.Sprintf(`%d`, subject.ID)}, Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusUnprocessableEntity, response.StatusCode)
}

func TestUpdateTaskMismatchID(test *testing.T) {
	//-- Shared Variables ----------
	var input Request

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	var subject task.Task

	//-- Test Parameters ----------
	var name = `Test API update task`
	var details = `Testing updating a task through the event handler API`
	var resolvedAt = time.Now()

	var updatedName = `Test API updated task name`

	//-- Pre-conditions ----------
	subject = task.Task{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
	}
	insertTask(test, &subject)

	ctx = context.Background()

	input = Request{
		ID:         subject.ID,
		Name:       updatedName,
		Details:    subject.Details,
		ResolvedAt: subject.ResolvedAt,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{PathParameters: map[string]string{`id`: fmt.Sprintf(`%d`, subject.ID+1)}, Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusNotAcceptable, response.StatusCode)
}
