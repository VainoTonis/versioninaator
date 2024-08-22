package main

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
	URL        string                         `yaml:"url"`
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
	ApiVersion string                                  `yaml:"apiVersion"`
	Entries    map[string][]helmRepositoryIndexEntries `yaml:"entries"`
}

type helmRepositoryIndexEntries struct {
	ApiVersion string `yaml:"apiVersion"`
	AppVersion string `yaml:"appVersion"`
	Version    string `yaml:"version"`
	Name       string `yaml:"name"`
	Created    string `yaml:"created"`
}
