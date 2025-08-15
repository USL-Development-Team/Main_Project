# USL Templates & Routes Refactoring Plan - 2025-08-15 (Fresh Start)

## Problem Analysis

**CRITICAL ISSUE**: Template system is fundamentally broken due to misunderstanding of Go's template patterns.

### Current Broken State:
1. **Template Name Mismatch**: `usl_admin_dashboard.html` defines `{{define "usl_admin_dashboard"}}` but handler calls `ExecuteTemplate(w, "usl_admin_dashboard.html", data)`
2. **Conflicting Inheritance**: Templates use `{{template "usl_base" .}}` but also have conflicting `{{define "navigation"}}` and `{{define "content"}}` blocks across multiple files
3. **Template Pollution**: All templates are loaded together, causing name conflicts
4. **Inconsistent Patterns**: Mix of standalone templates and inheritance attempts

### Root Cause:
Go templates work differently than expected:
- `ExecuteTemplate(w, "templateName", data)` looks for a `{{define "templateName"}}` block
- When using `{{template "base" .}}`, the template name should match the filename (minus .html)
- All templates are parsed together, so `{{define}}` names must be unique across ALL files

## Go-Appropriate Solutions

### Option A: Standalone Templates (Recommended for USL)
Each template is self-contained with shared CSS via a separate file or embedded styles.
- **Pros**: Simple, no conflicts, easy to debug, appropriate for temporary migration code
- **Cons**: Some CSS duplication (acceptable for temporary code)

### Option B: Proper Template Inheritance
Use unique template names and proper inheritance patterns.
- **Pros**: DRY principle, shared layouts
- **Cons**: More complex, overkill for temporary migration code

## Recommended Approach: Option A (Standalone Templates)

Since USL is explicitly temporary migration code, we should prioritize:
1. **Simplicity**: Easy to understand and debug
2. **Reliability**: No complex inheritance that can break
3. **Speed**: Fast to implement and maintain
4. **Clarity**: Each template is independent

## Implementation Plan

### Phase 1: Fix Template System (HIGH PRIORITY)
1. **Convert to standalone templates**: Each USL template becomes self-contained
2. **Shared CSS approach**: Extract common CSS to a shared snippet or embed in each template
3. **Fix template names**: Ensure template names match what handlers expect
4. **Test each template**: Verify each page renders correctly

### Phase 2: Validate Routes & Handlers
1. **Check all USL routes**: Ensure they're properly registered and working
2. **Test form submissions**: Verify all CRUD operations work
3. **Check authentication**: Ensure Discord OAuth protection is working
4. **Validate data flow**: Test user/tracker creation, editing, etc.

### Phase 3: Handler Organization (OPTIONAL)
1. **Keep single handler**: Since it's temporary code, don't over-engineer
2. **Add error handling**: Improve error messages and validation
3. **Add logging**: Better debugging for issues

## Files to Fix

### Templates (7 files):
- `usl_admin_dashboard.html` - Broken template name
- `usl_users.html` - Template inheritance conflicts
- `usl_user_form.html` - Template inheritance conflicts  
- `usl_user_edit_form.html` - Template inheritance conflicts
- `usl_trackers.html` - Template inheritance conflicts
- `usl_tracker_form.html` - Template inheritance conflicts
- `usl_import.html` - Template inheritance conflicts

### Routes to Test:
- `/usl/admin` - Dashboard
- `/usl/users` - User list with search
- `/usl/users/new` - User creation form
- `/usl/users/edit` - User edit form
- `/usl/users/create` - User creation handler
- `/usl/users/update` - User update handler
- `/usl/trackers` - Tracker list
- `/usl/trackers/new` - Tracker creation form
- `/usl/trackers/create` - Tracker creation handler
- `/usl/import` - Data import interface

## Success Criteria

### Template Validation:
- [ ] All USL pages render without errors
- [ ] Navigation works correctly
- [ ] Forms submit successfully
- [ ] Search functionality works
- [ ] Styling is consistent
- [ ] Mobile responsive (basic)

### Route Validation:
- [ ] All routes respond correctly
- [ ] Authentication required for all routes
- [ ] CRUD operations work (Create, Read, Update)
- [ ] Error handling works
- [ ] Redirects work properly

### Data Flow Validation:
- [ ] User creation works
- [ ] User editing works (name and status only)
- [ ] Tracker creation works
- [ ] Search functionality works
- [ ] Dashboard shows correct stats

## Implementation Strategy

1. **Start with admin dashboard**: Fix the most broken template first
2. **Fix one template at a time**: Test each individually
3. **Use simple standalone pattern**: No complex inheritance
4. **Test immediately**: Check each template in browser after fixing
5. **Keep backup files**: Don't delete working backups until confirmed

## File Changes Summary

**Remove files**:
- `usl_base.html` (complex inheritance causing conflicts)

**Convert to standalone**:
- All 7 USL templates become self-contained
- Each includes its own CSS (or shared CSS file)
- Each has proper `{{define "templateName"}}` matching filename

**No handler changes needed**:
- Handlers are calling templates correctly
- Issue is in template definitions, not handler calls

---

**PHILOSOPHY**: Keep it simple for temporary migration code. Optimize for reliability and ease of debugging, not for perfect code architecture.