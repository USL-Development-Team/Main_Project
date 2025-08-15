# USL-TrueSkill Integration Refactor Plan
*Generated: 2025-08-15*

## Overview
Integrate the USL (Ultra Soccer League) migration system with the existing sophisticated TrueSkill service to enable automatic TrueSkill updates when USL trackers are created or modified.

## Current State Analysis

### USL System
- **MigrationHandler**: Handles USL CRUD operations (users/trackers) but **NO TrueSkill integration**
- **USL Models**: Similar to core but with field naming differences and `*string` vs `time.Time` types
- **USL Repository**: Direct Supabase access, returns `[]*USLUserTracker` format
- **Current Flow**: Create tracker → Store in database → Return response (**Missing TrueSkill step**)

### TrueSkill System
- **UserTrueSkillService**: Sophisticated percentile-based TrueSkill calculation engine
- **Expected Input**: Core `*models.Tracker` format via `TrackerRepository` interface
- **Current Integration**: Used by main app, but **NOT by USL system**
- **Service Dependencies**: TrackerRepo, UserRepo, MMRCalculator, UncertaintyCalculator, DataTransformation, Config

### Integration Gap
USL operates in isolation from TrueSkill system - trackers created but TrueSkill never updated.

## Integration Tasks

### Step 1: Service Dependencies Integration ⏱️ 30min | 🔧 Low
**Problem**: USL MigrationHandler lacks access to TrueSkill services

**Current Constructor** (internal/usl/handlers/migration_handler.go:24):
```go
func NewMigrationHandler(uslRepo *usl.USLRepository, templates *template.Template) *MigrationHandler
```

**Required Changes**:
```go
func NewMigrationHandler(
    uslRepo *usl.USLRepository,
    templates *template.Template,
    trueskillService *services.UserTrueSkillService,  // Add this
    config *config.Config,                            // Add this
) *MigrationHandler
```

**Files to Update**:
- `internal/usl/handlers/migration_handler.go:19-30` - Update struct and constructor
- `cmd/server/main.go:243` - Update dependency injection in setupUSLRoutes

---

### Step 2: Data Model Mapping Implementation ⏱️ 45min | 🔧 Medium
**Problem**: USL models don't map directly to TrueSkill service input formats

**Field Mapping Required**:
```go
// USL Format → Core/TrackerData Format
DiscordID                       → DiscordID
URL                             → URL  
OnesCurrentSeasonPeak           → OnesCurrentPeak
OnesCurrentSeasonGamesPlayed    → OnesCurrentGames
OnesPreviousSeasonPeak          → OnesPreviousPeak
OnesPreviousSeasonGamesPlayed   → OnesPreviousGames
// ... (similar for Twos/Threes)
*string LastUpdated             → time.Time LastUpdated
```

**New Mapping Function**:
```go
func (h *MigrationHandler) mapUSLTrackerToTrackerData(uslTracker *usl.USLUserTracker) *services.TrackerData {
    var lastUpdated time.Time
    if uslTracker.LastUpdated != nil && *uslTracker.LastUpdated != "" {
        if parsed, err := time.Parse(time.RFC3339, *uslTracker.LastUpdated); err == nil {
            lastUpdated = parsed
        } else {
            lastUpdated = time.Now()
        }
    } else {
        lastUpdated = time.Now()
    }

    return &services.TrackerData{
        DiscordID:           uslTracker.DiscordID,
        URL:                uslTracker.URL,
        OnesCurrentPeak:     uslTracker.OnesCurrentSeasonPeak,
        OnesCurrentGames:    uslTracker.OnesCurrentSeasonGamesPlayed,
        OnesPreviousPeak:    uslTracker.OnesPreviousSeasonPeak,
        OnesPreviousGames:   uslTracker.OnesPreviousSeasonGamesPlayed,
        TwosCurrentPeak:     uslTracker.TwosCurrentSeasonPeak,
        TwosCurrentGames:    uslTracker.TwosCurrentSeasonGamesPlayed,
        TwosPreviousPeak:    uslTracker.TwosPreviousSeasonPeak,
        TwosPreviousGames:   uslTracker.TwosPreviousSeasonGamesPlayed,
        ThreesCurrentPeak:   uslTracker.ThreesCurrentSeasonPeak,
        ThreesCurrentGames:  uslTracker.ThreesCurrentSeasonGamesPlayed,
        ThreesPreviousPeak:  uslTracker.ThreesPreviousSeasonPeak,
        ThreesPreviousGames: uslTracker.ThreesPreviousSeasonGamesPlayed,
        LastUpdated:        lastUpdated,
    }
}
```

---

### Step 3: Repository Bridge Pattern ⏱️ 1hour | 🔧 High
**Problem**: TrueSkill service expects core repository interfaces, USL has direct Supabase access

**Current TrueSkill Dependencies**:
- `*repositories.TrackerRepository` (core interface)
- `*repositories.UserRepository` (core interface)

