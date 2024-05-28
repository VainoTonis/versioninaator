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
	Repository helmRepositoryDetails `yaml:"repository"`
}

type helmRepositoryDetails struct {
	URL        string                  `yaml:"URL"`
	Dependency []applicationDependency `yaml:"dependencies"`
}

type applicationDependency struct {
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	UsedChart     string `yaml:"usedChart"`
	UsedChartName string `yaml:"usedChartName"`
	UsedChartPath string `yaml:"UsedChartPath"`
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
	fmt.Print(targetDependencies)
}

func readConfiguration(configurationFile string) versioninaatorTargets {
	var targetConfigs versioninaatorTargets
	// Read the YAML file
	data, err := os.ReadFile(configurationFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the data into the Config struct
	if err := yaml.Unmarshal(data, &targetConfigs); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	log.Println(targetConfigs)
	return targetConfigs
}

func getTargetDependencies(targetConfigs versioninaatorTargets) []helmRepository {

	// Find and sort every dependency by repository URL
	depenenciesByRepository := []helmRepository{}

	for _, targetConfig := range targetConfigs.Targets {
		targetHelmChart := readChart(targetConfig.Path)

		for _, dependency := range targetHelmChart.Dependencies {
			repositoryExists := false
			applicationDepenencyDetails := []applicationDependency{{
				Name:          dependency.Name,
				Version:       dependency.Version,
				UsedChart:     targetConfig.URL,
				UsedChartName: targetHelmChart.Name,
				UsedChartPath: targetConfig.Path,
			}}

			// Check for already existing repositories

			for repositoryIndex, existingDependencyRepository := range depenenciesByRepository {

				if dependency.Repository == existingDependencyRepository.Repository.URL {
					depenenciesByRepository[repositoryIndex].Repository.Dependency = append(depenenciesByRepository[repositoryIndex].Repository.Dependency, applicationDepenencyDetails...)
					repositoryExists = true
					break
				}
			}

			// To avoid copy pasta same code a check goes through to either initialize the array or add a new entry to that array
			if len(depenenciesByRepository) == 0 || !repositoryExists {
				newDependencyRepository := helmRepository{
					Repository: helmRepositoryDetails{
						URL:        dependency.Repository,
						Dependency: applicationDepenencyDetails,
					},
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
