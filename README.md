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

### [**fs**](https://pkg.go.dev/github.com/codeshelldev/goplater/internals/template/funcs/fs)


<table align="left" width="100%">
<tr><th align="left"><h4><code>read(path string) string</code></h4></th>
</tr><tr><td>
Reads and templates file.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ read "/path/to/file" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>readRaw(path string) string</code></h4></th>
</tr><tr><td>
Reads file (without templating).

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ readRaw "/path/to/file" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>readArgs(path string, args any) string</code></h4></th>
</tr><tr><td>
Reads file and passes arguments to it.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ readArgs "/path/to/file" (sliceCreate "my" "args") }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>write(path string, content string) throws error</code></h4></th>
</tr><tr><td>
Writes to a file path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ write "/path/to/file" "Hello" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>mkdir(path string) throws error</code></h4></th>
</tr><tr><td>
Creates directory at path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ mkdir "/path/to/somewhere/" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>appendFile(path string, content string) throws error</code></h4></th>
</tr><tr><td>
Appends to an existing file.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ appendFile "/path/to/somewhere/" "Goodbye" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>fsExists(path string) bool</code></h4></th>
</tr><tr><td>
Returns if a file or folder exists at a given path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ if (fsExists "/path/to/somewhere") }}
Path exists!
+{{ end }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>isFile(path string) bool</code></h4></th>
</tr><tr><td>
Returns if a file exists at a given path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ if (isFile "/path/to/file") }}
Path is a file!
+{{ end }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>isDir(path string) bool</code></h4></th>
</tr><tr><td>
Returns if a folder exists at a given path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ if (isDir "/path/to/somewhere") }}
Path is a folder!
+{{ end }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>listDir(path string) ([]string, throws error)</code></h4></th>
</tr><tr><td>
Returns all files and folders inside of a given directory.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ join ", " (listDir "/path/to") }}
</code></pre>
**Output**:
<pre><code>
/path/to/file1.txt, /path/to/file2.txt, /path/to/folder
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>walkDir(path string) ([]string, throws error)</code></h4></th>
</tr><tr><td>
Returns all files and folders recursively under a given directory.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ join ", " (walkDir "/path/to") }}
</code></pre>
**Output**:
<pre><code>
/path/to/file1.txt, /path/to/file2.txt, /path/to/folder, /path/to/folder/file3.txt
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>fsRemove(path string) throws error</code></h4></th>
</tr><tr><td>
Removes a file or folder at a given path.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ fsRemove "/path/to/file/or/folder" }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>joinPath(paths []string) string</code></h4></th>
</tr><tr><td>
Joins paths together.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>basePath(path string) string</code></h4></th>
</tr><tr><td>
Returns the last element of a path.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>pathDir(path string) string</code></h4></th>
</tr><tr><td>
Returns all but the last element of a path.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>fileExt(path string) string</code></h4></th>
</tr><tr><td>
Returns the file name extension of a file path.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>absPath(path string) (string, throws error)</code></h4></th>
</tr><tr><td>
Returns an absolute representation of a given path.
If the path is not absolute it will be joined with the current working directory to turn it into an absolute path.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>relPath(basePath string, targetPath string) (string, throws error)</code></h4></th>
</tr><tr><td>
Returns a relative path that is lexically equivalent to targetPath when joined to basePath with an intervening separator.
The returned path will always be relative to basePath, even if basePath and targetPath share no elements.

</td></tr></table>
### [**outputto**](https://pkg.go.dev/github.com/codeshelldev/goplater/internals/template/funcs/outputto)


<table align="left" width="100%">
<tr><th align="left"><h4><code>outputTo(path string) throws error</code></h4></th>
</tr><tr><td>
Changes the output path for the current file.
Errors if path is not allowed.

</td></tr></table>
### [**base64**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/base64/base64)


