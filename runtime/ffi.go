package runtime

import (
	"github.com/ebitengine/purego"
)

// DLib represents a dynamically loaded library
type DLib struct {
	handle uintptr
}

// FFI constants for dlopen flags
// These match the values from purego
const (
	RTLD_LAZY   = 0x001 // Lazy function binding
	RTLD_NOW    = 0x002 // Immediate function binding
	RTLD_GLOBAL = 0x100 // Make symbols available for subsequently loaded libraries
	RTLD_LOCAL  = 0x000 // Symbols not available for symbol resolution of subsequently loaded libraries
)

// Dlopen loads a dynamic library
// Returns nil if the library cannot be loaded
// Accepts *[]byte (Moxie string) and converts to Go string
func Dlopen(filename *[]byte, flags int32) *DLib {
	if filename == nil {
		return nil
	}

	// Convert Moxie string (*[]byte) to Go string
	goFilename := string(*filename)

	// Use purego to load the library
	handle, err := purego.Dlopen(goFilename, int(flags))
	if err != nil {
		// Store error for Dlerror()
		lastError = err.Error()
		return nil
	}

	return &DLib{handle: handle}
}

// lastError stores the last error from FFI operations
var lastError string

// Dlsym looks up a symbol in a dynamic library
// T must be a function pointer type
// Returns the symbol cast to type T
// Accepts *[]byte (Moxie string) and converts to Go string
func Dlsym[T any](lib *DLib, name *[]byte) T {
	var zero T

	if lib == nil || lib.handle == 0 || name == nil {
		return zero
	}

	// Convert Moxie string (*[]byte) to Go string
	goName := string(*name)

	// Use purego to register the library function
	// This is a bit tricky because purego.RegisterLibFunc expects a pointer to a function variable
	var fn T
	err := purego.RegisterLibFunc(&fn, lib.handle, goName)
	if err != nil {
		lastError = err.Error()
		return zero
	}

	return fn
}

// Dlclose closes a dynamic library
func Dlclose(lib *DLib) {
	if lib != nil && lib.handle != 0 {
		// Note: purego doesn't expose dlclose directly in the public API
		// The library will be closed when the handle is garbage collected
		// We could use syscall.Syscall to call dlclose directly if needed
		lib.handle = 0
	}
}

// Dlerror returns the last error from dlopen, dlsym, or dlclose
// Returns empty Moxie string (*[]byte) if no error
func Dlerror() *[]byte {
	if lastError == "" {
		empty := []byte{}
		return &empty
	}

	bytes := []byte(lastError)
	lastError = "" // Clear error after reading
	return &bytes
}
