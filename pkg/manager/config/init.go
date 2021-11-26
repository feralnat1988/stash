package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

var (
	initOnce     sync.Once
	instanceOnce sync.Once
)

type FlagStruct struct {
	Host           net.IP
	Port           int
	ConfigFilePath string
	CpuProfilePath string
	NoBrowser      bool
}

func GetInstance() *Instance {
	instanceOnce.Do(func() {
		instance = &Instance{
			main:      viper.New(),
			overrides: viper.New(),
		}
	})
	return instance
}

func Initialize(flags FlagStruct) (*Instance, error) {
	var err error
	initOnce.Do(func() {
		overrides := makeOverrideConfig()

		_ = GetInstance()
		instance.overrides = overrides
		instance.cpuProfilePath = flags.CpuProfilePath

		if err = initConfig(instance, flags); err != nil {
			return
		}

		if instance.isNewSystem {
			if instance.Validate() == nil {
				// system has been initialised by the environment
				instance.isNewSystem = false
			}
		}

		if !instance.isNewSystem {
			err = instance.setExistingSystemDefaults()
			if err == nil {
				err = instance.SetInitialConfig()
			}
		}
	})
	return instance, err
}

func initConfig(instance *Instance, flags FlagStruct) error {
	v := instance.main

	// The config file is called config.  Leave off the file extension.
	v.SetConfigName("config")

	v.AddConfigPath(".")                                // Look for config in the working directory
	v.AddConfigPath(filepath.FromSlash("$HOME/.stash")) // Look for the config in the home directory

	configFile := ""
	envConfigFile := os.Getenv("STASH_CONFIG_FILE")

	if flags.ConfigFilePath != "" {
		configFile = flags.ConfigFilePath
	} else if envConfigFile != "" {
		configFile = envConfigFile
	}

	if configFile != "" {
		v.SetConfigFile(configFile)

		// if file does not exist, assume it is a new system
		if exists, _ := utils.FileExists(configFile); !exists {
			instance.isNewSystem = true

			// ensure we can write to the file
			if err := utils.Touch(configFile); err != nil {
				return fmt.Errorf(`could not write to provided config path "%s": %s`, configFile, err.Error())
			} else {
				// remove the file
				os.Remove(configFile)
			}

			return nil
		}
	}

	err := v.ReadInConfig() // Find and read the config file
	// if not found, assume its a new system
	var notFoundErr viper.ConfigFileNotFoundError
	if errors.As(err, &notFoundErr) {
		instance.isNewSystem = true
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func initEnvs(viper *viper.Viper) {
	viper.SetEnvPrefix("stash")     // will be uppercased automatically
	bindEnv(viper, "host")          // STASH_HOST
	bindEnv(viper, "port")          // STASH_PORT
	bindEnv(viper, "external_host") // STASH_EXTERNAL_HOST
	bindEnv(viper, "generated")     // STASH_GENERATED
	bindEnv(viper, "metadata")      // STASH_METADATA
	bindEnv(viper, "cache")         // STASH_CACHE
	bindEnv(viper, "stash")         // STASH_STASH
}

func bindEnv(viper *viper.Viper, key string) {
	if err := viper.BindEnv(key); err != nil {
		panic(fmt.Sprintf("unable to set environment key (%v): %v", key, err))
	}
}

func makeOverrideConfig() *viper.Viper {
	viper := viper.New()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logger.Infof("failed to bind flags: %s", err.Error())
	}

	initEnvs(viper)

	return viper
}