<table align="left" width="100%">
<tr><th align="left"><h4><code>base64Encode(encode string) string</code></h4></th>
</tr><tr><td>
Encodes strings with base64.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>base64Decode(decode string) (string, throws error)</code></h4></th>
</tr><tr><td>
Decodes a base64 string.

</td></tr></table>
### [**container**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/container/container)


<table align="left" width="100%">
<tr><th align="left"><h4><code>delete(container map[any]any|[]any, key any) (any, throws error)</code></h4></th>
</tr><tr><td>
Deletes a container entry with the given key.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ $newSlice := delete (sliceCreate "Apple" "Banana") 0 }}
+{{ echo $newSlice }}
</code></pre>
**Output**:
<pre><code>
[Banana]
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>set(container map[any]any|[]any, key any, value any) (any, throws error)</code></h4></th>
</tr><tr><td>
Set a container entry value with the given key.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ $newSlice := set (sliceCreate "Apple" "Banana") 0 "Strawberry" }}
+{{ echo $newSlice }}
</code></pre>
**Output**:
<pre><code>
[Strawberry Banana]
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>has(container map[any]any|[]any, key any) (bool, throws error)</code></h4></th>
</tr><tr><td>
Returns if an entry with the given key exists in container.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ if (has (mapCreate "version" "v1") "version") }}
Container has "version"!
+{{ end }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>slicePush(container []any, value any) []any</code></h4></th>
</tr><tr><td>
Appends a value to a slice.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ $newSlice := sliceCreate "Apple" "Banana" }}
+{{ $newSlice = slicePush $newSlice "Strawberry" }}
+{{ echo $newSlice }}
</code></pre>
**Output**:
<pre><code>
[Apple Banana Strawberry]
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>sliceCreate(values []any) []any</code></h4></th>
</tr><tr><td>
Creates a new slice, optionally with values.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ $newSlice := sliceCreate "Apple" "Banana" }}
+{{ echo $newSlice }}
</code></pre>
**Output**:
<pre><code>
[Apple Banana]
</code></pre>


</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>mapCreate(entries []any) map[string]any</code></h4></th>
</tr><tr><td>
Creates a new map, optionally with key value pairs.
Pairs are constructed one by one, a key followed by the given value.
Currently only supports string keys!

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ $newMap := mapCreate "key1" "value1" "key2" "value2" }}
+{{ echo $newMap }}
</code></pre>
**Output**:
<pre><code>
map[key1:value1 key2:value2]
</code></pre>


</details>

</td></tr></table>
### [**conversion**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/conversion/conversion)


