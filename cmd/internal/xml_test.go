package internal

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAddNewStringSorted(t *testing.T) {
	testCases := []struct {
		name             string
		initialResource  Resources
		newString        String
		expectedResource Resources
	}{
		{
			name:            "Add to empty Resources",
			initialResource: Resources{},
			newString:       String{Key: "test_key", Value: "Test Value"},
			expectedResource: Resources{
				Strings: []String{{Key: "test_key", Value: "Test Value"}},
			},
		},
		{
			name: "Insert into a sorted Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "d", Value: "D"},
				},
			},
			newString: String{Key: "c", Value: "C"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
				},
			},
		},
		{
			name: "Insert into an unsorted Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "d", Value: "D"},
					{Key: "b", Value: "B"},
				},
			},
			newString: String{Key: "c", Value: "C"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
					{Key: "b", Value: "B"},
				},
			},
		},
		{
			name: "Insert into a sorted Resources with duplicated keys",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
				},
			},
			newString: String{Key: "c", Value: "C"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "c", Value: "C"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
				},
			},
		},
		{
			name: "Insert into an unsorted Resources with duplicated keys",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "d", Value: "D"},
					{Key: "c", Value: "C"},
				},
			},
			newString: String{Key: "c", Value: "C"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
					{Key: "c", Value: "C"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.AddNewStringSorted(tc.newString)

			if !reflect.DeepEqual(result, tc.expectedResource) {
				t.Errorf("AddNewStringSorted() = %v, want %v", result, tc.expectedResource)
			}
		})
	}
}

func TestIndexToAddSorted(t *testing.T) {
	testCases := []struct {
		name            string
		initialResource Resources
		newString       String
		expectedIndex   int
	}{
		{
			name:            "Add to empty Resources",
			initialResource: Resources{},
			newString:       String{Key: "test_key", Value: "Test Value"},
			expectedIndex:   0,
		},
		{
			name: "Insert into a sorted Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "d", Value: "D"},
				},
			},
			newString:     String{Key: "c", Value: "C"},
			expectedIndex: 2,
		},
		{
			name: "Insert into an unsorted Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "d", Value: "D"},
					{Key: "b", Value: "B"},
				},
			},
			newString:     String{Key: "c", Value: "C"},
			expectedIndex: 1,
		},
		{
			name: "Insert into an unsorted with duplicated key",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "d", Value: "D"},
					{Key: "c", Value: "C"},
				},
			},
			newString:     String{Key: "c", Value: "C"},
			expectedIndex: 2,
		},
		{
			name: "Insert into a sorted with duplicated key",
			initialResource: Resources{
				Strings: []String{
					{Key: "a", Value: "A"},
					{Key: "b", Value: "B"},
					{Key: "c", Value: "C"},
					{Key: "d", Value: "D"},
				},
			},
			newString:     String{Key: "c", Value: "C"},
			expectedIndex: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.IndexToAddSorted(tc.newString)

			if !reflect.DeepEqual(result, tc.expectedIndex) {
				t.Errorf("IndexToAddSorted() = %v, want %v", result, tc.expectedIndex)
			}
		})
	}
}

func TestAppendNewString(t *testing.T) {
	testCases := []struct {
		name             string
		initialResource  Resources
		newString        String
		expectedResource Resources
	}{
		{
			name:            "Append to empty Resources",
			initialResource: Resources{},
			newString:       String{Key: "test_key", Value: "Test Value"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "test_key", Value: "Test Value"},
				},
			},
		},
		{
			name: "Append to non-empty Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "existing_key", Value: "Existing Value"},
				},
			},
			newString: String{Key: "new_key", Value: "New Value", Translatable: "false"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "existing_key", Value: "Existing Value"},
					{Key: "new_key", Value: "New Value", Translatable: "false"},
				},
			},
		},
		{
			name: "Append a duplicated value in Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "existing_key", Value: "Existing Value"},
				},
			},
			newString: String{Key: "existing_key", Value: "Existing Value"},
			expectedResource: Resources{
				Strings: []String{
					{Key: "existing_key", Value: "Existing Value"},
					{Key: "existing_key", Value: "Existing Value"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.AppendNewString(tc.newString)

			if !reflect.DeepEqual(result, tc.expectedResource) {
				t.Errorf("AppendNewString() = %v, want %v", result, tc.expectedResource)
			}
		})
	}
}

