# Contributing to USL Management System

## Development Workflow

### Branch Strategy
- **`main`** - Production releases (tagged with semantic versions)
- **`develop`** - Integration branch for features
- **`feature/*`** - Feature development branches

### Getting Started

1. **Clone and setup:**
   ```bash
   git clone https://github.com/USL-Development-Team/Main_Project.git
   cd Main_Project
   git checkout develop
   ```

2. **Create feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make changes and commit:**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

4. **Push and create PR:**
   ```bash
   git push -u origin feature/your-feature-name
   # Create PR to develop branch on GitHub
   ```

### Commit Message Convention

We use [Conventional Commits](https://www.conventionalcommits.org/) for automated versioning:

- **`feat:`** - New features (minor version bump)
- **`fix:`** - Bug fixes (patch version bump)
- **`refactor:`** - Code refactoring
- **`docs:`** - Documentation changes
- **`test:`** - Test additions/changes
- **`chore:`** - Maintenance tasks
- **`BREAKING CHANGE:`** - Breaking changes (major version bump)

### Release Process

1. **Feature complete on develop:**
   - All features tested and working
   - PR checks passing

2. **Create release PR:**
   - Merge `develop` into `main`
   - This triggers automated release process

3. **Automated release:**
   - Semantic version calculated from commits
   - CHANGELOG.md updated automatically
   - Git tag created
   - GitHub release published
   - Render deployment triggered

### Code Quality

- All PRs require review before merging
- Automated checks run on every PR:
  - Go formatting (`go fmt`)
  - Go vet checks
  - Build verification
  - Commit message validation

### Local Development

1. **Setup environment:**
   ```bash
   cp .env.develop .env
   supabase start
   go run cmd/server/main.go
   ```

2. **Run checks locally:**
   ```bash
   go fmt ./...
   go vet ./...
   go build ./cmd/server
   ```

### Questions?

Open an issue or contact the development team on Discord.