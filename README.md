# USL Management System

A comprehensive Discord-based management system for competitive gaming leagues, featuring TrueSkill rating calculations, user tracking, and automated workflows.

## Features

### Core Functionality
- **Discord OAuth Authentication** - Secure login via Supabase Auth
- **User Management** - Track players and their gaming profiles
- **TrueSkill Rating System** - Advanced skill rating calculations
- **Tracker Integration** - Link external tracking platforms
- **Web Interface** - Clean, responsive UI for management tasks

### Development & Deployment
- **Automated Releases** - Semantic versioning with conventional commits
- **Environment Management** - Development, staging, and production configs
- **Docker Support** - Containerized deployment ready for Render
- **Branch Protection** - Enforced git flow with develop → main workflow
- **Professional Templates** - GitHub issue and PR templates

## Quick Start

### Local Development

1. **Clone and setup:**
   ```bash
   git clone https://github.com/USL-Development-Team/Main_Project.git
   cd Main_Project
   cp .env.develop .env
   ```

2. **Start services:**
   ```bash
   supabase start
   go run cmd/server/main.go
   ```

3. **Access the application:**
   - Local: `http://localhost:8080`
   - Health check: `http://localhost:8080/health`

### Production Deployment

The application automatically deploys to Render when releases are published from the main branch.

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL via Supabase
- **Authentication**: Discord OAuth through Supabase Auth
- **Frontend**: HTML templates with Tailwind CSS
- **Deployment**: Docker containers on Render
- **CI/CD**: GitHub Actions with automated versioning

## Environment Configuration

Environment-specific configurations support:
- Development (`localhost:8080`)
- Staging (`staging-usl.render.com`)
- Production (`usl.render.com`)

Each environment has its own OAuth redirect configuration and database settings.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow, commit conventions, and release process.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
