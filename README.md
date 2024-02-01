# Installation

First, make sure you have Go properly installed and setup.

```bash
go install github.com/gomig/mig
```

## Create New Project

```bash
mig new myApp
```

Enter your repository and configure what you want. By default mig use [gomig/boilerplate](https://github.com/gomig/boilerplate) template.

## Create Project From Custom Branch

```bash
mig new myApp -b v1.2.3
```

## Create From Private Repository

For creating project from private github repository you need to store your github auth information. You must get [personal access token](https://github.com/settings/tokens) from your github account.

```bash
mig new myApp -a main
```

### Add New Credential Or Override

```bash
# Command
mig auth [key] [user] [access token]

# Usage
mig auth main johndoe asdfvzxcvq123asdfz
```

### Delete Credential

```bash
# Command
mig unauth [key]

# Usage
mig unauth main
```

### List Registered Credentials

```bash
mig users

# Key                           User
# +----------------------------+----
# main                          johndoe
```

## Template Repository Guide

Mig compile repository template base on [Go Text Template Library](https://pkg.go.dev/text/template). Each repository must contain `mig.json` configuration file in th root of repository.

### Extra Template Signs

Some times you want to comment code in template by default and compiled to code under some conditions.

- **`// {{`** This line translated to `{{` on template compile time. You could use template sign with comment.
- **`//-`** This comment sign remove from start of line on compile time.

```js
// Create new instance of {{ .appName }}
var app = new App("__APP__");

// This code only uncommented if web parameters equal y
// {{ if .web eq "y" }}
//- var client = new HttpClient();
//- client.initialize();
//- app.Client = client;
// {{ end }}
```

### Configuration

- **scripts:** A list of scripts to run after project template compile done.
- **conditions:** A list of questions to ask from user in creation time to customize template. Conditions answer will passed to output on.
  - **name:** Template key name. You can use parameter answer in your template by parameter name.
  - **desc:** Parameter description to show in creation time.
  - **default:** Default parameter value.
  - **placeholder:** A placeholder text to replace with parameter value.
  - **valids:** List of valid value for parameter.
  - **falsy:** False value of parameter.
  - **files:** List of files to ignore if parameter value equal `falsy` value.

```json
// mig.json
{
  "name": "Template Name",
  "intro": "Introduction text to show before start configuration",
  "message": "Final message to show after project creation",
  "scripts": [
    "go mod tidy",
    "go fmt ./...",
    "git init",
    "git add .",
    "git commit -m init"
  ],
  "conditions": [
    {
      "name": "appName",
      "desc": "app description",
      "default": "",
      "placeholder": "__APP__"
    },
    {
      "name": "web",
      "desc": "contain web part",
      "default": "n",
      "valids": ["y", "n"],
      "falsy": "n",
      "files": ["views", "src/web.js"]
    }
  ]
}
```
