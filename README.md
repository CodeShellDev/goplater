# Goplater

Goplater is a Go commandline programm that helps you template your files.

## Getting Started

Download the `latest` binaries from the Release Page.
Make it executable with `chmod +x goplater` and run it for the first time.

Use the `goplater template` command to template files:

```bash
./goplater template TEMPLATE.md -o README.md
```

This will create a new file called `README.md` in your current working directory.

## Usage

Take a look at this `TEMPLATE.md` file:

```markdown
Wow, look at this incredible file... 🥳

`cmd/root.go`:

```go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goplater",
    Short: "Go Template CLI",
    Long:  `Go CLI Programm to Template files.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	
}
```

You want the `docker-compose.yaml` file for Secured Signal API?
Really? Here you go:

```yaml
services:
  signal-api:
    image: bbernhard/signal-cli-rest-api:latest
    container_name: signal-api
    environment:
      - MODE=normal
    volumes:
      - ./data:/home/.local/share/signal-cli
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - signal-api

  secured-signal:
    image: ghcr.io/codeshelldev/secured-signal-api:latest
    container_name: secured-signal
    environment:
      API__URL: http://signal-api:8080
      SETTINGS__VARIABLES__RECIPIENTS:
        '[+123400002, +123400003, +123400004]'
      SETTINGS__VARIABLES__NUMBER: "+123400001"
      API__TOKENS: '[LOOOOOONG_STRING]'
    ports:
      - "8880:8880"
    restart: unless-stopped
    networks:
      backend:
        aliases:
          - secured-signal-api

networks:
  backend:
```
```

Notice the `{{ #://... }}}` and `{{ @://... }}}`, these are used to include local and remot files in your Template respectively.
This Template will then include `examples/fs/INCREDIBLE.md` and `docker-compose.yaml` (from [Secured Signal API](https://github.com/CodeShellDev/secured-signal-api/blob/main/docker-compose.yaml)) in its File Content.

Which results in:

```markdown
Wow, look at this incredible file... 🥳

{{{ #://examples/fs/INCREDIBLE.md }}}

You want the `docker-compose.yaml` file for Secured Signal API?
Really? Here you go:

```yaml
{{{ @://https://raw.githubusercontent.com/CodeShellDev/secured-signal-api/refs/heads/main/docker-compose.yaml }}}
```
```

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

[MIT](https://choosealicense.com/licenses/mit/)
