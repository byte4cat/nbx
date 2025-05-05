package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yimincai/nbx/internal/enumgo/generator"
)

func TestEnumGeneration_Golden(t *testing.T) {
	tests := []struct {
		name       string
		yamlFile   string
		goldenFile string
	}{
		{
			name:       "fruit",
			yamlFile:   "fruit.yml",
			goldenFile: "fruit.go.golden",
		},
		{
			name:       "animal",
			yamlFile:   "animal.yml",
			goldenFile: "animal.go.golden",
		},
		{
			name:       "car",
			yamlFile:   "car.yml",
			goldenFile: "car.go.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			yamlPath := filepath.Join("..", "..", "..", "tests", "enumgo", "input", tt.yamlFile)
			expectedPath := filepath.Join("..", "..", "..", "tests", "enumgo", "golden", tt.goldenFile)
			outputFile := filepath.Join(tmpDir, tt.name+".go")

			generator.NewEnumFile(yamlPath, tmpDir, "example")

			out, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}
			exp, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("Failed to read golden file: %v", err)
			}

			assert.Equal(t, string(exp), string(out), "Generated file does not match golden file")
		})
	}
}
