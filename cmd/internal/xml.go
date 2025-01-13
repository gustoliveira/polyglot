package internal

import (
	"encoding/xml"
	"fmt"
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
		fmt.Printf("Error reading file: %v\n", err)
		return resources, err
	}

	// Unmarshal the file contents into the Resources struct.
	if err := xml.Unmarshal(fileBytes, &resources); err != nil {
		fmt.Printf("Error unmarshaling XML: %v\n", err)
		return resources, err
	}

	return resources, nil
}

func (r Resources) AppendNewString(newString String) Resources {
	r.Strings = append(r.Strings, newString)
	return r
}

func (r Resources) RemoveStringByKey(key string) Resources {
	for index, s := range r.Strings {
		if s.Key == key {
			r.Strings = append(r.Strings[:index], r.Strings[index+1:]...)
		}
	}
	return r
}

func (r Resources) ContainsStringByKey(key string) bool {
	for _, s := range r.Strings {
		if s.Key == key {
			return true
		}
	}

	return false
}

func (r Resources) CreateOrSubstituteStringByKey(key string, value string) Resources {
	for index, s := range r.Strings {
		if s.Key == key {
			r.Strings[index].Value = value
			return r
		}
	}

	r = r.AppendNewString(String{
		XMLName: xml.Name{Local: "string"},
		Key:     key,
		Value:   value,
	})

	return r
}

// Marshal the updated Resources struct back to XML
func (r Resources) UpdateResourcesToXMLFile(path string) error {
	output, err := xml.MarshalIndent(r, "", "    ")
	if err != nil {
		fmt.Printf("Error marshaling XML: %v\n", err)
		return err
	}

	// Overwrite the same file (or create a new one).
	err = os.WriteFile(path, output, 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return err
	}

	return nil
}

func (r Resources) IsSortedByKey() bool {
	for i := range r.Strings {
		if i == 0 {
			continue
		}

		if r.Strings[i-1].Key > r.Strings[i].Key {
			return false
		}
	}

	return true
}
