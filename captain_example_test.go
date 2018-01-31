// external example tests for captain package
package captain_test

import (
	"fmt"
	"github.com/cyberhck/captain"
	"strconv"
	"time"
)

func ExampleConfig_WithRuntimeProcessor() {
	job := captain.CreateJob()
	job.WithRuntimeProcessor(func(tick time.Time, message string, startTime time.Time) {
		fmt.Println("ticked")
	})
	job.WithRuntimeProcessingFrequency(50 * time.Millisecond)
	job.SetWorker(func(channels captain.CommChan) {
		time.Sleep(120 * time.Millisecond)
	})
	job.Run()
	// Output:
	// ticked
	// ticked
}

func ExampleConfig_CallsResultProcessorAfterJobIsDone() {
	job := captain.CreateJob()
	job.SetWorker(func(channels captain.CommChan) {
		channels.Result <- "Total Items: " + strconv.Itoa(80)
	})
	job.WithResultProcessor(func(results []string) {
		fmt.Printf("%+v\n", results[0])
	})
	job.Run()
	// Output:
	// Total Items: 80
}
