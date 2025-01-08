package internal

import (
	"encoding/xml"
	"log"
	"os"
)

type Resources struct {
	XMLName xml.Name `xml:"resources"`
	Strings []String `xml:"string"`
}

type String struct {
	XMLName      xml.Name
	Key          string `xml:"name,attr"`
	Value        string `xml:",innerxml"`
	Translatable string `xml:"translatable,attr,omitempty"`
}

// Open, read and marshal a existing XML file to Resources
func GetResourcesFromPathXML(path string) (Resources, error) {
	var resources Resources

	// Open and read the existing XML file
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return resources, err
	}

	// Unmarshal the file contents into the Resources struct.
	if err := xml.Unmarshal(fileBytes, &resources); err != nil {
		log.Fatalf("Error unmarshaling XML: %v", err)
		return resources, err
	}

	return resources, nil
}

func (r Resources) AppendNewString(newString String) Resources {
	r.Strings = append(r.Strings, newString)
	return r
}

// Marshal the updated Resources struct back to XML
func (r Resources) UpdateResourcesToXMLFile(path string) error {
	output, err := xml.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling XML: %v", err)
		return err
	}

	// Overwrite the same file (or create a new one).
	err = os.WriteFile(path, output, 0o644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
		return err
	}

	return nil
}