//revive:disable:add-constant
package main

import (
	"fmt"
	"log"
	"sync"
)

func allMetrics(
	actionPages map[string]string, client *gatewayClient, configs configs, flags *flags, model string,
) {
	returnFact := empty

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Channel to capture errors from goroutines
	errChan := make(chan error, len(metricActions))

	for _, action := range metricActions {
		// Increment the WaitGroup counter for each goroutine
		wg.Add(1)

		// Launch each action retrieval in a new goroutine
		go func(action string) {
			defer wg.Done() // Decrement the counter when the goroutine finishes

			_, err := client.retrieveAction(action, actionPages, configs, flags, model, returnFact)
			if err != nil {
				// Send the error to the error channel
				errChan <- fmt.Errorf("failed to get %s: %v", action, err)
			}
		}(action)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the error channel after all goroutines are done
	close(errChan)

	// Check for errors
	for err := range errChan {
		log.Println("Error:", err)
	}
}
