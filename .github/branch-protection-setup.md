# Branch Protection Setup Guide

## GitHub Repository Settings

To prevent merging to `main` from any branch other than `develop`, configure these settings:

### 1. Navigate to Repository Settings
1. Go to your GitHub repository
2. Click **Settings** tab
3. Click **Branches** in the left sidebar

### 2. Add Branch Protection Rule for `main`

Click **Add rule** and configure:

#### Branch name pattern:
```
main
```

#### Protection Settings:
- ✅ **Restrict pushes that create files larger than 100 MB**
- ✅ **Require a pull request before merging**
  - ✅ **Require approvals**: 1
  - ✅ **Dismiss stale PR approvals when new commits are pushed**
  - ✅ **Require review from code owners** (optional)
  - ✅ **Restrict reviews to users with write access**

- ✅ **Require status checks to pass before merging**
  - ✅ **Require branches to be up to date before merging**
  - **Required status checks:**
    - `test`
    - `validate-commits`

- ✅ **Require conversation resolution before merging**
- ✅ **Require signed commits** (optional but recommended)
- ✅ **Require linear history** (prevents merge commits)
- ✅ **Require deployments to succeed before merging** (optional)

#### Advanced Settings:
- ✅ **Restrict pushes that create files larger than 100 MB**
- ✅ **Lock branch** (prevents any direct pushes)
- ✅ **Do not allow bypassing the above settings**
- ✅ **Allow force pushes** → **Specify who can force push** → Nobody
- ✅ **Allow deletions** → Disabled

### 3. Add Branch Protection Rule for `develop`

Click **Add rule** and configure:

#### Branch name pattern:
```
develop
```

#### Protection Settings:
- ✅ **Require a pull request before merging**
  - ✅ **Require approvals**: 1
  - ✅ **Dismiss stale PR approvals when new commits are pushed**

- ✅ **Require status checks to pass before merging**
  - **Required status checks:**
    - `test`

- ✅ **Require conversation resolution before merging**
- ✅ **Require linear history**

### 4. Repository-Level Settings

In **Settings > General**:

#### Pull Requests:
- ✅ **Allow merge commits** → Disabled
- ✅ **Allow squash merging** → Enabled (recommended)
- ✅ **Allow rebase merging** → Enabled
- ✅ **Always suggest updating pull request branches**
- ✅ **Allow auto-merge**
- ✅ **Automatically delete head branches**

## Result

With these settings:

1. **Direct pushes to `main`** → ❌ Blocked
2. **PRs to `main` from feature branches** → ❌ Blocked  
3. **PRs to `main` from `develop`** → ✅ Allowed (with checks)
4. **All PRs require approval** → ✅ Required
5. **Status checks must pass** → ✅ Required

## Workflow After Setup

```bash
# ❌ This will be blocked
git checkout main
git push origin feature/something

# ✅ This is the only allowed path to main
feature/branch → develop → main
```

## Emergency Override

Repository admins can still override protections if needed, but it will be logged and visible in the audit trail.

## Testing the Protection

After setup, try:
```bash
git checkout main
echo "test" > test.txt
git add test.txt
git commit -m "test: direct push to main"
git push origin main
```

You should see: `remote: error: GH006: Protected branch update failed`