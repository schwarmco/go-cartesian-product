package cartesian

import (
	"maps"
	"slices"
)

// Iter takes slices and returns a channel, receiving cartesian products
func Iter[T any](sets ...[]T) chan []T {
	if len(sets) == 0 {
		c := make(chan []T)
		close(c)
		return c
	}

	buffer := 1
	for _, set := range sets {
		buffer *= len(set)
	}
	if buffer > 1000 {
		buffer = 1000
	}

	c := make(chan []T, buffer)
	go func() {
		defer close(c)
		iterate(c, sets[0], make([]T, 0, len(sets)), sets[1:]...)
	}()
	return c
}

func iterate[T any](c chan []T, topLevel, result []T, needUnpacking ...[]T) {
	if len(needUnpacking) == 0 {
		for _, val := range topLevel {
			c <- append(slices.Clip(result), val)
		}
		return
	}
	for _, val := range topLevel {
		iterate(c, needUnpacking[0], append(result, val), needUnpacking[1:]...)
	}
}

// IterMap takes maps to slices and returns a channel, receiving cartesian products
func IterMap[K comparable, V any](maps ...map[K][]V) chan map[K]V {
	if len(maps) == 0 {
		c := make(chan map[K]V)
		close(c)
		return c
	}

	// Flatten the maps so that each map has only 1 key
	combined := make([]keyedSlice[K, V], 0, len(maps))
	buffer := 1
	for _, m := range maps {
		for k, v := range m {
			buffer *= len(v)
			combined = append(combined, keyedSlice[K, V]{
				Key:    k,
				Values: v,
			})
		}
	}
	if buffer > 1000 {
		buffer = 1000
	}

	c := make(chan map[K]V, buffer)
	go func() {
		defer close(c)
		iterateMap(c, combined[0], make(map[K]V, len(combined)), combined[1:]...)
	}()
	return c
}

func iterateMap[K comparable, V any](c chan map[K]V, topLevel keyedSlice[K, V], result map[K]V, needUnpacking ...keyedSlice[K, V]) {
	if len(needUnpacking) == 0 {
		for _, val := range topLevel.Values {
			finalResult := maps.Clone(result)
			finalResult[topLevel.Key] = val
			c <- finalResult
		}
		return
	}
	for _, val := range topLevel.Values {
		nextResult := maps.Clone(result)
		nextResult[topLevel.Key] = val
		iterateMap(c, needUnpacking[0], nextResult, needUnpacking[1:]...)
	}
}

type keyedSlice[K comparable, V any] struct {
	Key    K
	Values []V
}
