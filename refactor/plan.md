# USL Tracker Validation System Refactoring Plan
**Session ID**: validation-system-2025-08-16  
**Created**: 2025-08-16T20:30:00Z  
**Focus**: Implement comprehensive validation for USL tracker CRUD operations

## Initial State Analysis

### Current Validation State
- **Basic HTML5 validation**: Only `required` attribute on Discord ID field
- **Minimal server validation**: Single check for empty Discord ID in handler
- **No error feedback system**: Users get generic HTTP errors or silent failures
- **No business rule validation**: No MMR range checks, logical validation, or format validation
- **No client-side validation**: Beyond basic HTML5 attributes

### Problem Areas Identified
1. **🚨 CRITICAL**: No user feedback for validation errors
2. **🚨 CRITICAL**: Silent failure on invalid input (parseInt returns 0)
3. **⚠️ HIGH**: No business rule validation for MMR fields
4. **⚠️ HIGH**: No URL format validation for tracker URLs
5. **⚠️ HIGH**: No Discord ID format validation
6. **📈 MEDIUM**: No client-side validation for better UX

### Current Architecture
- Handler functions: Direct form parsing with minimal validation
- Templates: Basic HTML forms with no error display capability
- Error handling: Generic HTTP errors without field-specific feedback

## Refactoring Tasks

### Phase 1: Structured Validation Foundation (Sprint 1)

#### Task 1.1: Add Validation Types and Constants
**Priority**: 🚨 CRITICAL  
**Risk**: Low (adding new types)  
**Impact**: Foundation for all validation logic

**Actions**:
- [ ] Add `ValidationError` struct with Field, Message, Code
- [ ] Add `ValidationResult` struct with IsValid and Errors slice
- [ ] Add business rule constants (MinMMR, MaxMMR, MaxGames)
- [ ] Add validation error codes for consistent categorization

**Files**: 
- `internal/usl/handlers/migration_handler.go`

#### Task 1.2: Implement Core Validation Functions
**Priority**: 🚨 CRITICAL  
**Risk**: Low (pure functions)  
**Impact**: Enables comprehensive validation

**Actions**:
- [ ] Implement `validateTracker()` main validation function
- [ ] Add `validatePlaylistMMR()` for business rule validation
- [ ] Add `validateGamesPlayed()` for games count validation
- [ ] Add `isValidDiscordID()` for Discord ID format validation
- [ ] Add `isValidTrackerURL()` for URL format validation
- [ ] Add `hasNoPlaylistData()` for business rule validation

**Files**:
- `internal/usl/handlers/migration_handler.go`

#### Task 1.3: Add Error Display System
**Priority**: 🚨 CRITICAL  
**Risk**: Medium (template changes)  
**Impact**: Users can see validation errors

**Actions**:
- [ ] Add `renderFormWithErrors()` function
- [ ] Add `buildErrorMap()` helper for template lookup
- [ ] Update template data structures to include errors
- [ ] Update CreateTracker handler to use validation

**Files**:
- `internal/usl/handlers/migration_handler.go`

### Phase 2: Template Integration (Sprint 1 continuation)

#### Task 2.1: Update New Tracker Template
**Priority**: 🚨 CRITICAL  
**Risk**: Medium (UI changes)  
**Impact**: Error display for create tracker form

**Actions**:
- [ ] Add error display HTML for each form field
- [ ] Add conditional CSS classes for error states
- [ ] Add value preservation on validation failure
- [ ] Update template data structure expectations

**Files**:
- `templates/tracker-new.html`

#### Task 2.2: Update Edit Tracker Template
**Priority**: 🚨 CRITICAL  
**Risk**: Medium (UI changes)  
**Impact**: Error display for edit tracker form

**Actions**:
- [ ] Mirror validation error display from new template
- [ ] Ensure value preservation works with existing data
- [ ] Update UpdateTracker handler to use validation

**Files**:
- `templates/tracker-edit.html`
- `internal/usl/handlers/migration_handler.go`

### Phase 3: Enhanced Validation Rules (Sprint 2)

#### Task 3.1: Advanced Business Logic Validation
**Priority**: ⚠️ HIGH  
**Risk**: Low (expanding existing validation)  
**Impact**: Catches more user input errors

**Actions**:
- [ ] Add logical validation (all-time >= previous >= current)
- [ ] Add cross-field validation rules
- [ ] Add "at least one playlist" requirement
- [ ] Add reasonable upper bounds for all fields

**Files**:
- `internal/usl/handlers/migration_handler.go`

#### Task 3.2: Enhanced URL Validation
**Priority**: ⚠️ HIGH  
**Risk**: Low (URL parsing)  
**Impact**: Prevents invalid tracker URLs

