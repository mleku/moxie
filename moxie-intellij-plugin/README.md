# Moxie IntelliJ Plugin

IntelliJ IDEA plugin for Moxie language support (.mx files).

## Features

- **Syntax Highlighting**: Full syntax highlighting for Moxie language
  - Keywords (func, var, const, etc.)
  - Built-in functions (append, len, make, grow, clone, free)
  - Built-in types (int, string, bool, etc.)
  - Comments (line and block)
  - String literals (including raw strings)
  - Numbers (integers, floats, imaginary)
  - Operators and delimiters

- **File Type Recognition**: Automatically recognizes .mx files as Moxie

- **Editor Features**:
  - Brace matching ((), {}, [])
  - Code folding (blocks and comments)
  - Comment/uncomment support (Ctrl+/ for line, Ctrl+Shift+/ for block)
  - Automatic quote pairing

- **Customizable Colors**: Configure syntax colors in Settings → Editor → Color Scheme → Moxie

## Installation

### From Source

1. Clone the repository
2. Open in IntelliJ IDEA
3. Run `./gradlew buildPlugin`
4. Install the plugin from `build/distributions/moxie-intellij-plugin-*.zip` via Settings → Plugins → Install Plugin from Disk

### Building

```bash
./gradlew buildPlugin
```

The plugin will be generated in `build/distributions/`.

## Development

### Requirements

- IntelliJ IDEA 2023.2 or later
- JDK 17 or later
- Gradle 8.5 or later

### Project Structure

```
src/main/
├── kotlin/com/moxie/lang/
│   ├── MoxieLanguage.kt              # Language definition
│   ├── MoxieFileType.kt              # File type definition
│   ├── MoxieLexer.kt                 # Lexical analyzer
│   ├── MoxieTokenTypes.kt            # Token definitions
│   ├── MoxieParser.kt                # Parser
│   ├── MoxieParserDefinition.kt      # Parser configuration
│   ├── MoxieSyntaxHighlighter.kt     # Syntax highlighter
│   ├── MoxieColorSettingsPage.kt     # Color settings UI
│   ├── MoxieCommenter.kt             # Comment handler
│   ├── MoxieBraceMatcher.kt          # Brace matching
│   ├── MoxieFoldingBuilder.kt        # Code folding
│   └── MoxieQuoteHandler.kt          # Quote handling
└── resources/
    ├── META-INF/plugin.xml            # Plugin configuration
    └── icons/moxie-file.svg           # File icon
```

### Running in Development

```bash
./gradlew runIde
```

This will launch a new IntelliJ IDEA instance with the plugin loaded.

## Moxie Language Features Supported

- All Go keywords and syntax
- Moxie-specific built-in functions:
  - `grow()` - Grow slice capacity
  - `clone()` - Clone slice
  - `free()` - Free slice memory
- Explicit pointer syntax for slices and maps: `*[]T`, `*map[K]V`
- Mutable strings (string = *[]byte in Moxie)
- FFI functions (dlopen, dlsym, dlclose, dlerror)

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
