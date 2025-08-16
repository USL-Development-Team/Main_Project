# Branch Protection via GitHub CLI

## Option 1: Make Repository Public (Free)

If you're comfortable making the repository public:

```bash
# First make repo public
gh repo edit USL-Development-Team/Main_Project --visibility public

# Then add branch protection for main
gh api repos/USL-Development-Team/Main_Project/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"checks":[{"context":"test"},{"context":"validate-commits"}]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true,"require_code_owner_reviews":false}' \
  --field restrictions='{"users":[],"teams":[],"apps":[]}' \
  --field allow_force_pushes=false \
  --field allow_deletions=false \
  --field required_linear_history=true \
  --field required_conversation_resolution=true

# Add branch protection for develop  
gh api repos/USL-Development-Team/Main_Project/branches/develop/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"checks":[{"context":"test"}]}' \
  --field enforce_admins=false \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field restrictions='{"users":[],"teams":[],"apps":[]}' \
  --field allow_force_pushes=false \
  --field allow_deletions=false \
  --field required_linear_history=true
```

## Option 2: Workflow-Based Protection (Current Setup)

Our `.github/workflows/branch-policy.yml` provides similar protection:
- Blocks PRs to main from non-develop branches
- Shows clear error messages
- Works with free private repos

## Option 3: Manual Setup (Web UI)

1. Go to: https://github.com/USL-Development-Team/Main_Project/settings/branches
2. Click "Add rule"
3. Configure protection settings manually

## Option 4: Upgrade to GitHub Pro

```bash
# This would enable advanced features
# Costs $4/month per user for private repos
```

## Current Status

✅ **Workflow enforcement** is active (works on free tier)
❌ **Native branch protection** requires public repo or GitHub Pro

## Recommendation

For now, our workflow-based protection is sufficient and gives you:
- PR blocking from wrong branches
- Status check requirements  
- Clear error messages
- Auto-labeling of release PRs

The main difference is that admins can still force-push to protected branches with native protection, but workflow protection blocks PRs at the GitHub Actions level.