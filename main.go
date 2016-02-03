package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/bitrise-io/bitrise-webhooks/config"
	"github.com/bitrise-io/bitrise-webhooks/metrics"
)

func stringFlagOrEnv(flagValue *string, envKey string) string {
	if flagValue != nil && *flagValue != "" {
		return *flagValue
	}
	return os.Getenv(envKey)
}

func main() {
	var (
		portFlag        = flag.String("port", "", `Use port`)
		newRelicKeyFlag = flag.String("newrelic", "", `NewRelic license key`)
	)
	flag.Parse()

	port := stringFlagOrEnv(portFlag, "PORT")
	if port == "" {
		log.Fatal("Port must be set")
	}
	config.SetupServerEnvMode()

	// Monitoring
	if config.GetServerEnvMode() == config.ServerEnvModeProd {
		newRelicKey := stringFlagOrEnv(newRelicKeyFlag, "NEW_RELIC_LICENSE_KEY")
		metrics.SetupNewRelic("BitriseWebhooksProcessor", newRelicKey)
	} else {
		log.Println(" (!) Skipping NewRelic setup - environment is not 'production'")
	}

	// Routing
	setupRoutes()

	log.Println("Starting - using port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to ListenAndServe: %s", err)
	}
}