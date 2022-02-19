package main

import (
	"context"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
)

var proxyPass string
var proxyMethod string
var ginLambda *ginadapter.GinLambda

func init() {

	//flags parse
	logLevel := "info"
	httpServiceInitTimeout := int64(3000)
	flag.StringVar(&logLevel, "logLevel", "info", "Log level")
	flag.StringVar(&proxyPass, "proxyPass", "http://localhost:80", "Endpoint of the service that will handle the request")
	flag.StringVar(&proxyMethod, "proxyMethod", "POST", "HTTP method to use for the proxy")
	flag.Int64Var(&httpServiceInitTimeout, "httpServiceInitTimeout", 15, "HTTP service bridged init timeout in seconds")
	flag.Parse()

	l, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("Invalid loglevel")
	}
	logrus.SetLevel(l)

	// defaults
	if logLevel == "" {
		logLevel = "info"
	}
	if proxyPass == "" {
		proxyPass = "http://localhost:80"
	}
	if httpServiceInitTimeout == 0 {
		httpServiceInitTimeout = 15
	}

	logrus.Infof("logLevel=%s", logLevel)
	logrus.Infof("proxyPass=%s", proxyPass)
	logrus.Infof("httpServiceInitTimeout=%d", httpServiceInitTimeout)

	//setup gin routes
	logrus.Debug("Initializing gin server")

	router := gin.Default()

	// setup CORS to allow everthing because we are running as a totally transparent proxy
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST",
		RequestHeaders:  "Origin, Content-Type, Authorization",
		ExposedHeaders:  "",
		MaxAge:          24 * 3600 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	//Catch all route
	router.Any("/*anything", proxy)

	httpServiceInitTimeoutNano := httpServiceInitTimeout * int64(1000) * int64(1000) * int64(1000)
	startAttemptTime := time.Now().UnixNano()
	for {
		_, err := http.Get(proxyPass)
		if err != nil {
			if time.Now().UnixNano()-startAttemptTime > httpServiceInitTimeoutNano {
				logrus.Errorf("Timeout waiting for the HTTP service at % to be ready", proxyPass)
				break
			}
			logrus.Warnf("HTTP service expected at %s not ready yet. Waiting...", proxyPass)
			time.Sleep(200 * time.Millisecond)
		} else {
			break
		}
	}

	ginLambda = ginadapter.New(router)
}

func proxy(c *gin.Context) {
	urlProxyPass, err := url.Parse(proxyPass)
	if err != nil {
		logrus.Errorf("Error parsing --proxyPass flag. Details: %s", err)
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(urlProxyPass)
	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, e error) {
		logrus.Errorf("Error proxying the request, quitting function. Details: %s", e)
		os.Exit(1)
	}
	proxy.Director = func(req *http.Request) {
		logrus.Debugf("Function invoked. Proxying to %s. Request data: %s", proxyPass, req)
		req.Header = c.Request.Header
		req.Host = urlProxyPass.Host
		if proxyMethod != "" {
			req.Method = proxyMethod
		}
		req.URL.Scheme = urlProxyPass.Scheme
		req.URL.Host = urlProxyPass.Host
		req.URL.Path = c.Param("anything")
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
func GinProxyHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	logrus.Infof("Starting HTTP Lambda Bridge")
	lambda.Start(GinProxyHandler)
}
