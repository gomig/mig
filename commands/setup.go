package commands

import (
	"fmt"
	"os"

	"github.com/gomig/crypto"
	"github.com/gomig/mig/internal/helpers"
	"github.com/gomig/mig/internal/questions"
	"github.com/google/uuid"
)

func setup(name string, w *questions.Wizard) {
	// helpers
	pathResolver := func(p string) string {
		return fmt.Sprintf("./%s/%s", name, p)
	}

	// init global data
	data := make(helpers.TemplateData)
	data["name"] = name
	data["description"] = w.Result("description")
	data["namespace"] = w.Result("namespace")
	data["locale"] = w.Result("locale")
	data["config"] = w.Result("config")
	data["cache"] = w.Result("cache")
	data["database"] = w.Result("database")
	data["translator"] = w.Result("translator")
	data["web"] = w.Result("web")
	data["template"] = w.Result("template")

	// set app key
	c := crypto.NewCryptography(uuid.New().String())
	appKey, err := c.Hash(uuid.New().String(), crypto.SHA3256)
	helpers.Handle(err)
	data["appKey"] = appKey
	mcKey, err := c.Hash(uuid.New().String(), crypto.SHA3256)
	helpers.Handle(err)
	data["mcKey"] = mcKey

	// Clean and compile
	helpers.Handle(helpers.CompileTemplate(pathResolver("go.mod"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("main.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("app.tpl.service"), data))
	helpers.Handle(os.Rename(pathResolver("app.service"), pathResolver(name+".service")))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/boot.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/cache.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/config.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/crypto.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/helpers.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/logger.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/translator.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/utils.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/validator.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/clear.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/down.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/hash.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/up.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/version.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/config/strings.go"), data))

	// config
	switch w.Result("config") {
	case "env":
		helpers.Handle(helpers.CompileTemplate(pathResolver("config/config.tpl.env"), data))
	case "json":
		helpers.Handle(helpers.CompileTemplate(pathResolver("config/config.tpl.json"), data))
	}
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/config/config.tpl.go"), data))
	if w.Result("config") == "memory" && w.Result("translator") == "memory" {
		os.RemoveAll(pathResolver("config"))
	}

	// cache
	if w.Result("cache") != "redis" {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/cleanup.tpl.go"), data))
	}

	// translator
	if w.Result("translator") == "memory" {
		os.RemoveAll(pathResolver("config/strings"))
	} else {
		os.Rename(pathResolver("config/strings/locale"), pathResolver("config/strings/")+w.Result("locale"))
	}

	// database
	if w.Result("database") != "mysql" {
		os.RemoveAll(pathResolver("database"))
	}
	if w.Result("database") == "mysql" {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/database.tpl.go"), data))
	}
	if w.Result("database") == "mongo" {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/mongo.tpl.go"), data))
	}

	// web
	if w.Result("web") == "n" {
		os.RemoveAll(pathResolver("public"))
		os.RemoveAll(pathResolver("src/http"))
	} else {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/serve.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/web.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/http/controllers.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/http/errors.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/http/middlewares.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/http/routes.tpl.go"), data))
	}

	// template
	if w.Result("template") == "n" {
		os.RemoveAll(pathResolver("views"))
	} else {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/app/template.tpl.go"), data))
	}
}
