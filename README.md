<p align="center">
    <img width="256" height="256" alt="Goplater Logo" src="https://github.com/codeshelldev/goplater/raw/refs/heads/main/logo/goplater.png">
</p>

<h1 align="center">Goplater</h1>

<p align="center"><strong>Goplater</strong> is a Go commandline program that helps you template your files</p>

## Contents

<!-- prettier-ignore -->
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Standard Library](#standard-library)
	- [fs](#fs)
	- [outputto](#outputto)
	- [base64](#base64)
	- [container](#container)
	- [globals](#globals)
	- [import](#import)
	- [debug](#debug)
	- [html](#html)
	- [http](#http)
	- [json](#json)
	- [markdown](#markdown)
	- [funcs](#funcs)
	- [operators](#operators)
	- [regex](#regex)
	- [strings](#strings)
	- [yaml](#yaml)
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)
- [Legal](#legal)

Need help? Come join our [Matrix Server](https://matrix.to/#/#codeshelldev.oss.goplater:matrix.org)!

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
File Content: +​{​​{ read "./myfile.txt" }​}​
```

<details>
  <summary>Toggle me!</summary>

## Standard Library

### [`fs`](https://pkg.go.dev/github.com/codeshelldev/goplater/internals/template/funcs/fs)

#### `read(path string) string`

Reads and templates file.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ read "/path/to/file" }}
```
</details>
#### `readRaw(path string) string`

Reads file (without templating).

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ readRaw "/path/to/file" }}
```
</details>
#### `readArgs(path string, args any) string`

Reads file and passes arguments to it.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ read "/path/to/file" }}
```
</details>
#### `write(path string, content string) throws error`

Writes to a file path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ write "/path/to/file" "Hello" }}
```
</details>
#### `mkdir(path string) throws error`

Creates directory at path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ read "/path/to/somewhere/" }}
```
</details>
#### `appendFile(path string, content string) throws error`

Appends to an existing file.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ appendFile "/path/to/somewhere/" "Goodbye" }}
```
</details>
#### `fsExists(path string) bool`

Returns if a file or folder exists at a given path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ if (fsExists "/path/to/somewhere") }}
Path exists!
+{{ end }}
```
</details>
#### `isFile(path string) bool`

Returns if a file exists at a given path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ if (isFile "/path/to/file") }}
Path is a file!
+{{ end }}
```
</details>
#### `isDir(path string) bool`

Returns if a folder exists at a given path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ if (fsExists "/path/to/somewhere") }}
Path is a folder!
+{{ end }}
```
</details>
#### `listDir(path string) ([]string, throws error)`

Returns all files and folders inside of a given directory.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ join ", " (listDir "/path/to") }}
```
**Output**:
```
/path/to/file1.txt, /path/to/file2.txt, /path/to/folder
```

</details>
#### `walkDir(path string) ([]string, throws error)`

Returns all files and folders recursively under a given directory.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ join ", " (walkDir "/path/to") }}
```
**Output**:
```
/path/to/file1.txt, /path/to/file2.txt, /path/to/folder, /path/to/folder/file3.txt
```

</details>
#### `fsRemove(path string) throws error`

Removes a file or folder at a given path.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ fsRemove "/path/to/file/or/folder" }}
```
</details>
#### `joinPath(paths []string) string`

Joins paths together.

#### `basePath(path string) string`

Returns the last element of a path.

#### `pathDir(path string) string`

Returns all but the last element of a path.

#### `fileExt(path string) string`

Returns the file name extension of a file path.

#### `absPath(path string) (string, throws error)`

Returns an absolute representation of a given path.
If the path is not absolute it will be joined with the current working directory to turn it into an absolute path.

#### `relPath(basePath string, targetPath string) (string, throws error)`

Returns a relative path that is lexically equivalent to targetPath when joined to basePath with an intervening separator.
The returned path will always be relative to basePath, even if basePath and targetPath share no elements.

### [`outputto`](https://pkg.go.dev/github.com/codeshelldev/goplater/internals/template/funcs/outputto)

#### `outputTo(path string) throws error`

Changes the output path for the current file.
Errors if path is not allowed.

### [`base64`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/base64/base64)

#### `base64Encode(encode string) string`

Encodes strings with base64.

#### `base64Decode(decode string) (string, throws error)`

Decodes a base64 string.

### [`container`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/container/container)

#### `delete(container map[any]any|[]any, key any) (any, throws error)`

Deletes a container entry with the given key.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ $newSlice := delete (sliceCreate "Apple" "Banana") 0 }}
+{{ echo $newSlice }}
```
**Output**:
```
[Banana]
```

</details>
#### `set(container map[any]any|[]any, key any, value any) (any, throws error)`

Set a container entry value with the given key.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ $newSlice := set (sliceCreate "Apple" "Banana") 0 "Strawberry" }}
+{{ echo $newSlice }}
```
**Output**:
```
[Strawberry Banana]
```

</details>
#### `has(container map[any]any|[]any, key any) (bool, throws error)`

Returns if an entry with the given key exists in container.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ if (has (mapCreate "version" "v1") "version") }}
Container has "version"!
+{{ end }}
```
</details>
#### `slicePush(container []any, value any) []any`

Appends a value to a slice.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ $newSlice := sliceCreate "Apple" "Banana" }}
+{{ $newSlice = slicePush $newSlice "Strawberry" }}
+{{ echo $newSlice }}
```
**Output**:
```
[Apple Banana Strawberry]
```

</details>
#### `sliceCreate(values []any) []any`

Creates a new slice, optionally with values.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ $newSlice := sliceCreate "Apple" "Banana" }}
+{{ echo $newSlice }}
```
**Output**:
```
[Apple Banana]
```

</details>
#### `mapCreate(entries []any) map[string]any`

Creates a new map, optionally with key value pairs.
Pairs are constructed one by one, a key followed by the given value.
Currently only supports string keys!

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ $newMap := mapCreate "key1" "value1" "key2" "value2" }}
+{{ echo $newMap }}
```
**Output**:
```
map[key1:value1 key2:value2]
```

</details>### [`conversion`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/conversion/conversion)

#### `toString(value string) string`

Returns a given value as string (using [fmt.Sprint](https://pkg.go.dev/fmt#Sprint)).

#### `toInt(str string) (int, throws error)`

Attempts to parse a string as int.

#### `toFloat64(str string) (float64, throws error)`

Attempts to parse a string as float64.

#### `toFloat32(str string) (float32, throws error)`

Attempts to parse a string as float32.

#### `toBool(str string) (bool, throws error)`

Attempts to parse a string as a bool (using [strconv.ParseBool](https://pkg.go.dev/strconv#ParseBool)).

### [`globals`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/globals)

#### `globalSet(key string, value any)`

Defines a global value with a key.

#### `globalGet(key string) (any, throws error)`

Tries to retrieve a global value with a key.

#### `globalHas(key string) bool`

Returns whether a given global variable exists.

### [`import`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/import)

#### `import(alias string?, path string) throws error`

Tries to import a template file with func definitions.
Optionally supply an alias for easier access.
The fs resolver is only available as an internal module and does not ship with the `pkg/templating/modules` directory.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ import "https://domain.com/functions" }}
```
**Output**:
```
+{{ call "functions.greet" "John" }}
```

