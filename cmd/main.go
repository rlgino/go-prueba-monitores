package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"rlgino/go-prueba-datadog/internal/handler"
	"rlgino/go-prueba-datadog/internal/logs"
)

func main() {
	grafanaURL := os.Getenv("GRAFANA_URL")
	if len(grafanaURL) == 0 {
		grafanaURL = "http://localhost:3100"
	}
	logger := logs.NewLogger(grafanaURL)
	// ########## DATADOG ##########

	//Override de JAEGER_DISABLED en caso de tracing no habilitado
	// New
	tracer.Start(
		tracer.WithService("api-golang-dd"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()

	// use JSONFormatter
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// log an event as usual with logrus
	logrus.WithFields(logrus.Fields{"instance": "instance-01", "int": 1, "float": 1.1}).Info("My first event from golang to stdout")
	// ############################

	tracer.Start(
		tracer.WithService("api-golang-dd"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()

	// Create a traced mux router
	mux := httptrace.NewServeMux()

	handlerV1 := handler.NewGreetingHandler(logger, uuid.New(), "v1")
	handlerV2 := handler.NewGreetingHandler(logger, uuid.New(), "v2")

	// Continue using the router as you normally would.
	mux.HandleFunc(handlerV1.GetURI(), handlerV1.Handle)
	mux.HandleFunc(handlerV2.GetURI(), handlerV2.Handle)

	log.Println("HTTP Server running")
	portNumber := os.Getenv("PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", portNumber), mux)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
