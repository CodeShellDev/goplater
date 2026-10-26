+{{- import "composer" "../composer" -}}

**`TEMPLATE.md`**

```markdown
+{{ call "composer.exampleRaw" "app-readme/template" }}
```

```bash
goplater template TEMPLATE.md -o README.md
```

**`README.md`**

```markdown
+{{ call "composer.example" "app-readme/template" }}
```

[This very README is generated with Goplater.](./github/templates/README.template.md)