**USL Repository**:
- `*usl.USLRepository` (direct Supabase client)

**Solution Options**:
1. **Option A**: Create adapter pattern (quick fix)
2. **Option B**: Modify TrueSkill service to accept TrackerData directly (better long-term)

**Recommended: Option B - Direct TrackerData Integration**:
```go
// Add new method to UserTrueSkillService
func (s *UserTrueSkillService) UpdateUserTrueSkillFromTrackerData(trackerData *TrackerData) *TrueSkillUpdateResult {
    // Validate tracker data
    if err := s.dataTransformationService.ValidateTrackerData(trackerData); err != nil {
        return &TrueSkillUpdateResult{
            Success: false,
            Error:   fmt.Sprintf("invalid tracker data: %v", err),
        }
    }

    // Calculate TrueSkill values
    trueSkillResult, err := s.calculateTrueSkillValues(trackerData)
    if err != nil {
        return &TrueSkillUpdateResult{
            Success: false,
            Error:   fmt.Sprintf("calculation failed: %v", err),
        }
    }

    // Update user with new TrueSkill values
    err = s.updateUserWithTrueSkillValues(trackerData.DiscordID, trueSkillResult)
    if err != nil {
        return &TrueSkillUpdateResult{
            Success: false,
            Error:   fmt.Sprintf("update failed: %v", err),
        }
    }

    return &TrueSkillUpdateResult{
        Success:         true,
        HadTrackers:     true,
        TrueSkillResult: trueSkillResult,
    }
}
```

---

### Step 4: TrueSkill Update Integration ⏱️ 45min | 🔧 Medium
**Problem**: USL CreateTracker flow doesn't trigger TrueSkill updates

**Current Flow** (internal/usl/handlers/migration_handler.go:278-334):
```go
func (h *MigrationHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
    // 1. Parse form data
    // 2. Create USL tracker  
    // 3. Store in database
    // 4. Return response
    // ❌ NO TRUESKILL UPDATE!
}
```

**New Integrated Flow**:
```go
func (h *MigrationHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
    // 1. Parse form data
    tracker := parseTrackerFromForm(r)

    // 2. Store tracker in USL tables
    created, err := h.uslRepo.CreateTracker(tracker)
    if err != nil {
        http.Error(w, "Failed to create tracker", http.StatusInternalServerError)
        return
    }

    // 3. ✅ NEW: Trigger TrueSkill update
    trackerData := h.mapUSLTrackerToTrackerData(created)
    result := h.trueskillService.UpdateUserTrueSkillFromTrackerData(trackerData)
    
    if !result.Success {
        log.Printf("TrueSkill update failed for %s: %s", created.DiscordID, result.Error)
        // Continue anyway - tracker was created successfully
    }

    // 4. ✅ NEW: Sync TrueSkill results back to USL user table
    if result.Success {
        err = h.uslRepo.UpdateUserTrueSkill(
            created.DiscordID,
            result.TrueSkillResult.Mu,
            result.TrueSkillResult.Sigma,
        )
        if err != nil {
            log.Printf("Failed to sync TrueSkill to USL user table: %v", err)
        }
    }

    // 5. Return enhanced response
    h.renderTemplate(w, "usl_tracker_form.html", map[string]interface{}{
        "Tracker":          created,
        "TrueSkillUpdated": result.Success,
        "TrueSkillResult":  result.TrueSkillResult,
        "Message":          "Tracker created and TrueSkill updated successfully",
    })
}
```

---

### Step 5: Data Storage Integration ⏱️ 1hour | 🔧 High
**Problem**: TrueSkill service updates core tables, USL expects USL tables

**Core TrueSkill Updates**:
- `player_effective_mmr(user_id, guild_id, mmr, trueskill_mu, trueskill_sigma)`
- `player_historical_mmr(user_id, guild_id, mmr_before, mmr_after, change_reason)`

**USL Tables**:
- `usl_users(trueskill_mu, trueskill_sigma, trueskill_last_updated)`

**Required: Dual Update Strategy**
1. **Core System**: TrueSkill service updates core tables (existing functionality)
2. **USL Sync**: New method to sync TrueSkill results to USL tables

**New USL Repository Method**:
```go
func (r *USLRepository) UpdateUserTrueSkill(discordID string, mu, sigma float64) error {
    updateData := map[string]interface{}{
        "trueskill_mu":           mu,
        "trueskill_sigma":        sigma,
        "trueskill_last_updated": time.Now().Format(time.RFC3339),
    }

    _, err := r.client.From("usl_users").
        Update(updateData, "", "").
        Eq("discord_id", discordID).
        Execute()

    if err != nil {
        return fmt.Errorf("failed to update TrueSkill for user %s: %w", discordID, err)
    }

    return nil
}
```

---

### Step 6: Configuration Integration ⏱️ 30min | 🔧 Low
**Problem**: TrueSkill service uses core config, USL might need custom settings

