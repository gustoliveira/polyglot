package internal

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

type Resources struct {
	XMLName     xml.Name    `xml:"resources"`
	Strings     []String    `xml:"string"`
	Translation Translation `xml:"-"`
}

type String struct {
	XMLName      xml.Name
	Key          string `xml:"name,attr"`
	Value        string `xml:",innerxml"`
	Translatable string `xml:"translatable,attr,omitempty"`
}

type AllResources struct {
	existentResourcesPaths []string
	stringKeys             map[string][]string
}

type ListResources []Resources

func (ar AllResources) CheckMissingTranslationsRelatory() string {
	result := ""

	relatory := []string{}

	for key, paths := range ar.stringKeys {
		diff := StringSlicesDiff(ar.existentResourcesPaths, paths)
		if len(diff) == 0 {
			continue
		}

		defined_in := fmt.Sprintf("DEFINED IN: [%v]", strings.Join(paths, ", "))
		missing_from := fmt.Sprintf("MISSING FROM: [%v]", strings.Join(diff, ", "))
		relatory = append(relatory, fmt.Sprintf("\t%v:\n\t\t%v\n\t\t%v\n", key, defined_in, missing_from))
	}

	if len(relatory) == 0 {
		return "No missing translations found."
	}

	result += fmt.Sprintf("\nFound %v possible missing translations:\n", len(relatory))
	for _, l := range relatory {
		result += l
	}

	return result
}

func (lr ListResources) CheckMissingTranslations() AllResources {
	allResources := AllResources{
		existentResourcesPaths: []string{},
		stringKeys:             make(map[string][]string),
	}

	for _, r := range lr {
		allResources.existentResourcesPaths = append(allResources.existentResourcesPaths, r.Translation.Language)

		for _, s := range r.Strings {
			if s.Translatable == "false" {
				continue
			}

			if _, ok := allResources.stringKeys[s.Key]; !ok {
				allResources.stringKeys[s.Key] = []string{}
			}

			allResources.stringKeys[s.Key] = append(allResources.stringKeys[s.Key], r.Translation.Language)
		}
	}

	return allResources
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

	translation, err := GetTranslationFromFileName(path)
	if err != nil {
		return resources, err
	}

	resources.Translation = translation

	return resources, nil
}

func (r Resources) AppendNewString(newString String) Resources {
	r.Strings = append(r.Strings, newString)
	return r
}

func (r Resources) AddNewStringSorted(newString String) Resources {
	index := r.IndexToAddSorted(newString)
	r.Strings = slices.Insert(r.Strings, index, newString)
	return r
}

func (r Resources) IndexToAddSorted(newString String) int {
	return sort.Search(len(r.Strings), func(i int) bool {
		return r.Strings[i].Key >= newString.Key
	})
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

	if r.IsSortedByKey() {
		r = r.AddNewStringSorted(String{
			XMLName: xml.Name{Local: "string"},
			Key:     key,
			Value:   value,
		})
	} else {
		r = r.AppendNewString(String{
			XMLName: xml.Name{Local: "string"},
			Key:     key,
			Value:   value,
		})
	}

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

func (r Resources) SortByKey() {
	if r.IsSortedByKey() {
		return
	}

	sort.SliceStable(r.Strings, func(i, j int) bool {
		return r.Strings[i].Key < r.Strings[j].Key
	})
}
