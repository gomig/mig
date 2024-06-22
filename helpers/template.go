package helpers

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gomig/crypto"
	"github.com/google/uuid"
)

var pipes = make(template.FuncMap, 0)

func init() {
	pipes["iif"] = func(cond bool, yes, no any) any {
		if cond {
			return yes
		}
		return no
	}

	pipes["oneOf"] = func(v any, valids string) bool {
		for _, valid := range strings.Split(valids, "|") {
			if valid == fmt.Sprint(v) {
				return true
			}
		}
		return false
	}

	pipes["uuid"] = func() string {
		return uuid.NewString()
	}

	pipes["key32"] = func() string {
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.MD5); err != nil {
			panic(err)
		} else {
			return key
		}
	}

	pipes["key64"] = func() string {
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.SHA3256); err != nil {
			panic(err)
		} else {
			return key
		}
	}

	pipes["key96"] = func() string {
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.SHA3384); err != nil {
			panic(err)
		} else {
			return key
		}
	}

	pipes["key128"] = func() string {
		if key, err := crypto.NewCryptography(uuid.New().String()).Hash(uuid.New().String(), crypto.SHA3512); err != nil {
			panic(err)
		} else {
			return key
		}
	}
}

// CompileTemplate compile template file
//
// @param name file name
// @param content file content
// @param data template data
// @param replacements string based replacement
func CompileTemplate(name, content string, data map[string]string, replacements map[string]string) (string, error) {
	content = ResolvePlaceholders(content, replacements)
	content = normalizeTemplate(content)
	tpl, err := template.
		New(name).
		Delims(`<%`, `%>`).
		Funcs(pipes).
		Parse(content)
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

func ResolvePlaceholders(content string, replacements map[string]string) string {
	oldNew := make([]string, 0)
	for search, replace := range replacements {
		oldNew = append(oldNew, search)
		oldNew = append(oldNew, replace)
	}
	oldNew = append(oldNew, `//- `, "", `//-`, "")
	return strings.NewReplacer(oldNew...).Replace(content)
}

func normalizeTemplate(content string) string {
	return strings.NewReplacer(`// <%`, "<%", `//<%`, "<%", `"-<%`, "<%", `%>-"`, "%>").Replace(content)
}
