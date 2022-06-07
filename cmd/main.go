package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{"string": "foo", "int": 1, "float": 1.1}).Info("My first event from golang to stdout")
	// ############################

	handlerV1 := handler.NewGreetingHandler(logger, uuid.New(), "v1")
	handlerV2 := handler.NewGreetingHandler(logger, uuid.New(), "v2")

	http.HandleFunc(handlerV1.GetURI(), handlerV1.Handle)
	http.HandleFunc(handlerV2.GetURI(), handlerV2.Handle)

	log.Println("HTTP Server running")
	portNumber := os.Getenv("PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", portNumber), nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
