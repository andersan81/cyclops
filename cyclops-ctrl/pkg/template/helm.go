package template

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	json "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
	helmchart "helm.sh/helm/v3/pkg/chart"

	cyclopsv1alpha1 "github.com/andersan81/cyclops/cyclops-ctrl/api/v1alpha1"
	"github.com/andersan81/cyclops/cyclops-ctrl/internal/mapper"
	"github.com/andersan81/cyclops/cyclops-ctrl/internal/models"
	"github.com/andersan81/cyclops/cyclops-ctrl/internal/models/helm"
)

// extractPropertyOrder extracts property names from JSON schema in the order they appear
func extractPropertyOrder(jsonBytes []byte) []string {
	return extractPropertyOrderFromJSON(string(jsonBytes), "properties")
}

// extractPropertyOrderFromJSON extracts property names from a JSON string for a given properties key
func extractPropertyOrderFromJSON(jsonStr, propertiesKey string) []string {
	// Find the properties section
	searchKey := `"` + propertiesKey + `"`
	propertiesStart := strings.Index(jsonStr, searchKey)
	if propertiesStart == -1 {
		return nil
	}

	// Find the opening brace after "properties":
	braceStart := strings.Index(jsonStr[propertiesStart:], "{")
	if braceStart == -1 {
		return nil
	}
	braceStart += propertiesStart

	// Extract property names by finding quoted keys at the root level of properties
	var order []string
	remaining := jsonStr[braceStart+1:]
	depth := 0
	inString := false
	escaped := false
	i := 0

	for i < len(remaining) {
		char := remaining[i]

		if escaped {
			escaped = false
			i++
			continue
		}

		if char == '\\' {
			escaped = true
			i++
			continue
		}

		if char == '"' {
			inString = !inString
			if inString && depth == 0 {
				// Look for the end of this quoted string
				start := i + 1
				j := start
				for j < len(remaining) && remaining[j] != '"' {
					if remaining[j] == '\\' {
						j++ // skip escaped character
					}
					j++
				}
				if j < len(remaining) {
					propName := remaining[start:j]
					// Check if this is followed by a colon (indicating it's a property key)
					colonPos := j + 1
					for colonPos < len(remaining) && (remaining[colonPos] == ' ' || remaining[colonPos] == '\t' || remaining[colonPos] == '\n') {
						colonPos++
					}
					if colonPos < len(remaining) && remaining[colonPos] == ':' {
						order = append(order, propName)
					}
				}
			}
			i++
			continue
		}

		if !inString {
			if char == '{' || char == '[' {
				depth++
			} else if char == '}' || char == ']' {
				depth--
				if depth < 0 {
					break // End of properties object
				}
			}
		}

		i++
	}

	return order
}

// injectPropertyOrderRecursive recursively injects property order into schema and all nested schemas
func injectPropertyOrderRecursive(schema *helm.Property, jsonBytes []byte) {
	// Set the order for the current schema if it wasn't already specified
	if len(schema.Order) == 0 {
		propertyOrder := extractPropertyOrder(jsonBytes)
		if len(propertyOrder) > 0 {
			schema.Order = propertyOrder
		}
	}

	// Recursively process nested properties
	jsonStr := string(jsonBytes)
	for propName, property := range schema.Properties {
		if property.Type == "object" && len(property.Properties) > 0 {
			// Find this property's definition in the JSON to extract its nested order
			nestedOrder := findNestedPropertyOrder(jsonStr, propName)
			if len(nestedOrder) > 0 && len(property.Order) == 0 {
				// Update the property with the extracted order
				updatedProperty := property
				updatedProperty.Order = nestedOrder
				schema.Properties[propName] = updatedProperty

				// Recursively process this nested object
				injectPropertyOrderRecursive(&updatedProperty, jsonBytes)
			}
		}
	}
}

// findNestedPropertyOrder finds the property order for a nested object property
func findNestedPropertyOrder(jsonStr, propertyName string) []string {
	// Find the property definition
	propStart := strings.Index(jsonStr, `"`+propertyName+`"`)
	if propStart == -1 {
		return nil
	}

	// Find the opening brace of this property's object
	remaining := jsonStr[propStart:]
	bracePos := strings.Index(remaining, "{")
	if bracePos == -1 {
		return nil
	}

	// Find the "properties" key within this object
	objectStart := propStart + bracePos
	objectSection := jsonStr[objectStart:]

	// Look for "properties" within this object
	propertiesPos := strings.Index(objectSection, `"properties"`)
	if propertiesPos == -1 {
		return nil
	}

	// Extract from this specific properties section
	return extractPropertyOrderFromJSON(objectSection[propertiesPos:], "properties")
}

