// Package go_captain is used to run your job, and monitor during runtime
package captain

import (
	"errors"
	"sync"
	"time"
)

// Config represents a configuration of a job.
// You can either create your own, or use CreateJob function which will initialize basic configuration.
type Config struct {
	ResultProcessor            ResultProcessor
	RuntimeProcessor           RuntimeProcessor
	RuntimeProcessingFrequency time.Duration
	LockProvider               LockProvider
	Worker                     Worker
	SummaryBuffer              int
}

// CommChan is a basic struct containing channels used for communication between your worker,
// runtime and result processor
type CommChan struct {
	Logs   chan string
	Result chan string
}

// ResultProcessor is called after execution of worker.
type ResultProcessor func(results []string)

// RuntimeProcessor gets called every `RuntimeProcessingFrequency` duration.
type RuntimeProcessor func(tick time.Time, message string, startTime time.Time)

// Worker is called and is expected to do the real work.
type Worker func(channels CommChan)

// LockProvider is used if we need to make sure that two jobs aren't running at the same time.
type LockProvider interface {
	Acquire() error
	Release() error
}

// CreateJob creates a basic empty configuration with some defaults.
func CreateJob() Config {
	return Config{
		LockProvider:               nil,
		RuntimeProcessor:           nil,
		ResultProcessor:            nil,
		RuntimeProcessingFrequency: 200 * time.Millisecond,
		SummaryBuffer:              1,
	}
}

// WithLockProvider is used to set LockProvider.
func (config *Config) WithLockProvider(lockProvider LockProvider) {
	config.LockProvider = lockProvider
}

// WithResultProcessor is used to set ResultProcessor.
func (config *Config) WithResultProcessor(processor ResultProcessor) {
	config.ResultProcessor = processor
}

// WithRuntimeProcessingFrequency is used to set how frequently RuntimeProcessor is called.
func (config *Config) WithRuntimeProcessingFrequency(frequency time.Duration) {
	config.RuntimeProcessingFrequency = frequency
}

// WithRuntimeProcessor is used to set the RuntimeProcessor.
func (config *Config) WithRuntimeProcessor(processor RuntimeProcessor) {
	config.RuntimeProcessor = processor
}

// SetWorker is used to set Worker.
func (config *Config) SetWorker(worker Worker) {
	config.Worker = worker
}

// Run starts the job
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

func (config *Config) getCommunicationChannel() CommChan {
	return CommChan{
		Logs:   make(chan string),
		Result: make(chan string, config.SummaryBuffer),
	}
}

func (config *Config) runWorker() error {
	if config.Worker == nil {
		return errors.New("worker not set")
	}
	channels := config.getCommunicationChannel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go config.reportRuntimeProcessors(channels, time.Now())
	go config.invokeWorker(channels, &wg)
	wg.Wait()
	return nil
}

func (config *Config) invokeWorker(commChan CommChan, group *sync.WaitGroup) {
	defer group.Done()
	defer commChan.close()
	config.Worker(commChan)
	summary := getSummary(commChan.Result)
	config.invokeResultProcessor(summary)
}

func (config *Config) invokeResultProcessor(summary []string) {
	if config.ResultProcessor == nil {
		return
	}
	config.ResultProcessor(summary)
}

func (config *Config) reportRuntimeProcessors(commChan CommChan, startTime time.Time) {
	ticker := time.NewTicker(config.RuntimeProcessingFrequency)
	defer ticker.Stop()
	for t := range ticker.C {
		message := getCommunicationMessages(commChan)
		if config.RuntimeProcessor == nil {
			return
		}
		config.invokeRuntimeProcessor(t, message, startTime)
	}
}

func (config *Config) invokeRuntimeProcessor(t time.Time, message string, startTime time.Time) {
	config.RuntimeProcessor(t, message, startTime)
}

func (config *Config) ensureLock() error {
	if config.LockProvider == nil {
		return nil
	}
	return config.LockProvider.Acquire()
}

func getCommunicationMessages(ch CommChan) string {
	return getString(ch.Logs)
}

func getSummary(ch chan string) []string {
	var summary []string
	for {
		msg := getString(ch)
		if msg == "" {
			break
		}
		summary = append(summary, msg)
	}
	return summary
}

func getString(ch chan string) string {
	select {
	case msg := <-ch:
		return msg
	default:
		return ""
	}
}

func (channels *CommChan) close() {
	close(channels.Logs)
	close(channels.Result)
}
