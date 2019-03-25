//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports ----------------------------------------------------------------------------------------------------------
import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JustonDavies/go_serverless_api/cmd/shared/responses"
	"github.com/JustonDavies/go_serverless_api/pkg/services/task"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/json-iterator/go"
)

//-- Constants ---------------------------------------------------------------------------------------------------------
var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

//-- Structs -----------------------------------------------------------------------------------------------------------
type Response struct {
	Success bool `json:"success"`
}

//-- Event Handler -----------------------------------------------------------------------------------------------------
func Handler() (events.APIGatewayProxyResponse, error) {
	//-- Ignore Warm-Ups ----------
	{
		//Not configured for periodic warming
	}

	//-- Authenticate ----------
	{
		//No authentication required / implemented at this time
	}

	//-- Authorize ----------
	{
		//No authorization required / implemented at this time
	}

	//-- Shared variables ----------
	var start = time.Now().Unix()
	var logger = log.New(os.Stdout, `task_service: `, 3)

	var store task.Store

	var response *Response

	//-- Parse event ----------
	{

	}
	logger.Printf(`Event parsed: %d seconds`, time.Now().Unix()-start)

	//-- Connect Service ----------
	{
		store = task.NewPostgresStore()
		if err := store.Open(os.Getenv(`DATABASE_CONNECTION_PARAMETERS`)); err != nil {
			logger.Println(`Get fucked `, err)
			return responses.APIGatewayProxyError(responses.InternalServerErr(err))
		}

		defer func() {
			if err := store.Close(); err != nil {
				logger.Printf(`an unrecoverable error has occured while tring to close the store: %s`, err)
			}
		}()
	}
	logger.Printf(`Service started: %d seconds`, time.Now().Unix()-start)

	//-- Action ---------
	{
		if err := store.Prepare(`up`, `file://pkg/services/task/migrations`); err != nil {
			return responses.APIGatewayProxyError(responses.InternalServerErr(err))
		}

		response = &Response{
			Success: true,
		}
	}
	logger.Printf(`Action finished: %d seconds`, time.Now().Unix()-start)

	//-- Response ----------
	{
		if output, err := json.Marshal(response); err != nil {
			return responses.APIGatewayProxyError(responses.InternalServerErr(err))
		} else {
			logger.Printf(`Completed: %d seconds(%d bytes)`, time.Now().Unix()-start, len(output))

			return events.APIGatewayProxyResponse{
				Body:       string(output),
				StatusCode: http.StatusOK,
			}, nil
		}
	}

}

//-- Main --------------------------------------------------------------------------------------------------------------
func main() {
	lambda.Start(Handler)
}
