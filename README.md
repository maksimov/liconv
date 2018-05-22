# liconv

## Description
A small utility which helps to convert a dependencies report generated with **license-maven-plugin** to a CSV file which is handy for reports that are often required by a legal team in large enterprises, etc. The utility will gather the basic information (component name, groupId, version) available in the original report itself and also consult local and remote .pom files for further information (component description, URL, license information).

The list of fields that is currently generated:
- Component
  - Name
  - Version
  - Description
  - Use _(hardcoded to "Dynamically Linked")_
  - Type
  - Language _(hardcoded to "Java")_
  - Website URL
- License
  - Name
  - Version _(manual checking is required)_
  - Type _(hardcoded to "Open Source")_
  - Website URL
 
## Run
Run `license-maven-plugin` on your Maven project: 
```
mvn org.codehaus.mojo:license-maven-plugin:1.12:aggregate-download-licenses
```
This will generate a `licenses.xml` in `target/generated-resources` directory. 

Run `liconv` like so:
```
liconv target/generated-resources/licenses.xml
```
This will generate a `results.csv` file

## Build

```
go get github.com/gregjones/httpcache/...
go get golang.org/x/text/encoding/charmap

go install
```
