# USL Application Deployment Guide

## Render Deployment

This application is configured for deployment on Render using Docker containers.

### Files Created for Deployment

1. **Dockerfile** - Multi-stage build with Tailwind CSS compilation and Go binary
2. **render.yaml** - Infrastructure as Code configuration
3. **.dockerignore** - Excludes unnecessary files from Docker build context
4. **DEPLOYMENT.md** - This deployment guide

### Pre-deployment Setup

#### 1. Update render.yaml Repository URL
Edit `render.yaml` and update the repository URL:
```yaml
repo: https://github.com/YOUR_USERNAME/YOUR_REPO_NAME
```

#### 2. Set Environment Variables in Render Dashboard
The following environment variables need to be set in the Render dashboard (marked as `sync: false` in render.yaml):

**Required Secret Variables:**
- `SUPABASE_URL` - Your production Supabase project URL
- `SUPABASE_PUBLIC_URL` - Your production Supabase public URL (usually same as SUPABASE_URL)
- `SUPABASE_ANON_KEY` - Supabase anonymous key
- `SUPABASE_SERVICE_ROLE_KEY` - Supabase service role key
- `USL_ADMIN_DISCORD_IDS` - Comma-separated Discord IDs for admin access

**OAuth Configuration:**
After deployment, you'll need to update your Discord OAuth app and Supabase Auth settings with the production URLs.

### Deployment Steps

#### Option 1: Using Render Dashboard
1. Connect your GitHub repository to Render
2. Create a new Web Service
3. Select "Docker" as the runtime
4. Render will automatically detect the Dockerfile
5. Set the environment variables listed above
6. Deploy

#### Option 2: Using render.yaml Blueprint
1. Push the render.yaml file to your repository
2. In Render dashboard, create a new Blueprint
3. Connect to your repository
4. Set the required environment variables
5. Deploy

### Post-deployment Configuration

#### 1. Update Discord OAuth App
Add your Render production URL to Discord OAuth redirect URIs:
- `https://your-app-name.onrender.com/auth/callback`

#### 2. Update Supabase Auth Settings
In Supabase dashboard > Authentication > URL Configuration:
- Add your production domain to "Site URL"
- Add redirect URLs for OAuth callbacks

#### 3. Test OAuth Flow
1. Visit your production URL
2. Test Discord authentication
3. Verify admin access works correctly

### Environment Configuration

The application uses these environment variables (with defaults):

```bash
# Server (configured automatically by Render)
PORT=8080                    # Render provides this
SERVER_HOST=0.0.0.0         # Set in render.yaml

# Supabase (set in Render dashboard)
SUPABASE_URL=               # Required
SUPABASE_PUBLIC_URL=        # Required
SUPABASE_ANON_KEY=          # Required
SUPABASE_SERVICE_ROLE_KEY=  # Required

# TrueSkill (defaults in render.yaml)
TRUESKILL_INITIAL_MU=1000.0
TRUESKILL_INITIAL_SIGMA=8.333333
TRUESKILL_SIGMA_MIN=0.25
TRUESKILL_SIGMA_MAX=8.333333
TRUESKILL_GAMES_FOR_MAX_CERTAINTY=50

# MMR (defaults in render.yaml)
MMR_ONES_WEIGHT=0.3
MMR_TWOS_WEIGHT=0.5
MMR_THREES_WEIGHT=0.2

# USL (set in Render dashboard)
USL_ADMIN_DISCORD_IDS=      # Required - comma-separated Discord IDs
```

### Build Process

The Docker build process:

1. **Stage 1 (Tailwind)**: Builds CSS from source using Tailwind CSS
2. **Stage 2 (Go Build)**: Compiles Go application with optimized flags
3. **Stage 3 (Runtime)**: Creates minimal Alpine Linux container with compiled binary

### Scaling and Performance

- **Auto-scaling**: Configured for 1-2 instances based on CPU usage
- **Health checks**: Configured to check application health
- **Zero-downtime deploys**: Supported by Render
- **Static assets**: Served efficiently from compiled CSS and templates

### Monitoring and Logs

- Access logs through Render dashboard
- Application logs written to stdout (captured by Render)
- Health check endpoint: `/` (redirects to login if not authenticated)

### Troubleshooting

#### Common Issues:
1. **OAuth redirect errors**: Check Discord app and Supabase Auth settings
2. **Database connection issues**: Verify Supabase credentials and URL
3. **Permission errors**: Check USL_ADMIN_DISCORD_IDS includes your Discord ID
4. **Build failures**: Check Dockerfile and ensure all dependencies are available

#### Debug Steps:
1. Check Render logs for detailed error messages
2. Verify all environment variables are set correctly
3. Test Supabase connection from Render environment
4. Validate Discord OAuth configuration

### Security Considerations

- All sensitive credentials stored as environment variables
- Non-root user in Docker container
- Static asset compilation during build (no runtime dependencies)
- HTTPS enforced by Render platform
- Supabase handles database security and auth