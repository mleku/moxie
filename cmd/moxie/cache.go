// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BuildCache manages cached transpiled files
type BuildCache struct {
	cacheDir string
	enabled  bool
}

// CacheMetadata stores metadata about cached files
type CacheMetadata struct {
	SourceFile   string    `json:"source_file"`
	SourceHash   string    `json:"source_hash"`
	TranspiledTo string    `json:"transpiled_to"`
	Timestamp    time.Time `json:"timestamp"`
	Dependencies []string  `json:"dependencies,omitempty"`
}

// NewBuildCache creates a new build cache
func NewBuildCache(enabled bool) (*BuildCache, error) {
	if !enabled {
		return &BuildCache{enabled: false}, nil
	}

	// Use .moxie-cache in current directory
	cacheDir := ".moxie-cache"

	// Create cache directories
	if err := os.MkdirAll(filepath.Join(cacheDir, "transpiled"), 0755); err != nil {
		return nil, fmt.Errorf("creating cache dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(cacheDir, "metadata"), 0755); err != nil {
		return nil, fmt.Errorf("creating metadata dir: %w", err)
	}

	return &BuildCache{
		cacheDir: cacheDir,
		enabled:  true,
	}, nil
}

// hashFile computes SHA256 hash of a file
func hashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetCached checks if a cached version exists and is valid
func (bc *BuildCache) GetCached(sourceFile string) (cachedFile string, valid bool, err error) {
	if !bc.enabled {
		return "", false, nil
	}

	// Compute source file hash
	sourceHash, err := hashFile(sourceFile)
	if err != nil {
		return "", false, err
	}

	// Check metadata
	metadataFile := filepath.Join(bc.cacheDir, "metadata", sourceHash+".json")
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		// No cached version
		return "", false, nil
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return "", false, nil
	}

	// Verify hash matches
	if metadata.SourceHash != sourceHash {
		return "", false, nil
	}

	// Check if cached file exists
	cachedFile = filepath.Join(bc.cacheDir, "transpiled", sourceHash+".go")
	if _, err := os.Stat(cachedFile); os.IsNotExist(err) {
		return "", false, nil
	}

	return cachedFile, true, nil
}

// Store saves a transpiled file to cache
func (bc *BuildCache) Store(sourceFile, transpiledFile string) error {
	if !bc.enabled {
		return nil
	}

	// Compute source file hash
	sourceHash, err := hashFile(sourceFile)
	if err != nil {
		return fmt.Errorf("hashing source: %w", err)
	}

	// Copy transpiled file to cache
	cachedFile := filepath.Join(bc.cacheDir, "transpiled", sourceHash+".go")
	if err := copyFile(transpiledFile, cachedFile); err != nil {
		return fmt.Errorf("copying to cache: %w", err)
	}

	// Store metadata
	metadata := CacheMetadata{
		SourceFile:   sourceFile,
		SourceHash:   sourceHash,
		TranspiledTo: cachedFile,
		Timestamp:    time.Now(),
	}

	metadataFile := filepath.Join(bc.cacheDir, "metadata", sourceHash+".json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling metadata: %w", err)
	}

	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return fmt.Errorf("writing metadata: %w", err)
	}

	return nil
}

// Clear removes all cached files
func (bc *BuildCache) Clear() error {
	if !bc.enabled {
		return nil
	}

	return os.RemoveAll(bc.cacheDir)
}

// Stats returns cache statistics
func (bc *BuildCache) Stats() (files int, size int64, err error) {
	if !bc.enabled {
		return 0, 0, nil
	}

	transpiledDir := filepath.Join(bc.cacheDir, "transpiled")
	err = filepath.Walk(transpiledDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files++
			size += info.Size()
		}
		return nil
	})

	return files, size, err
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
