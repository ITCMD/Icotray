package channel

import (
	"time"
)

// Capacity transform the in channel to a filtered channel, which will only emit
// if the capacity during the duration is met
func Capacity(in <-chan interface{}, capacity int, duration time.Duration, allowOverlap bool) <-chan interface{} {
	mapToTimeNow := func(in <-chan interface{}) <-chan interface{} {
		return Map(in, func(value interface{}) interface{} {
			return time.Now()
		})
	}

	bufferCapacity := func(in <-chan interface{}) <-chan interface{} {
		return BufferCountEvery(in, capacity, 1)
	}

	filterExpired := func(in <-chan interface{}) <-chan interface{} {
		usedTimes := make([]time.Time, 0)

		removeExpiredUsedTimes := func() {
			for i := 0; i < len(usedTimes); i++ {
				if time.Now().Sub(usedTimes[i]) <= duration {
					usedTimes[i] = usedTimes[len(usedTimes)-1]
					usedTimes = usedTimes[:len(usedTimes)-1]
				}
			}
		}

		isTimeUsed := func(value time.Time) bool {
			for _, usedTime := range usedTimes {
				if usedTime == value {
					return true
				}
			}
			return false
		}

		addToUsed := func(values []interface{}) {
			for _, val := range values {
				timeVal := val.(time.Time)
				usedTimes = append(usedTimes, timeVal)
			}
		}

		return Filter(in, func(value interface{}) bool {
			sliceVal := value.([]interface{})
			if len(sliceVal) > 0 {
				initialTime := sliceVal[0].(time.Time)
				isWithinAllowed := time.Now().Sub(initialTime) <= duration

				if !allowOverlap && isWithinAllowed {
					if isTimeUsed(initialTime) {
						return false
					}

					removeExpiredUsedTimes()
					addToUsed(sliceVal)
				}

				return isWithinAllowed
			}

			return true
		})
	}

	mapToCount := func(in <-chan interface{}) <-chan interface{} {
		return Count(in)
	}

	return mapToCount(filterExpired(bufferCapacity(mapToTimeNow(in))))
}

func Filter(in <-chan interface{}, passFilterFn func(value interface{}) bool) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		for value := range in {
			if passFilterFn(value) {
				out <- value
			}
		}
		close(out)
	}()

	return out
}
