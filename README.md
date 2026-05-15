# pass;D

`passman` is a lightweight command-line password manager written in Go.  
It provides a simple way to securely store, retrieve, and manage credentials locally from the terminal.

The project is built for learning backend and systems programming concepts such as:
- CLI application development
- file handling
- JSON serialization
- encryption
- modular architecture
- secure credential storage


## Features

- Add credentials
- Retrieve stored passwords
- List saved services
- Delete entries
- Local vault storage
- Simple terminal interface
- Cross-platform support


## Project Structure

```text
pass;D/
├── cmd/
│   └── main.go
├── internal/
│   ├── cli/
│   ├── crypto/
│   ├── models/
│   └── storage/
├── data/
├── go.mod
└── README.md
```


## Installation

Clone the repository:

```bash
git clone <your-repo-url>
cd passman
```

Initialize dependencies:

```bash
go mod tidy
```

Run the project:

```bash
go run cmd/main.go
```


## Usage

### Add a password

```bash
go run cmd/main.go add
```

### List saved entries

```bash
go run cmd/main.go list
```

### Retrieve a password

```bash
go run cmd/main.go get github
```

### Delete an entry

```bash
go run cmd/main.go delete github
```


## Planned Features

- AES encrypted vault
- Master password authentication
- Password generator
- Clipboard support
- Search functionality
- TUI interface
- Export/import support

## Technologies Used

- Go
- JSON
- AES-GCM encryption
- File-based storage


## Learning Goals

This project is intended to improve understanding of:
- Go project structure
- secure coding practices
- data serialization
- CLI design
- cryptography basics
- modular software architecture

## Features
- [ ] Master password authentication
- [ ] AES-256-GCM encryption
- [ ] Password hashing with scrypt or argon2
- [ ] Secure memory wiping
- [ ] Auto-lock vault after inactivity
- [ ] Encrypted backups
- [ ] Vault integrity verification

## License

MIT License
