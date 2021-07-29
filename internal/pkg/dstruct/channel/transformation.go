package channel

func Count(in <-chan interface{}) <-chan interface{} {
	count := 0

	return Map(in, func(_ interface{}) interface{} {
		return count
	})
}

func Map(in <-chan interface{}, mapFn func(value interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func(channel <-chan interface{}) {
		for value := range channel {
			out <- mapFn(value)
		}
		close(out)
	}(in)

	return out
}

func MapVoid(in <-chan struct{}, mapFn func() interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func(channel <-chan struct{}) {
		for range channel {
			out <- mapFn()
		}
		close(out)
	}(in)

	return out
}

func BufferCount(in <-chan interface{}, bufferSize int) <-chan interface{} {
	return BufferCountEvery(in, bufferSize, -1)
}

func BufferCountEvery(in <-chan interface{}, bufferSize int, startBufferEvery int) <-chan interface{} {
	out := make(chan interface{})
	if startBufferEvery < 0 {
		startBufferEvery = bufferSize
	}

	buffers := make([][]interface{}, 0)
	inCount := 0

	go func() {
		for value := range in {
			// add a new buffer if the inCount reaches a multiple of startBufferEvery
			if inCount % startBufferEvery == 0 {
				buffers = append(buffers, make([]interface{}, 0))
			}

			// iterate from the last buffer to the first
			for i := len(buffers) - 1; i >= 0; i-- {
				// add the value to the buffer
				buffers[i] = append(buffers[i], value)

				// if the buffer has reached the desired size, emit and remove it
				if len(buffers[i]) >= bufferSize {
					out <- buffers[i]

					// remove the buffer at i by replacing it with the last element
					// then assigning the buffer to the range until the moved element
					buffers[i] = buffers[len(buffers) - 1]
					buffers = buffers[:len(buffers) - 1]
				}
			}

			inCount++
		}
		close(out)
	}()

	return out
}