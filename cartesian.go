package cartesian

func Iter(params ...[]interface{}) chan []interface{} {
	// create channel
	c := make(chan []interface{})
	if len(params) == 0 {
		close(c)
		return c // Return a safe value for nil/empty params.
	}
	go func() {
		iterate(c, params[0], []interface{}{}, params[1:]...)
		close(c)
	}()
	return c
}

func iterate(channel chan []interface{}, topLevel, result []interface{}, needUnpacking ...[]interface{}) {
	for _, p := range topLevel {
		newResult := append(append([]interface{}{}, result...), p)
		if len(needUnpacking) == 0 {
			channel <- newResult
			continue
		}
		iterate(channel, needUnpacking[0], newResult, needUnpacking[1:]...)
	}
}