**Example**:

```
+{{ import "funcs" "https://domain.com/functions" }}
```

**Output**:

```
+{{ call "funcs.greet" "John" }}
```

</details>### [`return`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/return)

#### `return(i int, value any) throws error`

Sets an output value at the given index.

#### `returnNext(value any)`

Appends a value to the output.

#### `returnAll(value []any)`

Sets the output object as a whole.

<details>
  <summary>Examples</summary>
**Example**:
```
+{{ returnAll "1" "2" "3" }}
+{{ returnAll (sliceCreate "1" "2" "3") }}
```
</details>
#### `getOutputs() []any`

Returns current outputs.

### [`debug`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/debug/debug)

#### `echo(data []any)`

Outputs data to stdout.

### [`html`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/html/html)

#### `htmlDecode(html string) (*goquery.Document, throws error)`

Decodes html into a document.

#### `htmlDocFindAll(doc *goquery.Document, selector string) ([]*goquery.Selection, throws error)`

Finds elements in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

#### `htmlDocFind(doc *goquery.Document, selector string) (*goquery.Selection, throws error)`

Finds element in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

#### `htmlFindAll(doc *goquery.Document, selector string) ([]*goquery.Selection, throws error)`

Finds elements in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

#### `htmlFind(doc *goquery.Document, selector string) (*goquery.Selection, throws error)`

Finds element in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

#### `htmlText(el *goquery.Selection) string`

Returns the elements text content.

#### `htmlAttr(el *goquery.Selection, attr string) string`

Returns an elements attribute content.

#### `htmlInner(el *goquery.Selection) string`

Returns the elements inner html.

### [`http`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/http/http)

#### `fetch(url string) string`

Fetches a resource from a given url.
Outputs errors as string.

#### `fetchThrow(url string) (string, throws error)`

Fetches a resource from a given url.
Throws error instead of outputting it.

### [`json`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/json/json)

#### `jsonEncode(obj any) (string, throws error)`

Encodes an object into a json string.

#### `jsonDecode(json string) (any, throws error)`

Decodes a json string into an object.

### [`markdown`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown/markdown)

#### `markdownDecode(md string) *Document`

Decodes markdown into a document.

#### `markdownEncode(doc *Document) string`

Encodes a document into markdown.

#### `markdownFindAll(target Searchable, selector string) []*Node`

Finds elements in target by selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

#### `markdownFind(target Searchable, selector string) *Node`

Finds element in target by selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

#### `markdownIs(node *Node, selector string) bool`

Returns whether node matches selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

#### `markdownHTML(doc *Document) (string, throws error)`

Returns html markdown representation.

#### `markdownText(node *Node) string`

Returns all raw markdown text in a node recursively.

#### `markdownHeadings(doc *Document) []*Heading`

Returns all markdown headings.

#### `markdownLinks(doc *Document) []*Link`

Returns all markdown links.

#### `markdownLinkURL(link *Link) string`

