Roadmap
- Get repo url from dependencies
  - Future: Handle situations where the repo is only referenced by name
- get the index.yaml from dependencies
  - Doublecheck if that index.yaml already exists
- get latest chart version and date index.yaml (it contains all versions)
- Compare requested version against latest
  - Optional: Some releases can be weekly, get the latest stable version
- TODO! Figure out how and in what format the data is going to be shown in grafana, this is crucial to understand how the data is formated
- Helm chart for deployment
