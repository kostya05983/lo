package lo

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

// Filter iterates over elements of collection, returning an array of all elements predicate returns truthy for.
// Play: https://go.dev/play/p/Apjg3WeSi7K
func Filter[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := make([]V, 0, len(collection))

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// FilterNotNil iterates over elements of collection, returning an array of all elements that are not nil
// Play: https://go.dev/play/p/qTJ69EgONdU
func FilterNotNil[V any](collection []V) []V {
	result := make([]V, 0, len(collection))

	for _, item := range collection {
		if !IsNil(item) {
			result = append(result, item)
		}
	}

	return result
}

// Map manipulates a slice and transforms it to a slice of another type.
// Play: https://go.dev/play/p/OkPcYAhBo0D
func Map[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}

// FilterMap returns a slice which obtained after both filtering and mapping using the given callback function.
// The callback function should return two values:
//   - the result of the mapping operation and
//   - whether the result element should be included or not.
//
// Play: https://go.dev/play/p/-AuYXfy7opz
func FilterMap[T any, R any](collection []T, callback func(item T, index int) (R, bool)) []R {
	result := []R{}

	for i, item := range collection {
		if r, ok := callback(item, i); ok {
			result = append(result, r)
		}
	}

	return result
}

// FlatMap manipulates a slice and transforms and flattens it to a slice of another type.
// The transform function can either return a slice or a `nil`, and in the `nil` case
// no value is added to the final slice.
// Play: https://go.dev/play/p/YSoYmQTA8-U
func FlatMap[T any, R any](collection []T, iteratee func(item T, index int) []R) []R {
	result := make([]R, 0, len(collection))

	for i, item := range collection {
		result = append(result, iteratee(item, i)...)
	}

	return result
}

// Reduce reduces collection to a value which is the accumulated result of running each element in collection
// through accumulator, where each successive invocation is supplied the return value of the previous.
// Play: https://go.dev/play/p/R4UHXZNaaUG
func Reduce[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}

// ReduceRight helper is like Reduce except that it iterates over elements of collection from right to left.
// Play: https://go.dev/play/p/Fq3W70l7wXF
func ReduceRight[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i := len(collection) - 1; i >= 0; i-- {
		initial = accumulator(initial, collection[i], i)
	}

	return initial
}

// ForEach iterates over elements of collection and invokes iteratee for each element.
// Play: https://go.dev/play/p/oofyiUPRf8t
func ForEach[T any](collection []T, iteratee func(item T, index int)) {
	for i, item := range collection {
		iteratee(item, i)
	}
}

// Times invokes the iteratee n times, returning an array of the results of each invocation.
// The iteratee is invoked with index as argument.
// Play: https://go.dev/play/p/vgQj3Glr6lT
func Times[T any](count int, iteratee func(index int) T) []T {
	result := make([]T, count)

	for i := 0; i < count; i++ {
		result[i] = iteratee(i)
	}

	return result
}

// Uniq returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array.
// Play: https://go.dev/play/p/DTzbeXZ6iEN
func Uniq[T comparable](collection []T) []T {
	result := make([]T, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))

	for _, item := range collection {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// UniqBy returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is
// invoked for each element in array to generate the criterion by which uniqueness is computed.
// Play: https://go.dev/play/p/g42Z3QSb53u
func UniqBy[T any, U comparable](collection []T, iteratee func(item T) U) []T {
	result := make([]T, 0, len(collection))
	seen := make(map[U]struct{}, len(collection))

	for _, item := range collection {
		key := iteratee(item)

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}

// GroupBy returns an object composed of keys generated from the results of running each element of collection through iteratee.
// Play: https://go.dev/play/p/XnQBd_v6brd
func GroupBy[T any, U comparable](collection []T, iteratee func(item T) U) map[U][]T {
	result := map[U][]T{}

	for _, item := range collection {
		key := iteratee(item)

		result[key] = append(result[key], item)
	}

	return result
}

// Chunk returns an array of elements split into groups the length of size. If array can't be split evenly,
// the final chunk will be the remaining elements.
// Play: https://go.dev/play/p/EeKl0AuTehH
func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([][]T, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}

	return result
}

// PartitionBy returns an array of elements split into groups. The order of grouped values is
// determined by the order they occur in collection. The grouping is generated from the results
// of running each element of collection through iteratee.
// Play: https://go.dev/play/p/NfQ_nGjkgXW
func PartitionBy[T any, K comparable](collection []T, iteratee func(item T) K) [][]T {
	result := [][]T{}
	seen := map[K]int{}

	for _, item := range collection {
		key := iteratee(item)

		resultIndex, ok := seen[key]
		if !ok {
			resultIndex = len(result)
			seen[key] = resultIndex
			result = append(result, []T{})
		}

		result[resultIndex] = append(result[resultIndex], item)
	}

	return result

	// unordered:
	// groups := GroupBy[T, K](collection, iteratee)
	// return Values[K, []T](groups)
}

// Flatten returns an array a single level deep.
// Play: https://go.dev/play/p/rbp9ORaMpjw
func Flatten[T any](collection [][]T) []T {
	totalLen := 0
	for i := range collection {
		totalLen += len(collection[i])
	}

	result := make([]T, 0, totalLen)
	for i := range collection {
		result = append(result, collection[i]...)
	}

	return result
}

// Interleave round-robin alternating input slices and sequentially appending value at index into result
// Play: https://go.dev/play/p/-RJkTLQEDVt
func Interleave[T any](collections ...[]T) []T {
	if len(collections) == 0 {
		return []T{}
	}

	maxSize := 0
	totalSize := 0
	for _, c := range collections {
		size := len(c)
		totalSize += size
		if size > maxSize {
			maxSize = size
		}
	}

	if maxSize == 0 {
		return []T{}
	}

	result := make([]T, totalSize)

	resultIdx := 0
	for i := 0; i < maxSize; i++ {
		for j := range collections {
			if len(collections[j])-1 < i {
				continue
			}

			result[resultIdx] = collections[j][i]
			resultIdx++
		}
	}

	return result
}

// Shuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.
// Play: https://go.dev/play/p/Qp73bnTDnc7
func Shuffle[T any](collection []T) []T {
	rand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})

	return collection
}

