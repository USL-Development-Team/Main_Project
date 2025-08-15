# Refactor Plan - TrueSkill Auto-Update Integration

## Initial State Analysis

**Current Problem:**
- TrackerHandler creates/updates tracker data but doesn't trigger TrueSkill recalculation
- TrueSkillService exists but requires manual HTTP calls to update skills
- Users update trackers but see stale skill ratings until admin runs batch update
- Google Sheets likely has immediate skill updates, we don't

**Current Architecture:**
```go
// TrackerHandler - missing TrueSkill integration
type TrackerHandler struct {
    trackerRepo *repositories.TrackerRepository
    templates   *template.Template
    // MISSING: trueSkillService *services.UserTrueSkillService
}

// Current flow (broken):
tracker := CreateTracker(data)
// NO TrueSkill update here
redirect("/trackers")
```

**Dependencies Available:**
- `TrueSkillService` exists in AppDependencies
- `UpdateUserTrueSkillFromTrackers()` method available
- Fast synchronous operation (10-50ms estimated)

## Refactoring Goals

**Target Architecture:**
```go
// TrackerHandler - with TrueSkill integration
type TrackerHandler struct {
    trackerRepo      *repositories.TrackerRepository
    trueSkillService *services.UserTrueSkillService  // ADD THIS
    templates        *template.Template
}

// Target flow (fixed):
tracker := CreateTracker(data)
result := h.trueSkillService.UpdateUserTrueSkillFromTrackers(tracker.DiscordID)
if !result.Success {
    log.Printf("TrueSkill update failed: %s", result.Error)
}
redirect("/trackers")
```

**User Experience Goal:**
User updates tracker → TrueSkill recalculates immediately → Google Sheets parity

## Refactoring Tasks

### Phase 1: Dependency Injection (15 minutes)
- [x] **Task 1.1**: Update TrackerHandler struct to include TrueSkillService
- [x] **Task 1.2**: Update NewTrackerHandler constructor signature  
- [x] **Task 1.3**: Update TrackerHandler instantiation in main.go (2 locations)

### Phase 2: Auto-Update Integration (30 minutes)  
- [x] **Task 2.1**: Add TrueSkill update after CreateTracker (line ~128)
- [x] **Task 2.2**: Add TrueSkill update after UpdateTracker (line ~248)
- [x] **Task 2.3**: Add TrueSkill update after DeleteTracker (line ~277)

### Phase 3: Error Handling & Logging (15 minutes)
- [x] **Task 3.1**: Add error handling for TrueSkill update failures
- [x] **Task 3.2**: Add success logging with skill change information
- [x] **Task 3.3**: Ensure tracker operations don't fail if TrueSkill update fails

## Implementation Details

### Code Changes Required:

**File 1: internal/handlers/tracker_handler.go**
```go
// BEFORE:
type TrackerHandler struct {
    trackerRepo *repositories.TrackerRepository
    templates   *template.Template
}

func NewTrackerHandler(trackerRepo *repositories.TrackerRepository, templates *template.Template) *TrackerHandler

// AFTER:
type TrackerHandler struct {
    trackerRepo      *repositories.TrackerRepository
    trueSkillService *services.UserTrueSkillService
    templates        *template.Template
}

func NewTrackerHandler(trackerRepo *repositories.TrackerRepository, 
                      trueSkillService *services.UserTrueSkillService,
                      templates *template.Template) *TrackerHandler
```

**File 2: cmd/server/main.go** (2 locations)
```go
// BEFORE:
trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.Templates)

// AFTER:
trackerHandler := handlers.NewTrackerHandler(deps.TrackerRepo, deps.TrueSkillService, deps.Templates)
```

**File 3: Add TrueSkill updates in tracker operations**
```go
// Add after successful tracker create/update/delete:
result := h.trueSkillService.UpdateUserTrueSkillFromTrackers(tracker.DiscordID)
if !result.Success {
    log.Printf("TrueSkill auto-update failed for %s: %s", tracker.DiscordID, result.Error)
    // Continue anyway - don't fail the tracker operation
} else {
    log.Printf("Auto-updated TrueSkill for %s: μ=%.1f", 
        tracker.DiscordID, result.TrueSkillResult.Mu)
}
```

## Risk Assessment

**Low Risk:**
- Dependency injection (constructor signature change)
- Adding optional TrueSkill calls (don't break on failure)
- Synchronous operation is fast enough (<50ms)

**Very Low Risk:**
- TrueSkill calculations are isolated and well-tested
- No changes to existing TrueSkill logic
- No database schema changes

**Mitigation:**
- TrueSkill update failures don't break tracker operations
- Detailed error logging for troubleshooting
- Can be easily rolled back if issues arise

## Performance Considerations

**Impact Assessment:**
- TrueSkill calculation: ~10-50ms (percentile math + single DB query)
- User frequency: ~1 tracker update per user per day
- Total overhead: Negligible for current user base

**No Async Needed Because:**
- Operations are fast enough to be imperceptible
- Low frequency of use
- Users expect to see updated skill rating immediately
- YAGNI principle applies - don't over-engineer

## Success Criteria

**Functional:**
- [x] Tracker create/update/delete triggers automatic TrueSkill recalculation
- [x] Users see updated skill ratings immediately after tracker changes
- [x] Google Sheets parity achieved for user experience

**Technical:**
- [x] All tracker operations complete successfully
- [x] TrueSkill service integration works without breaking existing functionality
- [x] Error handling prevents TrueSkill failures from breaking tracker operations
- [x] Comprehensive logging for monitoring and debugging

## Validation Checklist

- [x] TrackerHandler constructor updated with TrueSkillService parameter
- [x] All TrackerHandler instantiations updated in main.go
- [x] CreateTracker method includes TrueSkill update
- [x] UpdateTracker method includes TrueSkill update  
- [x] DeleteTracker method includes TrueSkill update
- [x] Error handling prevents failures from breaking tracker operations
- [x] Build successful with no compilation errors
- [x] Logging provides clear success/failure information
- [x] No orphaned code or broken imports

## Timeline

**Total Estimated Time: 1 hour**
- Phase 1 (Dependency Injection): 15 minutes
- Phase 2 (Auto-Update Integration): 30 minutes  
- Phase 3 (Error Handling): 15 minutes

**Simple, focused refactoring with immediate business value and Google Sheets parity.**