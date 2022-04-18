# switch

Sets the current Wiki that you're working with, given a Wiki ID.

```sh
$ wikicmd switch my_wiki
```

Your wikicmd configuration file can contain multiple Wiki Configuration Objects.

```json
{
  "config": [
    {
      "id": "my_wiki"
      ...
    },
    {
      "id": "another_wiki"
      ...
    }
  ]
}
```

Given the configuration structure above, we could switch to either `my_wiki` or `another_wiki`.

If you're using wikicmd to work with a single Wiki only, you won't be using this command.
