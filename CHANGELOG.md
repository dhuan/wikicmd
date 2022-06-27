# Changelog

## wikicmd 0.2.1

This update contains only a minor bugfix - editor programs, opened when using wikicmd's `edit` command, were receiving filenames with line breaks. Some programs have issues with filenames formatted as such, causing wikicmd to crash. This version fixes that.

## wikicmd 0.2.0

- Config param renamed from `config` to `wikis`.
- New configuration parameter `editor`.
- "txt" files can be imported.
- `import` can now have change summary, with the `-m` flag.
- Better error messages - exit gracefully instead of panicking when there are no Wikis configured.
