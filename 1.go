package example

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Catalog struct {
	XMLName    xml.Name    `xml:"archetype-catalog"`
	Archetypes []Archetype `xml:"archetypes>archetype"`
}

type Archetype struct {
	GroupID     string `xml:"groupId"`
	ArtifactID  string `xml:"artifactId"`
	Version     string `xml:"version"`
	Description string `xml:"description"`
}

func GetArchetypeCatalog() {
	resp, err := http.Get("https://repo1.maven.org/maven2/archetype-catalog.xml")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
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
			fmt.Println(archetype)
		}
	}
}
