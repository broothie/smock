# `smock`
A server mock

## Installation
```bash
$ brew tap broothie/smock
$ brew install smock
```

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