func (r Repo) LoadHelmChart(repo, chart, version, resolvedVersion string) (*models.Template, error) {
	var err error
	strictVersion := version
	if len(resolvedVersion) > 0 {
		strictVersion = resolvedVersion
	} else if !isValidVersion(version) {
		strictVersion, err = getRepoStrictVersion(repo, chart, version)
		if err != nil {
			return nil, err
		}
	}

	cached, ok := r.cache.GetTemplate(repo, chart, strictVersion, string(cyclopsv1alpha1.TemplateSourceTypeHelm))
	if ok {
		return cached, nil
	}

	tgzData, err := r.loadFromHelmChartRepo(repo, chart, version)
	if err != nil {
		return nil, err
	}

	extractedFiles, err := unpackTgzInMemory(tgzData)
	if err != nil {
		return nil, err
	}

	template, err := r.mapHelmChart(chart, extractedFiles)
	if err != nil {
		return nil, err
	}

	template.Version = version
	template.ResolvedVersion = strictVersion

	r.cache.SetTemplate(repo, chart, strictVersion, string(cyclopsv1alpha1.TemplateSourceTypeHelm), template)

	return template, nil
}

func (r Repo) LoadHelmChartInitialValues(repo, chart, version string) (map[string]interface{}, error) {
	var err error
	strictVersion := version
	if !isValidVersion(version) {
		strictVersion, err = getRepoStrictVersion(repo, chart, version)
		if err != nil {
			return nil, err
		}
	}

	cached, ok := r.cache.GetTemplateInitialValues(repo, chart, strictVersion, string(cyclopsv1alpha1.TemplateSourceTypeHelm))
	if ok {
		return cached, nil
	}

	tgzData, err := r.loadFromHelmChartRepo(repo, chart, version)
	if err != nil {
		return nil, err
	}

	extractedFiles, err := unpackTgzInMemory(tgzData)
	if err != nil {
		return nil, err
	}

	initial, err := r.mapHelmChartInitialValues(extractedFiles)
	if err != nil {
		return nil, err
	}

	r.cache.SetTemplateInitialValues(repo, chart, strictVersion, string(cyclopsv1alpha1.TemplateSourceTypeHelm), initial)

	return initial, nil
}

func IsHelmRepo(repo string) (bool, error) {
	indexURL, err := url.JoinPath(repo, "index.yaml")
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodHead, indexURL, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

func (r Repo) loadFromHelmChartRepo(repo, chart, version string) ([]byte, error) {
	tgzURL, err := getTarUrl(repo, chart, version)
	if err != nil {
		return nil, err
	}

	return downloadFile(tgzURL)
}

func (r Repo) mapHelmChart(chartName string, files map[string][]byte) (*models.Template, error) {
	metadataBytes := []byte{}
	schemaBytes := []byte{}
	chartFiles := make([]*helmchart.File, 0)

	templateFiles := make([]*helmchart.File, 0)
	crdFiles := make([]*helmchart.File, 0)

	dependenciesFromChartsDir := make(map[string]map[string][]byte, 0)

	for name, content := range files {
		parts := strings.Split(name, "/")

		if len(parts) == 2 && parts[1] == "Chart.yaml" {
			metadataBytes = content
			continue
		}

		if len(parts) == 2 && parts[1] == "values.schema.json" {
			schemaBytes = content
			continue
		}

		if len(parts) > 2 && parts[1] == "templates" &&
			(parts[2] != "Notes.txt" && parts[2] != "NOTES.txt" && parts[2] != "tests") {
			templateFiles = append(templateFiles, &helmchart.File{
				Name: path.Join(parts[1:]...),
				Data: content,
			})
			continue
		}

		if len(parts) > 2 && parts[1] == "crds" &&
			(parts[2] != "Notes.txt" && parts[2] != "NOTES.txt" && parts[2] != "tests") {
			crdFiles = append(crdFiles, &helmchart.File{
				Name: path.Join(parts[1:]...),
				Data: content,
			})
			continue
		}

		if len(parts) > 3 && parts[1] == "charts" {
			depName := parts[2]
			if _, ok := dependenciesFromChartsDir[depName]; !ok {
				dependenciesFromChartsDir[depName] = make(map[string][]byte)
			}

			dependenciesFromChartsDir[depName][path.Join(parts[2:]...)] = content
			continue
		}

		chartFiles = append(chartFiles, &helmchart.File{
			Name: name,
			Data: content,
		})

	}

	var schema helm.Property
	// unmarshal values schema only if present
	if len(schemaBytes) > 0 {
		if err := json.Unmarshal(schemaBytes, &schema); err != nil {
			fmt.Println("error on schema bytes", chartName)
			return &models.Template{}, err
		}

		// Inject property order recursively into the schema
		injectPropertyOrderRecursive(&schema, schemaBytes)
	}

	var metadata *helm.Metadata
	if err := yaml.Unmarshal(metadataBytes, &metadata); err != nil {
		fmt.Println("error on meta unm", chartName)
		return &models.Template{}, err
	}

	// region load dependencies
	dependencies, err := r.loadDependencies(metadata)
	if err != nil {
		return &models.Template{}, err
	}

	for depName, files := range dependenciesFromChartsDir {
		if dependencyExists(depName, dependencies) {
			continue
		}

		dep, err := r.mapHelmChart(depName, files)
		if err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}
	// endregion

	return &models.Template{
		Name:              chartName,
		RootField:         mapper.HelmSchemaToFields("", schema, schema.Definitions, dependencies),
		Files:             chartFiles,
		Templates:         templateFiles,
		CRDs:              crdFiles,
		Dependencies:      dependencies,
		HelmChartMetadata: metadata,
		RawSchema:         schemaBytes,
		IconURL:           metadata.Icon,
	}, nil
}

func (r Repo) mapHelmChartInitialValues(files map[string][]byte) (map[string]interface{}, error) {
	metadataBytes := []byte{}
	valuesBytes := []byte{}
	dependenciesFromChartsDir := make(map[string]map[string][]byte, 0)

	for name, content := range files {
		parts := strings.Split(name, "/")

		if len(parts) == 2 && parts[1] == "Chart.yaml" {
			metadataBytes = content
			continue
		}

		if len(parts) == 2 && parts[1] == "values.yaml" {
			valuesBytes = content
			continue
		}

		if len(parts) > 3 && parts[1] == "charts" {
			depName := parts[2]
			if _, ok := dependenciesFromChartsDir[depName]; !ok {
				dependenciesFromChartsDir[depName] = make(map[string][]byte)
			}

			dependenciesFromChartsDir[depName][path.Join(parts[2:]...)] = content
			continue
		}
	}

	values := make(map[string]interface{})
	if err := yaml.Unmarshal(valuesBytes, &values); err != nil {
		return nil, err
	}

	var metadata *helm.Metadata
	if err := yaml.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, err
	}

	// region load dependencies
	for depName, files := range dependenciesFromChartsDir {
		dep, err := r.mapHelmChartInitialValues(files)
		if err != nil {
			return nil, err
		}

		if values[depName] == nil {
			values[depName] = map[string]interface{}{}
		}

		values[depName] = overlayValues(values[depName], dep)
	}

	dependenciesFromMeta, err := r.loadDependenciesInitialValues(metadata)
	if err != nil {
		return nil, err
	}

	for depName, depValues := range dependenciesFromMeta {
		if values[depName] == nil {
			values[depName] = map[string]interface{}{}
		}

		values[depName] = overlayValues(values[depName], depValues)
	}
	// endregion

	return values, nil
}

