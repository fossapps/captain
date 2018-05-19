// external example tests for captain package
package captain_test

import (
	"fmt"
	"github.com/fossapps/captain"
	"strconv"
)

func ExampleConfig_WithResultProcessor() {
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
