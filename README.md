# `smock`
A server mock. I find it useful for mocking external service calls or inspecting inter-service calls while still allowing them to go through (via proxying).

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
Start a smock:
```bash
$ smock
[smock] ui @ http://localhost:9091
[smock] mock server @ http://localhost:9090
```
then, in another terminal, hit the server:
```bash
$ curl -i localhost:9090
HTTP/1.1 200 OK
Date: Tue, 04 Feb 2020 05:57:03 GMT
Content-Length: 0

```

### Customizing the mock response
Example: status 201, a couple headers, and a body. Start smock:
```bash
$ smock -c 201 -h 'Authorization: asdf' -h 'Content-Type: text/plain' -b ok
```
then:
```bash
$ curl -i localhost:9090
HTTP/1.1 201 Created
Authorization: asdf
Content-Length: 2
Content-Type: text/plain
Date: Tue, 04 Feb 2020 06:02:54 GMT

ok
```
You can also respond with the contents of a file. Start smock pointed at a file:
```bash
$ echo '{"key": "value"}' > response.json
$ smock file response.json
```
then:
```bash
$ curl localhost:9090
{"key": "value"}
```

### Proxying
Start smock proxied toward a uri:
```bash
$ smock proxy http://www.mocky.io
```
then:
```bash
$ curl localhost:9090/v2/5185415ba171ea3a00704eed
{"hello": "world"}
```

## All Options
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
