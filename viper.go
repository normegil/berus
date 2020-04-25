package berus

import (
	"fmt"
	"github.com/spf13/viper"
)

// Initialize represent a object that will initialize some viper options for the given instance.
type Initializer interface {
	Initialize(viper *viper.Viper) error
}

// Configuration will load viper configuration using different combinaisons of Initializers
type Configuration struct {
	viper        *viper.Viper
	// Option to force the usage of a specific file. Leave empty to ignore.
	ForcedFile   string
	// Initializers will initialize options of viper for some type of configuration (file, envvar, cli, ...)
	Initializers []Initializer
}

func NewConfiguration(initializers []Initializer) *Configuration {
	return &Configuration{
		viper:        viper.New(),
		ForcedFile:   "",
		Initializers: initializers,
	}
}

// ReadConfiguration will initialize viper configuration and read values to populate the initialized viper instance.
func (c *Configuration) ReadConfiguration() error {
	if err := c.Initialize(c.viper); nil != err {
		return err
	}

	if err := c.viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); !isNotFound {
			return fmt.Errorf("reading configuration: %w", err)
		}
	}
	return nil
}

// Initialize will setup given viper instance with all registered initializers
func (c Configuration) Initialize(vipercfg *viper.Viper) error {
	if c.ForcedFile != "" {
		vipercfg.SetConfigFile(c.ForcedFile)
	}
	for _, initializer := range c.Initializers {
		if err := initializer.Initialize(vipercfg); nil != err {
			return fmt.Errorf("initialization: %w", err)
		}
	}
	return nil
}

// Return viper instance associated with this Configuration
func (c *Configuration) ViperInstance() *viper.Viper {
	return c.viper
}