func overlayValues(existing interface{}, overlay interface{}) interface{} {
	existingMap, existingOk := existing.(map[string]interface{})
	overlayMap, overlayOk := overlay.(map[string]interface{})

	if !existingOk || !overlayOk {
		return existing
	}

	for key, overlayValue := range overlayMap {
		if existingValue, exists := existingMap[key]; exists {
			switch existingValueTyped := existingValue.(type) {
			case map[string]interface{}:
				if overlayValueTyped, ok := overlayValue.(map[string]interface{}); ok {
					existingMap[key] = overlayValues(existingValueTyped, overlayValueTyped)
				} else {
					existingMap[key] = overlayValue
				}
			}
		} else {
			existingMap[key] = overlayValue
		}
	}
	return existingMap
}

func getTarUrl(repo, chart, version string) (string, error) {
	indexURL, err := url.JoinPath(repo, "index.yaml")
	if err != nil {
		return "", err
	}

	response, err := http.Get(indexURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data helm.Index
	err = yaml.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if _, ok := data.Entries[chart]; !ok {
		return "", errors.New(fmt.Sprintf("chart %v not found in repo %v", chart, repo))
	}

	v, err := resolveVersion(data.Entries[chart], version)
	if err != nil {
		return "", err
	}

	for _, entry := range data.Entries[chart] {
		if entry.Version == v {
			if len(entry.URLs) == 0 {
				return "", errors.New(fmt.Sprintf("no URL on version %v of chart %v and repo %v", version, chart, repo))
			}

			return entry.URLs[0], nil
		}
	}

	return "", errors.New(fmt.Sprintf("version %v not found in chart %v and repo %v", version, chart, repo))
}

func getRepoStrictVersion(repo, chart, version string) (string, error) {
	indexURL, err := url.JoinPath(repo, "index.yaml")
	if err != nil {
		return "", err
	}

	response, err := http.Get(indexURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data helm.Index
	err = yaml.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if _, ok := data.Entries[chart]; !ok {
		return "", errors.New(fmt.Sprintf("chart %v not found in repo %v", chart, repo))
	}

	return resolveVersion(data.Entries[chart], version)
}

func resolveVersion(indexEntries []helm.IndexEntry, version string) (string, error) {
	if isValidVersion(version) {
		return version, nil
	}

	versions := make([]string, 0, len(indexEntries))
	for _, entry := range indexEntries {
		versions = append(versions, entry.Version)
	}

	return resolveSemver(version, versions)
}

func downloadFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func unpackTgzInMemory(tgzData []byte) (map[string][]byte, error) {
	// Create a gzip reader
	gzipReader, err := gzip.NewReader(bytes.NewReader(tgzData))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	files := make(map[string][]byte)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var fileBuffer bytes.Buffer
		if _, err := io.Copy(&fileBuffer, tarReader); err != nil {
			return nil, err
		}

		files[header.Name] = fileBuffer.Bytes()
	}

	return files, nil
}

func dependencyExists(name string, existing []*models.Template) bool {
	for _, ed := range existing {
		if ed.Name == name {
			return true
		}
	}

	return false
}
