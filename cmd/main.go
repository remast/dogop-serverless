package main

import (
	"os"

	// Blank-import the function package so the init() runs
	_ "crossnative.com/dogop"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	// 1. Port lesen (Default 8080)
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// 2. Hostname lesen (default 127.0.0.1)
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	// 3, Function Framework mit Host/Port starten
	funcframework.StartHostPort(hostname, port)
}
