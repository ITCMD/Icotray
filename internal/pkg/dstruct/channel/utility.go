package channel

func Tap(in <-chan interface{}, tapFn func(value interface{}) interface{}) {
	go func() {
		for value := range in {
			tapFn(value)
		}
	}()
}
