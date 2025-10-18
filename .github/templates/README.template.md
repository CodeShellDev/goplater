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
{{{ #://examples/TEMPLATE.md }}}
```

Notice the `{{{ #://... }}}` and `{{{ @://... }}}`, these are used to include local and remot files in your Template respectively.
This Template will then include `examples/data.json` and `docker-compose.yaml` (from [Secured Signal API](https://github.com/CodeShellDev/secured-signal-api/blob/main/docker-compose.yaml)) in its File Content.

Which results in:

```markdown
{{{ #://examples/TEMPLATE.md }}}
```

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

[MIT](https://choosealicense.com/licenses/mit/)
