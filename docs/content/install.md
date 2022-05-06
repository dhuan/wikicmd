# Installation

## Download it

This is the easiest installation method.

```sh
$ wget https://github.com/dhuan/wikicmd/releases/download/v0.1.0-beta-5/wikicmd_v0.1.0-beta-5_linux-386.zip
$ unzip wikicmd_v0.1.0-beta-5_linux-386.zip
```

> If you're using MacOS, you'll need to replace the download link above with the one for your platform. [Check the releases page](https://github.com/dhuan/wikicmd/releases) to see all supported platforms and download the one that matches your system.

A `wikicmd` binary should then be available to you at your current directory. You can just use it:

```sh
$ ./wikicmd
```

Optionally, you may move it to your user's programs folder in order to make it executable without specifying the program's path: 

```sh
$ mv ./wikicmd ~/bin/.
$ wikicmd
```

## Installing from source

Make sure that your system meets the following requirements before proceeding:

- [Go 1.18 or more recent](https://go.dev/)
- [GNU Make](https://www.gnu.org/software/make/)
- [Git](https://git-scm.com/)

```sh
$ git clone https://github.com/dhuan/wikicmd.git
$ cd wikicmd
$ make build
```

You should then have an executable named `wikicmd` located in the `bin` folder of the repository, which you can execute:

```sh
$ ./bin/wikicmd
```
