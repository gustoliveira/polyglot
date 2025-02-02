package internal

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestCheckCurrentDirectoryIsAndroidProject(t *testing.T) {
	tests := []struct {
		name       string
		indicators []string
		want       bool
	}{
		{
			name:       "No Indicators",
			indicators: []string{},
			want:       false,
		},
		{
			name:       "Has build.gradle",
			indicators: []string{"build.gradle"},
			want:       true,
		},
		{
			name:       "Has settings.gradle",
			indicators: []string{"settings.gradle"},
			want:       true,
		},
		{
			name:       "Has settings.gradle.kts",
			indicators: []string{"settings.gradle.kts"},
			want:       true,
		},
		{
			name:       "Has app/",
			indicators: []string{"app/"},
			want:       true,
		},
		{
			name:       "Has multiples indicators",
			indicators: []string{"build.gradle", "settings.gradle"},
			want:       true,
		},
		{
			name:       "Has multiples indicators but one that is not",
			indicators: []string{"palmeiras", "build.gradle", "settings.gradle"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, indicator := range tt.indicators {
				path := filepath.Join(tmpDir, indicator)

				if len(indicator) > 0 && indicator[len(indicator)-1] == '/' {
					os.MkdirAll(path, 0o755)
				} else {
					os.WriteFile(path, []byte("test"), 0o644)
				}
			}

			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			got := CheckCurrentDirectoryIsAndroidProject()
			if got != tt.want {
				t.Errorf("CheckCurrentDirectoryIsAndroidProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlockIfNotAndroidProject(t *testing.T) {
	tests := []struct {
		name       string
		setupFiles []string
		wantExit   int // Expected exit code (0 = no exit, 1 = exited)
	}{
		{
			name:       "Valid Android Project With build.gradle",
			setupFiles: []string{"build.gradle"},
			wantExit:   0,
		},
		{
			name:       "Valid Android Project With app/",
			setupFiles: []string{"app/"},
			wantExit:   0,
		},
		{
			name:       "Valid directory no mix of indicators",
			setupFiles: []string{"palmeiras", "build.gradle"},
			wantExit:   0,
		},
		{
			name:       "Invalid directory no indicators",
			setupFiles: []string{},
			wantExit:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for _, file := range tt.setupFiles {
				path := filepath.Join(tmpDir, file)
				if filepath.Ext(file) == "/" {
					if err := os.MkdirAll(path, 0o755); err != nil {
						t.Fatalf("failed to create directory %s: %v", path, err)
					}
					continue
				}

				if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
					t.Fatalf("failed to create file %s: %v", path, err)
				}
			}

			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			exitCode := 0

			error := BlockIfNotAndroidProject()

			if error != nil {
				exitCode = 1
			}

			if exitCode != tt.wantExit {
				t.Errorf("expected exit code %d, got %d", tt.wantExit, exitCode)
			}
		})
	}
}

func TestGetTranslationFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expected     Translation
		expectLocale string
		expectRegion string
		expectError  bool
	}{
		{
			name: "Valid values-en directory",
			path: "res/values-en/strings.xml",
			expected: Translation{
				Path:       "res/values-en/strings.xml",
				LocaleCode: "en",
				RegionCode: "",
				Language:   "English",
			},
			expectLocale: "en",
			expectRegion: "",
			expectError:  false,
		},
		{
			name: "Valid values-es directory with region",
			path: "res/values-es-rMX/strings.xml",
			expected: Translation{
				Path:       "res/values-es-rMX/strings.xml",
				LocaleCode: "es",
				RegionCode: "MX",
				Language:   "Mexican Spanish",
			},
			expectLocale: "es",
			expectRegion: "MX",
			expectError:  false,
		},
		{
			name: "Fallback to default locale (en)",
			path: "res/values/strings.xml",
			expected: Translation{
				Path:       "res/values/strings.xml",
				LocaleCode: "en",
				RegionCode: "",
				Language:   "English",
			},
			expectLocale: "en",
			expectRegion: "",
			expectError:  false,
		},
		{
			name: "Valid values-nn-rNO directory",
			path: "res/values-nn-rNO/strings.xml",
			expected: Translation{
				Path:       "res/values-nn-rNO/strings.xml",
				LocaleCode: "nn",
				RegionCode: "NO",
				Language:   "Norwegian Nynorsk",
			},
			expectLocale: "nn",
			expectRegion: "NO",
			expectError:  false,
		},
		{
			name:         "Invalid language code",
			path:         "res/values-asdf/strings.xml",
			expected:     Translation{},
			expectLocale: "",
			expectRegion: "",
			expectError:  true,
		},
		{
			name:         "Invalid region code",
			path:         "res/values-pt-rBrasilMeuBrasilBrasileiro/strings.xml",
			expected:     Translation{},
			expectLocale: "",
			expectRegion: "",
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetTranslationFromFileName(tc.path)

			if tc.expectError && err == nil {
				t.Errorf("Expected error, got nil - Translation %v", got)
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if got.Path != tc.expected.Path {
				t.Errorf("Path mismatch. Got %v, want %v", got.Path, tc.expected.Path)
			}
			if got.LocaleCode != tc.expected.LocaleCode {
				t.Errorf("Locale mismatch. Got %v, want %v", got.LocaleCode, tc.expected.LocaleCode)
			}
			if got.RegionCode != tc.expected.RegionCode {
				t.Errorf("Region mismatch. Got %v, want %v", got.RegionCode, tc.expected.RegionCode)
			}
			if got.Language != tc.expected.Language {
				t.Errorf("Language mismatch. Got %v, want %v", got.Language, tc.expected.Language)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name       string
		dirName    string
		wantLocale string
		wantRegion string
	}{
		{
			name:       "Valid locale only (en)",
			dirName:    "values-en",
			wantLocale: "en",
			wantRegion: "",
		},
		{
			name:       "Valid locale and region (es-rMX)",
			dirName:    "values-es-rMX",
			wantLocale: "es",
			wantRegion: "MX",
		},
		{
			name:       "No locale or region (default values)",
			dirName:    "values",
			wantLocale: "",
			wantRegion: "",
		},
		{
			name:       "Invalid random string",
			dirName:    "endrick",
			wantLocale: "",
			wantRegion: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotLocale, gotRegion := extract(tc.dirName)

			if gotLocale != tc.wantLocale {
				t.Errorf("Locale mismatch. Got %v, want %v", gotLocale, tc.wantLocale)
			}
			if gotRegion != tc.wantRegion {
				t.Errorf("Region mismatch. Got %v, want %v", gotRegion, tc.wantRegion)
			}
		})
	}
}

func TestGetTranslationsFromAllModules(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string
		expectedTrans []Translation
		expectError   bool
		errorContains string
	}{
		{
			name: "Single module with multiple translations",
			setupFiles: map[string]string{
				"app/src/main/res/values/strings.xml":    `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
				"app/src/main/res/values-es/strings.xml": `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
				"app/src/main/res/values-pt/strings.xml": `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
			},
			expectedTrans: []Translation{
				{
					Path:       "app/src/main/res/values/strings.xml",
					Language:   "English",
					LocaleCode: "en",
					RegionCode: "",
				},
				{
					Path:       "app/src/main/res/values-es/strings.xml",
					Language:   "Spanish",
					LocaleCode: "es",
					RegionCode: "",
				},
				{
					Path:       "app/src/main/res/values-pt/strings.xml",
					Language:   "Portuguese",
					LocaleCode: "pt",
					RegionCode: "",
				},
			},
			expectError: false,
		},
		{
			name: "Multiple modules with translations",
			setupFiles: map[string]string{
				"app/src/main/res/values/strings.xml":     `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
				"feature/src/main/res/values/strings.xml": `<?xml version="1.0" encoding="utf-8"?><resources></resources>`,
			},
			expectedTrans: []Translation{
				{
					Path:       "app/src/main/res/values/strings.xml",
					Language:   "English",
					LocaleCode: "en",
					RegionCode: "",
				},
				{
					Path:       "feature/src/main/res/values/strings.xml",
					Language:   "English",
					LocaleCode: "en",
					RegionCode: "",
				},
			},
			expectError: false,
		},
		{
			name:          "No resource directories found",
			setupFiles:    map[string]string{},
			expectedTrans: nil,
			expectError:   true,
			errorContains: "no android resource directories found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get current directory: %v", err)
			}
			defer os.Chdir(oldDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(tmpDir, path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
				if err != nil {
					t.Fatalf("failed to create directories: %v", err)
				}
				err = os.WriteFile(fullPath, []byte(content), 0o644)
				if err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}

			result, err := GetTranslationsFromAllModules()

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			t.Logf("LOG: %v", result)
			if len(result) != len(tt.expectedTrans) {
				t.Errorf("got %d translations, want %d", len(result), len(tt.expectedTrans))
				return
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i].Path < result[j].Path
			})
			sort.Slice(tt.expectedTrans, func(i, j int) bool {
				return tt.expectedTrans[i].Path < tt.expectedTrans[j].Path
			})

			for i := range result {
				if result[i].Language != tt.expectedTrans[i].Language {
					t.Errorf("translation[%d].Language = %v, want %v", i, result[i].Language, tt.expectedTrans[i].Language)
				}
				if result[i].LocaleCode != tt.expectedTrans[i].LocaleCode {
					t.Errorf("translation[%d].LocaleCode = %v, want %v", i, result[i].LocaleCode, tt.expectedTrans[i].LocaleCode)
				}
				if result[i].RegionCode != tt.expectedTrans[i].RegionCode {
					t.Errorf("translation[%d].RegionCode = %v, want %v", i, result[i].RegionCode, tt.expectedTrans[i].RegionCode)
				}
				if !strings.HasSuffix(result[i].Path, tt.expectedTrans[i].Path) {
					t.Errorf("translation[%d].Path = %v, want suffix %v", i, result[i].Path, tt.expectedTrans[i].Path)
				}
			}
		})
	}
}
