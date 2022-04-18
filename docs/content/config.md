# Configuring

wikicmd looks for `/home/myuser/.wikicmd.json` in order to read your configuration. You can either create this file by hand, or you can run the [config command](cmd_config.md) where a command-line wizard will help you to accomplish the same.

```sh
$ wikicmd config
Wiki address: (https://wikipedia.org) 
Login: myuser
Password: mypassword

Next, a configuration file will be created for you and saved as /home/myuser/.wikicmd.json

Is this OK? (yes):    
Done!
```

## Configuration Structure

A configuration file is formatted as follows:

```json
{
  "config": [
    {
      "id": "my_wiki",
      "address": "https://wikipedia.org",
      "user": "myuser",
      "password": "mypassword"
    }
  ]
}
```

The `config` field takes a list of "Wiki Configuration Objects". In the example above, we have only one Wiki, `my_wiki`, that we want to manage with wikicmd. If you want to configure wikicmd to be able to use multiple wikis, make sure to read about the [switch command](cmd_switch.md).

## Configuration Parameters

**id**: An ID to identify a Wiki with.

**address**: A Wiki URL. For example `https://wikipedia.org`.

**user**: A username that you can login with.

**password**: Your password.
