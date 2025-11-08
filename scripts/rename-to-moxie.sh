#!/bin/bash
# Moxie Phase 0: Systematic renaming script
# This script updates user-facing strings from "go" to "moxie"

set -e

echo "Moxie Renaming Script - Phase 0"
echo "================================"
echo ""
echo "This script will update user-facing strings from 'Go' to 'Moxie'"
echo "WARNING: This modifies source files. Ensure you have backups!"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# Function to replace in files
replace_in_file() {
    local file=$1
    local pattern=$2
    local replacement=$3

    if [ -f "$file" ]; then
        sed -i "s/$pattern/$replacement/g" "$file"
        echo "Updated: $file"
    fi
}

# Update version command output
echo ""
echo "Step 1: Updating version command..."
replace_in_file "src/cmd/go/internal/version/version.go" 'print Go version' 'print Moxie version'
replace_in_file "src/cmd/go/internal/version/version.go" 'go version' 'moxie version'
replace_in_file "src/cmd/go/internal/version/version.go" 'Go version' 'Moxie version'

# Update main go command
echo ""
echo "Step 2: Updating go command help text..."
if [ -f "src/cmd/go/main.go" ]; then
    sed -i 's/Go is a tool for managing Go source code\./Moxie is a tool for managing Moxie source code./g' src/cmd/go/main.go
    sed -i 's/\tgo <command>/\tmoxie <command>/g' src/cmd/go/main.go
    echo "Updated: src/cmd/go/main.go"
fi

# Update build info
echo ""
echo "Step 3: Updating build info..."
if [ -f "src/runtime/debug/mod.go" ]; then
    sed -i 's/Go toolchain/Moxie toolchain/g' src/runtime/debug/mod.go
    echo "Updated: src/runtime/debug/mod.go"
fi

# Update compiler version strings
echo ""
echo "Step 4: Updating compiler version strings..."
if [ -f "src/cmd/compile/internal/base/flag.go" ]; then
    sed -i 's/"Go compiler"/"Moxie compiler"/g' src/cmd/compile/internal/base/flag.go
    echo "Updated: src/cmd/compile/internal/base/flag.go"
fi

echo ""
echo "================================"
echo "Renaming complete!"
echo ""
echo "Next steps:"
echo "1. Test the build: cd src && ./make.bash"
echo "2. Review changes: git diff"
echo "3. Run tests: ./all.bash"
echo ""
echo "Note: This script only updates user-facing strings."
echo "Additional binary renaming and environment variable updates needed."
