package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type versioninaator struct {
	ApiVersion string                  `yaml:"apiVersion"`
	Targets    []versioninaatorTargets `yaml:"targets"`
}

type versioninaatorTargets struct {
	URL    string `yaml:"URL"`
	Path   string `yaml:"path"`
	Branch string `yaml:"branch"`
}

type helmChart struct {
	ApiVersion   string                  `yaml:"apiVersion"`
	Name         string                  `yaml:"name"`
	Dependencies []helmChartDependencies `yaml:"dependencies"`
}
type helmChartDependencies struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Repository string `yaml:"repository"`
}

type inUseHelmRepositories struct {
	Repository string                         `yaml:"repository"`
	Dependency []inUseApplicationDependencies `yaml:"dependencies"`
}

type inUseApplicationDependencies struct {
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	UsedChart     string `yaml:"usedChart"`
	UsedChartName string `yaml:"usedChartName"`
	UsedChartPath string `yaml:"UsedChartPath"`
}

type helmRepositoryIndex struct {
	ApiVersion string                                `yaml:"apiVersion"`
	Entries    map[string]helmRepositoryIndexEntries `yaml:"entries"`
}

type helmRepositoryIndexEntries struct {
	ApiVersion string `yaml:"apiVersion"`
	AppVersion string `yaml:"appVersion"`
	Version    string `yaml:"version"`
	Name       string `yaml:"name"`
	Created    string `yaml:"created"`
}

func main() {
	// baseChartYaml := "Chart.yaml"
	configurationFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	if *configurationFile == "" {
		envVariableFile, found := os.LookupEnv("versioninaatorConfiguration")
		if !found {
			log.Fatalf("You must set a configuration file either by env variable (versioninaatorConfiguration) or cli argument (-config <configName>)")
		}
		*configurationFile = envVariableFile
	}
	targetConfigs := readConfiguration(*configurationFile)
	targetDependencies := getTargetDependencies(targetConfigs)
	getLatestDependencies(targetDependencies)

	/*
		get every repository and every unique dependecy name
		get the latest versions of every dependency and when they were released
		Create a list of charts that use outdated dependencies
		create custom metrics to show the results
	*/
}

func readConfiguration(configurationFile string) versioninaator {
	var targetConfigs versioninaator
	// Read the YAML file
	data, err := os.ReadFile(configurationFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the data into the Config struct
	if err := yaml.Unmarshal(data, &targetConfigs); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	return targetConfigs
}

func getTargetDependencies(targetConfigs versioninaator) []inUseHelmRepositories {

	// Find and sort every dependency by repository URL
	depenenciesByRepository := []inUseHelmRepositories{}

	for _, targetConfig := range targetConfigs.Targets {
		targetHelmChart := readChart(targetConfig.Path)

		for _, dependency := range targetHelmChart.Dependencies {
			repositoryExists := false
			applicationDepenencyDetails := []inUseApplicationDependencies{{
				Name:          dependency.Name,
				Version:       dependency.Version,
				UsedChart:     targetConfig.URL,
				UsedChartName: targetHelmChart.Name,
				UsedChartPath: targetConfig.Path,
			}}

			// Check for already existing repositories

			for repositoryIndex, existingDependencyRepository := range depenenciesByRepository {

				if dependency.Repository == existingDependencyRepository.Repository {
					depenenciesByRepository[repositoryIndex].Dependency = append(depenenciesByRepository[repositoryIndex].Dependency, applicationDepenencyDetails...)
					repositoryExists = true
					break
				}
			}

			// To avoid copy pasta same code a check goes through to either initialize the array or add a new entry to that array
			if len(depenenciesByRepository) == 0 || !repositoryExists {
				newDependencyRepository := inUseHelmRepositories{
					Repository: dependency.Repository,
					Dependency: applicationDepenencyDetails,
				}
				depenenciesByRepository = append(depenenciesByRepository, newDependencyRepository)
			}
		}
	}

	return depenenciesByRepository
}

func readChart(pathToLocalChart string) helmChart {
	rawChart, err := os.ReadFile(pathToLocalChart)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var unmarshaledChart helmChart
	if err := yaml.Unmarshal(rawChart, &unmarshaledChart); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	return unmarshaledChart
}

func getLatestDependencies(inUseHelmRepositories []inUseHelmRepositories) {
	fmt.Println()

	marshaled, err := yaml.Marshal(inUseHelmRepositories)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}
	fmt.Println(string(marshaled))

	
}
