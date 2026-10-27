Keep your functions in one file and import them from the filesystem or from a URL:

+{{ GOPLATER_IGNORE }}

<!-- prettier-ignore -->
```handlebars
+{{ import "funcs" "https://raw.githubusercontent.com/CodeShellDev/goplater/refs/heads/main/examples/funcs.gplt" }}
+{{ call "funcs.greet" "John" }}
```

+{{ GOPLATER_ENDIGNORE }}

Pin remote imports to a tag or a commit and treat them like any other dependency: imported functions run inside your template.