<table align="left" width="100%">
<tr><th align="left"><h4><code>toString(value string) string</code></h4></th>
</tr><tr><td>
Returns a given value as string (using [fmt.Sprint](https://pkg.go.dev/fmt#Sprint)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>toInt(str string) (int, throws error)</code></h4></th>
</tr><tr><td>
Attempts to parse a string as int.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>toFloat64(str string) (float64, throws error)</code></h4></th>
</tr><tr><td>
Attempts to parse a string as float64.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>toFloat32(str string) (float32, throws error)</code></h4></th>
</tr><tr><td>
Attempts to parse a string as float32.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>toBool(str string) (bool, throws error)</code></h4></th>
</tr><tr><td>
Attempts to parse a string as a bool (using [strconv.ParseBool](https://pkg.go.dev/strconv#ParseBool)).

</td></tr></table>
### [**globals**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/globals)


<table align="left" width="100%">
<tr><th align="left"><h4><code>globalSet(key string, value any)</code></h4></th>
</tr><tr><td>
Defines a global value with a key.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>globalGet(key string) (any, throws error)</code></h4></th>
</tr><tr><td>
Tries to retrieve a global value with a key.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>globalHas(key string) bool</code></h4></th>
</tr><tr><td>
Returns whether a given global variable exists.

</td></tr></table>
### [**import**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/import)


<table align="left" width="100%">
<tr><th align="left"><h4><code>import(alias string?, path string) throws error</code></h4></th>
</tr><tr><td>
Tries to import a template file with func definitions.
Optionally supply an alias for easier access.
The fs resolver is only available as an internal module and does not ship with the `pkg/templating/modules` directory.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ import "https://domain.com/functions" }}
</code></pre>
**Output**:
<pre><code>
+{{ call "functions.greet" "John" }}
</code></pre>


**Example**:
<pre><code>
+{{ import "funcs" "https://domain.com/functions" }}
</code></pre>
**Output**:
<pre><code>
+{{ call "funcs.greet" "John" }}
</code></pre>


</details>

</td></tr></table>
### [**return**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/core/return)


<table align="left" width="100%">
<tr><th align="left"><h4><code>return(i int, value any) throws error</code></h4></th>
</tr><tr><td>
Sets an output value at the given index.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>returnNext(value any)</code></h4></th>
</tr><tr><td>
Appends a value to the output.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>returnAll(value []any)</code></h4></th>
</tr><tr><td>
Sets the output object as a whole.

<br/>
<details>
  <summary>Examples</summary>
  <br/>

**Example**:
<pre><code>
+{{ returnAll "1" "2" "3" }}
+{{ returnAll (sliceCreate "1" "2" "3") }}
</code></pre>

</details>

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>getOutputs() []any</code></h4></th>
</tr><tr><td>
Returns current outputs.

</td></tr></table>
### [**debug**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/debug/debug)


<table align="left" width="100%">
<tr><th align="left"><h4><code>echo(data []any)</code></h4></th>
</tr><tr><td>
Outputs data to stdout.

</td></tr></table>
### [**html**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/html/html)


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlDecode(html string) (*goquery.Document, throws error)</code></h4></th>
</tr><tr><td>
Decodes html into a document.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlDocFindAll(doc *goquery.Document, selector string) ([]*goquery.Selection, throws error)</code></h4></th>
</tr><tr><td>
Finds elements in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlDocFind(doc *goquery.Document, selector string) (*goquery.Selection, throws error)</code></h4></th>
</tr><tr><td>
Finds element in document by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlFindAll(doc *goquery.Document, selector string) ([]*goquery.Selection, throws error)</code></h4></th>
</tr><tr><td>
Finds elements in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlFind(doc *goquery.Document, selector string) (*goquery.Selection, throws error)</code></h4></th>
</tr><tr><td>
Finds element in element by selector (using [goquery.Find](https://pkg.go.dev/github.com/PuerkitoBio/goquery#Selection.Find)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlText(el *goquery.Selection) string</code></h4></th>
</tr><tr><td>
Returns the elements text content.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlAttr(el *goquery.Selection, attr string) string</code></h4></th>
</tr><tr><td>
Returns an elements attribute content.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>htmlInner(el *goquery.Selection) string</code></h4></th>
</tr><tr><td>
Returns the elements inner html.

</td></tr></table>
### [**http**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/http/http)


<table align="left" width="100%">
<tr><th align="left"><h4><code>fetch(url string) string</code></h4></th>
</tr><tr><td>
Fetches a resource from a given url.
Outputs errors as string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>fetchThrow(url string) (string, throws error)</code></h4></th>
</tr><tr><td>
Fetches a resource from a given url.
Throws error instead of outputting it.

</td></tr></table>
### [**json**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/json/json)


<table align="left" width="100%">
<tr><th align="left"><h4><code>jsonEncode(obj any) (string, throws error)</code></h4></th>
</tr><tr><td>
Encodes an object into a json string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>jsonDecode(json string) (any, throws error)</code></h4></th>
</tr><tr><td>
Decodes a json string into an object.

</td></tr></table>
### [**markdown**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown/markdown)


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownDecode(md string) *Document</code></h4></th>
</tr><tr><td>
Decodes markdown into a document.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownEncode(doc *Document) string</code></h4></th>
</tr><tr><td>
Encodes a document into markdown.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownFindAll(target Searchable, selector string) []*Node</code></h4></th>
</tr><tr><td>
Finds elements in target by selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownFind(target Searchable, selector string) *Node</code></h4></th>
</tr><tr><td>
Finds element in target by selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownIs(node *Node, selector string) bool</code></h4></th>
</tr><tr><td>
Returns whether node matches selector (using [custom](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/markdown#TypeMatches)).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownHTML(doc *Document) (string, throws error)</code></h4></th>
</tr><tr><td>
Returns html markdown representation.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownText(node *Node) string</code></h4></th>
</tr><tr><td>
Returns all raw markdown text in a node recursively.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownHeadings(doc *Document) []*Heading</code></h4></th>
</tr><tr><td>
Returns all markdown headings.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownLinks(doc *Document) []*Link</code></h4></th>
</tr><tr><td>
Returns all markdown links.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownLinkURL(link *Link) string</code></h4></th>
</tr><tr><td>
Returns the url of a markdown link node.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownLinkSetURL(link *Link, url string)</code></h4></th>
</tr><tr><td>
Sets a link node's url.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownLinkSetTitle(link *Link, title string)</code></h4></th>
</tr><tr><td>
Sets a link node's title.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownLinkText(link *Link) string</code></h4></th>
</tr><tr><td>
Returns the text of a markdown link node.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownImages(doc *Document) []*Image</code></h4></th>
</tr><tr><td>
Returns all markdown images.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownImageURL(image *Image) string</code></h4></th>
</tr><tr><td>
Returns the url of a markdown image node.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownImageSetURL(image *Image, url string)</code></h4></th>
</tr><tr><td>
Sets a image node's url.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownImageSetAlt(image *Image, alt string)</code></h4></th>
</tr><tr><td>
Sets a image node's alt.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownParagraphs(doc *Document) []*Paragraph</code></h4></th>
</tr><tr><td>
Returns all markdown paragraphs.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownCodeBlocks(doc *Document) []*CodeBlock</code></h4></th>
</tr><tr><td>
Returns all markdown code blocks.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownBlockquotes(doc *Document) []*Blockquote</code></h4></th>
</tr><tr><td>
Returns all markdown block quotes.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownRemove(node *Node)</code></h4></th>
</tr><tr><td>
Removes a markdown node.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownAppend(doc *Document, str string)</code></h4></th>
</tr><tr><td>
Appends text to a markdown document.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>markdownPrepend(doc *Document, str string)</code></h4></th>
</tr><tr><td>
Prepends text to a markdown document.

</td></tr></table>
### [**funcs**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/math/funcs)


<table align="left" width="100%">
<tr><th align="left"><h4><code>abs(numbers int|float64) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Returns the abs value of a number.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>min(a int|float64, b int|float64) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Returns the smaller value of two numbers.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>max(a int|float64, b int|float64) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Returns the bigger value of two numbers.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>clamp(value int|float64, min int|float64, max int|float64) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Clamps a value between a minimum and maximum value.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>pow(base int|float64, exponent int|float64) (float64, throws error)</code></h4></th>
</tr><tr><td>
Performs pow on base with exponent.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>sqrt(number int|float64) (float64, throws error)</code></h4></th>
</tr><tr><td>
Returns the square root of a number.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>round(number int|float64) (int, throws error)</code></h4></th>
</tr><tr><td>
Rounds a number to the nearest integer.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>floor(number int|float64) (int, throws error)</code></h4></th>
</tr><tr><td>
Returns the greatest integer value less than or equal to number.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>ceil(number int|float64) (int, throws error)</code></h4></th>
</tr><tr><td>
Returns the least integer value greater than or equal to number.

</td></tr></table>
### [**operators**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/math/operators)


<table align="left" width="100%">
<tr><th align="left"><h4><code>add(numbers [](int|float64)) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Adds numbers together.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>sub(numbers [](int|float64)) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Subtracts numbers from each other.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>mult(numbers [](int|float64)) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Multiplies numbers together.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>divd(numbers [](int|float64)) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Divides numbers from each other.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>mod(numbers [](int|float64)) (int|float64, throws error)</code></h4></th>
</tr><tr><td>
Performs modulo on numbers.

</td></tr></table>
### [**regex**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/regex/regex)


<table align="left" width="100%">
<tr><th align="left"><h4><code>regexMatch(regex string) (bool, throws error)</code></h4></th>
</tr><tr><td>
Returns whether a string matches the given regex pattern.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>regexFind(regex string) ([]string, throws error)</code></h4></th>
</tr><tr><td>
Returns all matches against the given regex pattern in a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>regexFindGroups(regex string) ([][]string, throws error)</code></h4></th>
</tr><tr><td>
Returns all matches and submatches from the given regex pattern in a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>regexFindGroupsIndex(regex string) ([][]int, throws error)</code></h4></th>
</tr><tr><td>
Returns all match and submatch positions from the given regex pattern in a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>regexReplace(regex string) (string, throws error)</code></h4></th>
</tr><tr><td>
Replaces all matches from the given regex pattern in a string with the replaceWith param (which allows using '$1', etc.).

</td></tr></table>
### [**strings**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/strings/strings)


<table align="left" width="100%">
<tr><th align="left"><h4><code>trimSpace(str string) string</code></h4></th>
</tr><tr><td>
Trims space from a string (leading and trailing).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>trim(str string, cutset string) string</code></h4></th>
</tr><tr><td>
Trims all chars defined in the cutset param from a string (leading and trailing).

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>upper(str string) string</code></h4></th>
</tr><tr><td>
Uppercases a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>lower(str string) string</code></h4></th>
</tr><tr><td>
Lowercases a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>contains(str string, sub string) bool</code></h4></th>
</tr><tr><td>
Returns whether a string contains the given substring.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>count(str string, sub string) int</code></h4></th>
</tr><tr><td>
Counts all occurences of a substring in the given string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>startsWith(str string, prefix string) bool</code></h4></th>
</tr><tr><td>
Returns whether a string starts with the given prefix.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>endsWith(str string, suffix string) bool</code></h4></th>
</tr><tr><td>
Returns whether a string ends with the given suffix.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>isEmpty(str string) bool</code></h4></th>
</tr><tr><td>
Returns whether the given string is empty.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>replace(str string, sub string, replaceWith string) string</code></h4></th>
</tr><tr><td>
Replaces all occurences of a substring in the given string with the replaceWith param.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>split(str string, sub string, sep string) []string</code></h4></th>
</tr><tr><td>
Split a string by the given seperator.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>after(str string, sub string) string</code></h4></th>
</tr><tr><td>
Returns the string after the given substring in a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>before(str string, sub string) string</code></h4></th>
</tr><tr><td>
Returns the string before the given substring in a string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>between(str string, startSub string, endSub string) string</code></h4></th>
</tr><tr><td>
Returns the substring between the start and end substring.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>cutPrefix(str string, sub string) string</code></h4></th>
</tr><tr><td>
Removes a substring at the start of the given string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>cutSuffix(str string, sub string) string</code></h4></th>
</tr><tr><td>
Removes a substring at the end of the given string.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>slice(str string, start int, end int) string</code></h4></th>
</tr><tr><td>
Slices a string based on start and end index.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>join(sep string, strings []string) string</code></h4></th>
</tr><tr><td>
Joins multiple strings together by the given separator.

</td></tr></table>


<table align="left" width="100%">
<tr><th align="left"><h4><code>repeat(str string, count int) string</code></h4></th>
</tr><tr><td>
Repeat a string n times.

</td></tr></table>
### [**yaml**](https://pkg.go.dev/github.com/codeshelldev/goplater/pkg/templating/modules/yaml/yaml)


<table align="left" width="100%">
<tr><th align="left"><h4><code>yamlDecode(yaml string) (any, throws error)</code></h4></th>
</tr><tr><td>
Decodes a yaml string into an object.

</td></tr></table>


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
