package berus

import (
	"github.com/spf13/viper"
	"strings"
)

// EnvironmentVariableConfiguration allow to configure environment variable configuration for an application.
// Environment variables are base on requested configuration keys where "." are replaced by "_"
type EnvironmentVariableConfiguration struct {
	// Prefix will be used as prefix for environment variables
	Prefix string
}

// NewDefaultEnvironmentVariableConfiguration create a EnvironmentVariableConfiguration where prefix is based on
// application name ending with an underscore
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
