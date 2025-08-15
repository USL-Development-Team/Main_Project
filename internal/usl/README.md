# USL Migration - Temporary Code

⚠️ **THIS DIRECTORY CONTAINS TEMPORARY CODE** ⚠️

## Purpose
This directory contains USL-specific migration code to get USL off Google Sheets ASAP. 

**This code is intentionally temporary and will be deleted once USL is migrated to the multi-guild platform.**

## Contents

### `handlers/migration_handler.go`
- Simplified USL-only web interface 
- No multi-guild logic or complexity
- Hardcoded USL Discord Guild ID: `1390537743385231451`
- Server-side rendering with basic forms

### `templates/`
- USL-specific HTML templates
- Simple, functional interface for USL admins
- No fancy styling - just get the job done

### `scripts/`
- `export-usl-data.gs` - Google Apps Script to export USL data from sheets
- `import.go` - Go script to import USL data into new database schema

## Routes Created
- `/usl/admin` - Admin dashboard  
- `/usl/users` - User management
- `/usl/trackers` - Tracker management
- `/usl/import` - Data import tools

## Migration Timeline
1. **Phase 1:** Use this code to migrate USL off Google Sheets (2-3 weeks)
2. **Phase 2:** USL uses this temporarily while multi-guild platform is built (2-3 months)  
3. **Phase 3:** Migrate USL to multi-guild platform and **DELETE THIS DIRECTORY**

## Important Notes
- **No tests needed** - this is temporary migration code
- **Keep it simple** - don't over-engineer 
- **Mark clearly** - anyone reading this code should know it's temporary
- **Easy to remove** - designed for deletion, not maintenance

## When to Delete This Directory
Delete this entire directory when:
- [ ] Multi-guild platform is ready for production
- [ ] USL has been migrated to multi-guild platform  
- [ ] USL admins are comfortable with new interface
- [ ] All USL data has been verified in new system

---

**Remember: This is a bridge, not a destination.**