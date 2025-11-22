<p align="center"> 
    <img alt="immortal" src="./.assets/images/immo.svg" width="150" height="150" />
</p>

<h1 align="center">
Immortal
</h1>

<br/>


The Immortal is a [Nostr](https://nostr.com) relay implementation in Rust.

Immortal is aimed and designed to be **scalable**, **high-performance**, and **configurable**. It's a good choice for paid relays or big community relays, and not a good choice for a personal relay.

## Installation & Running

### Prerequisites

- Rust 1.70 or later
- Cargo (comes with Rust)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/dezh-tech/immortal.git
cd immortal

# Build the project
cargo build --release

# Run the relay
cargo run --release
```

### Configuration

The relay can be configured using a `config.toml` file. See the example configuration in the repository.

## Updates

Updates, changes, or community discussions can be followed on the GitHub issue/discussion or the [Dezh Technologies Nostr profile](https://njump.me/dezh.tech).

## Contribution

All kinds of contributions are welcome!
Read the [Contribution guides](./CONTRIBUTING.md) before any code contribution.

## Donation

Donations and financial support for the development process are possible using Bitcoin and Lightning:

**on-chain**:

```
bc1qfw30k9ztahppatweycnll05rzmrn6u07slehmc
```

**lightning**: 

```
donate@dezh.tech
```

## License

The Immortal software is published under the [MIT License](./LICENSE), and contributing to and using this code means you agree with the license.
