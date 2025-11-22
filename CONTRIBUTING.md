# Contributing

Thank you for considering contributing to the Immortal relay!
Please read these guidelines before submitting a pull request or opening an issue.

> [!NOTE]
> This code guideline must be followed for both contributors and maintainers to review the PRs.

## Code Guidelines

We strive to maintain clean, readable, and maintainable code.
Please follow these guidelines when contributing to the project:

- Follow the [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/) and [The Rust Programming Language](https://doc.rust-lang.org/book/) best practices.
- Use `rustfmt` for code formatting and `clippy` for linting.
- Follow the principles of clean code as outlined in
  Robert C. Martin's "[Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)" book.
- Write comprehensive tests and benchmarks for new code or changes to existing code.
- Ensure all tests pass before submitting a pull request.
- Use meaningful variable and function names.
- Prefer explicit error handling over panic.

### Development Commands

Use these Cargo commands for development:

- `cargo check`: quickly check your code for compile errors
- `cargo fmt`: format the code using rustfmt
- `cargo clippy`: run the Clippy linter for additional checks
- `cargo test`: run all tests
- `cargo bench`: run benchmarks
- `cargo build`: build the project
- `cargo run`: build and run the project
- `cargo doc --open`: generate and open documentation

### Code Style

- Use `snake_case` for variable and function names.
- Use `PascalCase` for types and traits.
- Use `SCREAMING_SNAKE_CASE` for constants.
- Document public APIs with `///` doc comments.
- Use `#[must_use]` attribute for functions that return values that should be used.

### Error and Log Messages

Error and log messages should not start with a capital letter (unless it's a proper noun or acronym).

#### Examples

- Correct ✅: "unable to connect to client"
- Incorrect ❌: "Unable to connect to client"

### Testing

All changes to the core must contain proper and well-defined unit tests. Previous tests must continue to pass.
This codebase uses Rust's built-in testing framework:

- Use `#[test]` attribute for unit tests.
- Use `#[cfg(test)]` module for test-only code.
- Use `assert!`, `assert_eq!`, `assert_ne!` macros for assertions.
- Use `#[should_panic]` attribute for tests that should panic.
- Place integration tests in the `tests/` directory.

### Benchmarking

Use Rust's built-in benchmarking framework or the `criterion` crate for benchmarks:

- Use `#[bench]` attribute for benchmarks (requires nightly Rust).
- For stable Rust, use the `criterion` crate for more detailed benchmarking.
- Run benchmarks with `cargo bench`.

### Help Messages

Follow these rules for help messages for CLI commands and flags:

- Help string should not start with a capital letter.
- Don't include the default value in the help string.
- Include the acceptable range for the flags that accept a range of values.

## Commit Guidelines

Please follow these rules when committing changes to the Immortal:

- Each commit should represent a single, atomic change to the codebase.
  Avoid making multiple unrelated changes in a single commit.
- Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for commit messages and
  Pull Request titles.

### Commit type

List of conventional commit [types](https://github.com/commitizen/conventional-commit-types/blob/master/index.json):

| Types    | Description                                                                       |
| -------- | --------------------------------------------------------------------------------- |
| fix      | A big fix                                                                         |
| feat     | A new feature                                                                     |
| docs     | Documentation only changes                                                        |
| test     | Adding missing tests or correcting existing tests                                 |
| build    | Changes that affect the build system or external dependencies                     |
| ci       | Changes to our CI configuration files and scripts                                 |
| perf     | A code change that improves performance                                           |
| refactor | A code change that neither fixes a bug nor adds a feature                         |
| style    | Changes that do not affect the meaning of the code (white-space, formatting, etc) |
| chore    | Other changes that don't modify src or test files                                 |

### Commit Scope

The scope helps specify which part of the code is affected by your commit.
It must be included in the commit message to provide clarity.
Multiple scopes can be used if the changes impact several areas.

### Commit Description

- Keep the commit message under 50 characters.
- Start the commit message with a lowercase letter and do not end with punctuation.
- Write commit messages in the imperative: "fix bug", not "fixed bug" or "fixes bug".

### Examples

  - Correct ✅: "feat(ws): close stale connections."
  - Correct ✅: "feat(ws, config): max_wss connection limit"
  - Incorrect ❌: 'feat(config): Blacklist npubs"
  - Incorrect ❌: 'feat(config): blacklisted npubs"

-------------------------------------------------

Thank you for your contributions to the Immortal!