// Reverse reverses array so that the first element becomes the last, the second element becomes the second to last, and so on.
// Play: https://go.dev/play/p/fhUMLvZ7vS6
func Reverse[T any](collection []T) []T {
	length := len(collection)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}

	return collection
}

// Fill fills elements of array with `initial` value.
// Play: https://go.dev/play/p/VwR34GzqEub
func Fill[T Clonable[T]](collection []T, initial T) []T {
	result := make([]T, 0, len(collection))

	for range collection {
		result = append(result, initial.Clone())
	}

	return result
}

// Repeat builds a slice with N copies of initial value.
// Play: https://go.dev/play/p/g3uHXbmc3b6
func Repeat[T Clonable[T]](count int, initial T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, initial.Clone())
	}

	return result
}

// RepeatBy builds a slice with values returned by N calls of callback.
// Play: https://go.dev/play/p/ozZLCtX_hNU
func RepeatBy[T any](count int, predicate func(index int) T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, predicate(i))
	}

	return result
}

// KeyBy transforms a slice or an array of structs to a map based on a pivot callback.
// Play: https://go.dev/play/p/mdaClUAT-zZ
func KeyBy[K comparable, V any](collection []V, iteratee func(item V) K) map[K]V {
	result := make(map[K]V, len(collection))

	for _, v := range collection {
		k := iteratee(v)
		result[k] = v
	}

	return result
}

// Associate returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func Associate[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for _, t := range collection {
		k, v := transform(t)
		result[k] = v
	}

	return result
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
// Alias of Associate().
// Play: https://go.dev/play/p/WHa2CfMO3Lr
func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return Associate(collection, transform)
}

// Drop drops n elements from the beginning of a slice or array.
// Play: https://go.dev/play/p/JswS7vXRJP2
func Drop[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return make([]T, 0)
	}

	result := make([]T, 0, len(collection)-n)

	return append(result, collection[n:]...)
}

// DropRight drops n elements from the end of a slice or array.
// Play: https://go.dev/play/p/GG0nXkSJJa3
func DropRight[T any](collection []T, n int) []T {
	if len(collection) <= n {
		return []T{}
	}

	result := make([]T, 0, len(collection)-n)
	return append(result, collection[:len(collection)-n]...)
}

