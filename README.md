# Installation

First, make sure you have Go installed.

```bash
go install github.com/gomig/mig
```

## Update

Update cli tools to latest version.

```bash
mig update
```

## Create New Project

```bash
mig new myApp
```

Enter your repository and configure what you want. By default mig use [gomig/boilerplate](https://github.com/gomig/boilerplate) template.

## Create Project From Custom Branch Or Tag

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

# Output
┌───┬──────────────────┬──────────┐
│ # │ AUTHENTICATE KEY │ USERNAME │
├───┼──────────────────┼──────────┤
│ 1 │ main             │ johndoe  │
└───┴──────────────────┴──────────┘
```

## Template Repository Guide

Mig compile repository template base on [Go Text Template Library](https://pkg.go.dev/text/template) with `<% %>` delimeter. Each repository must contain `mig.json` configuration file in th root of repository.

**NOTE:** You can use user entered app name with `.name` variable or `__name__` placeholder in your template syntax.

**NOTE:** You can use a unique project key with `.key` variable or `__key__` placeholder. Each time you create new project mig generate new key for your project.

**NOTE:** You can use a unique project access token with `.token` variable or `__token__` placeholder. Each time you create new project mig generate new access token for your project.

**NOTE:** You can use a unique project identifier (uuid) with `.uuid` variable or `__uuid__` placeholder. Each time you create new project mig generate new uuid for your project.

### Extra Template Signs

Sometimes you want to comment code in template by default and compiled to code under some conditions.

- **`// <%`** This line translated to `<%` on compile time. You could use template sign with comment.
- **`//-`** This comment sign remove from start of line on compile time.

```js
// Create new instance of <% .appName %>
var app = new App("__APP__");

// This code only uncommented if web parameters equal y
// <% if eq .web "y" %>
//- var client = new HttpClient();
//- client.initialize();
//- app.Client = client;
// <% end %>
```

### Configuration

To configure your project template you need to put `mig.json` file in the root of your project. All configuration file `rules` value accessable in template syntax with rule name.

- **name:** Template name.
- **intro:** Introduction text to show before start configuration.
- **message:** Final message to show after project creation.
- **rules:** List of question to ask from user to configure template.

  - **name:** Name of rule. You can use rule answer in your template with rule name (`.web`, `.app_name`, ...).
  - **default:** Default rule value if user leave answer empty.
  - **placeholder:** A placeholder text to replace with rule answer. You can use placeholder in file name or string quotes.
  - **desc:** Rule description.
  - **options:** List of valid options for rule. If options not defined any text allowed.
  - **files:** File proccess conditions. You can define which file must included on final result based on rule value.

- **statics:** List of static assets to not compile with template engine.
- **ignores:** List of repository files to ignore.
- **scripts:** List of scripts to run after project create. Script must splited to array of args.

```json
{
  "name": "My Template",
  "intro": "You need git and node installed on your system.",
  "message": "Visit config/config.json to configure app port and information",
  "rules": [
    {
      "name": "app_name",
      "placeholder": "__APP__",
      "desc": "application name"
    },
    {
      "name": "web",
      "default": "n",
      "desc": "include web assets",
      "options": ["y", "n"],
      "files": {
        "n": ["console.md"], // console.md only included if web answer was n
        "y": ["web.md", "web.js", "views", "public", "src/app.js"] // this files and directory not listed if answer was n
      }
    },
    {
      "name": "locale",
      "default": "en",
      "placeholder": "__LOCALE__"
    }
  ],
  "statics": ["public/favicon.ico", "assets/icons"],
  "ignores": ["readme.md"],
  "scripts": [
    ["npm", "install"],
    ["git", "init"],
    ["git", "add", "."],
    ["git", "commit", "-m", "'initialize app'"]
  ]
}
```
