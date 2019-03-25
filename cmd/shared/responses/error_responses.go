//-- Package Declaration -----------------------------------------------------------------------------------------------
package responses

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/jsonapi"
)

//-- Constants ---------------------------------------------------------------------------------------------------------

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Exported Functions ------------------------------------------------------------------------------------------------
func BadPathParameterErr(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusBadRequest),
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: `Path parameter was invalid or unable to parse/process`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func NotAcceptableErr(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusNotAcceptable),
		Title:  http.StatusText(http.StatusNotAcceptable),
		Detail: `Request parameter was not acceptabe`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func UnprocessableEntryErr(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusUnprocessableEntity),
		Title:  http.StatusText(http.StatusUnprocessableEntity),
		Detail: `Request parameter was invalid or unable to parse/process`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func MalformedRequestErr(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusBadRequest),
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: `Request parameter was invalid or unable to parse/process`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func InternalServerErr(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusInternalServerError),
		Title:  http.StatusText(http.StatusInternalServerError),
		Detail: `an unrecoverable error has occurred`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func NotFound(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusNotFound),
		Title:  http.StatusText(http.StatusNotFound),
		Detail: `Record not found`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func Unauthorized(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf(`%d`, http.StatusUnauthorized),
		Title:  http.StatusText(http.StatusUnauthorized),
		Detail: `Unauthorized`,
		Meta:   &map[string]interface{}{`error`: err.Error()},
	}
}

func APIGatewayProxyError(err *jsonapi.ErrorObject) (events.APIGatewayProxyResponse, error) {
	var errs []*jsonapi.ErrorObject
	errs = append(errs, err)

	return APIGatewayProxyErrors(errs)
}

func APIGatewayProxyErrors(errs []*jsonapi.ErrorObject) (events.APIGatewayProxyResponse, error) {
	//-- Response ----------
	var writer = new(bytes.Buffer)

	if err := jsonapi.MarshalErrors(writer, errs); err != nil {
		log.Println(`Error Serializing Errors: `, err)
	}

	var status int
	{
		if parsed, err := strconv.Atoi(errs[0].Status); err != nil {
			status = http.StatusInternalServerError
		} else {
			status = parsed
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       writer.String(),
	}, nil
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
