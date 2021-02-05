package http

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func AuthenticateMiddleware(authAddr string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := req.Header{
			"Accept":        "application/json",
			"Authorization": findToken(c),
		}

		type Request struct {
			Resource string `json:"resource"`
		}

		body := Request{
			Resource: c.Request.URL.Path,
		}

		r, err := req.Post("http://"+authAddr+"/auth/authenticate", header, req.BodyJSON(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())

			return
		}

		var response authenticateResponse
		if err = r.ToJSON(&response); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)

			return
		}

		if !response.IsAuthenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Permission Denied")

			return
		}

		c.Next()
	}
}

func findToken(c *gin.Context) string {
	if token := c.Request.Header.Get("Authorization"); token != "" {
		return token
	}

	v := c.Request.URL.Query()
	if token, exist := v["token"]; exist {
		return token[0]
	}

	return ""
}

func TracerMiddleware(spanName string, tag opentracing.Tag) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := opentracing.GlobalTracer().StartSpan(
			spanName,
			tag,
		)
		defer span.Finish()

		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.HTTPMethod.Set(span, c.Request.Method)

		opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header))

		c.Next()
	}
}
