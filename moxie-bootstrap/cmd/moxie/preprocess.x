// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"regexp"
	"strings"
)

// preprocessMoxieSyntax performs text-level preprocessing of Moxie syntax
// before parsing with the Go parser. This handles syntax that Go's parser
// doesn't natively support, like channel literals and endianness tuple syntax.
func preprocessMoxieSyntax(source string) string {
	// 1. Channel literal preprocessing
	// Replace &chan T{ with a marker structure that can be parsed
	// &chan T{...} → &__MoxieChan[T]{...}
	// This allows the parser to handle it as a composite literal

	// Pattern to match: &chan <type> {
	// We need to handle various channel types:
	// - &chan int{
	// - &chan string{
	// - &chan <-int{  (receive-only)
	// - &chan<- int{  (send-only)

	// Match bidirectional channels: &chan <type>{
	chanPattern := regexp.MustCompile(`&chan\s+([^{}<>]+)\s*\{`)
	source = chanPattern.ReplaceAllString(source, "&__MoxieChan[$1]{")

	// Match send-only channels: &chan<- <type>{
	sendOnlyPattern := regexp.MustCompile(`&chan<-\s+([^{}<>]+)\s*\{`)
	source = sendOnlyPattern.ReplaceAllString(source, "&__MoxieChanSend[$1]{")

	// Match receive-only channels: &<-chan <type>{
	recvOnlyPattern := regexp.MustCompile(`&<-chan\s+([^{}<>]+)\s*\{`)
	source = recvOnlyPattern.ReplaceAllString(source, "&__MoxieChanRecv[$1]{")

	// 2. Endianness tuple syntax preprocessing
	// Transform: (*[]T, LittleEndian)(expr) → __MoxieCoerceLE[T](expr)
	// Transform: (*[]T, BigEndian)(expr) → __MoxieCoerceBE[T](expr)
	// Transform: (*[]T, NativeEndian)(expr) → (*[]T)(expr) (no change needed)

	// Match: (*[]<type>, LittleEndian)(<expr>)
	// We need to be careful with the type pattern - it can be complex
	littleEndianPattern := regexp.MustCompile(`\(\*\[\]([^,)]+),\s*LittleEndian\)\s*\(`)
	source = littleEndianPattern.ReplaceAllString(source, "__MoxieCoerceLE[$1](")

	// Match: (*[]<type>, BigEndian)(<expr>)
	bigEndianPattern := regexp.MustCompile(`\(\*\[\]([^,)]+),\s*BigEndian\)\s*\(`)
	source = bigEndianPattern.ReplaceAllString(source, "__MoxieCoerceBE[$1](")

	// Match: (*[]<type>, NativeEndian)(<expr>) → (*[]<type>)(<expr>)
	// This is actually optional since NativeEndian is the default,
	// but we handle it for explicitness
	nativeEndianPattern := regexp.MustCompile(`\(\*\[\]([^,)]+),\s*NativeEndian\)\s*\(`)
	source = nativeEndianPattern.ReplaceAllString(source, "(*[]$1)(")

	return source
}

// postprocessMoxieSyntax reverses preprocessing for error messages
// This ensures error messages show the original Moxie syntax
func postprocessMoxieSyntax(message string) string {
	// Reverse channel literal markers
	message = strings.ReplaceAll(message, "__MoxieChan[", "chan ")
	message = strings.ReplaceAll(message, "__MoxieChanSend[", "chan<- ")
	message = strings.ReplaceAll(message, "__MoxieChanRecv[", "<-chan ")

	// Reverse endianness markers
	// Note: We need to be careful with the replacements to reconstruct the original syntax
	message = regexp.MustCompile(`__MoxieCoerceLE\[([^\]]+)\]`).ReplaceAllString(message, "(*[]$1, LittleEndian)")
	message = regexp.MustCompile(`__MoxieCoerceBE\[([^\]]+)\]`).ReplaceAllString(message, "(*[]$1, BigEndian)")

	// Clean up remaining brackets from channel markers
	message = strings.ReplaceAll(message, "]", "")

	return message
}