func TestContainsStringByKey(t *testing.T) {
	testCases := []struct {
		name             string
		initialResource  Resources
		searchedString   string
		expectedContains bool
	}{
		{
			name: "Contains in Resources",
			initialResource: Resources{
				Strings: []String{{Key: "test_key", Value: "Test Value"}},
			},
			searchedString:   "test_key",
			expectedContains: true,
		},
		{
			name: "Dont Contains in Resources",
			initialResource: Resources{
				Strings: []String{{Key: "palmeiras", Value: "Quando surge o alviverde imponente"}},
			},
			searchedString:   "test_key",
			expectedContains: false,
		},
		{
			name: "Empty Resources",
			initialResource: Resources{
				Strings: []String{},
			},
			searchedString:   "test_key",
			expectedContains: false,
		},
		{
			name: "Empty Resources with empty string",
			initialResource: Resources{
				Strings: []String{},
			},
			searchedString:   "",
			expectedContains: false,
		},
		{
			name: "With duplicated values",
			initialResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Quando surge o alviverde imponente"},
					{Key: "palmeiras", Value: "No gramado em que a luta o aguarda"},
				},
			},
			searchedString:   "palmeiras",
			expectedContains: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.ContainsStringByKey(tc.searchedString)

			if !reflect.DeepEqual(result, tc.expectedContains) {
				t.Errorf("ContainsStringByKey() = %v, want %v", result, tc.expectedContains)
			}
		})
	}
}

func TestRemoveStringByKey(t *testing.T) {
	testCases := []struct {
		name             string
		initialResource  Resources
		removedStringKey string
		expectedResource Resources
	}{
		{
			name: "Remove from Resources",
			initialResource: Resources{
				Strings: []String{{Key: "test_key", Value: "Test Value"}},
			},
			removedStringKey: "test_key",
			expectedResource: Resources{
				Strings: []String{},
			},
		},
		{
			name: "Remove into empty Resources",
			initialResource: Resources{
				Strings: []String{},
			},
			removedStringKey: "test_key",
			expectedResource: Resources{
				Strings: []String{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.RemoveStringByKey(tc.removedStringKey)

			if !reflect.DeepEqual(result, tc.expectedResource) {
				t.Errorf("RemoveNewString() = %v, want %v", result, tc.expectedResource)
			}
		})
	}
}

func TestCreateOrSubstituteStringByKey(t *testing.T) {
	testCases := []struct {
		name             string
		initialResource  Resources
		addedKey         string
		addedValue       string
		expectedResource Resources
	}{
		{
			name: "Add to empty Resources",
			initialResource: Resources{
				Strings: []String{},
			},
			addedKey:   "palmeiras",
			addedValue: "Quando surge o alviverde imponente",
			expectedResource: Resources{
				Strings: []String{
					{
						XMLName: xml.Name{Local: "string"},
						Key:     "palmeiras",
						Value:   "Quando surge o alviverde imponente",
					},
				},
			},
		},
		{
			name: "Add to non-empty Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Quando surge o alviverde imponente"},
				},
			},
			addedKey:   "palestra_italia",
			addedValue: "Cantando em coro a victoria",
			expectedResource: Resources{
				Strings: []String{
					{
						XMLName: xml.Name{Local: "string"},
						Key:     "palestra_italia",
						Value:   "Cantando em coro a victoria",
					},
					{
						Key:   "palmeiras",
						Value: "Quando surge o alviverde imponente",
					},
				},
			},
		},
		{
			name: "Substitute in Resources",
			initialResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Quando surge o alviverde imponente"},
				},
			},
			addedKey:   "palmeiras",
			addedValue: "Sabe bem o que vem pela frente",
			expectedResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Sabe bem o que vem pela frente"},
				},
			},
		},
		{
			name: "Substitute in Resources with duplicated values",
			initialResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Quando surge o alviverde imponente"},
				},
			},
			addedKey:   "palmeiras",
			addedValue: "Quando surge o alviverde imponente",
			expectedResource: Resources{
				Strings: []String{
					{Key: "palmeiras", Value: "Quando surge o alviverde imponente"},
				},
			},
		},
		{
			name: "Substitute in Resources with for one that contains another type of xml.Name",
			initialResource: Resources{
				Strings: []String{
					{
						XMLName: xml.Name{Local: "time"},
						Key:     "palmeiras",
						Value:   "Quando surge o alviverde imponente",
					},
				},
			},
			addedKey:   "palmeiras",
			addedValue: "No gramado em que a luta o aguarda",
			expectedResource: Resources{
				Strings: []String{
					{
						XMLName: xml.Name{Local: "time"},
						Key:     "palmeiras",
						Value:   "No gramado em que a luta o aguarda",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.initialResource.CreateOrSubstituteStringByKey(tc.addedKey, tc.addedValue)

			if !reflect.DeepEqual(result, tc.expectedResource) {
				t.Errorf("CreateOrSubstituteStringByKey() = %v, want %v", result, tc.expectedResource)
			}
		})
	}
}

