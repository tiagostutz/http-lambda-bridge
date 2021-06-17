package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
)

var proxyPass string

func init() {
	logLevel := "info"
	flag.StringVar(&logLevel, "logLevel", "info", "Log level")
	flag.StringVar(&proxyPass, "proxyPass", "http://localhost:80", "Endpoint of the service that will handle the request")
	flag.Parse()
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type proxy struct {
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (interface{}, error) {
	logrus.Debugf("Function invoked. Proxying to: %s", proxyPass)
	urlProxy, err := url.Parse(proxyPass)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing --proxyPass URL. Details: %s", err)
		logrus.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	if urlProxy.Scheme != "http" && urlProxy.Scheme != "https" {
		errMsg := fmt.Sprintf("Unsupported protocal scheme: %s", urlProxy.Scheme)
		logrus.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// create the client request
	var httpClient = &http.Client{Timeout: 30 * time.Second}

	// wrap the received body into a io.Reader to send the request body
	payloadBuf := strings.NewReader(request.Body)
	req, err := http.NewRequest(request.HTTPMethod, proxyPass, payloadBuf)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating the request to the Proxied endpoint. --proxyPass: %s. Error: %s", proxyPass, err)
		logrus.Errorf(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// copy the incoming headers to the proxied endpoint
	for k, vv := range request.Headers {
		req.Header.Add(k, vv)
	}

	//send the request to --proxyPass
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("Error invoking --proxyPass endpoint. Details: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Body, nil

}

func main() {
	lambda.Start(HandleRequest)
}
