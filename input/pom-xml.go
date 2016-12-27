package input

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/text/encoding/charmap"
)

// Project represents the <project> element of pom.xml
type Project struct {
	ArtifactID  string `xml:"artifactId"`
	Name        string `xml:"name"`
	URL         string `xml:"url"`
	Description string `xml:"description"`
	Packaging   string `xml:"packaging"`
}

var client http.Client

var remoteBase = "http://search.maven.org/remotecontent?filepath="
var localBase string

func init() {
	cache := diskcache.New("pom")
	t := httpcache.NewTransport(cache)
	client = http.Client{Transport: t}
	localBase = userHomeDir() + "/.m2/repository/"
}

// ReadPOMFile reads a POM file from local $HOME/.m2/repository and if it fails fetches one from maven.org
func ReadPOMFile(uri string) (*Project, error) {
	pom, err := readLocal(uri)
	if err != nil {
		pom, err = readRemote(uri)
	}

	if err != nil {
		return nil, err
	}
	var project Project
	decoder := xml.NewDecoder(strings.NewReader(pom))
	decoder.CharsetReader = makeCharsetReader
	if err := decoder.Decode(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

func readLocal(uri string) (string, error) {
	f, err := os.Open(localBase + uri)
	defer f.Close()

	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func readRemote(uri string) (string, error) {
	fmt.Printf("-> Reading remote cached " + remoteBase + uri + "\n")
	req, err := http.NewRequest("GET", remoteBase+uri, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%d while reading %s", resp.StatusCode, uri)
		return "", err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var isWindows = func() bool {
	return runtime.GOOS == "windows"
}

func userHomeDir() string {
	if isWindows() {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	charset = strings.ToLower(charset)
	if charset == "iso-8859-1" || charset == "windows-1252" {
		// Windows-1252 is a superset of ISO-8859-1, so should do here
		return charmap.Windows1252.NewDecoder().Reader(input), nil
	}
	return nil, fmt.Errorf("Unknown charset: %s", charset)
}
