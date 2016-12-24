package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"strconv"

	"github.com/maksimov/liconv/input"
)

/**
TODO:
1. try and read POM file contents as text and try and parse the header comment if one is present, look for Copyrights and License info
	(ex. org/wildfly/common/wildfly-common/1.1.0.Final/wildfly-common-1.1.0.Final.pom)
**/
var wg sync.WaitGroup

const versionRegexp string = `(\d+\.)?(\d+\.)?(\*|\d+)`

var versionRegexpCompiled *regexp.Regexp

// escape to avoid conversion to number in excel, e.g. make versions return as "=""1.0"""
func escapeNumberString(s string) string {
	return `"="` + strconv.Quote(s) + `""`
}

func worker(wn int, in <-chan input.Dependency, out chan<- Component) {
	defer wg.Done()

	for dep := range in {
		group := strings.Replace(dep.GroupID, ".", "/", -1)
		artifact := dep.ArtifactID

		version := dep.Version

		uri := fmt.Sprintf("%s/%s/%s/%s-%s.pom", group, artifact, version, artifact, version)
		proj, err := input.ReadPOMFile(uri)
		fmt.Printf("[%d] %s\n\t%v\n", wn, uri, proj)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, license := range dep.Licenses.Licenses {
			c := Component{}
			c.Name = proj.Name
			if c.Name == "" || c.Name == "null" {
				c.Name = artifact
			}
			c.Name = strconv.Quote(strings.TrimSpace(c.Name))
			c.GroupID = dep.GroupID
			c.Language = "Java"
			c.Use = "Dynamically Linked"
			if proj.Description != "null" {
				// replace double quotes with single quotes in description to avoid csv breakage
				c.Description = strconv.Quote(strings.Replace(strings.TrimSpace(proj.Description), "\"", "'", -1))
			}
			c.SiteURL = proj.URL
			c.Version = escapeNumberString(version)
			c.ArtifactID = artifact
			c.Type = proj.Packaging
			c.LicenseName = strconv.Quote(strings.TrimSpace(license.Name))
			c.LicenseSiteURL = license.URL
			c.LicenseMisc = license.Distribution
			c.LicenseType = "Open Source"
			matches := versionRegexpCompiled.FindStringSubmatch(license.Name)
			if matches != nil {
				c.LicenseVersion = matches[0]
			}
			if c.LicenseVersion == "" {
				matches := versionRegexpCompiled.FindStringSubmatch(license.URL)
				if matches != nil {
					c.LicenseVersion = matches[0]
				}
			}
			c.LicenseVersion = escapeNumberString(c.LicenseVersion)
			out <- c
		}
	}
}

func writeResults(out chan Component, writer io.Writer) int {
	written := 0
	fmt.Fprintf(writer, "Component Name,Version,Description,Use,Type,Language,Website,License Name,Lic.Version,Lic.Type,Lic. Website URL\n")
	for c := range out {
		fmt.Fprintf(writer, fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			c.Name, c.Version, c.Description, c.Use, c.Type, c.Language, c.SiteURL, c.LicenseName, c.LicenseVersion, c.LicenseType, c.LicenseSiteURL))
		written++
	}
	return written
}

func init() {
	versionRegexpCompiled = regexp.MustCompile(versionRegexp)
}

func main() {
	var numWorkers int
	var outFilename string

	flag.IntVar(&numWorkers, "w", 8, "Number of simultaneous workers running at once")
	flag.StringVar(&outFilename, "o", "results.csv", "Output filename")
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("\t%s OPTIONS filename\n", os.Args[0])
		fmt.Printf("\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	arg := flag.Args()
	if len(arg) != 1 {
		flag.Usage()
		os.Exit(-1)
	}
	dependencies, err := input.ReadLicenses(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-3)
	}

	in, out := make(chan input.Dependency), make(chan Component)

	outFile, err := os.Create(outFilename)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}

	defer outFile.Close()

	go func() {
		numResults := writeResults(out, outFile)
		fmt.Printf("Written %d component licenses to %s\n", numResults, outFilename)
	}()

	// start numWorkers workers
	for n := 0; n < numWorkers; n++ {
		wg.Add(1)
		go worker(n, in, out)
	}

	// send every dependency into the "in" channel
	go func() {
		for _, dep := range dependencies {
			in <- dep
		}
		close(in)
	}()

	wg.Wait()
	close(out)
}
