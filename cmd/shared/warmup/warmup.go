//-- Package Declaration -----------------------------------------------------------------------------------------------
package warmup

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

//-- Constants ---------------------------------------------------------------------------------------------------------

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Exported Functions ------------------------------------------------------------------------------------------------
func IsScheduledWarmupEvent(event events.APIGatewayProxyRequest) bool {
	if len(event.Resource) < 1 {
		return true
	} else {
		return false
	}
}

func DefaultAPIGatewatResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `Ignoring warm up invocation`,
	}, nil
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
