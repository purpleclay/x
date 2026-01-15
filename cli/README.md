# cli

A cobra-infused CLI starter kit with custom help rendering inspired by Rust's clap.

```sh
go get github.com/purpleclay/x/cli
```

## Features

- **Custom Help Rendering**: clap-inspired help output with clean formatting and text wrapping
- **Theming**: fully customizable styles for commands, flags, headers, and more
- **Environment Variable Binding**: associate env vars with flags using cobra annotations, values are displayed in help if set
- **Flag Grouping**: organize related flags into named sections
- **Enum Flags**: type-safe enums with optional help text for each value
- **Version Flag**: automatic `--version` flag and `version` subcommand support
- **Shell Completion**: enhanced completions for 11 shells powered by [carapace](https://github.com/carapace-sh/carapace)

## Example

```
Deploy your application to the cloud with configurable options for environment,
replicas, and resource limits.

USAGE
  demo deploy [FLAGS]

EXAMPLES
  # Deploy to staging with defaults
  demo deploy --env staging

  # Deploy to production with custom replicas
  demo deploy --env production --replicas 5

FLAGS
  -e, --env <string>
          target environment (default: "staging")

  -r, --replicas <int>
          number of replicas (default: 3)

      --token <string>  [env: DEPLOY_TOKEN=sk-abc123]
          authentication token

  -l, --log-level <string>
          set the logging verbosity (default: "info")

          Possible values:
          - debug: Enable debug logging
          - info: Standard logging
          - warn: Warnings only
          - error: Errors only

AUTHENTICATION
      --url <string>
          API endpoint URL

GLOBAL FLAGS
  -c, --config <string>  [env: DEMO_CONFIG]
          path to config file
```
