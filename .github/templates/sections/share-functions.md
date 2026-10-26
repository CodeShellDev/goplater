Keep your functions in one file and import them from the filesystem or from a URL:

```
+{{ import "funcs" "https://github.com/CodeShellDev/goplater/blob/main/examples/importing-funcs/funcs" }}
+{{ call "funcs.greet" "John" }}
```

Pin remote imports to a tag or a commit and treat them like any other dependency: imported functions run inside your template.
