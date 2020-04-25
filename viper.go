package berus

import (
	"fmt"
	"github.com/spf13/viper"
)

type Initializer interface {
	Initialize(viper *viper.Viper) error
}

type Configuration struct {
	viper        *viper.Viper
	ForcedFile   string
	Initializers []Initializer
}

func NewConfiguration(initializers []Initializer) *Configuration {
	return &Configuration{
		viper:        viper.New(),
		ForcedFile:   "",
		Initializers: initializers,
	}
}

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

func (c *Configuration) GetViperInstance() *viper.Viper {
	return c.viper
}
