# advancedplayer

advanceDplayer (alternatively advancedplayer) is a music player daemon written in Go.

![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white) ![Version](https://img.shields.io/github/v/release/noorhossam312/advancedplayer?include_prereleases)

## Table of Contents

- [Features](#features)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Local file playback
- YouTube Music search and playback
- File looping
- Pause and resume playback
- Playback position tracking
- Forward/backward seeking

## Prerequisites
- `libmpv`
- `go 1.26+` (if building from source)

## Usage
Download binaries from the [Releases](https://github.com/noorhossam312/advancedplayer/releases) page, or build from source:
- Compile
```bash
CGO_ENABLED=0 go build -o bin/ ./cmd/...
```

- Start the daemon:
```bash
./bin/advancedplayerd &
```
or if you downloaded the binaries:

```bash
<path_to_advancedplayerd> &
```
- Control the daemon:
```bash
./bin/adpctl <command> [args]
```
or if you downloaded the binaries:

```bash
<path_to_adpctl> <command> [args]
```

## Contributing

### Prerequisites
- `libmpv`
- `go 1.26+`

### Build
`go mod download`

`CGO_ENABLED=0 go build -o bin/ ./cmd/...`

> [!NOTE]
> You don't need a C compiler, libmpv is loaded dynamically at runtime.

### Checks before opening a PR
`go vet ./...`

> [!NOTE]
> `-race` requires cgo to be enabled, make sure you have a C compiler and that CGO_ENABLED=1

`CGO_ENABLED=1 go build -race -o bin/ ./cmd/...`

`GORACE="halt_on_error=1" ./bin/advancedplayerd`

Then in another terminal, run this:

`./bin/adpctl play <audio_file>`

`seq 1 200 | xargs -P 16 -I{} sh -c './bin/adpctl togglepause; ./bin/adpctl seek 10; ./bin/adpctl timepos'`


### Opening a PR
> [!NOTE]
> advanceDplayer uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).
1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'feat: add amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤ by [noorhossam312](https://github.com/noorhossam312)