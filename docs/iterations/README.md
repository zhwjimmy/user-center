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

- [user-features-v2025-07-26-154029.md](./user-features-v2025-07-26-154029.md) - user-features feature list

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
