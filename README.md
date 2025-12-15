# Goplater

Goplater is a Go commandline programm that helps you template your files.

## Getting Started

Download the latest binary from the Release page.
Make it executable with `chmod +x goplater` and run it for the first time.

Use the `goplater template` command to template files:

```bash
./goplater template TEMPLATE.md -o README.md
```

This will create a new file called `README.md` in your current working directory.

## Usage

Take a look at this `TEMPLATE.md` file:

---

```md
Wow, look at this incredible file... 🥳
```

`{{{ #://examples/fs/INCREDIBLE.md }}}`

```md
You want the `docker-compose.yaml` file for Secured Signal API?
Really? Here you go:
```

```yaml
{{{ @://https://raw.githubusercontent.com/CodeShellDev/secured-signal-api/refs/heads/docs/docker-compose.yaml }}}
```

---

Notice the `{­{­{ #://... }­}­}` and `{­{­{ @://... }­}­}`, these are used to include local and remote files in your Template respectively.
This Template will then include `examples/fs/INCREDIBLE.md` and `docker-compose.yaml` (from [Secured Signal API](https://github.com/CodeShellDev/secured-signal-api/blob/main/docker-compose.yaml)) in its File Content.

Which results in:

---

```md
Wow, look at this incredible file... 🥳
```

````go
// src: cmd/root.go

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
````

```md
You want the `docker-compose.yaml` file for Secured Signal API?
Really? Here you go:
```

```yaml
404: Not Found
```

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

[MIT](https://choosealicense.com/licenses/mit/)
