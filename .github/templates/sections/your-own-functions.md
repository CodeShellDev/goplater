Define a function in the template and call it:

+{{ GOPLATER_IGNORE }}

<!-- prettier-ignore -->
```handlebars
+{{ define "greet(name, surname?)" }}
    +{{ $name := .name }}
    +{{ if .surname }}
        +{{ $name = concat .name " " .surname }}
    +{{ end }}
    +{{ return (printf "Hello, %s!" $name) }}
+{{ end }}

+{{ call "greet" "Gopher" }}
+{{ call "greet" "John" "Doe" }}
```

+{{ GOPLATER_ENDIGNORE }}

Output:

<!-- prettier-ignore -->
```handlebars
+{{ define "greet(name, surname?)" }}
    +{{ $name := .name }}
    +{{ if .surname }}
        +{{ $name = concat .name " " .surname }}
    +{{ end }}
    +{{ return (printf "Hello, %s!" $name) }}
+{{ end -}}

+{{ call "greet" "Gopher" }}
+{{ call "greet" "John" "Doe" }}
```

The `define` name is the function signature. Parameters are available as `.name` inside the function:

- `name` is required
- `name?` is optional (`nil` if not given)
- `...name` collects all remaining arguments into a slice

Optional and variadic parameters must come last. A `define` without a signature (e.g. `define "test"`) is a plain template and can't be called with `call`.

Prefer Go? Register custom Go functions as well.
