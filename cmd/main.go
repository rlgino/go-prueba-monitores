package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	otr "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	tr "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"rlgino/go-prueba-datadog/internal/handler"
	"rlgino/go-prueba-datadog/internal/logs"
)

func main() {
	logger := logs.NewLogger("http://localhost:3100")
	// ########## DATADOG ##########

	var tracerCloser io.Closer

	//Override de JAEGER_DISABLED en caso de tracing no habilitado
	// New
	ddHost := "127.0.0.1"
	addr := net.JoinHostPort(ddHost, "8126")

	tracer := otr.New(tr.WithAgentAddr(addr))

	defer tr.Stop()

	opentracing.SetGlobalTracer(tracer)
	if tracerCloser != nil {
		defer tracerCloser.Close()
	}

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
