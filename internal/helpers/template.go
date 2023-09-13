package helpers

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData structure
type TemplateData map[string]any

// CompileTemplate compile template file
func CompileTemplate(filePath string, maps ...TemplateData) error {
	var data TemplateData
	if len(maps) == 0 {
		data = make(TemplateData)
	} else {
		data = maps[0]
		for i, m := range maps {
			if i == 0 {
				continue
			}
			for k, v := range m {
				data[k] = v
			}
		}
	}

	// Read file
	var writer bytes.Buffer
	fBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(fBytes)
	content = setNamespace(content, data["namespace"].(string), data["name"].(string))
	dest := strings.ReplaceAll(filePath, ".tpl.", ".")

	tpl, err := template.New(filepath.Base(filePath)).Delims(`// {{`, "}}").Parse(content)
	if err != nil {
		return err
	}
	if err := tpl.Execute(&writer, data); err != nil {
		return err
	}
	return os.WriteFile(dest, writer.Bytes(), 0644)
}

func setNamespace(code, namemspace, name string) string {
	res := code
	replacer := strings.NewReplacer("mekramy/__boiler", namemspace+"/"+name)
	return replacer.Replace(res)
}
