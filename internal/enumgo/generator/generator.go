package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	etpl "github.com/byte4cat/nbx/v2/internal/enumgo/templates"
	"github.com/goccy/go-yaml"
)

type EnumValue struct {
	Name    string            `yaml:"name"`
	Value   any               `yaml:"value"` // support string and int
	Label   map[string]string `yaml:"label"`
	Comment string            `yaml:"comment"`
}

type EnumDefinition struct {
	Type        string      `yaml:"type"`
	Description string      `yaml:"description"`
	Values      []EnumValue `yaml:"values"`
}

func NewEnumFile(inputPath, outputDir, pkg string) error {
	yamlData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", inputPath, err)
	}

	var def EnumDefinition
	if err := yaml.Unmarshal(yamlData, &def); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %w", inputPath, err)
	}

	if len(def.Values) == 0 {
		return fmt.Errorf("no values found in %s", inputPath)
	}

	// type of the first value determines the template to use
	kind := reflect.TypeOf(def.Values[0].Value).Kind()

	var tpl *template.Template
	switch kind {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Float64:
		tpl = template.Must(template.New("enum").Funcs(template.FuncMap{
			"lower":                  strings.ToLower,
			"formatMultilineComment": formatMultilineComment,
		}).Parse(etpl.EnumTplInt))
	case reflect.String:
		tpl = template.Must(template.New("enum").Funcs(template.FuncMap{
			"lower":                  strings.ToLower,
			"formatMultilineComment": formatMultilineComment,
		}).Parse(etpl.EnumTplString))
	default:
		fmt.Fprintf(os.Stderr, "Unsupported enum value type %s in file %s\n", kind.String(), inputPath)
		os.Exit(1)
	}

	// filename = input file name without extension
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputFilePath := filepath.Join(outputDir, baseName+".go")

	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outputFilePath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file %s: %v\n", outputFilePath, err)
		}
	}()

	var tplOutput bytes.Buffer

	err = tpl.Execute(&tplOutput, map[string]any{
		"Package": pkg,
		"Def":     def,
	})
	if err != nil {
		return fmt.Errorf("error executing template for file %s: %w", inputPath, err)
	}

	formattedCode, err := formatCode(tplOutput.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting code: %w", err)
	}

	err = os.WriteFile(outputFilePath, formattedCode, 0644)
	if err != nil {
		return fmt.Errorf("error writing file %s: %w", outputFilePath, err)
	}

	fmt.Printf("Generated %s -> %s\n", inputPath, outputFilePath)

	return nil
}

func formatMultilineComment(comment string) string {
	// Split comment into lines
	lines := strings.Split(comment, "\n")
	var formattedComment []string

	// Iterate over each line and add "//" only to non-empty lines
	for _, line := range lines {
		// Skip empty lines to avoid adding extra "//"
		line = strings.TrimSpace(line)
		if line != "" {
			formattedComment = append(formattedComment, "// "+line)
		}
	}
	// Join all lines back into a single string with proper line breaks
	return strings.Join(formattedComment, "\n")
}