Returns the url of a markdown link node.

#### `markdownLinkSetURL(link *Link, url string)`

Sets a link node's url.

#### `markdownLinkSetTitle(link *Link, title string)`

Sets a link node's title.

#### `markdownLinkText(link *Link) string`

Returns the text of a markdown link node.

#### `markdownImages(doc *Document) []*Image`

Returns all markdown images.

#### `markdownImageURL(image *Image) string`

Returns the url of a markdown image node.

#### `markdownImageSetURL(image *Image, url string)`

Sets a image node's url.

#### `markdownImageSetAlt(image *Image, alt string)`

Sets a image node's alt.

#### `markdownParagraphs(doc *Document) []*Paragraph`

Returns all markdown paragraphs.

#### `markdownCodeBlocks(doc *Document) []*CodeBlock`

Returns all markdown code blocks.

#### `markdownBlockquotes(doc *Document) []*Blockquote`

Returns all markdown block quotes.

#### `markdownRemove(node *Node)`

Removes a markdown node.

#### `markdownAppend(doc *Document, str string)`

Appends text to a markdown document.

#### `markdownPrepend(doc *Document, str string)`

Prepends text to a markdown document.

### [`funcs`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/math/funcs)

#### `abs(numbers int|float64) (int|float64, throws error)`

Returns the abs value of a number.

#### `min(a int|float64, b int|float64) (int|float64, throws error)`

Returns the smaller value of two numbers.

#### `max(a int|float64, b int|float64) (int|float64, throws error)`

Returns the bigger value of two numbers.

#### `clamp(value int|float64, min int|float64, max int|float64) (int|float64, throws error)`

Clamps a value between a minimum and maximum value.

#### `pow(base int|float64, exponent int|float64) (float64, throws error)`

Performs pow on base with exponent.

#### `sqrt(number int|float64) (float64, throws error)`

Returns the square root of a number.

#### `round(number int|float64) (int, throws error)`

Rounds a number to the nearest integer.

#### `floor(number int|float64) (int, throws error)`

Returns the greatest integer value less than or equal to number.

#### `ceil(number int|float64) (int, throws error)`

Returns the least integer value greater than or equal to number.

### [`operators`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/math/operators)

#### `add(numbers [](int|float64)) (int|float64, throws error)`

Adds numbers together.

#### `sub(numbers [](int|float64)) (int|float64, throws error)`

Subtracts numbers from each other.

#### `mult(numbers [](int|float64)) (int|float64, throws error)`

Multiplies numbers together.

#### `divd(numbers [](int|float64)) (int|float64, throws error)`

Divides numbers from each other.

#### `mod(numbers [](int|float64)) (int|float64, throws error)`

Performs modulo on numbers.

### [`regex`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/regex/regex)

#### `regexMatch(regex string) (bool, throws error)`

Returns whether a string matches the given regex pattern.

#### `regexFind(regex string) ([]string, throws error)`

Returns all matches against the given regex pattern in a string.

#### `regexFindGroups(regex string) ([][]string, throws error)`

Returns all matches and submatches from the given regex pattern in a string.

#### `regexFindGroupsIndex(regex string) ([][]int, throws error)`

Returns all match and submatch positions from the given regex pattern in a string.

#### `regexReplace(regex string) (string, throws error)`

Replaces all matches from the given regex pattern in a string with the replaceWith param (which allows using '$1', etc.).

### [`strings`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/strings/strings)

#### `trimSpace(str string) string`

Trims space from a string (leading and trailing).

#### `trim(str string, cutset string) string`

Trims all chars defined in the cutset param from a string (leading and trailing).

#### `upper(str string) string`

Uppercases a string.

#### `lower(str string) string`

Lowercases a string.

#### `contains(str string, sub string) bool`

Returns whether a string contains the given substring.

#### `count(str string, sub string) int`

Counts all occurences of a substring in the given string.

#### `startsWith(str string, prefix string) bool`

Returns whether a string starts with the given prefix.

#### `endsWith(str string, suffix string) bool`

Returns whether a string ends with the given suffix.

#### `isEmpty(str string) bool`

Returns whether the given string is empty.

#### `replace(str string, sub string, replaceWith string) string`

Replaces all occurences of a substring in the given string with the replaceWith param.

#### `split(str string, sub string, sep string) []string`

Split a string by the given seperator.

#### `after(str string, sub string) string`

Returns the string after the given substring in a string.

#### `before(str string, sub string) string`

Returns the string before the given substring in a string.

#### `between(str string, startSub string, endSub string) string`

Returns the substring between the start and end substring.

#### `cutPrefix(str string, sub string) string`

Removes a substring at the start of the given string.

#### `cutSuffix(str string, sub string) string`

Removes a substring at the end of the given string.

#### `slice(str string, start int, end int) string`

Slices a string based on start and end index.

#### `join(sep string, strings []string) string`

Joins multiple strings together by the given separator.

#### `repeat(str string, count int) string`

Repeat a string n times.

### [`yaml`](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/yaml/yaml)

#### `yamlDecode(yaml string) (any, throws error)`

Decodes a yaml string into an object.

</details>

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
