Define a function in the template and call it:

```
+{{ define "greet(name)" }}
    +{{ returnNext (printf "Hello, %s!" $name) }}
+{{ end }}

+{{ call "greet" "Gopher" }}
```

Output:

```
Hello, Gopher!
```

Prefer Go? Register custom Go functions as well.
