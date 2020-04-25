package berus

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CommandLineKey struct {
	Command *cobra.Command
	Key     string
}

type PFlagConfiguration struct {
	Bindings map[string]*pflag.Flag
}

func (c PFlagConfiguration) Initialize(viper *viper.Viper) error {
	for cfgKey, flag := range c.Bindings {
		if err := viper.BindPFlag(cfgKey, flag); nil != err {
			return fmt.Errorf("bind flags {cli:%s,cfg:%s}: %w", flag, cfgKey, err)
		}
	}
	return nil
}
