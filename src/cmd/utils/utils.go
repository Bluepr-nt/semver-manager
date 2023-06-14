package utils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = "ccs"
	envPrefix             = "CCS"
)

func InitializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(defaultConfigFilename)

	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)

	v.AutomaticEnv()

	// Bind the current command's flags to viper
	// BindFlags(cmd, v)
	// v.BindPFlags(cmd.Flags())

	return nil
}
