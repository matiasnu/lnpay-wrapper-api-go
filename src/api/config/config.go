package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration estructura
type Configuration struct {
	APIRestServerHost string `mapstructure:"api_host"`
	APIRestServerPort string `mapstructure:"api_port"`
	LoggingPath       string `mapstructure:"api_logpath"`
	LoggingFile       string `mapstructure:"api_logfile"`
	LoggingLevel      string `mapstructure:"api_loglevel"`
}

// Config is package struct containing conf params
var ConfMap Configuration

func Load(path string, name string, ext string) {

	// name := "parameters"
	// ext := "yml"
	// path := "./config"
	fmt.Printf("Loading configuration %s/%s.%s\n", path, name, ext)
	viper.SetConfigType(ext)
	viper.SetConfigName(name)
	viper.AddConfigPath(path)

	// Setting defaults if the config not read
	// API
	viper.SetDefault("api_host", "127.0.0.1")
	viper.SetDefault("api_port", ":8080")
	// LOG
	viper.SetDefault("api_logpath", "/var/log/")
	viper.SetDefault("api_logfile", "lnpay_wrapper_api_go.log")
	viper.SetDefault("api_loglevel", "trace")

	if _, err := os.Stat(filepath.Join(path, name+"."+ext)); err == nil {
		err = viper.ReadInConfig()
		if err == nil {
			viper.WatchConfig()
			viper.OnConfigChange(func(e fsnotify.Event) {
				// TODO: load new config values ...
				log.Println("Config file changed: ", e.Name)
			})
		} else {
			log.Errorln(err)
		}
	} else {
		log.Warningf("File parameters.yml not found. Working with default config: %s \n", err)
	}

	err := viper.Unmarshal(&ConfMap)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %+v", err)
	}
	fmt.Printf("Load configuration : \n")
	spew.Dump(ConfMap)
}
