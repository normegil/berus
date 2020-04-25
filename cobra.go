package berus

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type CobraConfiguration struct {
	ApplicationName string
	RootCommand     *cobra.Command
	Bindings        []CobraBinding
}

func NewCobraDefaultConfiguration(cfg CobraConfiguration) (*Configuration, error) {
	bindings, err := CobraBindings(cfg.RootCommand, cfg.Bindings)
	if err != nil {
		return nil, fmt.Errorf("loading cobra bindings: %w", err)
	}

	return NewConfiguration([]Initializer{
		PFlagConfiguration{
			Bindings: bindings,
		},
		NewDefaultFileConfiguration(cfg.ApplicationName),
		NewDefaultEnvironmentVariableConfiguration(cfg.ApplicationName),
	}), nil
}

type CobraBinding struct {
	ConfigurationKey    string
	CobraCommandLineKey CobraCommandLineKey
}

func CobraBindings(root *cobra.Command, configurations []CobraBinding) (map[string]*pflag.Flag, error) {
	bindings := make(map[string]*pflag.Flag)
	for _, configuration := range configurations {
		flag, err := configuration.CobraCommandLineKey.AssociatedFlag(root)
		if err != nil {
			return nil, err
		}
		bindings[configuration.ConfigurationKey] = flag
	}
	return bindings, nil
}

type CobraCommandLineKey struct {
	Key string
}

func (k CobraCommandLineKey) AssociatedFlag(root *cobra.Command) (*pflag.Flag, error) {
	return findFlag(root, k.Key)
}

func findFlag(cmd *cobra.Command, flagPath string) (*pflag.Flag, error) {
	pathParts := strings.Split(flagPath, ".")
	if len(pathParts) == 1 {
		return cmd.Flags().Lookup(flagPath), nil
	}
	subCmd := findSubCommand(cmd, pathParts[0])
	if nil == subCmd {
		return nil, fmt.Errorf("not found %s sub command", pathParts[0])
	}
	subPath := strings.Join(pathParts[1:len(pathParts)], ".")
	flag, err := findFlag(subCmd, subPath)
	if err != nil {
		return nil, fmt.Errorf("search %s: %w", pathParts[0], err)
	}
	return flag, nil
}

func findSubCommand(root *cobra.Command, use string) *cobra.Command {
	var subCommand *cobra.Command
	subCmds := root.Commands()
	for _, subCmd := range subCmds {
		if subCmd.Use == use {
			return subCmd
		}
	}
	return subCommand
}
