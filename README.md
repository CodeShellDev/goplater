<p align="center">
    <img width="256" height="256" alt="Goplater Logo" src="https://github.com/codeshelldev/goplater/raw/refs/heads/main/logo/goplater.png">
</p>

<h1 align="center">Goplater</h1>

<p align="center"><strong>Goplater</strong> is a Go commandline program that helps you template your files</p>

## Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
  - [File Functions](#file-functions)
  - [String Functions](#string-functions)
  - [Math Functions](#math-functions)
  - [Container Functions](#container-functions)
  - [Parser Functions](#parser-functions)
  - [Advanced Functions](#advanced-functions)
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)
- [Legal](#legal)

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
File Content: +​{​{​{ read "./myfile.txt" }​}​}
```

### File Functions

As you saw in the example above `read` is used for reading and output file contents.
But there are more as you will see in the following…

#### `read`

Reads from absolute or relative file path (depending on input),
where relative paths are relative to the invoker.

```
read "path"
```

#### `readArgs`

Same as [`read`], but allows to supply additional arguments to the file for further processing.

```
read "path" arg1 arg2 arg3
```

Arguments are accessible in under `.args` with `{{{ index .args 0 }}}`.

#### `fetch`

Performs a get http request to the specified url.

```
fetch "url"
```

### String Functions

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

#### `count`

Returns amount of times that substring is present in string.

#### `repeat`

Repeats string n times.

```
repeat "*" 5
```

#### `startsWith`

Returns wether string starts with prefix.

```
startsWith "Sunflower" "Sun"
```

#### `endsWith`

Returns wether string ends with prefix.

```
endsWith "Sunflower" "flower"
```

#### `isEmpty`

Returns wether string is empty (`""`).

#### `indexOf`

Returns starting index of substring in string.

```
indexOf "Apple Banana Strawberry" "Apple"
```

#### `replace`

Replaces substring with new string in old string.

```
replace "Sunflower" "flower" "shine"
```

#### `split`

Split string by separator and return slice of all parts.

```
split "Apple, Banana, Strawberry" ", "
```

#### `before`

Outputs string before substring in string.

```
before "Apple Banana Strawberry" "Banana"
```

#### `after`

Outputs string after substring in string.

```
after "Apple Banana Strawberry" "Apple"
```

#### `between`

Outputs inbetween of starting and ending substring.

```
between "Apple Banana Strawberry" "Apple" "Strawberry"
```

#### `slice`

Slice string by `start` and `end` bounds.

```
slice "  Apple  " 2 7
```

#### `join`

Joins strings with separator.

```
join "Apple" "Banana" "Strawberry" ", "
```

#### `concat`

Concat multiple strings.

```
concat "Apple" "Banana" "Strawberry"
```

#### `append`

Appends another string to string.

```
append "Hello " "World!"
```

#### `regexMatch`

Outputs wether string matches regex.

```
regexMatch "[0-9]" "0123456789"
```

#### `regexFind`

Returns a list of all regex matches.

```
index ( regexFind "[1-36-9]" "01234 56789" ) 0
```

#### `regexFindGroups`

Returns a nested list (`[][]string`) of all regex submatches (groups `(.*)`).

```
index ( index ( regexFindGroups "(_*)(\S+)(_*)" "__xyz__" 0 ) 0 )
```

#### `regexReplace`

Replaces substring via regex in string.

```
regexReplace "string" "replace_regex" "replace_with"
```

### Math Functions

#### `add`

Adds two numbers together.

#### `sub`

Subtraction for two numbers.

#### `mult`

Multiplies two numbers.

#### `divd`

Divides a through b.

#### `mod`

Performs a modulo b.

### Conversion Functions

Functions for converting types.

#### `toString`

Returns value as string (via `fmt.Sprint()`).

#### `toInt`

Parses string as int.

#### `toFloat64`

Parses string as float64.

#### `toFloat32`

Same as [`toFloat64`](#tofloat64), but for float32.

#### `toBool`

Parses string as bool.

### Container Functions

The following are functions for slices and maps.

#### `has`

Returns wether map or slice has key.

#### `includes`

Returns wether map or slice includes value.

#### `delete`

Deletes an entry from a map or slice.

```
delete map "key"
```

```
delete slice 0
```

#### `set`

Sets key in map or slice to value.

```
set map "key" value
```

```
set slice 0 value
```

#### `slicePush`

Pushes value on top of slice.

```
slicePush slice value
```

#### `sliceCreate`

Creates slice by using arguments as items.

```
sliceCreate "a" "b" "c"
```

### Parser Functions

The following section is dedicated to parser functions, for example json.

#### `jsonDecode`

Parses json string as map.

#### `jsonEncode`

Returns json string from object.

#### `yamlDecode`

Parses yaml string as map.

#### `yamlEncode`

Returns yaml string from object.

#### `base64Decode`

Decodes base64 into raw string.

#### `base64Encode`

Encodes raw string into base64.

#### `htmlDecode`

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

### Advanced Functions

In addition to the function in the [Simple Functions](#simple-functions) section, there are also some functions for advanced usage.

#### `import`

Imports a file and **executes** it as template, output is **discarded**.

```
import "functions.inc.gtmpl"
```

#### `globalSet`

Sets key globally to value.

```
globalSet "key" value
```

#### `globalGet`

Returns global value at key.

#### `funcDefine`

Defines a global function.

```
funcDefine "name" `
    {{{ return 0 "Hello World" }}}
`
```

The second argument is the template body, notice the `{{{ ... }}}` instead of `$​{​{​{ ... }​}​}`.

Raw output is discarded only output via [`return`](#return) persists.

##### `return`

> [!WARNING]
> This function is **only** accessible from within functions!

Sets return argument at **index** to **value**.

```
funcDefine "helloWorld" "{{{ return 0 "Hello World!" }}}"
```

**Overwriting** previous return arguments is possible.

##### `returnNext`

Same as [`return`](#return), but **appends** to output.

Shorthand for:

```
return i+1 value
```

```
returnNext value
```

##### `returnAll`

Sets return output slice directly with **multiple** arguments.

```
returnAll out1 out2 out3
```

##### `returnOutputs`

Sets return output slice directly with **one slice**.

```
returnOutputs slice
```

##### `getOutputs`

Returns the whole output slice.

#### `funcCall`

Calls a global function by its name (without passing any arguments).

```
funcCall "name"
```

Returns list of [`return`](#return) outputs in order of index, or if applicable only a single output.

#### `funcCallArgs`

Same as [`funcCall`](#funccall), but arguments can be passed.

```
funcCallArgs "name" arg1 arg2
```

Arguments are accessible in function body with `{{{ index .args 0 }}}`.

## Contributing

Found a bug or just want to change or add something?
Feel free to open up an issue or a PR!

## Support

Like this Project? Or just want to help?
Why not ⭐️ this Repo? :)

## License

This Project is licensed under the [MIT License](./LICENSE).

## Legal

Logo designed by [@CodeShellDev](https://github.com/codeshelldev) — All Rights Reserved. Go gopher mascot originally created by [Renée French](https://instagram.com/reneefrench/), used under the [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) license.
