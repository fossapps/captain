// external example tests for captain package
package go_captain_test

import (
	"fmt"
	"github.com/cyberhck/go-captain"
	"strconv"
)

func ExampleConfig_WithResultProcessor() {
	job := go_captain.CreateJob()
	job.SetWorker(func(channels go_captain.CommChan) {
		channels.Result <- "Total Items: " + strconv.Itoa(80)
	})
	job.WithResultProcessor(func(results []string) {
		fmt.Printf("%+v\n", results[0])
	})
	job.Run()
	// Output:
	// Total Items: 80
}
