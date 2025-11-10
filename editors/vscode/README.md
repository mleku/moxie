# Moxie Language Support for VS Code

Official VS Code extension for the Moxie programming language.

## Features

### Language Server Protocol (LSP) Support
- **Symbol Navigation**: Go to definition, find references, document symbols
- **Code Intelligence**: Hover information, code completion
- **Real-time Diagnostics**: Syntax errors and linter integration
- **Code Formatting**: Format on save with `moxie fmt`

### Syntax Highlighting
- Full syntax highlighting for Moxie language
- Support for Moxie-specific syntax:
  - Channel literals: `&chan T{}`
  - Endianness coercion: `(*[]T, LittleEndian)(bytes)`
  - Builtin functions: `clone()`, `free()`, `grow()`

### Commands
- **Moxie: Build** - Build the current project
- **Moxie: Run** - Run the current file
- **Moxie: Test** - Run tests
- **Moxie: Format Document** - Format the current file
- **Moxie: Run Linter** - Run static analysis
- **Moxie: Clean Cache** - Clear build cache

### Code Snippets
Includes snippets for:
- Function declarations
- Struct and interface definitions
- Control flow statements
- Channel operations
- Moxie-specific patterns (clone, free, grow)

## Requirements

- Moxie compiler must be installed and available in PATH
- Run `moxie version` to verify installation

## Installation

### From Source
1. Clone the Moxie repository
2. Navigate to `editors/vscode/`
3. Run `npm install`
4. Run `npm run compile`
5. Run `vsce package` to create `.vsix` file
6. In VS Code, run "Extensions: Install from VSIX..."

### From Marketplace (Future)
Search for "Moxie" in the VS Code extension marketplace.

## Configuration

Configure the extension in VS Code settings:

```json
{
  // Path to moxie executable
  "moxie.moxiePath": "moxie",

  // Enable Language Server
  "moxie.lsp.enabled": true,

  // LSP trace level (off, messages, verbose)
  "moxie.lsp.trace": "off",

  // Format document on save
  "moxie.formatOnSave": true,

  // Run linter on save
  "moxie.vetOnSave": false
}
```

## Usage

### Opening a Moxie Project
1. Open a folder containing `.mx` or `.x` files
2. The extension will automatically activate
3. LSP features will be available immediately

### Running Commands
- Press `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
- Type "Moxie" to see all available commands

### Keyboard Shortcuts
- Format Document: `Shift+Alt+F` (default VS Code shortcut)
- Go to Definition: `F12`
- Find References: `Shift+F12`
- Hover: Hold `Ctrl` (Windows/Linux) or `Cmd` (Mac) over symbol

## Development

### Building from Source
```bash
cd editors/vscode
npm install
npm run compile
```

### Packaging
```bash
npm run package
```

### Publishing
```bash
vsce publish
```

## Troubleshooting

### LSP Server Not Starting
- Check that `moxie` is in your PATH
- Run `moxie lsp` manually to test
- Check VS Code Output panel → "Moxie Language Server"

### Syntax Highlighting Not Working
- Verify file extension is `.mx` or `.x`
- Reload VS Code window

### Format on Save Not Working
- Check `moxie.formatOnSave` setting
- Ensure `moxie fmt` command works

## Contributing

Contributions are welcome! Please see the main Moxie repository for contribution guidelines.

## License

BSD 3-Clause License - see LICENSE file in the Moxie repository.

## Links

- [Moxie Language](https://github.com/mleku/moxie)
- [Documentation](https://github.com/mleku/moxie/blob/main/README.md)
- [Issue Tracker](https://github.com/mleku/moxie/issues)
