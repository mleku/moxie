// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watchCommand watches for file changes and rebuilds/reruns automatically
func watchCommand(args []string) error {
	// Parse flags
	var (
		runMode   bool   // --run flag: run after build
		testMode  bool   // --test flag: test after build
		execCmd   string // --exec flag: custom command to run
		verbose   bool   // --verbose flag: show detailed output
		clearTerm bool   // --clear flag: clear terminal before each build (default true)
	)

	clearTerm = true // default to clearing terminal

	var paths []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--run":
			runMode = true
		case "--test":
			testMode = true
		case "--exec":
			if i+1 < len(args) {
				execCmd = args[i+1]
				i++
			} else {
				return fmt.Errorf("--exec requires a command argument")
			}
		case "--verbose", "-v":
			verbose = true
		case "--clear":
			clearTerm = true
		case "--no-clear":
			clearTerm = false
		case "-h", "--help":
			watchUsage()
			return nil
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			paths = append(paths, args[i])
		}
	}

	// Validate flag combinations
	modeCount := 0
	if runMode {
		modeCount++
	}
	if testMode {
		modeCount++
	}
	if execCmd != "" {
		modeCount++
	}

	if modeCount > 1 {
		return fmt.Errorf("cannot use --run, --test, and --exec together")
	}

	// Default to current directory if no paths specified
	if len(paths) == 0 {
		paths = []string{"."}
	}

	// Create watcher
	watcher, err := newMoxieWatcher(paths, verbose)
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Set up the build/run configuration
	watcher.clearTerm = clearTerm
	watcher.runMode = runMode
	watcher.testMode = testMode
	watcher.execCmd = execCmd

	// Print initial message
	fmt.Println("🔍 Watching for changes...")
	if runMode {
		fmt.Println("   Mode: Build and run")
	} else if testMode {
		fmt.Println("   Mode: Build and test")
	} else if execCmd != "" {
		fmt.Printf("   Mode: Build and exec '%s'\n", execCmd)
	} else {
		fmt.Println("   Mode: Build only")
	}
	fmt.Printf("   Watching: %s\n", strings.Join(paths, ", "))
	fmt.Println()

	// Do initial build
	if err := watcher.rebuild(); err != nil {
		fmt.Fprintf(os.Stderr, "Initial build failed: %v\n", err)
	}

	// Start watching
	return watcher.Watch()
}

func watchUsage() {
	fmt.Print(`usage: moxie watch [flags] [path ...]

Watch watches for changes to Moxie source files and automatically rebuilds.

The flags are:

	--run
		Run the program after successful build.
	--test
		Run tests after successful build.
	--exec <command>
		Execute custom command after successful build.
	--verbose, -v
		Show detailed output including file change events.
	--clear
		Clear terminal before each build (default).
	--no-clear
		Don't clear terminal before each build.

When watching a directory, moxie watch monitors all .mx and .x files recursively.

Examples:

	moxie watch                           Watch current directory, build on changes
	moxie watch --run examples/hello      Watch and run on changes
	moxie watch --test ./...              Watch and test on changes
	moxie watch --exec "go run ."         Watch and execute custom command
	moxie watch --verbose .               Watch with detailed logging
`)
}

// moxieWatcher handles file watching and rebuilding
type moxieWatcher struct {
	watcher  *fsnotify.Watcher
	paths    []string
	verbose  bool
	debounce time.Duration

	// Rebuild configuration
	clearTerm bool
	runMode   bool
	testMode  bool
	execCmd   string

	// Debouncing
	mu           sync.Mutex
	pendingEvent bool
	timer        *time.Timer

	// Track last build
	lastBuild time.Time
	buildArgs []string
}

