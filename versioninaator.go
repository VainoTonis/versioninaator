package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type RSSFeedConfigurationApiInfo struct {
	Url       string                          `yaml:"url"`
	Uri       string                          `yaml:"uri"`
	ChartInfo []RSSFeedConfigurationChartInfo `yaml:"chartInfo"`
}

type RSSFeedConfigurationChartInfo struct {
	Repository string                          `yaml:"repository"`
	Chart      []RSSFeedConfigurationChartName `yaml:"chart"`
}

type RSSFeedConfigurationChartName struct {
	Name string `yaml:"name"`
}

type RSSFeedConfiguration struct {
	ApiInfo []RSSFeedConfigurationApiInfo `yaml:"apiInfo"`
}

func main() {
	// Read the YAML file
	data, err := os.ReadFile("rssFeedConfiguration.yaml")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the data into the Config struct
	var rssFeedConfiguration RSSFeedConfiguration
	if err := yaml.Unmarshal(data, &rssFeedConfiguration); err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}

	// Print the data
	fmt.Printf("%+v\n", rssFeedConfiguration)
}
