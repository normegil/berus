package berus

import (
	"github.com/spf13/viper"
	"os"
)

// FileConfiguration will be used to setup a file based configuration loading.
type FileConfiguration struct {
	// Name is the searched file name
	Name      string
	// Extension is the searched file type
	Extension string
	// Paths is a list of paths to search. Last path will tke precedence of previous path, first as the lowest priority.
	Paths     []string
}

// NewDefaultFileConfiguration load a default configuration, where application name is the file name
// and confifuration file will be a yaml file
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
