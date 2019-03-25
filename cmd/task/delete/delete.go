//-- Package Declaration -----------------------------------------------------------------------------------------------
package main

//-- Imports ----------------------------------------------------------------------------------------------------------
import (
	"context"
	logger2 "github.com/JustonDavies/go_serverless_api/cmd/shared/logger"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JustonDavies/go_serverless_api/cmd/shared/responses"
	"github.com/JustonDavies/go_serverless_api/cmd/shared/warmup"
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
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Details    *string    `json:"details,omitempty"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

//-- Event Handler -----------------------------------------------------------------------------------------------------
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//-- Ignore Warm-Ups ----------
	{
		if warmup.IsScheduledWarmupEvent(event) {
			log.Print(`Warmup event detected...ignoring...`)
			return warmup.DefaultAPIGatewatResponse()
		}
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
	var logger = logger2.NewLogger()

	var subjectID uint
	var service task.Service
	var subjectTask *task.Task

	var response *Response

	//-- Parse event ----------
	{
		if parsed, err := strconv.ParseUint(event.PathParameters[`id`], 10, 64); err != nil {
			return responses.APIGatewayProxyError(responses.BadPathParameterErr(err))
		} else {
			subjectID = uint(parsed)
		}
	}
	log.Printf(`Event parsed: %d seconds`, time.Now().Unix()-start)

	//-- Connect Service ----------
	{
		var middlewares []task.Middleware
		var store task.Store

		middlewares = append(middlewares, task.NewLogMiddleare(*logger))

		store = task.NewPostgresStore()
		if err := store.Open(os.Getenv(`DATABASE_CONNECTION_PARAMETERS`)); err != nil {
			return responses.APIGatewayProxyError(responses.InternalServerErr(err))
		}

		service = task.NewService(middlewares, store)

		defer func() {
			if err := service.Shutdown(); err != nil {
				log.Printf(`an unrecoverable error has occured while tring to shutdown the service: %s`, err)
			}
		}()
	}
	log.Printf(`Service started: %d seconds`, time.Now().Unix()-start)

	//-- Action ---------
	{
		if result, err := service.Delete(ctx, subjectID); err != nil {
			return responses.APIGatewayProxyError(responses.NotFound(err))
		} else {
			subjectTask = result
		}

		response = &Response{
			ID:         subjectTask.ID,
			Name:       subjectTask.Name,
			Details:    subjectTask.Details,
			ResolvedAt: subjectTask.ResolvedAt,
			CreatedAt:  subjectTask.CreatedAt,
			UpdatedAt:  subjectTask.UpdatedAt,
		}
	}
	log.Printf(`Action finished: %d seconds`, time.Now().Unix()-start)

	//-- Response ----------
	{
		if output, err := json.Marshal(response); err != nil {
			return responses.APIGatewayProxyError(responses.InternalServerErr(err))
		} else {
			log.Printf(`Completed: %d seconds(%d bytes)`, time.Now().Unix()-start, len(output))

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
