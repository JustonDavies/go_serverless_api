//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"fmt"
	"github.com/JustonDavies/go_serverless_api/pkg/services/task"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
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

func deleteTasks(test *testing.T) {
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

	if tasks, err := service.List(ctx, 1000000, 0); err != nil {
		test.Fatalf(`an unexpected error occured while fetching all tasks in the database: %s`, err)
	} else {
		for _, item := range tasks {
			if _, err := service.Delete(ctx, item.ID); err != nil {
				test.Fatalf(`an unexpected error occured while deleting all tasks the database: %s`, err)
			}
		}
	}
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestIndexTask(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API list task`
	var quantity = 10

	var limit uint = 25
	var offset uint = 0

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, quantity, len(output.Tasks))
	}
}

func TestIndexTaskLimit(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API limited list task`
	var quantity = 10

	var limit uint = 5
	var offset uint = 0

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, int(limit), len(output.Tasks))
	}
}

func TestIndexTaskZeroLimit(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API zero-limited list task`
	var quantity = 10

	var limit uint = 0
	var offset uint = 0

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, int(limit), len(output.Tasks))
	}
}

func TestIndexTaskOverLimit(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API zero-limited list task`
	var quantity = 10

	var limit uint = 50
	var offset uint = 0

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, int(quantity), len(output.Tasks))
	}
}

func TestIndexTaskOffset(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API zero-limited list task`
	var quantity = 10

	var limit uint = 10
	var offset uint = 5

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, 5, len(output.Tasks))
	}
}

func TestIndexTaskOverOffset(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API zero-limited list task`
	var quantity = 10

	var limit uint = 10
	var offset uint = 10

	//-- Pre-conditions ----------
	deleteTasks(test)
	for i := 0; i < quantity; i++ {
		var item = task.Task{}
		item.Name = fmt.Sprintf(`%s %d`, name, i)

		insertTask(test, &item)
	}

	ctx = context.Background()

	input = Request{
		Limit:  limit,
		Offset: offset,
	}

	if result, err := json.Marshal(input); err != nil {
		test.Fatalf(`unable to marshal request: %s`, err)
	} else {
		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
	}

	//-- Action ----------
	response, eventErr = Handler(ctx, request)

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.Equal(test, 0, len(output.Tasks))
	}
}

//func TestReadTaskNotFound(test *testing.T) {
//	//-- Shared Variables ----------
//	var input Request
//
//	var request events.APIGatewayProxyRequest
//	var response events.APIGatewayProxyResponse
//
//	var eventErr error
//
//	var ctx context.Context
//
//	//-- Test Parameters ----------
//
//	//-- Pre-conditions ----------
//	ctx = context.Background()
//
//	input = Request{
//		ID:  0,
//	}
//
//	if result, err := json.Marshal(input); err != nil {
//		test.Fatalf(`unable to marshal request: %s`, err)
//	} else {
//		request = events.APIGatewayProxyRequest{Body: string(result), Resource: `fake test resource`}
//	}
//
//	//-- Action ----------
//	response, eventErr = Handler(ctx, request)
//
//	//-- Post-conditions ----------
//	assert.Nil(test, eventErr)
//	assert.Equal(test, http.StatusNotFound, response.StatusCode)
//}
