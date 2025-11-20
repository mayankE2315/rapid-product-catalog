package testutils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

type TestServer struct {
	ginContext *gin.Context
	router     *gin.Engine
	recorder   *httptest.ResponseRecorder
}

func NewServer() *TestServer {
	server := new(TestServer)
	server.SetupTest()
	return server
}

func (suite *TestServer) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.recorder = httptest.NewRecorder()
	c, r := gin.CreateTestContext(suite.recorder)
	r.Use(SetupAuthMiddleware())
	suite.ginContext = c
	suite.router = r
}

func SetupAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("branchId", "branch-id")
	}
}

func (suite *TestServer) Start(req *http.Request) {
	suite.router.ServeHTTP(suite.recorder, req)
}

func (suite *TestServer) Router() *gin.Engine {
	return suite.router
}

func (suite *TestServer) Context() *gin.Context {
	return suite.ginContext
}

func (suite *TestServer) Recorder() *httptest.ResponseRecorder {
	return suite.recorder
}

func (suite *TestServer) PerformRequest(url, method string, requestBody interface{}) {
	buf := new(bytes.Buffer)
	if requestBody != nil {
		if err := json.NewEncoder(buf).Encode(requestBody); err != nil {
			panic(err)
		}
	}

	req, _ := http.NewRequest(strings.ToUpper(method), url, buf)
	suite.Start(req)
}

func (suite *TestServer) PerformRequestWithRequestBody(url, method string, requestBody string) {
	req, _ := http.NewRequest(strings.ToUpper(method), url, ioutil.NopCloser(bytes.NewBufferString(requestBody)))
	suite.Start(req)
}
