// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Package name mapping between Moxie and Go conventions.
// For stdlib packages, Moxie uses lowercase names, same as Go.
// This mapping handles any special cases or future divergence.

// PackageMapping defines bidirectional package name translation
type PackageMapping struct {
	moxieToGo map[string]string // Moxie package name -> Go package name
	goToMoxie map[string]string // Go package name -> Moxie package name
}

// NewPackageMapping creates a new package mapping with stdlib defaults
func NewPackageMapping() *PackageMapping {
	pm := &PackageMapping{
		moxieToGo: make(map[string]string),
		goToMoxie: make(map[string]string),
	}

	// Initialize with standard library mappings
	// For now, most packages have 1:1 mapping (same names)
	// This structure allows for future divergence
	pm.addMapping("fmt", "fmt")
	pm.addMapping("os", "os")
	pm.addMapping("io", "io")
	pm.addMapping("net", "net")
	pm.addMapping("http", "http") // net/http subpackage
	pm.addMapping("url", "url")   // net/url subpackage
	pm.addMapping("time", "time")
	pm.addMapping("sync", "sync")
	pm.addMapping("strings", "strings")
	pm.addMapping("bytes", "bytes")
	pm.addMapping("errors", "errors")
	pm.addMapping("context", "context")
	pm.addMapping("json", "json") // encoding/json subpackage
	pm.addMapping("xml", "xml")   // encoding/xml subpackage
	pm.addMapping("base64", "base64") // encoding/base64
	pm.addMapping("hex", "hex")   // encoding/hex
	pm.addMapping("sql", "sql")   // database/sql
	pm.addMapping("testing", "testing")
	pm.addMapping("reflect", "reflect")
	pm.addMapping("unsafe", "unsafe")
	pm.addMapping("math", "math")
	pm.addMapping("rand", "rand") // math/rand
	pm.addMapping("big", "big")   // math/big
	pm.addMapping("sort", "sort")
	pm.addMapping("regexp", "regexp")
	pm.addMapping("unicode", "unicode")
	pm.addMapping("utf8", "utf8") // unicode/utf8
	pm.addMapping("path", "path")
	pm.addMapping("filepath", "filepath") // path/filepath
	pm.addMapping("flag", "flag")
	pm.addMapping("log", "log")
	pm.addMapping("slog", "slog") // log/slog
	pm.addMapping("bufio", "bufio")
	pm.addMapping("scanner", "scanner") // text/scanner
	pm.addMapping("template", "template") // text/template or html/template
	pm.addMapping("tar", "tar")   // archive/tar
	pm.addMapping("zip", "zip")   // archive/zip
	pm.addMapping("gzip", "gzip") // compress/gzip
	pm.addMapping("zlib", "zlib") // compress/zlib
	pm.addMapping("crypto", "crypto")
	pm.addMapping("md5", "md5")   // crypto/md5
	pm.addMapping("sha1", "sha1") // crypto/sha1
	pm.addMapping("sha256", "sha256") // crypto/sha256
	pm.addMapping("sha512", "sha512") // crypto/sha512
	pm.addMapping("aes", "aes")   // crypto/aes
	pm.addMapping("des", "des")   // crypto/des
	pm.addMapping("rsa", "rsa")   // crypto/rsa
	pm.addMapping("ecdsa", "ecdsa") // crypto/ecdsa
	pm.addMapping("tls", "tls")   // crypto/tls
	pm.addMapping("x509", "x509") // crypto/x509
	pm.addMapping("runtime", "runtime")
	pm.addMapping("debug", "debug") // runtime/debug
	pm.addMapping("pprof", "pprof") // runtime/pprof
	pm.addMapping("trace", "trace") // runtime/trace

	return pm
}

// addMapping adds a bidirectional mapping between Moxie and Go package names
func (pm *PackageMapping) addMapping(moxie, goPkg string) {
	pm.moxieToGo[moxie] = goPkg
	pm.goToMoxie[goPkg] = moxie
}

// MoxieToGo converts a Moxie package name to Go package name
// Returns the input unchanged if no mapping exists
func (pm *PackageMapping) MoxieToGo(moxiePkg string) string {
	if goPkg, ok := pm.moxieToGo[moxiePkg]; ok {
		return goPkg
	}
	return moxiePkg
}

// GoToMoxie converts a Go package name to Moxie package name
// Returns the input unchanged if no mapping exists
func (pm *PackageMapping) GoToMoxie(goPkg string) string {
	if moxiePkg, ok := pm.goToMoxie[goPkg]; ok {
		return moxiePkg
	}
	return goPkg
}

// HasMoxieMapping returns true if a Moxie->Go mapping exists
func (pm *PackageMapping) HasMoxieMapping(moxiePkg string) bool {
	_, ok := pm.moxieToGo[moxiePkg]
	return ok
}

// HasGoMapping returns true if a Go->Moxie mapping exists
func (pm *PackageMapping) HasGoMapping(goPkg string) bool {
	_, ok := pm.goToMoxie[goPkg]
	return ok
}

// Global package mapping instance
var pkgMap = NewPackageMapping()
