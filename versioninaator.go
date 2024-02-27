package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type deploymentRepositories struct {
	DeploymentRepositories []struct {
		RepoURL string `yaml:"URL"`
		Path    string `yaml:"path"`
	} `yaml:"deploymentRepositories"`
}

type helmChart struct {
	ApiVersion   string `yaml:"apiVersion"`
	Name         string `yaml:"name"`
	Description  string `yaml:"description"`
	Version      string `yaml:""`
	AppVersion   string `yaml:"appVersion"`
	Dependencies []struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Repository string `yaml:"repository"`
	} `yaml:"dependencies"`
}

func loadlocalChart(pathToLocalChart string) {
	data, err := os.ReadFile(pathToLocalChart)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var helmConfig helmChart
	if err := yaml.Unmarshal(data, &helmConfig); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	fmt.Printf("%+v\n", helmConfig)
}

func main() {
	configurationFile := flag.String("config", "", "Path to the configuration file")
	debugChart := flag.String("debug", "", "Insert Local File for development")
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
	var deploymentConfig deploymentRepositories
	if err := yaml.Unmarshal(data, &deploymentConfig); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	fmt.Printf("%+v\n", deploymentConfig)

	loadlocalChart(*debugChart)

}
