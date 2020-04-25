package berus

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

// CobraConfiguration allow to customize a default viper configuration
type CobraConfiguration struct {
	// Name of the application, used as (uppercased) prefix for environment variables and file name
	ApplicationName string
	// RootCommand is the root cobra command
	RootCommand     *cobra.Command
	// Bindings refere all the bindings between cobra flags and viper configuration key
	Bindings        []CobraBinding
}

// NewCobraDefaultConfiguration allow to create a configuration with most used option.
// It will initialize it using given application name, and use
// - Command line options: using provided bindings and root command
// - Environment variables: see NewDefaultEnvironmentVariableConfiguration
// - Configuration file: see NewDefaultFileConfiguration
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

// CobraBinding is a class representing a binding between a viper configuration key and a cobra flag usage name.
type CobraBinding struct {
	// Viper configuration key. Will be used to access associated flag value.
	ConfigurationKey    string
	// Cobra flag identifier value to bind
	CobraCommandLineKey CobraCommandLineKey
}

// CobraBindings will load Flags for cobra commands (and subcommands), and associate them with a configuration key used by viper
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

// CobraCommandLineKey represent an identifier for a flag in a command or sub-command.
type CobraCommandLineKey struct {
	// Key is a string that will have the form <mycommand>.<mysubcommand>.<flag name>, where <mycommand> and <mysubcommand>
	// should be the the "use" field of the command and subcommand. <flag name> is the flag name to get and will be searched
	// using Lookup() method on pflag.Flagset returned by PersistentFlags() and Flags() method.
	Key string
}

// AssociatedFlag load flag associated with current key in the given command/subcommand.
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
	subPath := strings.Join(pathParts[1:], ".")
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
