package conf

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"time"
	"strconv"
	"github.com/kataras/golog"
)

// This is the struct that holds our application's configuration
// Note that it could be change depends on how the application configuration needs
type Config struct {
	Message string `yaml:"message"`
}

type basicConfig struct {
	configFileAddress string
	gRPCBindPort      string
	updateInterval    time.Duration
}

// Variables
var (
	// basicConfigInstance only use one time for getting necessary env variables
	basicConfigInstance *basicConfig
	// ConfigManagerInstance is an exported instance for using globally in project
	ConfigManagerInstance *MutexConfigManager
	//
	watcher *FileWatcher
	// callBackFunction is a function executed every time that config file updated
	callBackFunction = func() {
		golog.Info("ConfigFile Updated")
		conf := loadConfig(basicConfigInstance.configFileAddress)
		ConfigManagerInstance.Set(conf)
	}
)

// fatal only fatal error and print a fatal error if any needed
func fatal(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// Init initialize a new basicConfig struct
func (this *basicConfig) Init() {
	// Getting CONFIG_FILE_ADDR env variable and fatal for any error
	configFileAddress, err := GetEnvVariable("CONFIG_FILE_ADDR")
	fatal(err)
	// Getting GRPC_BIND_PORT env variable
	gRPCBindPort, err := GetEnvVariable("GRPC_BIND_PORT")
	fatal(err)
	// Getting UPDATE_INTERVAL env variable
	// Note: UPDATE_INTERVAL is in Second format
	updateInterval, err := GetEnvVariable("UPDATE_INTERVAL")
	fatal(err)
	// Converting string UPDATE_INTERVAL env variable to integer
	interval, err := strconv.Atoi(updateInterval)
	fatal(err)
	// Set basic configuration
	this.configFileAddress = configFileAddress
	this.gRPCBindPort = gRPCBindPort
	// convert interval to Time.Duration and than set
	this.updateInterval = time.Duration(interval) * time.Second
}

func init() {
	// Only for init function
	var err error
	// Initialize basic configuration
	basicConfigInstance.Init()
	// Initialize ConfigManagerInstance
	ConfigManagerInstance = NewMutexConfigManager(loadConfig(basicConfigInstance.configFileAddress))
	// Create new watcher goroutine for watching config change
	// CallBack function is executed when update event happens
	watcher, err = WatchFile(
		basicConfigInstance.configFileAddress,
		basicConfigInstance.updateInterval,
		callBackFunction)
	fatal(err)
}

// loadConfig load configuration file and unmarshal it to config struct
// Kubernetes default ConfigMap file format is 'yaml'
// Note that it is not bound to Kubernetes and can be use for any yaml files
func loadConfig(configFile string) *Config {
	conf := &Config{}
	// Read from file
	configData, err := ioutil.ReadFile(configFile)
	fatal(err)
	// Unmarshal it to config struct
	err = yaml.Unmarshal(configData, conf)
	fatal(err)
	return conf
}

// CleanUp close all goroutine instances and release the memory
func CleanUp() {
	ConfigManagerInstance.Close()
	watcher.Close()
}
