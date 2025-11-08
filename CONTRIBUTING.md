# Contributing to Moxie

Moxie is an open source project - a fork of Go with fundamental language revisions.

We appreciate your help in building a better systems programming language!

## About Moxie

Moxie is actively being developed according to the [implementation plan](go-to-moxie-plan.md). We are currently in **Phase 0: Foundation & Setup**.

## Before Contributing

1. Read the [language specification](go-language-revision.md) to understand Moxie's design goals
2. Review the [implementation plan](go-to-moxie-plan.md) to see what phase we're in
3. Check existing issues to avoid duplicate work
4. For major changes, open an issue first to discuss your approach

## How to Contribute

### Reporting Issues

When filing an issue, please include:

1. What version of Moxie are you using (`moxie version`)?
2. What operating system and processor architecture are you using?
3. What did you do?
4. What did you expect to see?
5. What did you see instead?
6. Which phase of the implementation does this relate to?

### Suggesting Features

Moxie's core language features are defined in [go-language-revision.md](go-language-revision.md). For new feature proposals:

1. Check if it aligns with Moxie's design principles (simplicity, explicitness, performance)
2. Open an issue with detailed rationale and use cases
3. Be prepared to discuss trade-offs

### Contributing Code

1. **Check the current phase**: We implement features in a specific order (see [go-to-moxie-plan.md](go-to-moxie-plan.md))
2. **Claim an issue**: Comment on an issue to indicate you're working on it
3. **Follow the plan**: Stick to the documented implementation approach
4. **Write tests**: All code changes must include tests
5. **Document changes**: Update relevant documentation

### Development Process

```bash
# Clone the repository
git clone https://github.com/mleku/moxie.git
cd moxie

# Build from source
cd src
./all.bash

# Run tests
./run.bash
```

### Code Style

- Follow the existing code style (based on Go's style)
- Run `moxiefmt` on your code (when available)
- Write clear commit messages explaining the "why"

### Testing Requirements

- All new features must have tests
- All bug fixes must have regression tests
- Tests must pass on all supported platforms
- Performance-sensitive code should include benchmarks

## Implementation Phases

Moxie is being built in phases:

- **Phase 0**: Foundation & Setup (Current)
- **Phase 1**: Type System Foundation
- **Phase 2**: Built-in Functions
- **Phase 3**: String & Byte Unification
- **Phase 4**: Const & Immutability
- **Phase 5**: Zero-Copy Type Coercion
- **Phase 6**: FFI & dlopen
- **Phase 7**: Standard Library Updates
- **Phase 8**: Testing & Validation
- **Phase 9**: Documentation & Tools
- **Phase 10**: Releases

See the [full implementation plan](go-to-moxie-plan.md) for details.

## Areas for Contribution

### High Priority (Phase 0)
- Build system updates
- Testing infrastructure
- Documentation improvements
- Example code

### Future Phases
- Core language features (follow the plan)
- Standard library updates
- Tool development
- Performance optimization

## Code of Conduct

Be respectful and constructive. We're all working toward the same goal: a better systems programming language.

## Questions?

- Open an issue for discussion
- Check the [documentation](go-language-revision.md)
- Review the [implementation plan](go-to-moxie-plan.md)

## License

Unless otherwise noted, the Moxie source files are distributed under the BSD-style license found in the LICENSE file.

---

## Acknowledgments

Moxie is built on top of Go, which is the work of hundreds of contributors. We are grateful for their foundational work.

Thank you for contributing to Moxie!
