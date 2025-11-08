// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package runtime provides runtime support for Moxie built-in functions
package runtime

import (
	"runtime"
	"unsafe"
)

// Grow pre-allocates capacity for a slice without changing its length
// Signature: grow(s *[]T, n int64) *[]T
func Grow[T any](s *[]T, n int64) *[]T {
	if s == nil {
		return nil
	}

	currentCap := int64(cap(*s))
	if currentCap >= n {
		return s // Already has enough capacity
	}

	// Create new slice with requested capacity
	newSlice := make([]T, len(*s), n)
	copy(newSlice, *s)
	*s = newSlice
	return s
}

// Clone creates a deep copy of a value
// Signature: clone(v *T) *T
func Clone[T any](v *T) *T {
	if v == nil {
		return nil
	}

	// For slices, maps, and other reference types,
	// we need special handling
	// This is a simplified version - full implementation would
	// need runtime type checking

	// Allocate new value
	result := new(T)
	*result = *v
	return result
}

// CloneSlice creates a deep copy of a slice
func CloneSlice[T any](s *[]T) *[]T {
	if s == nil {
		return nil
	}

	newSlice := make([]T, len(*s), cap(*s))
	copy(newSlice, *s)
	return &newSlice
}

// CloneMap creates a deep copy of a map
func CloneMap[K comparable, V any](m *map[K]V) *map[K]V {
	if m == nil {
		return nil
	}

	newMap := make(map[K]V, len(*m))
	for k, v := range *m {
		newMap[k] = v
	}
	return &newMap
}

// Free provides a hint to the garbage collector that memory can be freed
// Signature: free(p *T)
func Free[T any](p *T) {
	if p == nil {
		return
	}

	// Set pointer to zero value to help GC
	var zero T
	*p = zero

	// Trigger GC (this is aggressive, real implementation might be more subtle)
	runtime.GC()
}

// FreeSlice provides GC hint for slice memory
func FreeSlice[T any](s *[]T) {
	if s == nil {
		return
	}

	// Clear the slice
	*s = nil

	// Trigger GC
	runtime.GC()
}

// FreeMap provides GC hint for map memory
func FreeMap[K comparable, V any](m *map[K]V) {
	if m == nil {
		return
	}

	// Clear the map
	*m = nil

	// Trigger GC
	runtime.GC()
}

// SizeOf returns the size in bytes of the value pointed to
// This is a helper function, not a Moxie built-in
func SizeOf[T any](v *T) int64 {
	if v == nil {
		return 0
	}
	return int64(unsafe.Sizeof(*v))
}
