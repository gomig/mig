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
		for _, ignore := range q.Ignores(answer) {
			mig.ignores = append(mig.ignores, helpers.NormalizePath(ignore))
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

// Compile compile file
func (mig *Mig) Compile(path string, content []byte) error {
	path = helpers.ResolvePlaceholders(path, "//", mig.Replacements())
	if v, err := helpers.CompileTemplate(
		path,
		"//",
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

// AddIgnore append new ignore path
func (mig *Mig) AddIgnore(path ...string) {
	for _, p := range path {
		mig.ignores = append(mig.ignores, helpers.NormalizePath(p))
	}
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
