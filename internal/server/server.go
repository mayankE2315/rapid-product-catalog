package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roppenlabs/rapid-product-catalog/internal/config"
	logger "github.com/roppenlabs/rapido-logger-go"
)

type Server struct {
	config       config.Config
	engine       *gin.Engine
	routerGroups RouterGroups
}

type RouterGroups struct {
	rootRouter *gin.Engine
}

func NewServer(c config.Config) *Server {

	if c.IsProductionEnv() {
		logger.Info(logger.Format{Message: "Setting gin server to release mode for production environment"})
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			logger.Info(
				logger.Format{
					Message: fmt.Sprintf("Endpoint %s is declared via handler %s", absolutePath, handlerName),
					Data: map[string]string{
						"method": httpMethod,
					},
				})
		}
	}
	engine := gin.New()
	loggerConfig := gin.LoggerConfig{
		SkipPaths: []string{"/sanity", "/health"},
	}
	engine.Use(LoggerWithConfig(loggerConfig), gin.Recovery())

	return &Server{
		config: c,
		engine: engine,
		routerGroups: RouterGroups{
			rootRouter: engine,
		},
	}
}

func (s *Server) Run(h Handlers) {
	s.InitRoutes(h, s.config)
	srv := &http.Server{
		Addr:    s.config.Get().ListenAddress(),
		Handler: s.engine,
	}
	go listenServer(srv)
	waitForShutdown(srv)
}

func listenServer(server *http.Server) {
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Info(logger.Format{Message: "server shutting down"})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error(logger.Format{Message: fmt.Sprintf("Server forced to shutdown: %v", err)})
	}
	logger.Info(logger.Format{Message: "server shutdown complete"})
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {

	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.Debug(logger.Format{
				Message: fmt.Sprintf("Accessing %s", path),
				Data: map[string]string{
					"method":     param.Method,
					"clientIP":   param.ClientIP,
					"statusCode": strconv.Itoa(param.StatusCode),
					"error":      param.ErrorMessage,
					"latency":    param.Latency.String(),
					"size":       strconv.Itoa(param.BodySize),
				},
			})
		}
	}
}