**Actions**:
- [ ] Add whitelist of valid tracker domains
- [ ] Add URL format validation
- [ ] Add protocol validation (https required)
- [ ] Consider async URL verification (Phase 3)

**Files**:
- `internal/usl/handlers/migration_handler.go`

### Phase 4: Client-Side Validation (Sprint 2 continuation)

#### Task 4.1: JavaScript Validation Framework
**Priority**: 📈 MEDIUM  
**Risk**: Low (adding JS)  
**Impact**: Better user experience

**Actions**:
- [ ] Add `validateField()` JavaScript function
- [ ] Add real-time validation on field blur
- [ ] Add `showFieldErrors()` function for error display
- [ ] Integrate with existing MMR calculation

**Files**:
- `templates/tracker-new.html`
- `templates/tracker-edit.html`

#### Task 4.2: Enhanced Client-Side Rules
**Priority**: 📈 MEDIUM  
**Risk**: Low (JS validation)  
**Impact**: Immediate feedback to users

**Actions**:
- [ ] Add Discord ID format validation in JS
- [ ] Add MMR range validation in JS
- [ ] Add URL format validation in JS
- [ ] Add visual feedback for validation states

**Files**:
- `templates/tracker-new.html`
- `templates/tracker-edit.html`

### Phase 5: Advanced Features (Future)

#### Task 5.1: Async Validation
**Priority**: 🔧 LOW  
**Actions**:
- [ ] Add Discord user existence validation
- [ ] Add tracker URL accessibility verification
- [ ] Add duplicate tracker detection

#### Task 5.2: Validation Caching
**Priority**: 🔧 LOW  
**Actions**:
- [ ] Cache validation results for performance
- [ ] Add validation result expiration
- [ ] Add cache invalidation strategies

## Implementation Strategy

### Incremental Rollout
1. **Foundation First**: Add validation types and core functions
2. **Server-Side Integration**: Update handlers to use validation
3. **Template Updates**: Add error display capability
4. **Client-Side Enhancement**: Add JavaScript validation
5. **Advanced Features**: Add async validation and caching

### Risk Mitigation
- **Backward Compatibility**: All changes are additive, no breaking changes
- **Progressive Enhancement**: Forms work without JS, better with JS
- **Graceful Degradation**: Server-side validation always present
- **Testing Strategy**: Validate each field type and error case

## Validation Checklist

### After Each Phase
- [ ] All form submissions properly validated
- [ ] Error messages are user-friendly and specific
- [ ] Values preserved on validation failure
- [ ] No breaking changes to existing functionality
- [ ] All validation rules properly tested
- [ ] Error display works in all supported browsers

### Final Validation
- [ ] Complete validation coverage for all fields
- [ ] Consistent error handling across all forms
- [ ] Performance impact is negligible
- [ ] User experience significantly improved
- [ ] Security improved through input validation
- [ ] Maintainable validation code structure

## Business Rules Implemented

### MMR Validation
- **Range**: 0-3000 MMR (SSL ~1900, allowing buffer)
- **Logic**: All-time >= Previous >= Current (when values exist)
- **Required**: At least one playlist must have data

### Discord ID Validation  
- **Format**: 17-19 digit snowflake
- **Required**: Always required for tracker creation

### URL Validation
- **Format**: Valid URL format
- **Whitelist**: Known tracker sites (tracker.network, ballchasing.com, etc.)
- **Protocol**: HTTPS preferred

### Games Played Validation
- **Range**: 0-10000 reasonable upper limit
- **Logic**: Current season games should be reasonable

## Success Metrics

### Target Outcomes
1. **User Experience**: 90% reduction in form submission errors
2. **Data Quality**: 95% of submitted trackers have valid data
3. **Support Burden**: 75% reduction in invalid data support tickets
4. **Development Velocity**: Easier to add new validation rules

### Performance Targets
- **Validation Time**: <1ms for full tracker validation
- **Page Load Impact**: <100ms additional for validation JS
- **Memory Usage**: <1MB additional for validation state

## Rollback Strategy

### Safety Measures
- Git checkpoints before each major change
- Feature flags for validation enforcement
- Ability to disable client-side validation
- Graceful degradation to current behavior

### Emergency Rollback
1. Disable validation in handlers (return to current behavior)
2. Remove error display from templates
3. Revert to basic HTML5 validation only
4. Notify users of temporary simplified validation

**Estimated Total Effort**: 1-2 sprints for core implementation, 1 additional sprint for enhancements
**Next Session**: Begin Phase 1 - Structured Validation Foundation