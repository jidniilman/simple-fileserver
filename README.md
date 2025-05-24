# Simple File Server

A lightweight, self-contained file server written in Go using the Echo framework. This single-binary application serves files and directories from the location where it's executed, with a clean web interface.

## Features

- Single binary with no external dependencies
- Clean, responsive web interface
- Browse files and directories with ease
- View file sizes
- Mobile-friendly design
- No configuration needed - just run and go!

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jidniilman/simple-fileserver.git
   cd simple-fileserver
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the binary:
   ```bash
   go build -o fileserver
   ```

## Usage

1. Run the server:
   ```bash
   ./fileserver
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:4221
   ```

   The server will display all files and directories in the current working directory.

## Configuration

The server runs on port `4221` by default. To change the port, modify the `e.Start(":4221")` line in `main.go`.

## Building for Different Platforms

The server can be cross-compiled for various platforms:

### Linux (64-bit)
```bash
GOOS=linux GOARCH=amd64 go build -o fileserver-linux
```

### Windows (64-bit)
```bash
GOOS=windows GOARCH=amd64 go build -o fileserver.exe
```

### macOS (Intel 64-bit)
```bash
GOOS=darwin GOARCH=amd64 go build -o fileserver-macos
```

### macOS (Apple Silicon)
```bash
GOOS=darwin GOARCH=arm64 go build -o fileserver-m1
```

## Deployment

Simply copy the binary to your target machine and run it. The server is completely self-contained with all assets (including CSS) embedded in the binary.

## Security Note

By default, the server allows access to all files in the directory where it's run. Be cautious when running in directories containing sensitive information.

## License

This project is open source and available under the [MIT License](LICENSE).
