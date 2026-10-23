# advancedplayer

advanceDplayer (alternatively advancedplayer) is a music player daemon written in Go.

![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white) ![Version](https://img.shields.io/badge/version-0.1.0--alpha-green.svg)

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

### Tech Stack

- **Go**

## Prerequisites
- `libmpv`
- `go` (if building from source)

## Usage

Start the daemon:
```bash
./advancedplayerd
```
Control the daemon:
```bash
./adpctl <command> [args]
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'Add amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤ by [noorhossam312](https://github.com/noorhossam312)