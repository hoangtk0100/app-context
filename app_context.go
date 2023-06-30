package appctx

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	devEnv = "dev"
	stgEnv = "stg"
	prdEnv = "prd"

	envFileKey     = "ENV_FILE"
	defaultEnvFile = ".env"
)

type AppContext interface {
	GetPrefix() string
	GetName() string
	GetEnvName() string
	Get(id string) (interface{}, bool)
	MustGet(id string) interface{}
	Load() error
	Stop() error
	Logger(prefix string) Logger
	OutEnv()
}

type Component interface {
	ID() string
	InitFlags()
	Run(AppContext) error
	Stop() error
}

type appContext struct {
	prefix     string
	name       string
	env        string
	store      map[string]Component
	components []Component
	cmd        *appFlagSet
	logger     Logger
}

func NewAppContext(opts ...Option) AppContext {
	app := &appContext{
		store: make(map[string]Component),
	}

	app.components = []Component{defaultLogger}

	for _, opt := range opts {
		opt(app)
	}

	app.initFlags()
	app.cmd = newAppFlagSet(app.prefix, app.name, pflag.CommandLine)
	app.parseFlags()

	app.logger = defaultLogger.GetLogger(formatLogPrefix(app.prefix, app.name))

	return app
}

func formatLogPrefix(prefix, name string) string {
	if prefix == "" {
		prefix = name
	}

	prefix = strings.ToLower(strings.TrimSpace(prefix))
	prefix = strings.Replace(prefix, ".", "-", -1)
	prefix = strings.Replace(prefix, " ", "-", -1)

	return prefix
}

func (ac *appContext) initFlags() {
	pflag.StringVar(
		&ac.env,
		"app-env",
		devEnv,
		"Env (dev | stg | prd)",
	)

	for _, c := range ac.components {
		c.InitFlags()
	}
}

func (ac *appContext) parseFlags() {
	// Parse all flags have been defined
	err := ac.cmd.flagSet.Parse(viper.AllKeys())
	if err != nil {
		log.Fatal().Msg("Cannot parse command flags")
	}

	// Bind the command flags to the configuration variables
	if err := viper.BindPFlags(ac.cmd.flagSet); err != nil {
		log.Fatal().Err(err).Msg("Error binding flags")
	}

	envFile := os.Getenv(envFileKey)
	if envFile == "" {
		envFile = defaultEnvFile
	}

	_, err = os.Stat(envFile)
	if err == nil {
		viper.SetConfigFile(envFile)
		viper.SetConfigType("dotenv")

		// Automatically overrides values with the value of corresponding environment variable if they exist
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msgf("Load env(%s) failed", envFile)
		}
	} else if envFile != defaultEnvFile {
		log.Fatal().Err(err).Msgf("Load env(%s) failed", envFile)
	}

	// Set ENV variables to flags value if not passing flags in command
	if err := ac.cmd.ParseSet(); err != nil {
		log.Fatal().Msg("Cannot set flags")
	}
}

func (ac *appContext) GetPrefix() string {
	return ac.prefix
}

func (ac *appContext) GetName() string {
	return ac.name
}

func (ac *appContext) GetEnvName() string {
	return ac.env
}

func (ac *appContext) Get(id string) (interface{}, bool) {
	c, ok := ac.store[id]
	if ok {
		return c, true
	}

	return nil, false
}

func (ac *appContext) MustGet(id string) interface{} {
	c, ok := ac.Get(id)
	if ok {
		return c
	}

	panic(fmt.Sprintf("Cannot get %s\n", id))
}

func (ac *appContext) Load() error {
	for _, c := range ac.components {
		if err := c.Run(ac); err != nil {
			return err
		}
	}

	ac.logger.Info("Service context loaded")

	return nil
}

func (ac *appContext) Stop() error {
	for i := range ac.components {
		if err := ac.components[i].Stop(); err != nil {
			return err
		}
	}

	ac.logger.Info("Service context stopped")

	return nil
}

func (ac *appContext) Logger(prefix string) Logger {
	return defaultLogger.GetLogger(prefix)
}

func (ac *appContext) OutEnv() {
	ac.cmd.GetSampleEnvs()
}
