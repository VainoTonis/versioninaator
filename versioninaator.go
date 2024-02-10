package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type RemoteAPIInfo struct {
	RemoteAPIInfo []struct {
		URL       string `yaml:"url"`
		URI       string `yaml:"uri"`
		ChartInfo []struct {
			Repository string `yaml:"repository"`
			Chart      []struct {
				Name string `yaml:"name"`
			} `yaml:"chart"`
		} `yaml:"chartInfo"`
	} `yaml:"RemoteAPIInfo"`
}

func main() {
	configurationFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	if *configurationFile == "" {
		envVariableFile, found := os.LookupEnv("versioninaatorConfiguration")
		if !found {
			log.Fatalf("You must set a configuration file either by env variable (versioninaatorConfiguration) or cli argument (-config <configName>)")
		}
		*configurationFile = envVariableFile
	}

	// Read the YAML file
	data, err := os.ReadFile(*configurationFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the data into the Config struct
	var remoteConfig RemoteAPIInfo
	if err := yaml.Unmarshal(data, &remoteConfig); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	// Print the data
	fmt.Printf("%+v\n", remoteConfig)
}
