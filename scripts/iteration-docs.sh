#!/bin/bash

# Iteration Documentation Management Script
# This script provides utilities for managing iteration documents with versioning

set -e

# Configuration
ITERATION_DOCS_DIR="docs/iterations"
TIMESTAMP=$(date +"%Y-%m-%d-%H%M%S")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Create a new iteration document version
create_iteration_doc() {
    local name="$1"
    
    if [ -z "$name" ]; then
        log_error "Feature name is required"
        echo "Usage: $0 create <feature_name>"
        exit 1
    fi
    
    local filename="${name}-v${TIMESTAMP}.md"
    local filepath="${ITERATION_DOCS_DIR}/${filename}"
    local latest_link="${ITERATION_DOCS_DIR}/${name}-latest.md"
    
    log_info "Creating new iteration document: $filename"
    
    # Create directory if it doesn't exist
    mkdir -p "$ITERATION_DOCS_DIR"
    
    # Check if latest version exists
    if [ -f "$latest_link" ]; then
        # Copy from latest version
        cp "$latest_link" "$filepath"
        log_success "Created $filename from latest version"
    else
        # Create new document with template
        cat > "$filepath" << 'EOF'
# $name Feature List

## Version Information
- Created: $(date '+%Y-%m-%d %H:%M:%S')
- Version: v$TIMESTAMP

## Feature List

1. [Feature to be added]

## Implementation Details

### API Endpoints
- \`GET /api/v1/...\` - Description
- \`POST /api/v1/...\` - Description

### Database Tables
- \`table_name\` - Description

### Services
- \`ServiceName\` - Description

### Events
- \`EventType\` - Description

## Notes
> Add implementation details and notes here as the iteration progresses.
EOF
        log_success "Created new $filename with template"
    fi
    
    # Update latest symlink
    ln -sf "$filename" "$latest_link"
    log_success "Latest symlink updated to $filename"
}

# List all iteration documents
list_iteration_docs() {
    log_info "Iteration documents:"
    
    if [ ! -d "$ITERATION_DOCS_DIR" ]; then
        log_warning "Iteration docs directory not found: $ITERATION_DOCS_DIR"
        return
    fi
    
    if ls "$ITERATION_DOCS_DIR"/*.md >/dev/null 2>&1; then
        ls -la "$ITERATION_DOCS_DIR"/*.md | grep -v README.md || log_warning "No iteration documents found"
    else
        log_warning "No iteration documents found"
    fi
}

# Show the latest iteration document
show_latest_doc() {
    local name="$1"
    
    if [ -z "$name" ]; then
        log_error "Feature name is required"
        echo "Usage: $0 show-latest <feature_name>"
        exit 1
    fi
    
    local latest_link="${ITERATION_DOCS_DIR}/${name}-latest.md"
    
    if [ -f "$latest_link" ]; then
        log_info "Latest $name document:"
        cat "$latest_link"
    else
        log_warning "No latest document found for $name"
        log_info "Available documents:"
        ls "$ITERATION_DOCS_DIR"/*"$name"*.md 2>/dev/null || log_warning "No documents found"
    fi
}

# Clean old iteration documents (keep last 5 versions)
clean_old_docs() {
    log_info "Cleaning old iteration documents (keeping last 5 versions)..."
    
    if [ ! -d "$ITERATION_DOCS_DIR" ]; then
        log_warning "Iteration docs directory not found: $ITERATION_DOCS_DIR"
        return
    fi
    
    for doc in "$ITERATION_DOCS_DIR"/*-v*.md; do
        if [ -f "$doc" ]; then
            local base=$(basename "$doc" .md | sed 's/-v[0-9-]*$//')
            log_info "Processing $base..."
            
            # Keep only the last 5 versions
            ls -t "$ITERATION_DOCS_DIR/$base"-v*.md 2>/dev/null | tail -n +6 | xargs -r rm -f
        fi
    done
    
    log_success "Cleanup completed"
}

# Update iteration docs README
update_readme() {
    log_info "Updating iteration docs README..."
    
    cat > "$ITERATION_DOCS_DIR/README.md" << 'EOF'
# Iteration Documentation Management

## File Naming Convention

Iteration documents use the following naming format:
```
{feature_name}-v{version}.md
```

### Version Format
- Format: `YYYY-MM-DD-HHMMSS`
- Example: `2025-07-26-154029`
- Description: Year-Month-Day-HourMinuteSecond for easy sorting and identification

### Examples
- `user-features-v2025-07-26-154029.md` - User features list (created 2025-07-26 15:40:29)
- `auth-system-v2025-07-26-160000.md` - Authentication system features (created 2025-07-26 16:00:00)

## Version Management Guidelines

1. **Create new versions** for each major feature iteration
2. **Keep historical versions** for tracking feature evolution
3. **Document major changes** at the beginning of each document
4. **Use semantic versioning** (e.g., v1.0.0) as internal version identifiers

## Current Iteration Documents

EOF

    # Add current documents to README
    if [ -d "$ITERATION_DOCS_DIR" ]; then
        for doc in "$ITERATION_DOCS_DIR"/*-v*.md; do
            if [ -f "$doc" ]; then
                local filename=$(basename "$doc")
                local name=$(basename "$doc" .md | sed 's/-v[0-9-]*$//')
                echo "- [$filename](./$filename) - $name feature list" >> "$ITERATION_DOCS_DIR/README.md"
            fi
        done
    fi
    
    cat >> "$ITERATION_DOCS_DIR/README.md" << 'EOF'

## Quick Commands

```bash
# Create new iteration document
make docs-create-iteration name=feature_name

# List all iteration documents
make docs-list-iterations

# Show latest version
make docs-show-latest name=feature_name

# Clean old versions (keep last 5)
make docs-clean-old

# Update README
make docs-update-readme
```

## Script Usage

```bash
# Direct script usage
./scripts/iteration-docs.sh create <feature_name>
./scripts/iteration-docs.sh list
./scripts/iteration-docs.sh show-latest <feature_name>
./scripts/iteration-docs.sh clean-old
./scripts/iteration-docs.sh update-readme
```
EOF

    log_success "README updated successfully"
}

# Show help
show_help() {
    cat << EOF
Iteration Documentation Management Script

Usage: $0 <command> [options]

Commands:
  create <feature_name>    Create a new iteration document version
  list                     List all iteration documents
  show-latest <name>       Show the latest iteration document
  clean-old               Clean old iteration documents (keep last 5 versions)
  update-readme           Update iteration docs README with current documents
  help                    Show this help message

Examples:
  $0 create user-features
  $0 list
  $0 show-latest user-features
  $0 clean-old
  $0 update-readme

EOF
}

# Main script logic
main() {
    local command="$1"
    
    case "$command" in
        "create")
            create_iteration_doc "$2"
            ;;
        "list")
            list_iteration_docs
            ;;
        "show-latest")
            show_latest_doc "$2"
            ;;
        "clean-old")
            clean_old_docs
            ;;
        "update-readme")
            update_readme
            ;;
        "help"|"--help"|"-h"|"")
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@" 