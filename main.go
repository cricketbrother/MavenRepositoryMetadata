package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

//type Catalog struct {
//	XMLName    xml.Name    `xml:"archetype-catalog"`
//	Archetypes []Archetype `xml:"archetypes>archetype"`
//}

type Archetype struct {
	GroupID     string `xml:"groupId"`
	ArtifactID  string `xml:"artifactId"`
	Version     string `xml:"version"`
	Description string `xml:"description"`
}

func GetArchetypes() {
	resp, err := http.Get("https://repo1.maven.org/maven2/archetype-catalog.xml")
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	decoder := xml.NewDecoder(resp.Body)
	count := 0
	for {
		token, err := decoder.Token()
		if token == nil {
			break
		}
		if err != nil {
			panic(err)
		}
		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "archetype" {
			var archetype Archetype
			if err := decoder.DecodeElement(&archetype, &se); err != nil {
				panic(err)
			}
			fmt.Printf("%s, %s, %s, %s\n", archetype.GroupID, archetype.ArtifactID, archetype.Version, archetype.Description)
			count += 1
		}
	}
	fmt.Printf("Total: %d\n", count)
}

type ALevelMetadata struct {
	XMLName    xml.Name   `xml:"metadata"`
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
	Versioning Versioning `xml:"versioning"`
}

type Versioning struct {
	Latest      string   `xml:"latest"`
	Release     string   `xml:"release"`
	Versions    []string `xml:"versions>version"`
	LastUpdated string   `xml:"lastUpdated"`
}

func verifyMavenMetadata(prefixUrl string) (ALevelMetadata, error) {
	metadataUrl, _ := url.JoinPath(prefixUrl, "maven-metadata.xml")
	var aLevelMetadata ALevelMetadata
	resp, err := http.Head(metadataUrl)
	if err != nil {
		return aLevelMetadata, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode == 200 {
		resp1, err := http.Get(metadataUrl)
		if err != nil {
			return aLevelMetadata, err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp1.Body)
		err = xml.NewDecoder(resp1.Body).Decode(&aLevelMetadata)
		if err != nil {
			fmt.Println(metadataUrl + " is not an A Level Metadata")
			return aLevelMetadata, err
		}
		return aLevelMetadata, nil
	}
	fmt.Println(metadataUrl + " is not a valid metadata url")
	return aLevelMetadata, err
}

func getNextLevelDirectories(prefixUrl string) (string, []string) {
	var nextLevelDirectories []string
	resp, err := http.Get(prefixUrl)
	if err != nil {
		return prefixUrl, nextLevelDirectories
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	node, err := html.Parse(resp.Body)
	if err != nil || node == nil {
		return prefixUrl, nextLevelDirectories
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				nextLevelDirectories = append(nextLevelDirectories, attr.Val)
			}
		}
	}
	return prefixUrl, nextLevelDirectories
}

func getHash(hashUrl string) string {
	fmt.Println(hashUrl)
	resp, err := http.Get(hashUrl)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	hashBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(hashBytes)
}

func printMetadata(prefixUrl string, metadata ALevelMetadata) {
	for _, version := range metadata.Versioning.Versions {
		md5Url, _ := url.JoinPath(prefixUrl, version, fmt.Sprintf("%s-%s.jar.md5", metadata.ArtifactID, version))
		sha1Url, _ := url.JoinPath(prefixUrl, version, fmt.Sprintf("%s-%s.jar.sha1", metadata.ArtifactID, version))
		fmt.Printf("%s, %s, %s, %s, %s\n", metadata.GroupID, metadata.ArtifactID, version, getHash(md5Url), getHash(sha1Url))
	}
}

func getMavenMetadata(prefixUrl string) {
	metadata, err := verifyMavenMetadata(prefixUrl)
	if err != nil {
		currentUrl, nextLevelDirectories := getNextLevelDirectories(prefixUrl)
		for _, dir := range nextLevelDirectories {
			nextUrl, _ := url.JoinPath(currentUrl, dir)
			getMavenMetadata(nextUrl)
		}
	}
	printMetadata(prefixUrl, metadata)
}

func main() {
	getMavenMetadata("https://repo.maven.apache.org/maven2/org/apache/maven/plugins/maven-jar-plugin")
}
