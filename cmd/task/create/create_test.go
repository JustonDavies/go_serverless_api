//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestCreateTask(test *testing.T) {
	//-- Shared Variables ----------
	var input Request
	var output Response

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API create task`
	var details = `Testing creating a task through the event handler API`
	var resolvedAt = time.Now()

	//-- Pre-conditions ----------
	ctx = context.Background()

	input = Request{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
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
		assert.NotEqual(test, uint(0), output.ID)
		assert.NotEqual(test, time.Time{}, output.CreatedAt)
		assert.Nil(test, output.UpdatedAt)
	}
}

func TestCreateTaskInvalid(test *testing.T) {
	//-- Shared Variables ----------
	var input Request

	var request events.APIGatewayProxyRequest
	var response events.APIGatewayProxyResponse

	var eventErr error

	var ctx context.Context

	//-- Test Parameters ----------
	var name = `Test API create invalid task ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~`
	var details = `Testing creating am invalid task through the event handler API`
	var resolvedAt = time.Now()

	//-- Pre-conditions ----------
	ctx = context.Background()

	input = Request{
		Name:       name,
		Details:    &details,
		ResolvedAt: &resolvedAt,
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
	assert.Equal(test, http.StatusUnprocessableEntity, response.StatusCode)
}
