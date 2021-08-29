package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"moul.io/http2curl"
)

type LoggingRoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
}

func NewLoggingRoundTripper(
	next http.RoundTripper,
	w io.Writer,
) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		next:   next,
		logger: w,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(
	req *http.Request,
) (resp *http.Response, err error) {
	defer func(begin time.Time) {
		var msg string

		if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {

			msg = fmt.Sprintf(
				"method=%s host=%s path=%s status_code=%d took=%s\n",
				req.Method,
				req.URL.Host,
				req.URL.Path,
				resp.StatusCode,
				time.Since(begin),
			)
			if err != nil {
				fmt.Fprintf(rt.logger, "%s : %+v\n", msg, err)
			} else {
				fmt.Fprintf(rt.logger, "%s\n", msg)
			}
			command, _ := http2curl.GetCurlCommand(req)
			fmt.Println(command)
		}
	}(time.Now())

	return rt.next.RoundTrip(req)
}
