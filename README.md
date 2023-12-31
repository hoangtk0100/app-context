# App-Context

App Context is a library that manages common components such as configurations, database connections, caching, and more.

It simplifies the implementation of new services by abstracting away the complexities of component configurations,
allowing developers to focus on building functionality rather than dealing with intricate configuration details.

It was inspired by [go-micro](https://github.com/go-micro/go-micro) and [service-context](https://github.com/viettranx/service-context)

It provides the following features:

- Logger component using [zerolog](https://github.com/rs/zerolog).
- Dynamic management of environment variables and flag variables using [viper](https://github.com/spf13/viper) and
  it's [pflag](https://github.com/spf13/pflag) package (viper supports multiple configuration file formats and reading from
  remote config systems (etcd or Consul), and watching changes, ...).
- Ability to output environment variables and flag variables in `.env` format.
- Easy integration of additional components as plugins.

## Features

1. Base:

- [x] Context
- [x] Logger
- [x] Component

2. Components:

	- 2.1. Token
		- [x] JWT
		- [x] PASETO

	- 2.2. Server
		- [x] Gin
		- [x] gRPC

	- 2.3. Client
		- [x] gRPC

	- 2.4. Datastore
		- [x] GORM
			- [x] PostgreSQL
			- [x] MySQL
			- [x] MSSQL
			- [x] SQLite
		- [x] Redis

	- 2.5. Storage
		- [x] AWS S3
		- [x] Cloudflare R2

	- 2.6. Mail

	- 2.7. Cache
		- [x] Local
		- [x] Redis

	- 2.8. Pub/sub
		- [x] Local
		- [x] NATS

	- 2.9. Database migration
		- [x] PostgreSQL
		- [x] MySQL
		- [x] MSSL
		- [x] MongoDB

3. Utils:
- [x] Password
- [x] Job
- [x] UID


## Examples

- [Demo component](./examples/component)
- [Demo CLI](./examples/cli) (runs the service and outputs all environment variables)

## How to use it

### 1. Installation:

```shell
go get -u github.com/hoangtk0100/app-context
```

### 2. Define your component:

Your component can be anything but implements this interface:

```go
type Component interface {
	ID() string
	InitFlags()
	Run(AppContext) error
	Stop() error
}
```

Demo custom component:

```go
package main

import (
	"github.com/spf13/pflag"

	appctx "github.com/hoangtk0100/app-context"
)

type demoComponent struct {
	id     string
	data   string
	logger appctx.Logger
}

func NewDemoComponent(id string) *demoComponent {
	return &demoComponent{id: id}
}

func (c *demoComponent) ID() string {
	return c.id
}

func (c *demoComponent) InitFlags() {
	pflag.StringVar(&c.data, "component-data", "demo", "Data string")
}

func (c *demoComponent) Run(ac appctx.AppContext) error {
	c.logger = ac.Logger(c.id)
	return nil
}

func (c *demoComponent) Stop() error {
	return nil
}

func (c *demoComponent) GetData() string {
	return c.data
}

func (c *demoComponent) DoSomething() error {
	c.logger.Print("LOL")
	return nil
}
```

### 3. Use the component with App-Context:

```go
package main

import (
	appctx "github.com/hoangtk0100/app-context"
)

func main() {
	const cmpId = "abc"
	appCtx := appctx.NewAppContext(
		appctx.WithName("Demo Component"),
		appctx.WithComponent(NewDemoComponent(cmpId)),
	)

	log := appCtx.Logger("service")

	if err := appCtx.Load(); err != nil {
		log.Error(err)
	}

	type CanDoSomething interface {
		GetData() string
		DoSomething() error
	}

	cmp := appCtx.MustGet(cmpId).(CanDoSomething)

	log.Print(cmp.GetData())
	_ = cmp.DoSomething()

	_ = appCtx.Stop()
}
```

### 4. Run your code with ENV

Option 1: Command Line

```shell
go build -o app
COMPONENT_DATA="Hi There" ./app
```

Option 2: Environment Variable File (You should do it in a new terminal)

```shell
# Create a file named .env with the following content:
COMPONENT_DATA="Hi There"

# Run the application
./app
```

Option 3: Environment Variable File (with custom name - You should do it in a new terminal)

```shell
# Create a file named .env.dev with the following content:
COMPONENT_DATA="Hi There"

# Set the environment variable pointing to the file
ENV_FILE=".env.dev"

# Run the application
./app
```

You will see this row on your console.

```
## Case: Only run "./app"
{"level":"debug","prefix":"core.service","time":"2023-06-22T18:18:13+07:00","message":"demo"}

## Case: Use 1 in 3 options above
{"level":"debug","prefix":"core.service","time":"2023-06-22T18:21:35+07:00","message":"Hi There"}
```
