package channel

import "icotray/internal/pkg/dstruct"

func MapReceiving(channel <-chan dstruct.Any, mapFn func(value dstruct.Any) dstruct.Any) <-chan dstruct.Any {
	mappedChannel := make(chan dstruct.Any)

	go func(channel <-chan dstruct.Any) {
		for value := range channel {
			mappedChannel <- mapFn(value)
		}
		close(mappedChannel)
	}(channel)

	return mappedChannel
}

func MapReceivingVoid(channel <-chan dstruct.Void, mapFn func() dstruct.Any) <-chan dstruct.Any {
	mappedChannel := make(chan dstruct.Any)

	go func(channel <-chan dstruct.Void) {
		for _ = range channel {
			mappedChannel <- mapFn()
		}
		close(mappedChannel)
	}(channel)

	return mappedChannel
}

func MapReceivingVoidVendor(channel <-chan struct{}, mapFn func() dstruct.Any) <-chan dstruct.Any {
	mappedChannel := make(chan dstruct.Any)

	go func(channel <-chan struct{}) {
		for _ = range channel {
			mappedChannel <- mapFn()
		}
		close(mappedChannel)
	}(channel)

	return mappedChannel
}
