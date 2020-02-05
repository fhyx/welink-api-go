package welink

import (
	"github.com/fhyx/welink-api-go/log"
)

func logger() log.Logger {
	return log.GetLogger()
}
