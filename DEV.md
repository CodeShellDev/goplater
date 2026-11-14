## Common Mistakes & Solutions

- **Template funcs which only return an error must be avoided.** They will not throw **ANY** errors, use `(string, error)` with empty string returns instead!

## Code Standards

- **Do not shorten code that doesn't need to be.** There is no harm in leaving a bit of space between lines to improve readability (please don't `if err := ...; err == nil`).
- **Please use a code formatter for markdown documents.**
- ...
