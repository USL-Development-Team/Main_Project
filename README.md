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

### Required Environment Variables

**Core Application:**
- `APP_BASE_URL` - **Required for staging/production**. Your application's base URL (e.g., `https://rl-league-management.onrender.com`)
- `ENVIRONMENT` - Deployment environment (`development`, `staging`, `production`)
- `SERVER_HOST` - Server bind address (default: `0.0.0.0` for production)
- `SERVER_PORT` or `PORT` - Server port (default: `8080`)

**Database (Supabase):**
- `SUPABASE_URL` - **Required**. Your Supabase project URL
- `SUPABASE_ANON_KEY` - **Required**. Supabase anonymous/public key
- `SUPABASE_SERVICE_ROLE_KEY` - **Required**. Supabase service role key for server operations
- `DATABASE_URL` - PostgreSQL connection string (optional, uses Supabase by default)

**Discord OAuth:**
- `DISCORD_CLIENT_ID` - **Required**. Discord application client ID
- `DISCORD_CLIENT_SECRET` - **Required**. Discord application client secret

**Optional Configuration:**
- `SUPABASE_PUBLIC_URL` - Override for Supabase public URL (defaults to SUPABASE_URL)
- `USL_ADMIN_DISCORD_IDS` - Comma-separated Discord IDs for admin access
- TrueSkill configuration (`TRUESKILL_*`)
- MMR calculation weights (`MMR_*`)

**⚠️ Important**: Production and staging environments will fail to start if required variables are missing.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for development workflow, commit conventions, and release process.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
