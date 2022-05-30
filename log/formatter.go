package log

import (
	"encoding/json"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
)

type logMsg struct {
	ClientIP      string
	TimeStamp     string
	Method        string
	Path          string
	Status        int
	Latency       string
	ErrorMessage  string
	CorrelationID string
}

const (
	CORRELATION_HEADER = "Correlation-Id"
)

func Formatter(param gin.LogFormatterParams) string {

	correlationID := ""
	ch := param.Request.Header[CORRELATION_HEADER]
	if len(ch) > 0 {
		correlationID = ch[0]
	}

	m := logMsg{
		ClientIP:      param.ClientIP,
		TimeStamp:     param.TimeStamp.Format(time.RFC1123),
		Method:        param.Method,
		Path:          param.Path,
		Status:        param.StatusCode,
		Latency:       param.Latency.String(),
		ErrorMessage:  param.ErrorMessage,
		CorrelationID: correlationID,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "error formatting log:" + fmt.Sprintf("%+v\n", m) + "\n"
	}
	return string(b) + "\n"

}
