# Troubleshooting

If you're facing issues, it may help to use the verbose flag `-v`. It will output for you a bunch of useful log messages, including each HTTP request that's being sent, and their responses.

```sh
$ wikicmd -v edit 'Some page'
```

# F.A.Q

## I see the error message `Failed to login with user...`

A common mistake is trying to setup wikicmd with your actual login and password information, when you should instead setup a Bot Password for your account. [Check this guide for more details.](credentials.md)

In case you're still facing this issue even after making sure that your Bot Password is set correctly, [maybe it'd be better to report the issue.](https://github.com/dhuan/wikicmd/issues)

## I don't like the editor that opens up when I use wikicmd's `edit`. How do I use editor X?

```sh
$ export EDITOR=vscode
$ wikicmd edit 'Some page'
```
