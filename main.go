package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

func main() {
	GetArchetypes()
}
