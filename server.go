package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	// Import New Relic and Logrus integrations
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrlogrus"
	"github.com/sirupsen/logrus"
)

// Global variables
var (
	greeting = flag.String("g", "Hello", "Greet with `greeting`")
	addr     = flag.String("addr", "localhost:8080", "address to serve")
	app      *newrelic.Application // New Relic application instance
	logger   *logrus.Logger        // Logger integrated with New Relic
)

func main() {
	// Parse flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	// Initialize New Relic application
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("HelloServer"),                     // Application name for New Relic APM
		newrelic.ConfigLicense("YOUR_NEW_RELIC_LICENSE_KEY"),      // Replace with your New Relic license key
		newrelic.ConfigDistributedTracerEnabled(true),            // Enable distributed tracing
		newrelic.ConfigAppLogForwardingEnabled(true),             // Enable log forwarding to New Relic Logs
		newrelic.ConfigDebugLogger(os.Stdout),                    // Enable debug logging for troubleshooting
	)
	if err != nil {
		log.Fatalf("Failed to initialize New Relic: %v", err)
	}

	// Initialize logrus with New Relic formatter
	logger = logrus.New()
	logger.SetFormatter(nrlogrus.NewFormatter(app, &logrus.TextFormatter{}))

	// Test log forwarding to New Relic
	logger.WithFields(logrus.Fields{
		"startup": "log-test",
	}).Info("Testing log forwarding to New Relic")

	// Register HTTP handlers with New Relic transaction wrapping
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/", greet))      // Wrap the "/" endpoint
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/version", version)) // Wrap the "/version" endpoint

	// Start the HTTP server
	log.Printf("Serving at http://%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// version handler displays build information
func version(w http.ResponseWriter, r *http.Request) {
	// Start a New Relic transaction
	txn := app.StartTransaction("Version")
	defer txn.End()

	// Associate the HTTP request with the transaction
	txn.SetWebRequestHTTP(r)

	// Log access to the "/version" endpoint
	txnLogger := logger.WithContext(newrelic.NewContext(context.Background(), txn))
	txnLogger.WithFields(logrus.Fields{
		"endpoint": "/version",
		"method":   r.Method,
		"user_ip":  r.RemoteAddr,
	}).Info("Version endpoint accessed")

	// Respond with build information
	info, ok := debug.ReadBuildInfo()
	if !ok {
		http.Error(w, "no build information available", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<!DOCTYPE html>\n<pre>\n")
	fmt.Fprintf(w, "%s\n", html.EscapeString(info.String()))
}

// greet handler responds with a personalized greeting
func greet(w http.ResponseWriter, r *http.Request) {
	// Start a New Relic transaction
	txn := app.StartTransaction("Greet")
	defer txn.End()

	// Associate the HTTP request with the transaction
	txn.SetWebRequestHTTP(r)

	// Log access to the "/" endpoint
	txnLogger := logger.WithContext(newrelic.NewContext(context.Background(), txn))
	txnLogger.WithFields(logrus.Fields{
		"endpoint": "/",
		"method":   r.Method,
		"user_ip":  r.RemoteAddr,
	}).Info("Greeting endpoint accessed")

	// Respond with greeting
	name := strings.Trim(r.URL.Path, "/")
	if name == "" {
		name = "Gopher"
	}

	fmt.Fprintf(w, "<!DOCTYPE html>\n")
	fmt.Fprintf(w, "%s, %s!\n", *greeting, html.EscapeString(name))
}
