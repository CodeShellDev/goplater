Define a function in the template and call it:

+{{ GOPLATER_IGNORE }}

```
+{{ define "greet(name)" }}
    +{{ returnNext (printf "Hello, %s!" .name) }}
+{{ end }}

+{{ call "greet" "Gopher" }}
```

+{{ GOPLATER_ENDIGNORE }}

Output:

```
+{{ define "greet(name)" -}}
    +{{ returnNext (printf "Hello, %s!" .name) }}
+{{- end -}}

+{{ call "greet" "Gopher" }}
```

Prefer Go? Register custom Go functions as well.
