package middleware

import (
	"bytes"
	"io"
	"nunu-http-layout/pkg/helper/md5"
	"nunu-http-layout/pkg/helper/uuid"
	"nunu-http-layout/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLog(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// The configuration is initialized once per request
		trace := md5.Md5(uuid.GenUUID())
		logger.NewContext(ctx, zap.String("trace", trace))
		logger.NewContext(ctx, zap.String("request_method", ctx.Request.Method))
		logger.NewContext(ctx, zap.Any("request_headers", ctx.Request.Header))
		logger.NewContext(ctx, zap.String("request_url", ctx.Request.URL.String()))
		if ctx.Request.Body != nil {
			bodyBytes, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
			logger.NewContext(ctx, zap.String("request_params", string(bodyBytes)))
		}
		logger.WithContext(ctx).Info("Request")
		ctx.Next()
	}
}
func ResponseLog(logger *log.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime).String()
		ctx.Header("X-Response-Time", duration)
		logger.WithContext(ctx).Info("Response", zap.Any("response_body", blw.body.String()), zap.Any("time", duration))
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
