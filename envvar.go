package berus

import (
	"github.com/spf13/viper"
	"strings"
)

// EnvironmentVariableConfiguration allow to configure environment variable configuration for an application
type EnvironmentVariableConfiguration struct {
	Prefix string
}

func NewDefaultEnvironmentVariableConfiguration(applicationName string) *EnvironmentVariableConfiguration {
	return &EnvironmentVariableConfiguration{
		Prefix: strings.ToUpper(applicationName) + "_",
	}
}

func (c EnvironmentVariableConfiguration) Initialize(viper *viper.Viper) error {
	viper.SetEnvPrefix(c.Prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	return nil
}
