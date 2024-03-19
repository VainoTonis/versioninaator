package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type versioninaatorTargets struct {
	Targets []struct {
		URL    string `yaml:"URL"`
		Path   string `yaml:"path"`
		Branch string `yaml:"branch"`
	} `yaml:"targets"`
}

type helmChart struct {
	ApiVersion   string `yaml:"apiVersion"`
	Name         string `yaml:"name"`
	Dependencies []struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Repository string `yaml:"repository"`
	} `yaml:"dependencies"`
}

type helmRepository struct {
	Repository []struct {
		URL        string `yaml:"URL"`
		Dependency []struct {
			Name          string `yaml:"name"`
			Version       string `yaml:"version"`
			UsedChart     string `yaml:"usedChart"`
			UsedChartName string `yaml:"usedChartName"`
			UsedChartPath string `yaml:"usedPath"`
		} `yaml:"dependencies"`
	} `yaml:"repository"`
}

func main() {
	baseChartYaml := "Chart.yaml"
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
	var targetConfigs versioninaatorTargets
	if err := yaml.Unmarshal(data, &targetConfigs); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}
	for _, targetConfig := range targetConfigs.Targets {
		targetHelmChart := loadlocalChart(targetConfig.Path + baseChartYaml)

		fmt.Printf("%s\n", targetConfig)
		fmt.Printf("%s", targetHelmChart)
	}

}

func loadlocalChart(pathToLocalChart string) helmChart {
	data, err := os.ReadFile(pathToLocalChart)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var helmConfig helmChart
	if err := yaml.Unmarshal(data, &helmConfig); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	return helmConfig
}
