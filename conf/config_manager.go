package conf

import "sync"

// Simple interface that allows us to switch out both implementations of the Manager
type ConfigManager interface {
	Set(conf *Config)
	Get() *Config
	Close()
}

// This struct manages the configuration instance by
// Preforming locking around access to the Config struct.
type MutexConfigManager struct {
	conf  *Config
	mutex *sync.Mutex
}

// ChannelConfigManager manages the configuration instance by feeding a
// Pointer through a channel whenever the user calls Get()
type ChannelConfigManager struct {
	conf *Config
	get  chan *Config
	set  chan *Config
	done chan bool
}

// Set is config type setter
func (this *MutexConfigManager) Set(conf *Config) {
	this.mutex.Lock()
	this.conf = conf
	this.mutex.Unlock()
}

// Get is config type getter
func (this *MutexConfigManager) Get() *Config {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.conf
}

// Close does not do nothing for now
func (this *MutexConfigManager) Close() {
	// Do nothing for now
}

// Set is ChannelConfigManager type setter
func (self *ChannelConfigManager) Set(conf *Config) {
	self.set <- conf
}

// get is ChannelConfigManager type getter
func (self *ChannelConfigManager) Get() *Config {
	return <-self.get
}

// Close close the current channel and stop goroutine
func (self *ChannelConfigManager) Close() {
	self.done <- true
}

// NewMutexConfigManager return the new instance of MutexConfigManager type
func NewMutexConfigManager(conf *Config) *MutexConfigManager {
	return &MutexConfigManager{conf, &sync.Mutex{}}
}

func NewChannelConfigManager(conf *Config) *ChannelConfigManager {
	parser := &ChannelConfigManager{conf: conf, set: make(chan *Config), get: make(chan *Config), done: make(chan bool)}
	parser.start()
	return parser
}

// start starts a new goroutine for ChannelConfigManager
func (this *ChannelConfigManager) start() {
	go func() {
		defer func() {
			close(this.get)
			close(this.set)
			close(this.done)
		}()
		for {
			select {
			case this.get <- this.conf:
			case value := <-this.set:
				this.conf = value
			case <-this.done:
				return
			}
		}
	}()
}
