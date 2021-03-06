package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsPort    string
	dataSourceName string
)

const (
	metricsHost = "0.0.0.0"
	metricsPath = "/metrics"
	healthzPath = "/healthz"
	namespace   = "pgbouncer"
	indexHTML   = `
	<html>
		<head><title>PgBouncer Metrics Exporter</title></head>
		<body>
			<h1>PgBouncer Metrics Exporter</h1>
			<ul>
				<li><a href='` + metricsPath + `'>metrics</a></li>
				<li><a href='` + healthzPath + `'>healthz</a></li>
			</ul>
		</body>
	</html>`
)

func ParseEnv() {
	if dsn := os.Getenv("DATA_SOURCE_NAME"); len(dsn) != 0 {
		dataSourceName = dsn
	}
	if port := os.Getenv("PGB_EXPORTER_WEB_LISTEN_PORT"); len(port) != 0 {
		metricsPort = port
	}
}

func main() {
	flag.StringVar(&metricsPort, "p", "9127", "Port to listen on for web interface and telemetry")
	flag.StringVar(&dataSourceName, "d", "postgres://pgbouncer:@localhost:6432/pgbouncer?sslmode=disable", "PgBouncer connection url")
	flag.Parse()
	ParseEnv()

	// Connect to pgbouncer
	db, err := Connect(dataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to PgBouncer: ", err)
	}

	// Create new collector
	collector := NewCollector(db, namespace)
	defer collector.Close()

	// Register collector
	r := prometheus.NewRegistry()
	r.MustRegister(collector)

	listenAddress := net.JoinHostPort(metricsHost, fmt.Sprint(metricsPort))
	mux := http.NewServeMux()

	// Add metricsPath
	mux.Handle(metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	// Add healthzPath
	mux.HandleFunc(healthzPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("ok")); err != nil {
			log.Fatal("Unable to write to serve metrics: ", err)
		}
	})

	// Add index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(indexHTML))
		if err != nil {
			log.Fatal("Unable to write to serve metrics: ", err)
		}
	})

	err = http.ListenAndServe(listenAddress, mux)
	log.Fatal("Failed to serve metrics: ", err)
}
