package channel

import (
	"sync"
)

func Merge(ins ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	if len(ins) == 0 {
		return out
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(ins))

	// report to the out from a go routine per in
	for _, in := range ins {
		go func(c <-chan interface{}) {
			for v := range c {
				out <- v
			}
			waitGroup.Done()
		}(in)
	}

	// close out after all reporting ins have finished
	go func() {
		waitGroup.Wait()
		close(out)
	}()

	return out
}
