package channel

import (
	"icotray/internal/pkg/dstruct"
	"sync"
)

func MergeReceivingWithGoroutines(channels ...<-chan dstruct.Any) <-chan dstruct.Any {
	merged := make(chan dstruct.Any)

	if len(channels) == 0 {
		return merged
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(channels))

	// report to the merged channel from a go routine per channel
	for _, channel := range channels {
		go func(c <-chan dstruct.Any) {
			for v := range c {
				merged <- v
			}
			waitGroup.Done()
		}(channel)
	}

	// close the merged channel after all reporting channels have finished
	go func() {
		waitGroup.Wait()
		close(merged)
	}()

	return merged
}
