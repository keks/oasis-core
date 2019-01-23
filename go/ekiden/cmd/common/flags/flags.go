// Package flags implements common flags used across multiple commands
// and backends.
package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgVerbose = "verbose"
	cfgForce   = "force"
)

// Verbose returns true iff the verbose flag is set.
func Verbose() bool {
	return viper.GetBool(cfgVerbose)
}

// RegisterVerbose registers the verbose flag for the provided command.
func RegisterVerbose(cmd *cobra.Command) {
	if !cmd.Flags().Parsed() {
		cmd.Flags().BoolP(cfgVerbose, "v", false, "verbose output")
	}

	_ = viper.BindPFlag(cfgVerbose, cmd.Flags().Lookup(cfgVerbose))
}

// Force returns true iff the force flag is set.
func Force() bool {
	return viper.GetBool(cfgForce)
}

// RegisterForce registers the force flag for the provided command.
func RegisterForce(cmd *cobra.Command) {
	if !cmd.Flags().Parsed() {
		cmd.Flags().Bool(cfgForce, false, "force")
	}

	_ = viper.BindPFlag(cfgForce, cmd.Flags().Lookup(cfgForce))
}