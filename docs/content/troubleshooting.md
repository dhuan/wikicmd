# Troubleshooting

If you're facing issues, it may help to use the verbose flag `-v`. It will output for you a bunch of useful log messages, including each HTTP request that's being sent, and their responses.

```sh
$ wikicmd -v edit 'Some page'
```

# F.A.Q

## I don't like the editor that opens up when I use wikicmd's `edit`. How do I use editor X?

```sh
$ export EDITOR=vscode
$ wikicmd edit 'Some page'
```
