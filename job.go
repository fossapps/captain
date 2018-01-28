package job

import (
	"time"
	"errors"
	"sync"
)

type Config struct {
	ResultProcessor            ResultProcessor
	RuntimeProcessor           RuntimeProcessor
	RuntimeProcessingFrequency time.Duration
	LockProvider               LockProvider
	Worker                     Worker
}

type ResultProcessor interface{}
type RuntimeProcessor func(tick time.Time, message string, startTime time.Time) error
type Worker func(Channel chan string, WaitGroup *sync.WaitGroup)

type LockProvider interface {
	Acquire() error
	Release() error
}

func New() Config {
	return Config{
		LockProvider:               nil,
		RuntimeProcessor:           nil,
		ResultProcessor:            nil,
		RuntimeProcessingFrequency: 200 * time.Millisecond,
	}
}

func (config *Config) WithLockProvider(lockProvider LockProvider) {
	config.LockProvider = lockProvider
}

func (config *Config) WithResultProcessor(processor ResultProcessor) {
	config.ResultProcessor = processor
}

func (config *Config) WithRuntimeProcessingFrequency(frequency time.Duration) {
	config.RuntimeProcessingFrequency = frequency
}

func (config *Config) WithRuntimeProcessor(processor RuntimeProcessor) {
	config.RuntimeProcessor = processor
}

func (config *Config) SetWorker(worker Worker) {
	config.Worker = worker
}

func (config *Config) Run() {
	err := config.ensureLock()
	if err != nil {
		panic(err)
	}
	err = config.runWorker()
	if err != nil {
		panic(err)
	}
}

func (config *Config) runWorker() error {
	if config.Worker == nil {
		return errors.New("worker not set")
	}
	startTime := time.Now()
	ticker := time.NewTicker(config.RuntimeProcessingFrequency)
	defer ticker.Stop()
	commChan := make(chan string)
	defer close(commChan)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go config.Worker(commChan, &wg)
	go config.reportRuntimeProcessor(ticker, commChan, startTime, &wg)
	wg.Wait()
	return nil
}

func (config *Config) reportRuntimeProcessor(ticker *time.Ticker, commChan chan string, startTime time.Time, group *sync.WaitGroup) {
	for t := range ticker.C {
		message := getMessage(commChan)
		if config.RuntimeProcessor == nil {
			return
		}
		config.invokeRuntimeProcessor(t, message, startTime)
	}
}

func (config *Config) invokeRuntimeProcessor(t time.Time, message string, startTime time.Time) error {
	err := config.RuntimeProcessor(t, message, startTime)
	if err != nil {
		panic(err)
	}
	return nil
}

func (config *Config) ensureLock() error {
	if config.LockProvider == nil {
		return nil
	}
	return config.LockProvider.Acquire()
}

func getMessage(ch chan string) string {
	select {
	case msg := <-ch:
		return msg
	default:
		return ""
	}
}
