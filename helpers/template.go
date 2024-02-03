package helpers

import (
	"bytes"
	"strings"
	"text/template"
)

// CompileTemplate compile template file
//
// @param name file name
// @param content file content
// @param data template data
// @param replacements string based replacement
func CompileTemplate(name, commentSymbol, content string, data map[string]string, replacements map[string]string) (string, error) {
	content = ResolvePlaceholders(content, commentSymbol, replacements)
	content = normalizeTemplate(content, commentSymbol)
	tpl, err := template.New(name).Delims(`{{`, `}}`).Parse(content)
	if err != nil {
		return "", err
	}
	var writer bytes.Buffer
	if err := tpl.Execute(&writer, data); err != nil {
		return "", err
	}
	if writer.Len() == 0 {
		return "", nil
	} else {
		return writer.String(), nil
	}
}

func ResolvePlaceholders(content, commentSymbol string, replacements map[string]string) string {
	oldNew := make([]string, 0)
	for search, replace := range replacements {
		oldNew = append(oldNew, search)
		oldNew = append(oldNew, replace)
	}
	oldNew = append(oldNew, commentSymbol+`- `, "", commentSymbol+`-`, "")
	return strings.NewReplacer(oldNew...).Replace(content)
}

func normalizeTemplate(content, commentSymbol string) string {
	return strings.NewReplacer(commentSymbol+` {{`, `{{`, commentSymbol+`{{`, `{{`).Replace(content)
}
