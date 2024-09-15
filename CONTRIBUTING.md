# Contributing

Thank you for considering contributing to the Immortal relay!
Please read these guidelines before submitting a pull request or opening an issue.

> [!NOTE]
> This code guideline must be followed for both contributors and maintainers to review the PRs.

## Code Guidelines

We strive to maintain clean, readable, and maintainable code.
Please follow these guidelines when contributing to the project:

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
- Follow the [Go Doc Comments](https://go.dev/doc/comment) guidelines.
- Follow the principles of clean code as outlined in
  Robert C. Martin's "[Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)" book.
- Write tests/benchmarks for new code or changes to existing code, and make sure all tests pass before submitting a pull request.

### Makefile Targets

There is some make file targets which you can use when developing this codebase:

- `devtools`: will install all devtools you need in development proccess.
- `unit-test`, `test`, `test-race`: runs all existing tests.
- `fmt`: formats the code using gofumpt. (run before `check` target always.)
- `check`: runs golangci-lint linter based on its [config](./.golangci.yml).
- `build`: build an immortal binary on `build/immortal` path.
- `pre-commit`: executes formatter, linter and tests.
- `compose-up`: spin ups the development docker compose which runs all requiered third-party services for developement.
- `compose-down`: stops the development docker compose stuff.
- `models-generate`: generates the SQL tables using sqlboiler, only use it when you change the database.

### Error and Log Messages

Error and log messages should not start with a capital letter (unless it's a proper noun or acronym).

#### Examples

- Correct ✅: "unable to connect to client"
- Incorrect ❌: "Unable to connect to client"

### Testing

All changes on core must contain proper and well-defined unit-tests, also previous tests must be passed as well.
This codebase used `testify` for unit tests, make sure you follow these guide for tests:

- For panic cases make sure you use `assert.Panics`
- For checking err using `assert.ErrorIs` make sure you pass expected error as second argument.
- For checking equality using `assert.Equal` make sure you pass expected value as the first argument.

### Benchmarking

Make sure you follow [this guide](https://100go.co/89-benchmarks) when you write or change benchmarks to reach an accurate result.

### Help Messages

Follow these rules for help messages for CLI commands and flags:

- Help string should not start with a capital letter.
- Don't include default value in the help string.
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
- Write commit messages in the imperative: "fix bug" not "fixed bug" or "fixes bug".

### Examples

  - Correct ✅: "feat(ws): close stale connections."
  - Correct ✅: "feat(ws, config): max_wss connection limit"
  - Incorrect ❌: 'feat(config): Blacklist npubs"
  - Incorrect ❌: 'feat(config): blacklisted npubs"

-------------------------------------------------

Thank you for your contributions to the Immortal!
