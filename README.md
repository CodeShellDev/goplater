<p align="center">
    <img width="256" height="256" alt="Goplater Logo" src="https://github.com/codeshelldev/goplater/raw/refs/heads/main/logo/goplater.png">
</p>

<h1 align="center">Goplater</h1>

<p align="center">Goplater is a Go commandline programm that helps you template your files</p>
<br/>

## Getting Started

Download the latest binary from the Release page.
Make it executable with `chmod +x goplater` and run it for the first time.

Use the `goplater template` command to template files:

```bash
./goplater template TEMPLATE.md -o README.md
```

This will create a new file called `README.md` in your current working directory.

## Usage

### Format

Goplater uses Go's [builtin templating library](https://pkg.go.dev/text/template) therefor the syntax should be consistent with other projects.

**Example:**

```
File Content: $​{​{​{ read "./myfile.txt" }​}​}
```

### Functions

As you saw in the example above `read` is used for reading and output file contents.
But there are more as you will see in the following…

#### `read`

Reads from absolute or relative file path (depending on input),
where relative paths are relative to the invoker.

```
read "path"
```

##### `readOpts`

Same as `read` but with another parameter for additional arguments:

| Short | Long          | Type   | Note                        |
| ----- | ------------- | ------ | --------------------------- |
| `-r`  | `--recursive` | `bool` | tries templating read files |

```
read "path" "--flag1" "--flag2"
```

#### `fetch`

Performs a get http request to the specified url.

```
fetch "url"
```

#### `json`

Parses json string as dictionary (`map[string]any`).

#### `yaml`

Parses yaml string as dictionary (`map[string]any`).

#### `html`

Parses html string as html document.

##### `htmlDocFind`

Query element by selector in html document.

```
htmlDocFind ( html "html_string" ) "h3:contains['xyz']"
```

##### `htmlFind`

Query element by selector within another element.

```
htmlFind ( htmlDocFind document ) "h3:contains['xyz']"
```

##### `htmlText`

Outputs inner text of a html element.

```
htmlText ( htmlDocFind document "selector" )
```

##### `htmlAttr`

Outputs the value of the specified elements attribute.

```
htmlAttr ( htmlDocFind document "selector" ) "attribute"
```

##### `htmlInner`

Outputs element's inner html string.

```
htmlInner ( htmlDocFind document "selector" )
```

#### `trim`

Outputs trimmed string.

#### `upper`

Outputs uppercased string.

#### `lower`

Outputs lowercased string.

#### `contains`

Returns wether string contains substring.

```
contains "Sunflower" "flower"
```

#### `replace`

Replaces substring with new string in old string.

```
replace "Sunflower" "flower" "shine"
```

#### `split`

Split string by separator and return array of all parts.

```
split "Apple, Banana, Strawberry" ","
```

#### `delete`

Deletes an entry from a dictionary or array.

```
delete map "key"
```

```
delete array 0
```

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by @CodeShellDev — All Rights Reserved. Go gopher mascot originally created by [Renée French](https://instagram.com/reneefrench/), used under the [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) license.