**Current Config**: Core TrueSkill configuration exists in main config
**Required**: Verify USL can use existing TrueSkill configuration or add USL-specific overrides if needed

**Potential USL Config Extension**:
```go
type USLConfig struct {
    GuildID              string            `json:"guild_id"`
    UseCustomRanks       bool              `json:"use_custom_ranks"`
    TrueSkillOverrides   *TrueSkillConfig  `json:"trueskill_overrides,omitempty"`
}
```

**Note**: May not be needed if USL can use core TrueSkill config directly.

---

### Step 7: Error Handling & Rollback Strategy ⏱️ 30min | 🔧 Medium
**Problem**: What happens if TrueSkill update fails after tracker creation?

**Strategy**: **Continue-on-Error with Logging**
- Tracker creation succeeds even if TrueSkill fails
- Log TrueSkill failures for manual review
- Include TrueSkill status in user response
- Provide retry mechanisms for failed TrueSkill updates

**Error Handling Pattern**:
```go
func (h *MigrationHandler) CreateTracker(...) {
    // 1. Create tracker (must succeed)
    created, err := h.uslRepo.CreateTracker(tracker)
    if err != nil {
        return handleError(w, "Failed to create tracker", err)
    }

    // 2. Try TrueSkill update (log failure, don't abort)
    result := h.trueskillService.UpdateUserTrueSkillFromTrackerData(trackerData)
    if !result.Success {
        log.Printf("WARNING: Tracker created but TrueSkill update failed for %s: %s", 
                  created.DiscordID, result.Error)
        // Could queue for retry later
    }

    // 3. Success response includes TrueSkill status
    h.renderSuccessWithTrueSkillStatus(w, created, result)
}
```

## Implementation Order

1. **Step 1**: Service Dependencies (30min) - Foundation
2. **Step 2**: Data Model Mapping (45min) - Core transformation logic  
3. **Step 3**: Repository Bridge (1hr) - Service communication
4. **Step 4**: TrueSkill Integration (45min) - Main feature
5. **Step 5**: Data Storage Sync (1hr) - Dual table updates
6. **Step 6**: Configuration (30min) - Settings integration
7. **Step 7**: Error Handling (30min) - Production readiness

**Total Estimated Time**: 4.5 hours

## Validation Checklist

### Pre-Integration Tests
- [ ] Current USL tracker creation works
- [ ] Current TrueSkill service works independently
- [ ] Core TrueSkill integration works for main app

### Post-Integration Tests
- [ ] USL tracker creation triggers TrueSkill update
- [ ] TrueSkill values sync to both core and USL tables
- [ ] Error handling works (TrueSkill fails, tracker still created)
- [ ] Performance acceptable (additional service calls)
- [ ] Logging and monitoring in place

### Integration Test Scenarios
1. **Happy Path**: Create tracker → TrueSkill calculates → Both tables updated
2. **TrueSkill Failure**: Create tracker → TrueSkill fails → Tracker still created, error logged
3. **Data Validation**: Invalid tracker data → Proper error handling
4. **User Not Found**: TrueSkill update for non-existent user → Graceful handling

## Risk Assessment

### Low Risk
- Service dependency injection (Step 1)
- Configuration integration (Step 6)

### Medium Risk  
- Data model mapping (Step 2)
- TrueSkill update integration (Step 4)
- Error handling (Step 7)

### High Risk
- Repository bridge pattern (Step 3) - Complex service interactions
- Data storage integration (Step 5) - Dual table update coordination

## Success Criteria

1. ✅ **Functional**: USL tracker creation automatically updates TrueSkill
2. ✅ **Reliable**: System handles TrueSkill failures gracefully
3. ✅ **Performance**: Integration doesn't significantly slow USL operations
4. ✅ **Data Integrity**: Both core and USL tables stay synchronized
5. ✅ **Maintainable**: Code follows existing patterns and conventions

## File Change Summary

### New Files
- None (all integration within existing files)

### Modified Files
- `internal/usl/handlers/migration_handler.go` - Add TrueSkill integration
- `internal/services/trueskill_service.go` - Add TrackerData input method
- `internal/usl/repository.go` - Add UpdateUserTrueSkill method
- `cmd/server/main.go` - Update dependency injection

### Lines of Code
- **Added**: ~150 lines (mapping functions, integration logic, error handling)
- **Modified**: ~50 lines (constructors, method signatures)
- **Total Impact**: ~200 lines across 4 files

## Bottom Line

This is a **half-day integration task** that leverages the existing sophisticated TrueSkill system. The USL system simply needs to:

1. **Connect** to the TrueSkill service (dependency injection)
2. **Transform** its data format to TrueSkill input format  
3. **Trigger** TrueSkill calculations when trackers are created
4. **Sync** results back to USL tables for consistency

**No rewriting required** - just thoughtful integration of two existing, working systems.