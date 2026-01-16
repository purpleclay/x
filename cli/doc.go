// Package cli provides a cobra-infused CLI starter kit with custom help
// rendering inspired by Rust's clap. It also provides type-safe enum flag
// support via [Enum] with optional help text for each value.
//
// Shell completion is powered by [carapace], providing enhanced completions
// for 11 shells beyond Cobra's standard offering. Use [Completer] types like
// [Files], [Directories], and [Values] to define completions declaratively.
//
// [carapace]: https://github.com/carapace-sh/carapace
package cli
