<!--
THIS IS A GOPLATER COMMAND HEADER,
ALL INSTRUCTIONS HAVE BEEN CONSUMED BY IT!

# imports
+{{ import "commands" "./commands" }}
+{{ import "contents" "./contents" }}
+{{ import "composer" "./composer" }}

# library commands
+{{ $cmds := call "commands.renderCommandDocs" }}
+{{ $cmdContents := call "contents.indentTabs" (call "contents.renderContentSection" $cmds 3) 1 }}
-->

<p align="center">
    <img width="256" height="256" alt="Goplater Logo" src="https://github.com/CodeShellDev/goplater/raw/refs/heads/main/logo/goplater.png">
</p>

<h1 align="center">Goplater</h1>

<div align="center">
    <h3>Go templates with superpowers</h3>
    Read files, call APIs, parse JSON, YAML and Markdown and write your own template functions<br>
</div>

<br>

<div align="center">
    <a href="https://github.com/CodeShellDev/goplater/releases">
        <img 
            src="https://img.shields.io/github/v/release/CodeShellDev/goplater?sort=semver&logo=github&label=Release" 
            alt="GitHub release">
    </a>
    <a href="https://github.com/CodeShellDev/goplater/stargazers">
    <img 
		src="https://img.shields.io/github/stars/CodeShellDev/goplater?style=flat&logo=github&label=Stars" 
		alt="GitHub stars">
    </a>
    <a href="https://matrix.to/#/#codeshelldev.oss.goplater:matrix.org">
        <img 
            alt="Matrix" 
            src="https://img.shields.io/matrix/codeshelldev.oss.goplater%3Amatrix.org?style=flat&color=blue&logo=matrix&label=Chat">
    </a>
    <a href="./LICENSE">
        <img 
            src="https://img.shields.io/badge/License-MIT-green.svg"
            alt="License: MIT">
    </a>
</div>

## See it in action

+{{ call "composer.section" "see-it-in-action" }}

## Contents

<!-- prettier-ignore -->
- [Why Goplater?](#why-goplater)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Standard Library](#standard-library)
+{{ $cmdContents }}
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)
- [Legal](#legal)

Need help? Come join our [Matrix Server](https://matrix.to/#/#codeshelldev.oss.goplater:matrix.org)!

## Why Goplater?

- **Familiar.** It is Go's [`text/template`](https://pkg.go.dev/text/template) syntax. If you know it, you already know Goplater.
- **Batteries included.** Read and write files, fetch URLs, work with JSON, YAML, HTML and Markdown, run regex and string operations. See the [Standard Library](#standard-library).
- **Write your own functions, inside the template.** Define them with `define`, call them with `call`. Arguments and return values included.
- **Share and import functions.** Put your functions in a file, import them from disk or from a URL and reuse them everywhere.
- **Not just a CLI.** The engine behind it, `gplt-engine`, is a separate library you can embed in your own Go program.
- **Modular.** The library is split into independent modules. Load only what you need. Or let the CLI load them all.
- **Automated.** Use the [GitHub Action and Bot](https://github.com/CodeShellDev/goplater-action) to render templates and automate pull requests. See [Automation](#automate-with-github-actions).

## Getting Started

Install with Go:

```bash
go install github.com/codeshelldev/goplater@latest
```

Or download the latest binary from the [Releases page](https://github.com/CodeShellDev/goplater/releases), make it executable with `chmod +x goplater` and run it.

Template your first file:

```bash
goplater template examples/app-readme/template.md -o README.md
```

This creates `README.md` in your current working directory.

## Usage

### Templating files

+{{ GOPLATER_IGNORE }}

Goplater uses Go's [built-in templating library](https://pkg.go.dev/text/template), with actions written as `+{{ ... }}`. Files that already contain plain `{{ }}` (Helm charts, GitHub Actions, Mustache) pass through untouched.

```gotmpl
File Content: +{{ read "./myfile.txt" }}
```

+{{ GOPLATER_ENDIGNORE }}

### Whole folders

Point Goplater at a directory and it mirrors the structure into the output:

```bash
goplater template ./examples -o ./out -r
```

Useful flags:

| Flag                     | What it does                                                       |
| ------------------------ | ------------------------------------------------------------------ |
| `-o`                     | Output file or directory                                           |
| `-s`                     | Source directory for `~/` paths                                    |
| `-r`                     | Walk directories                                                   |
| `--match`                | Only template files matching a regex                               |
| `--flatten`              | Drop the directory structure                                       |
| `--root`                 | Choose the root that output paths are relative to                  |
| `--allowed-outputs`      | Limit where templates may write                                    |
| `--trim-max-blank-lines` | Experimental: Remove blank lines left behind by non-output actions |
| `--dry-run`              | Show what would be written, change nothing                         |
| `--ignore-errors`        | Continue processing when errors occur                              |

### Your own functions

+{{ call "composer.section" "your-own-functions" }}

### Share functions

+{{ call "composer.section" "share-functions" }}

### Use the engine instead

`gplt-engine` is not tied to the CLI.
Load the modules you want and embed templating in your own Go application.
Browse the packages on [pkg.go.dev](https://pkg.go.dev/github.com/CodeShellDev/goplater/pkg/gplt).

## Automate with GitHub Actions

Use Goplater in your CI/CD workflows to generate files automatically.

- **[Goplater GitHub Action](https://github.com/CodeShellDev/goplater-action)**: render templates directly in your GitHub Actions workflows.
- **[Goplater Bot](https://github.com/CodeShellDev/goplater-action/blob/main/.github/workflows/goplater-bot.yml)**: automate file generation in a repository and have the results submitted as a pull request or pushed as a commit.

We use [Goplater Bot](https://github.com/CodeShellDev/goplater-action/blob/main/.github/workflows/goplater-bot.yml) as well, check him out [PR #20](https://github.com/CodeShellDev/goplater/pull/20)!

## Standard Library

`gplt-engine` includes a large and growing selection of built-in functions for use in your templates.
Every module is independent, load all of them or only the ones you need.

Goplater ships with all implemented modules in addition to `fs`-specific and cli dependent functions.

| Module                        | Use it for                                         |
| ----------------------------- | -------------------------------------------------- |
| `fs`                          | Reading, writing and walking files and directories |
| `outputto`                    | Choosing the output path from inside a template    |
| `strings`, `regex`            | Trimming, splitting, replacing, matching           |
| `json`, `yaml`                | Decoding and encoding structured data              |
| `html`, `markdown`            | Querying and editing documents                     |
| `http`                        | Fetching remote resources                          |
| `math`                        | Operators and number helpers                       |
| `container`                   | Maps and slices                                    |
| `base64`, `conversion`        | Encoding and type conversion                       |
| `import`, `return`, `globals` | Sharing functions and passing values around        |
| `debug`                       | Printing while you develop                         |

Need something that is not there? Add your own with a custom module or your own Goplater function file and share it.

+{{ $cmds }}

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/CodeShellDev) — All Rights Reserved.

Go gopher mascot originally created by [Renée French](https://instagram.com/reneefrench/),
used under the [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) license.
