# Get your Wiki Credentials

wikicmd is meant to be used with [MediaWiki's Bot Passwords.](https://www.mediawiki.org/wiki/Manual:Bot_passwords) Getting your Bot Password is an easy process once you have your user account for Wikipedia or any MediaWiki website. This guide will show you how to do it.

> This guide links to Wikipedia but you can follow these steps in any MediaWiki instance.

## 1-) Get your account

[Signup to Wikipedia here](https://en.wikipedia.org/w/index.php?title=Special:CreateAccount) if you don't have an account yet. Then login with your account in Wikipedia.

## 2-) Setup your Bot Password

[Click here to go to the Bot Password page.](https://en.wikipedia.org/wiki/Special:BotPasswords)

Under **Bot Name**, type a bot name of your preference. You'll use it later when configuring wikicmd. Click **Create**.

Next, you'll be asked which privileges to assign to your new Bot Password. The following privileges are relevant to wikicmd - you will want to check them:

```
Import revisions
Edit existing pages
Edit protected pages
Create, edit, and move pages
Upload new files
Upload, replace, and move files
Delete pages
```

Click **Create**. MediaWiki creates a random password for you, which is shown immediately after you've created your Bot Password. Make sure not to lose the information in this page, as you'll be using it in the next step configuring wikicmd.

## 3-) Configure wikicmd with your new Bot Password

Supposing that your Wikipedia's Username is `JohnDoe` and your Bot Password is named `mybot`, and the randomly generated password is `123456789`, your wikicmd configuration should then be edited as such:

```json
{
  "wikis": [
    {
      "id": "my_wiki",
      "address": "https://en.wikipedia.org/w",
      "user": "JohnDoe@mybot",
      "password": "123456789"
    }
  ]
}
```

Having followed these steps, wikicmd should be set and ready to be used. In case of issues, [you can report bugs in this link.](https://github.com/dhuan/wikicmd/issues)