func TestGetResourcesFromPathXML(t *testing.T) {
	testCases := []struct {
		name           string
		xmlContent     string
		expectedResult Resources
		expectError    bool
	}{
		{
			name: "Valid XML file",
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
	<resources>
	  <string name="app_name">Test App</string>
	  <string name="welcome_message">Welcome to the app!</string>
	</resources>`,
			expectedResult: Resources{
				XMLName: xml.Name{Local: "resources"},
				Strings: []String{
					{XMLName: xml.Name{Local: "string"}, Key: "app_name", Value: "Test App"},
					{XMLName: xml.Name{Local: "string"}, Key: "welcome_message", Value: "Welcome to the app!"},
				},
			},
			expectError: false,
		},
		{
			name:           "Empty XML file",
			xmlContent:     `<?xml version="1.0" encoding="UTF-8"?><resources></resources>`,
			expectedResult: Resources{XMLName: xml.Name{Local: "resources"}},
			expectError:    false,
		},
		{
			name:           "Invalid XML file",
			xmlContent:     `This is not valid XML content`,
			expectedResult: Resources{},
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_resources_*.xml")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.Write([]byte(tc.xmlContent)); err != nil {
				t.Fatalf("Failed to write to temporary file: %v", err)
			}
			tmpFile.Close()

			result, err := GetResourcesFromPathXML(tmpFile.Name())

			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("GetResourcesFromPathXML() = %v, want %v", result, tc.expectedResult)
			}
		})
	}

	t.Run("Non-existent file", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "non_existent_file.xml")
		_, err := GetResourcesFromPathXML(nonExistentPath)
		fmt.Println(err)
		if err == nil {
			t.Errorf("Expected an error for non-existent file, but got none")
		}
	})
}

func TestUpdateResourcesToXMLFile(t *testing.T) {
	testCases := []struct {
		name           string
		resources      Resources
		expectedOutput string
		expectError    bool
	}{
		{
			name: "Valid Resources",
			resources: Resources{
				XMLName: xml.Name{Local: "resources"},
				Strings: []String{
					{XMLName: xml.Name{Local: "string"}, Key: "app_name", Value: "Test App"},
					{XMLName: xml.Name{Local: "string"}, Key: "welcome_message", Value: "Welcome to the app!"},
				},
			},
			expectedOutput: `<resources>
    <string name="app_name">Test App</string>
    <string name="welcome_message">Welcome to the app!</string>
</resources>`,
			expectError: false,
		},
		{
			name: "Empty Resources",
			resources: Resources{
				XMLName: xml.Name{Local: "resources"},
			},
			expectedOutput: `<resources></resources>`,
			expectError:    false,
		},
		{
			name: "Invalid Characters in String Key",
			resources: Resources{
				XMLName: xml.Name{Local: "resources"},
				Strings: []String{{Key: "invalid<>key", Value: "Can't"}},
			},
			expectedOutput: `<resources>
    <string name="invalid&lt;&gt;key">Can't</string>
</resources>`,
			expectError: false,
		},
		{
			name: "Invalid Characters in String Value",
			resources: Resources{
				XMLName: xml.Name{Local: "resources"},
				Strings: []String{{Key: "valid_key", Value: "Can<>t"}},
			},
			expectedOutput: `<resources>
    <string name="valid_key">Can<>t</string>
</resources>`,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_resources_*.xml")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			err = tc.resources.UpdateResourcesToXMLFile(tmpFile.Name())

			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			content, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read temporary file: %v", err)
			}

			if string(content) != tc.expectedOutput {
				t.Errorf("File content does not match expected output.\nGot:\n%s\nWant:\n%s", content, tc.expectedOutput)
			}
		})
	}

	t.Run("Write Permission Error", func(t *testing.T) {
		resources := Resources{XMLName: xml.Name{Local: "resources"}}
		readOnlyPath := filepath.Join(os.TempDir(), "readonly_test.xml")

		err := os.WriteFile(readOnlyPath, []byte(""), 0o400)
		if err != nil {
			t.Fatalf("Failed to create read-only file: %v", err)
		}
		defer os.Remove(readOnlyPath)

		err = resources.UpdateResourcesToXMLFile(readOnlyPath)
		if err == nil {
			t.Errorf("Expected an error due to write permissions, but got none")
		}
	})
}

func TestIsSortedByKey(t *testing.T) {
	testCases := []struct {
		name      string
		resources Resources
		expected  bool
	}{
		{
			name: "Empty list",
			resources: Resources{
				Strings: []String{},
			},
			expected: true,
		},
		{
			name: "Single element",
			resources: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
				},
			},
			expected: true,
		},
		{
			name: "Sorted list",
			resources: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
			expected: true,
		},
		{
			name: "Unsorted list",
			resources: Resources{
				Strings: []String{
					{Key: "b", Value: "value2"},
					{Key: "a", Value: "value1"},
					{Key: "c", Value: "value3"},
				},
			},
			expected: false,
		},
		{
			name: "Duplicate keys in sorted order",
			resources: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "a", Value: "value2"},
					{Key: "b", Value: "value3"},
				},
			},
			expected: true,
		},
		{
			name: "Duplicate keys but unsorted values must consider sorted",
			resources: Resources{
				Strings: []String{
					{Key: "a", Value: "value2"},
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value3"},
				},
			},
			expected: true,
		},
		{
			name: "Duplicate keys in unsorted order",
			resources: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "a", Value: "value3"},
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.resources.IsSortedByKey()
			if result != tc.expected {
				t.Errorf("Expected IsSortedByKey() to return %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestSortByKey(t *testing.T) {
	testCases := []struct {
		name     string
		input    Resources
		expected Resources
	}{
		{
			name:     "Empty list",
			input:    Resources{Strings: []String{}},
			expected: Resources{Strings: []String{}},
		},
		{
			name: "Single element",
			input: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
				},
			},
			expected: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
				},
			},
		},
		{
			name: "Already sorted list",
			input: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
			expected: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
		},
		{
			name: "Unsorted list",
			input: Resources{
				Strings: []String{
					{Key: "c", Value: "value3"},
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
				},
			},
			expected: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
		},
		{
			name: "List with duplicate keys",
			input: Resources{
				Strings: []String{
					{Key: "b", Value: "value2"},
					{Key: "a", Value: "value1"},
					{Key: "c", Value: "value3"},
					{Key: "a", Value: "value4"},
				},
			},
			expected: Resources{
				Strings: []String{
					{Key: "a", Value: "value1"},
					{Key: "a", Value: "value4"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
		},
		{
			name: "List with duplicate keys should keep the order and dont consider Value",
			input: Resources{
				Strings: []String{
					{Key: "b", Value: "value2"},
					{Key: "a", Value: "value4"},
					{Key: "c", Value: "value3"},
					{Key: "a", Value: "value1"},
				},
			},
			expected: Resources{
				Strings: []String{
					{Key: "a", Value: "value4"},
					{Key: "a", Value: "value1"},
					{Key: "b", Value: "value2"},
					{Key: "c", Value: "value3"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.SortByKey()
			if !reflect.DeepEqual(tc.input, tc.expected) {
				t.Errorf("Expected sorted Resources to be %v, but got %v", tc.expected, tc.input)
			}
		})
	}
}