func newMoxieWatcher(paths []string, verbose bool) (*moxieWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating watcher: %w", err)
	}

	mw := &moxieWatcher{
		watcher:   watcher,
		paths:     paths,
		verbose:   verbose,
		debounce:  300 * time.Millisecond, // Wait 300ms after last change
		lastBuild: time.Now(),
	}

	// Add all paths to watcher
	for _, path := range paths {
		if err := mw.addPath(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return mw, nil
}

func (mw *moxieWatcher) addPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if !info.IsDir() {
		// Watch the directory containing the file
		dir := filepath.Dir(path)
		if err := mw.watcher.Add(dir); err != nil {
			return fmt.Errorf("watching %s: %w", dir, err)
		}
		mw.buildArgs = []string{path}
		if mw.verbose {
			fmt.Printf("   Watching directory: %s (for %s)\n", dir, filepath.Base(path))
		}
		return nil
	}

	// Walk directory recursively
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Skip hidden directories and build artifacts
		name := info.Name()
		if strings.HasPrefix(name, ".") || name == "vendor" || strings.HasPrefix(name, "moxie-build-") {
			return filepath.SkipDir
		}

		if err := mw.watcher.Add(p); err != nil {
			return fmt.Errorf("watching %s: %w", p, err)
		}

		if mw.verbose {
			fmt.Printf("   Watching directory: %s\n", p)
		}

		return nil
	})
}

func (mw *moxieWatcher) Close() {
	if mw.timer != nil {
		mw.timer.Stop()
	}
	mw.watcher.Close()
}

func (mw *moxieWatcher) Watch() error {
	for {
		select {
		case event, ok := <-mw.watcher.Events:
			if !ok {
				return nil
			}

			// Only process relevant events
			if !mw.shouldProcess(event) {
				continue
			}

			if mw.verbose {
				fmt.Printf("   Event: %s %s\n", event.Op, event.Name)
			}

			// Debounce: schedule rebuild after delay
			mw.scheduleRebuild()

		case err, ok := <-mw.watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
		}
	}
}

func (mw *moxieWatcher) shouldProcess(event fsnotify.Event) bool {
	// Check if it's a Moxie source file
	ext := filepath.Ext(event.Name)
	if ext != ".mx" && ext != ".x" {
		return false
	}

	// Ignore certain operations
	if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return false
	}

	// Check if file still exists (might have been deleted)
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return true // Process deletions
	}

	if _, err := os.Stat(event.Name); os.IsNotExist(err) {
		return false
	}

	return true
}

func (mw *moxieWatcher) scheduleRebuild() {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	// Cancel existing timer
	if mw.timer != nil {
		mw.timer.Stop()
	}

	// Schedule new rebuild
	mw.timer = time.AfterFunc(mw.debounce, func() {
		mw.mu.Lock()
		mw.pendingEvent = false
		mw.mu.Unlock()

		if err := mw.rebuild(); err != nil {
			// Error already printed by rebuild()
		}
	})

	mw.pendingEvent = true
}

func (mw *moxieWatcher) rebuild() error {
	start := time.Now()

	// Clear terminal if requested
	if mw.clearTerm {
		// ANSI escape code to clear screen and move cursor to top
		fmt.Print("\033[2J\033[H")
	}

	// Print separator
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("🔨 Building... (%s)\n", start.Format("15:04:05"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Build the project
	buildErr := buildCommand(mw.buildArgs)

	// Calculate build time
	duration := time.Since(start)
	mw.lastBuild = time.Now()

	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))

	if buildErr != nil {
		// Build failed
		fmt.Printf("❌ Build failed (%.2fs)\n", duration.Seconds())
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println()
		return buildErr
	}

	// Build succeeded
	fmt.Printf("✅ Build succeeded (%.2fs)\n", duration.Seconds())
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println()

	// Execute post-build command if specified
	if mw.runMode {
		fmt.Println("▶️  Running...")
		fmt.Println()
		if err := runCommand(mw.buildArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Run failed: %v\n", err)
		}
	} else if mw.testMode {
		fmt.Println("🧪 Testing...")
		fmt.Println()
		if err := testCommand(mw.buildArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Tests failed: %v\n", err)
		}
	} else if mw.execCmd != "" {
		fmt.Printf("⚙️  Executing: %s\n", mw.execCmd)
		fmt.Println()
		// TODO: Execute custom command
		// This would require parsing and executing the command safely
		fmt.Println("(Custom command execution not yet implemented)")
	}

	fmt.Println()
	fmt.Println("🔍 Watching for changes...")
	fmt.Println()

	return nil
}
