package input

import (
	"encoding/xml"
	"os"
)

type _summary struct {
	Deps _deps `xml:"dependencies"`
}

type _deps struct {
	Dependencies []Dependency `xml:"dependency"`
}

// Dependency represents a <dependency> element of licenses.xml
type Dependency struct {
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Licenses   Licenses `xml:"licenses"`
}

// Licenses represents a <licenses> element of licenses.xml
type Licenses struct {
	Licenses []License `xml:"license"`
}

// License represents a <license> element of license.xml
type License struct {
	Name         string `xml:"name"`
	URL          string `xml:"url"`
	Distribution string `xml:"distribution"`
}

// ReadLicenses reads the licenses from a reader and return all dependencies
func ReadLicenses(uri string) ([]Dependency, error) {
	file, err := os.Open(uri)

	defer file.Close()

	if err != nil {
		return nil, err
	}

	var xmlSummary _summary
	if err := xml.NewDecoder(file).Decode(&xmlSummary); err != nil {
		return nil, err
	}
	return xmlSummary.Deps.Dependencies, nil
}