// DropWhile drops elements from the beginning of a slice or array while the predicate returns true.
// Play: https://go.dev/play/p/7gBPYw2IK16
func DropWhile[T any](collection []T, predicate func(item T) bool) []T {
	i := 0
	for ; i < len(collection); i++ {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make([]T, 0, len(collection)-i)
	return append(result, collection[i:]...)
}

// DropRightWhile drops elements from the end of a slice or array while the predicate returns true.
// Play: https://go.dev/play/p/3-n71oEC0Hz
func DropRightWhile[T any](collection []T, predicate func(item T) bool) []T {
	i := len(collection) - 1
	for ; i >= 0; i-- {
		if !predicate(collection[i]) {
			break
		}
	}

	result := make([]T, 0, i+1)
	return append(result, collection[:i+1]...)
}

// Reject is the opposite of Filter, this method returns the elements of collection that predicate does not return truthy for.
// Play: https://go.dev/play/p/YkLMODy1WEL
func Reject[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := []V{}

	for i, item := range collection {
		if !predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// Count counts the number of elements in the collection that compare equal to value.
// Play: https://go.dev/play/p/Y3FlK54yveC
func Count[T comparable](collection []T, value T) (count int) {
	for _, item := range collection {
		if item == value {
			count++
		}
	}

	return count
}

// CountBy counts the number of elements in the collection for which predicate is true.
// Play: https://go.dev/play/p/ByQbNYQQi4X
func CountBy[T any](collection []T, predicate func(item T) bool) (count int) {
	for _, item := range collection {
		if predicate(item) {
			count++
		}
	}

	return count
}

// CountValues counts the number of each element in the collection.
// Play: https://go.dev/play/p/-p-PyLT4dfy
func CountValues[T comparable](collection []T) map[T]int {
	result := make(map[T]int)

	for _, item := range collection {
		result[item]++
	}

	return result
}

// CountValuesBy counts the number of each element return from mapper function.
// Is equivalent to chaining lo.Map and lo.CountValues.
// Play: https://go.dev/play/p/2U0dG1SnOmS
func CountValuesBy[T any, U comparable](collection []T, mapper func(item T) U) map[U]int {
	result := make(map[U]int)

	for _, item := range collection {
		result[mapper(item)]++
	}

	return result
}

// Subset returns a copy of a slice from `offset` up to `length` elements. Like `slice[start:start+length]`, but does not panic on overflow.
// Play: https://go.dev/play/p/tOQu1GhFcog
func Subset[T any](collection []T, offset int, length uint) []T {
	size := len(collection)

	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}

	if offset > size {
		return []T{}
	}

	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}

	return collection[offset : offset+int(length)]
}

// Slice returns a copy of a slice from `start` up to, but not including `end`. Like `slice[start:end]`, but does not panic on overflow.
// Play: https://go.dev/play/p/8XWYhfMMA1h
func Slice[T any](collection []T, start int, end int) []T {
	size := len(collection)

	if start >= end {
		return []T{}
	}

	if start > size {
		start = size
	}
	if start < 0 {
		start = 0
	}

	if end > size {
		end = size
	}
	if end < 0 {
		end = 0
	}

	return collection[start:end]
}

// Replace returns a copy of the slice with the first n non-overlapping instances of old replaced by new.
// Play: https://go.dev/play/p/XfPzmf9gql6
func Replace[T comparable](collection []T, old T, new T, n int) []T {
	result := make([]T, len(collection))
	copy(result, collection)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}

	return result
}

// ReplaceAll returns a copy of the slice with all non-overlapping instances of old replaced by new.
// Play: https://go.dev/play/p/a9xZFUHfYcV
func ReplaceAll[T comparable](collection []T, old T, new T) []T {
	return Replace(collection, old, new, -1)
}

// Compact returns a slice of all non-zero elements.
// Play: https://go.dev/play/p/tXiy-iK6PAc
func Compact[T comparable](collection []T) []T {
	var zero T

	result := make([]T, 0, len(collection))

	for _, item := range collection {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// IsSorted checks if a slice is sorted.
// Play: https://go.dev/play/p/mc3qR-t4mcx
func IsSorted[T constraints.Ordered](collection []T) bool {
	for i := 1; i < len(collection); i++ {
		if collection[i-1] > collection[i] {
			return false
		}
	}

	return true
}

// IsSortedByKey checks if a slice is sorted by iteratee.
// Play: https://go.dev/play/p/wiG6XyBBu49
func IsSortedByKey[T any, K constraints.Ordered](collection []T, iteratee func(item T) K) bool {
	size := len(collection)

	for i := 0; i < size-1; i++ {
		if iteratee(collection[i]) > iteratee(collection[i+1]) {
			return false
		}
	}

	return true
}
