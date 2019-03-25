//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Decorators --------------------------------------------------------------------------------------------------------

//-- Helpers -----------------------------------------------------------------------------------------------------------

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestMigrate(test *testing.T) {
	//-- Shared Variables ----------
	var output Response

	var response events.APIGatewayProxyResponse

	var eventErr error

	//-- Test Parameters ----------

	//-- Pre-conditions ----------

	//-- Action ----------
	response, eventErr = Handler()

	//-- Post-conditions ----------
	assert.Nil(test, eventErr)
	assert.Equal(test, http.StatusOK, response.StatusCode)

	if err := json.Unmarshal([]byte(response.Body), &output); err != nil {
		test.Fatalf(`unable to marshal response: %s`, err)
	} else {
		assert.True(test, output.Success)
	}
}
