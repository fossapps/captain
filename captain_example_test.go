package captain_test

import (
	"github.com/cyberhck/captain"
	"time"
	"sync"
	"fmt"
)

func ExampleConfig_Run() {
	job := captain.CreateJob()
	job.WithRuntimeProcessor(func(tick time.Time, message string, startTime time.Time) {
		fmt.Println("ticked")
	})
	job.WithRuntimeProcessingFrequency(100 * time.Millisecond)
	job.SetWorker(func(Channel chan string, WaitGroup *sync.WaitGroup) {
		time.Sleep(250 * time.Millisecond)
		WaitGroup.Done()
	})
	job.Run()
	// Output:
	// ticked
	// ticked
}
