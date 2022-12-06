package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository/cache"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

var (
	RateLimitRepo cache.RateLimitRepo
	apiKeyRepo    repository.IOpenApiKeyRepo = new(cache.OpenApiKeyRepo)
)

const (
	signStrFmt                = "X-Timestamp: %d\nURI: %s\nBody: %s"
	defaultRateLimitFrequency = 100
	defaultRateLimitCycleTime = 60
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma, General, X-Timestamp, X-Signature, X-Api-Key")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Set("content-type", "application/json")

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

func SignatureVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		{
			//从header获取参数
			apiKey := c.Request.Header.Get(constant.HeaderApiKey)
			ok, apiSecret := verifyApiKey(apiKey)
			if !ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, response.FailBadRequest(fmt.Sprintf("Invalid %s", constant.HeaderApiKey)))
				return
			}

			timestamp := c.Request.Header.Get(constant.HeaderTimestamp)
			timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil || !checkTimeliness(timestampInt) {
				c.AbortWithStatusJSON(http.StatusBadRequest, response.FailBadRequest(fmt.Sprintf("Timeliness error. Please check the request header %s", constant.HeaderTimestamp)))
				return
			}

			signature := c.Request.Header.Get(constant.HeaderSignature)
			uri := c.Request.RequestURI
			var body string
			if c.Request.Method == http.MethodPost {
				var bz json.RawMessage
				if err = c.ShouldBindBodyWith(&bz, binding.JSON); err != nil {
					if err.Error() == "EOF" {
						body = ""
					} else {
						c.AbortWithStatusJSON(http.StatusBadRequest, response.FailBadRequest(err.Error()))
						return
					}
				}
				body = string(bz)
			} else {
				body = ""
			}

			if calculateSignature(uri, body, apiSecret, timestampInt) == signature {
				c.Next()
				return
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.FailBadRequest(fmt.Sprintf("Invalid %s", constant.HeaderSignature)))
			}
		}
	}
}

func checkTimeliness(timestamp int64) bool {
	now := time.Now().Unix()
	if now >= timestamp-constant.NetworkDelay && now <= timestamp+constant.NetworkDelay {
		return true
	}
	return false
}

func calculateSignature(uri, body, apiSecret string, timestamp int64) string {
	signStr := fmt.Sprintf(signStrFmt, timestamp, uri, body)
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(signStr))
	bytes := hash.Sum(nil)
	base64Str := base64.StdEncoding.EncodeToString(bytes)
	return base64Str
}

func verifyApiKey(apiKey string) (bool, string) {
	if apiKey == "" {
		return false, ""
	}

	res, err := apiKeyRepo.FindByApiKey(apiKey)
	if err != nil {
		return false, ""
	}

	return true, res.ApiSecret
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get(constant.HeaderApiKey)
		frequency, cycleTime := parseRateLimitPolicy()
		ok, err := RateLimitRepo.RateLimit(apiKey, frequency, cycleTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.FailSystemError())
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.FailBadRequest("Too many requests"))
			return
		}
		c.Next()
	}
}

func parseRateLimitPolicy() (frequency, cycleTime int) {
	split := strings.Split(global.Config.App.RateLimitPolicy, "/")
	if len(split) != 2 {
		return defaultRateLimitFrequency, defaultRateLimitCycleTime
	}

	var err error
	frequency, err = strconv.Atoi(split[0])
	if err != nil {
		return defaultRateLimitFrequency, defaultRateLimitCycleTime
	}

	cycleTime, err = strconv.Atoi(split[1])
	if err != nil {
		return defaultRateLimitFrequency, defaultRateLimitCycleTime
	}

	return frequency, cycleTime
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		if c.Writer.Status() != http.StatusOK {
			var resp vo.BaseResponse
			_ = json.Unmarshal(blw.Body.Bytes(), &resp)

			logrus.WithField("uri", c.Request.URL).WithField("req", c.Request.Body).WithField("resp", resp).
				Errorf("[%d]open api exception, msg: %s", c.Writer.Status(), resp.Message)

			//c.JSON().
		}
	}
}
