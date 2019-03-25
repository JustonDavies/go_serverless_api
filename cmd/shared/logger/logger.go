//-- Package Declaration -----------------------------------------------------------------------------------------------
package logger

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"bufio"
	"bytes"
	"log"
	"os"
)

//-- Constants ---------------------------------------------------------------------------------------------------------
var (
	environment = os.Getenv(`ENVIRONMENT`)
)

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Exported Functions ------------------------------------------------------------------------------------------------
func NewLogger() *log.Logger {
	if environment == `test` {
		var buffer bytes.Buffer
		var writer *bufio.Writer
		var logger log.Logger

		writer = bufio.NewWriter(&buffer)

		logger.SetOutput(writer)

		return &logger
	}

	return log.New(os.Stdout, `task_service: `, 3)
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
