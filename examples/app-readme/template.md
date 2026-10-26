+{{- $app := jsonDecode (readRaw "./package.json") -}}

# +{{ upper (index $app "name") }}

Version: +{{ index $app "version" }}

+{{ read "./install.md" }}