package main

import (
	"fmt"
	"os"
)

func ConfigFile() string {
	if s := os.Getenv("BUFFALAB_CONFIG"); s != "" {
		return s
	}
	return "buffalab.yml"
}

func main() {
	fmt.Printf("BuffaLab coming online...\n")

	fmt.Printf("reading configuration from %s...\n", ConfigFile())
	config, err := ReadConfig(ConfigFile())
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILED: %s\n", err)
		os.Exit(1)
	}

	api := NewAPIServer(config)
	api.Run()
	fmt.Printf("Buffalab going offline...\n")
}
