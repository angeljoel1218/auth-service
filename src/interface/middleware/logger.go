package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const CtxtKey string = "reqID"

type LoggerConfig struct {
	SkipPaths []string
	Logger    *zerolog.Logger
}

type loggerWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w loggerWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(conf LoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqID := uuid.New()
		c.Set(CtxtKey, reqID)

		var skip map[string]struct{}

		notlogged := conf.SkipPaths
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if length := len(notlogged); length > 0 {
			skip = make(map[string]struct{}, length)

			for _, path := range notlogged {
				skip[path] = struct{}{}
			}
		}

		if _, ok := skip[path]; !ok {
			var bodyBytes []byte

			if raw != "" {
				path = path + "?" + raw
			}

			if c.Request.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			lw := &loggerWriter{
				body:           bytes.NewBufferString(""),
				ResponseWriter: c.Writer,
			}

			c.Writer = lw

			c.Next()

			msg := c.Errors.ByType(gin.ErrorTypePrivate).String()

			conf.Logger.Info().
				Str("reqID", reqID.String()).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Str("method", c.Request.Method).
				Str("path", path).
				Int("status", c.Writer.Status()).
				Dur("latency", time.Now().Sub(start)).
				Str("request", string(bodyBytes)).
				Str("response", lw.body.String()).
				Msg(msg)
		} else {
			c.Next()
		}
	}
}
