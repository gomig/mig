package app

import (
	"encoding/json"
	"strings"

	"github.com/gomig/mig/helpers"
)

type MigConfig struct {
	Name    string     `json:"name"`
	Intro   string     `json:"intro"`
	Message string     `json:"message"`
	Rules   []Rule     `json:"rules"`
	Scripts [][]string `json:"scripts"`
	Statics []string   `json:"statics"`
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
	mig.config.Rules = make([]Rule, 0)
	mig.config.Scripts = make([][]string, 0)
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

// Start question wizard
func (mig *Mig) Start() {
	mig.answers = make(map[string]string)
	for _, q := range mig.config.Rules {
		answer := q.Ask()
		mig.ignores = append(mig.ignores, q.Ignores(answer)...)
		mig.answers[q.Name] = answer
	}
}

// ShouldIgnore check if path should ignore
func (mig Mig) ShouldIgnore(path string) bool {
	path = helpers.NormalizePath(path)
	for _, ignore := range mig.ignores {
		if helpers.IsPathOf(path, ignore) {
			return true
		}
	}
	return false
}

// ShouldCompile check if path should compile
func (mig Mig) ShouldCompile(path string) bool {
	path = helpers.NormalizePath(path)
	for _, static := range mig.config.Statics {
		if helpers.IsPathOf(path, static) {
			return true
		}
	}
	return false
}

// Compile compile file
func (mig *Mig) Compile(path string, content []byte) error {
	realPath := helpers.ResolvePlaceholders(path, "//", mig.Replacements())
	if mig.ShouldCompile(path) {
		if v, err := helpers.CompileTemplate(
			path,
			"//",
			string(content),
			mig.answers,
			mig.Replacements(),
		); err != nil {
			return err
		} else {
			mig.compiled[realPath] = v
		}
	} else {
		mig.compiled[realPath] = string(content)
	}
	return nil
}

// Getter and setters

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

// AddRule append rule
func (mig *Mig) AddRule(q ...Rule) {
	mig.config.Rules = append(mig.config.Rules, q...)
}

// AddScript append post script
func (mig *Mig) AddScript(s ...string) {
	mig.config.Scripts = append(mig.config.Scripts, s)
}

// AddStatic append static path
func (mig *Mig) AddStatic(path ...string) {
	mig.config.Statics = append(mig.config.Statics, path...)
}

// AddIgnore append new ignore path
func (mig *Mig) AddIgnore(path ...string) {
	mig.ignores = append(mig.ignores, path...)
}

// AddAnswer append new answer path
func (mig *Mig) AddAnswer(name, v string) {
	mig.answers[name] = v
}

// Scripts get post scripts
func (mig *Mig) Scripts() [][]string {
	res := make([][]string, 0)
	for _, script := range mig.config.Scripts {
		if len(script) > 0 {
			res = append(res, script)
		}
	}
	return res
}

// Replacements get replacements data
func (mig *Mig) Replacements() map[string]string {
	res := make(map[string]string)
	for _, q := range mig.config.Rules {
		if q.Placeholder != "" {
			res[q.Placeholder] = mig.Result(q.Name)
		}
	}
	return res
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
