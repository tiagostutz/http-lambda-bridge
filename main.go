package main

import (
	"context"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
)

var proxyPass string
var ginLambda *ginadapter.GinLambda

func init() {

	//flags parse
	logLevel := "info"
	flag.StringVar(&logLevel, "logLevel", "info", "Log level")
	flag.StringVar(&proxyPass, "proxyPass", "http://localhost:80", "Endpoint of the service that will handle the request")
	flag.Parse()

	l, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("Invalid loglevel")
	}
	logrus.SetLevel(l)

	logrus.Infof("logLevel=%s", logLevel)
	logrus.Infof("proxyPass=%s", proxyPass)

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

	time.Sleep(9 * time.Second)

	ginLambda = ginadapter.New(router)
}
func proxy(c *gin.Context) {
	urlProxyPass, err := url.Parse(proxyPass)
	if err != nil {
		logrus.Errorf("Error parsing --proxyPass flag. Details: %s", err)
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(urlProxyPass)
	proxy.Director = func(req *http.Request) {
		logrus.Debugf("Function invoked. Proxying to %s. Request data: %s", proxyPass, req)
		req.Header = c.Request.Header
		req.Host = urlProxyPass.Host
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
