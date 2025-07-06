#!/bin/bash

# GitHub Actions Setup Script for UserCenter
# This script helps you set up GitHub Actions for your repository

set -e

echo "🚀 Setting up GitHub Actions for UserCenter..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Check if we have the required files
if [ ! -f ".github/workflows/ci.yml" ]; then
    echo "❌ Error: GitHub Actions workflows not found"
    echo "Please ensure you have the .github/workflows/ directory with the required workflow files"
    exit 1
fi

echo "✅ GitHub Actions workflows found"

# Get repository information
REPO_URL=$(git remote get-url origin)
REPO_NAME=$(basename -s .git "$REPO_URL")

echo "📋 Repository: $REPO_NAME"
echo "🔗 URL: $REPO_URL"

# Check if this is a GitHub repository
if [[ $REPO_URL != *"github.com"* ]]; then
    echo "⚠️  Warning: This doesn't appear to be a GitHub repository"
    echo "GitHub Actions will only work with GitHub repositories"
fi

echo ""
echo "📝 Next steps to complete GitHub Actions setup:"
echo ""
echo "1. Push your code to GitHub:"
echo "   git add ."
echo "   git commit -m 'Add GitHub Actions workflows'"
echo "   git push origin main"
echo ""
echo "2. Configure repository secrets (Settings > Secrets and variables > Actions):"
echo "   - POSTGRES_HOST: Your PostgreSQL host"
echo "   - POSTGRES_PORT: Your PostgreSQL port (default: 5432)"
echo "   - POSTGRES_USER: Your PostgreSQL username"
echo "   - POSTGRES_PASSWORD: Your PostgreSQL password"
echo "   - POSTGRES_DB: Your PostgreSQL database name"
echo "   - REDIS_HOST: Your Redis host"
echo "   - REDIS_PORT: Your Redis port (default: 6379)"
echo "   - REDIS_PASSWORD: Your Redis password (if required)"
echo ""
echo "3. Set up environments (Settings > Environments):"
echo "   - staging: For staging deployments"
echo "   - production: For production deployments"
echo ""
echo "4. Enable GitHub Actions (Settings > Actions > General):"
echo "   - Allow all actions and reusable workflows"
echo ""
echo "5. Configure Dependabot (optional):"
echo "   - The .github/dependabot.yml file is already configured"
echo "   - Dependabot will automatically create PRs for dependency updates"
echo ""
echo "6. Set up Codecov integration (optional):"
echo "   - Visit https://codecov.io"
echo "   - Connect your GitHub repository"
echo "   - Add CODECOV_TOKEN secret if needed"
echo ""
echo "🎉 GitHub Actions setup complete!"
echo ""
echo "📚 For more information, see: docs/github-actions.md"
echo ""
echo "🔍 To monitor your workflows:"
echo "   https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/actions" 