# Go Game Server by AI

A game server developed with Go and AI technologies.

## Project Overview

This project is a game server built with Go programming language, integrated with AI capabilities for enhanced gaming experiences.

## Project Structure

```
.
├── api/                 # API definitions
├── cmd/                 # Main applications
├── docs/                # Documentation
├── internal/            # Private application and library code
├── pkg/                 # External facing library code
├── scripts/             # Scripts for build, install, analysis, etc.
├── .github/workflows/   # GitHub Actions workflows
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # Project documentation
└── LICENSE             # License information
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:xieyx/go-game-server-by-ai.git
   ```

2. Navigate to the project directory:
   ```bash
   cd go-game-server-by-ai
   ```

3. Build the project:
   ```bash
   go build ./cmd/...
   ```

### Running Tests

To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test -cover ./...
```

To run tests with verbose output:
```bash
go test -v ./...
```

### Running the Server

To run the game server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default. You can change the port by setting the PORT environment variable:
```bash
PORT=3000 go run cmd/server/main.go
```

### Running the Turn-Based Combat Demo

To run the turn-based combat demo:
```bash
go run cmd/battle_demo/main.go
```

This demo showcases the turn-based combat system with:
- Character creation (Warrior, Mage, Archer)
- Turn-based battle mechanics
- Skill usage
- Status effects
- Battle rewards

## Development

This project follows the GitFlow workflow with the following branch naming conventions:
- `main` - Production releases
- `develop` - Development branch
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Production hotfixes
- `experimental/*` - Experimental features

## Contributing

1. Create a feature branch from `develop`
2. Make your changes
3. Write tests for your changes
4. Update documentation as needed
5. Submit a pull request for review

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
