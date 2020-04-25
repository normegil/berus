package berus

import (
	"github.com/spf13/viper"
	"os"
)

type FileConfiguration struct {
	Name      string
	Extension string
	Paths     []string
}

func NewDefaultFileConfiguration(applicationName string) *FileConfiguration {
	return &FileConfiguration{
		Name:      applicationName,
		Extension: "yaml",
		Paths: []string{
			"/etc/" + applicationName,
			"$XDG_CONFIG_HOME" + string(os.PathSeparator) + applicationName,
			"$HOME" + string(os.PathSeparator) + "." + applicationName,
			".",
		},
	}
}

func (c FileConfiguration) Initialize(viper *viper.Viper) error {
	if nil != c.Paths {
		for _, path := range c.Paths {
			viper.AddConfigPath(path)
		}
	}

	viper.SetConfigType(c.Extension)
	viper.SetConfigName(c.Name)
	return nil
}
