# Installation

First, make sure you have Go properly installed and setup.

```bash
git clone https://github.com/gomig/mig
cd gomig
go install .
```

## Create New Project

```bash
gomig new myApp
```

And configure what you want!

## Usage

For library usage see [gomig](https://github.com/gomig) libraries docs.

### Access App Dependencies

For accessing app dependencies (config driver, cache driver, etc.) you must use `app` namespace functions.

```go
// github.com/myapp is your app namespace
import "github.com/myapp/src/app"
app.Config()
```
