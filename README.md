# Go Network Attached Storage (NAS)

This project implements a simple Network Attached Storage (NAS) system in Go, allowing users to store, retrieve, and delete files from the disk. It supports concurrency with mutex locks for file operations and offers path transformation functionality to store files based on custom or content-addressable paths.


## Getting Started

### Prerequisites
- Go 1.16 or higher
- Make (for MacOS/Linux)
- MinGW (for Windows)

### Installation
1. Clone the repository:
```sh
git clone https://github.com/sumit-behera-in/Go_Network_Attached_Storage
```

2. Install dependencies (if needed):
```sh
go mod tidy
```

#### Windows
```sh
mingw32-make run
```

#### MacOS / Linux
```sh
make run
```
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- `goLogger`: A custom logger used for logging operations.