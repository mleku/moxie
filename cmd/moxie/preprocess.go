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
// doesn't natively support, like channel literals.
func preprocessMoxieSyntax(source string) string {
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

	return source
}

// postprocessMoxieSyntax reverses preprocessing for error messages
// This ensures error messages show the original Moxie syntax
func postprocessMoxieSyntax(message string) string {
	message = strings.ReplaceAll(message, "__MoxieChan[", "chan ")
	message = strings.ReplaceAll(message, "__MoxieChanSend[", "chan<- ")
	message = strings.ReplaceAll(message, "__MoxieChanRecv[", "<-chan ")
	message = strings.ReplaceAll(message, "]", "")
	return message
}
