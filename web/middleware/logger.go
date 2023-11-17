package middleware

import (
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//start := time.Now()
		//requestBody, err := io.ReadAll(ctx.Request.Body)
		//if err != nil {
		//	response.Response(ctx, code.NewCode(code.RequestBodyInvalid))
		//	return
		//}
		//ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		//
		//raw, err := url.QueryUnescape(ctx.Request.URL.RawQuery)
		//if err != nil {
		//	response.Response(ctx, code.NewCode(code.RequestQueryInvalid))
		//	return
		//}
		//ctx.Next()
		//end := time.Now()
	}
}
