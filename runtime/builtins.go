// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package runtime provides runtime support for Moxie built-in functions
package runtime

import (
	"fmt"
	goruntime "runtime"
	"strconv"
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
	goruntime.GC()
}

// FreeSlice provides GC hint for slice memory
func FreeSlice[T any](s *[]T) {
	if s == nil {
		return
	}

	// Clear the slice
	*s = nil

	// Trigger GC
	goruntime.GC()
}

// FreeMap provides GC hint for map memory
func FreeMap[K comparable, V any](m *map[K]V) {
	if m == nil {
		return
	}

	// Clear the map
	*m = nil

	// Trigger GC
	goruntime.GC()
}

// SizeOf returns the size in bytes of the value pointed to
// This is a helper function, not a Moxie built-in
func SizeOf[T any](v *T) int64 {
	if v == nil {
		return 0
	}
	return int64(unsafe.Sizeof(*v))
}

// Concat concatenates two byte slices and returns a pointer to the result
// Used for string concatenation: s1 + s2 -> Concat(s1, s2)
func Concat(s1, s2 *[]byte) *[]byte {
	if s1 == nil && s2 == nil {
		empty := []byte{}
		return &empty
	}
	if s1 == nil {
		return s2
	}
	if s2 == nil {
		return s1
	}

	result := append(*s1, *s2...)
	return &result
}

// ConcatSlice concatenates two slices of any type and returns a pointer to the result
// Used for array concatenation: a1 + a2 -> ConcatSlice[T](a1, a2)
// Always allocates a new slice, never mutates operands
func ConcatSlice[T any](s1, s2 *[]T) *[]T {
	if s1 == nil && s2 == nil {
		empty := []T{}
		return &empty
	}
	if s1 == nil {
		// Clone s2
		result := make([]T, len(*s2))
		copy(result, *s2)
		return &result
	}
	if s2 == nil {
		// Clone s1
		result := make([]T, len(*s1))
		copy(result, *s1)
		return &result
	}

	// Concatenate - allocate new slice with exact size needed
	result := make([]T, len(*s1)+len(*s2))
	copy(result, *s1)
	copy(result[len(*s1):], *s2)
	return &result
}

// Print prints values to stdout, converting *[]byte to strings
// This provides better output for Moxie strings than fmt.Println
func Print(args ...any) {
	converted := convertArgs(args...)
	for i, arg := range converted {
		fmt.Print(arg)
		if i < len(converted)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

// Printf prints formatted output, converting *[]byte to strings
func Printf(format *[]byte, args ...any) {
	// Convert format string
	var formatStr string
	if format != nil {
		formatStr = string(*format)
	}
	fmt.Printf(formatStr, convertArgs(args...)...)
}

// convertArgs converts *[]byte arguments to strings for printing
func convertArgs(args ...any) []any {
	result := make([]any, len(args))
	for i, arg := range args {
		if s, ok := arg.(*[]byte); ok && s != nil {
			result[i] = string(*s)
		} else {
			result[i] = arg
		}
	}
	return result
}

// String conversion functions for Moxie type conversions

// IntToString converts an integer to a Moxie string (*[]byte)
// Used for string(intValue) conversions
func IntToString(n int) *[]byte {
	s := strconv.Itoa(n)
	b := []byte(s)
	return &b
}

// Int64ToString converts an int64 to a Moxie string (*[]byte)
func Int64ToString(n int64) *[]byte {
	s := strconv.FormatInt(n, 10)
	b := []byte(s)
	return &b
}

// Int32ToString converts an int32 to a Moxie string (*[]byte)
func Int32ToString(n int32) *[]byte {
	s := strconv.FormatInt(int64(n), 10)
	b := []byte(s)
	return &b
}

// RuneToString converts a rune to a Moxie string (*[]byte)
// Used for string(runeValue) conversions
func RuneToString(r rune) *[]byte {
	s := string(r)
	b := []byte(s)
	return &b
}

// RunesToString converts a rune slice to a Moxie string (*[]byte)
// Used for string(*[]rune) conversions
func RunesToString(runes *[]rune) *[]byte {
	if runes == nil {
		return &[]byte{}
	}
	s := string(*runes)
	b := []byte(s)
	return &b
}

// StringToRunes converts a Moxie string (*[]byte) to a rune slice (*[]rune)
// Used for []rune(string) or *[]rune(string) conversions
func StringToRunes(s *[]byte) *[]rune {
	if s == nil {
		return &[]rune{}
	}
	r := []rune(string(*s))
	return &r
}

// BytesToString converts a byte slice to a Moxie string (*[]byte)
// This is essentially an identity conversion since Moxie strings are *[]byte
func BytesToString(b *[]byte) *[]byte {
	if b == nil {
		return &[]byte{}
	}
	return b
}
