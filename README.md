# `smock`
A server mock

## Installation

### Mac
```bash
$ brew tap broothie/smock
$ brew install smock
```

### Releases
Releases are available on the [releases page](https://github.com/broothie/smock/releases).

### Source
You can also build from source of course if you have the Go toolchain installed and feel like doing that.

## Usage
```bash
$ smock --help-long
usage: main [<flags>] <command> [<args> ...]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=9090          port to run server mock on
  -u, --uiport=9091        port to run ui on
      --no-ui              disable ui
  -c, --code=200           response status code
  -h, --header=HEADER ...  response headers
  -b, --body=""            response body

Commands:
  help [<command>...]
    Show help.

  version
    print smock version


  mock [<flags>]
    mock response

    -c, --code=200           response status code
    -h, --header=HEADER ...  response headers
    -b, --body=""            response body

  file <filename>
    mock response from file


  proxy <target>
    reverse proxy to target url




```
