package app

import (
	"encoding/json"
	"strings"

	"github.com/gomig/mig/helpers"
)

type MigConfig struct {
	Name      string     `json:"name"`
	Intro     string     `json:"intro"`
	Message   string     `json:"message"`
	Comment   string     `json:"comment"`
	Scripts   []string   `json:"scripts"`
	Questions []Question `json:"conditions"`
}

type Mig struct {
	answers  map[string]string
	compiled map[string]string
	ignores  []string
	config   MigConfig
}

// Init initialized instance. this method clean all existed data
func (mig *Mig) Init() {
	mig.answers = make(map[string]string)
	mig.compiled = make(map[string]string)
	mig.ignores = make([]string, 0)
	mig.config = MigConfig{}
	mig.config.Scripts = make([]string, 0)
	mig.config.Questions = make([]Question, 0)
}

// Parse config from json data
func (mig *Mig) Parse(data []byte) error {
	config := new(MigConfig)
	err := json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	mig.config = *config
	return nil
}

// Name get or set app name
func (mig *Mig) Name(v ...string) string {
	if len(v) > 0 {
		mig.config.Name = v[0]
	}
	return strings.TrimSpace(mig.config.Name)
}

// Intro get or set app intro
func (mig *Mig) Intro(v ...string) string {
	if len(v) > 0 {
		mig.config.Intro = v[0]
	}
	return strings.TrimSpace(mig.config.Intro)
}

// Message get or set app complete message
func (mig *Mig) Message(v ...string) string {
	if len(v) > 0 {
		mig.config.Message = v[0]
	}
	return strings.TrimSpace(mig.config.Message)
}

// Comment get or set app comment symbol
func (mig *Mig) Comment(v ...string) string {
	if len(v) > 0 {
		mig.config.Comment = v[0]
	}
	if strings.TrimSpace(mig.config.Comment) == "" {
		mig.config.Comment = `//`
	}
	return strings.TrimSpace(mig.config.Comment)
}

// AddScript append post script
func (mig *Mig) AddScript(s ...string) {
	mig.config.Scripts = append(mig.config.Scripts, s...)
}

// AddQuestion append question
func (mig *Mig) AddQuestion(q ...Question) {
	mig.config.Questions = append(mig.config.Questions, q...)
}

// AddIgnore append new ignore path
func (mig *Mig) AddIgnore(path ...string) {
	for _, p := range path {
		mig.ignores = append(mig.ignores, helpers.NormalizePath(p))
	}
}

// Find find question
func (mig Mig) Find(name string) *Question {
	for _, rule := range mig.config.Questions {
		if rule.Name == name {
			return &rule
		}
	}
	return nil
}

// Start question wizard
func (mig *Mig) Start() {
	mig.answers = make(map[string]string)
	for _, q := range mig.config.Questions {
		answer := q.Ask()
		if answer == q.Falsy {
			for _, i := range q.Files {
				mig.ignores = append(mig.ignores, helpers.NormalizePath(i))
			}
		}
		mig.answers[q.Name] = answer
	}
}

// ShouldIgnore check if path should ignore
func (mig Mig) ShouldIgnore(path string) bool {
	path = helpers.NormalizePath(path)
	for _, ignore := range mig.ignores {
		if path == ignore || helpers.IsPathOf(path, ignore) {
			return true
		}
	}
	return false
}

// Result get wizard result
func (mig Mig) Result(name string) string {
	return mig.answers[name]
}

// Answers get wizard results
func (mig Mig) Answers() map[string]string {
	return mig.answers
}

// Compiled get compiled results
func (mig Mig) Compiled() map[string]string {
	return mig.compiled
}

// Scripts get post scripts
func (mig *Mig) Scripts() []string {
	res := make([]string, 0)
	for _, script := range mig.config.Scripts {
		script = strings.TrimSpace(script)
		if script != "" {
			res = append(res, script)
		}
	}
	return res
}

// Replacements get replacements data
func (mig *Mig) Replacements() map[string]string {
	res := make(map[string]string)
	for _, q := range mig.config.Questions {
		if q.Placeholder != "" {
			res[q.Placeholder] = mig.Result(q.Name)
		}
	}
	return res
}

// Compile compile file
func (mig *Mig) Compile(path string, content []byte) error {
	if v, err := helpers.CompileTemplate(
		path,
		mig.Comment(),
		string(content),
		mig.answers,
		mig.Replacements(),
	); err != nil {
		return err
	} else {
		mig.compiled[path] = v
	}
	return nil
}