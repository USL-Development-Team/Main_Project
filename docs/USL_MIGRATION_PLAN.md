# USL Migration Plan: Google Sheets → Web Application

## Overview
Convert USL's current Google Sheets-based MMR system to a web application while building infrastructure for future guilds to use the same system.

## Current State Analysis
- **USL Discord Guild ID**: `1390537743385231451`
- **Current System**: Google Sheets + Discord Bot
- **User Base**: Existing USL members with MMR data in sheets
- **MMR Tracking**: Rocket League tracker URLs, TrueSkill calculations

## Migration Strategy

### Phase 1: Data Export from Google Sheets
**Objective**: Extract all current USL data for migration

**Required Data**:
- User list (Discord ID, Name, Active status)
- Current MMR values (MMR, TrueSkill Mu, TrueSkill Sigma)
- Tracker URLs and validation status
- Historical MMR changes (if available)
- Admin/moderator roles

**Export Format**: CSV or JSON files for import

### Phase 2: Database Migration
**Objective**: Load USL data into new multi-guild schema

**Steps**:
1. **Create USL Guild Record**
   ```sql
   INSERT INTO guilds (discord_guild_id, name, active, config)
   VALUES ('1390537743385231451', 'Underrated Soccer League (USL)', true, {...});
   ```

2. **Import Users**
   - Create base user records (identity only)
   - Associate all users with USL guild via `user_guild_memberships`
   - Set appropriate permissions (admin, moderator, member)

3. **Import MMR Data**
   - Load into `player_effective_mmr` table with `guild_id = USL_GUILD_ID`
   - Create initial `player_historical_mmr` records for audit trail
   - Import tracker URLs into `user_trackers` table

### Phase 3: Bot Integration Update
**Objective**: Update Discord bot to work with new web app API

**Changes Required**:
- Bot commands now query web app API instead of sheets
- All operations scoped to USL guild (guild_id = USL_GUILD_ID)
- Maintain current command syntax for zero user disruption

### Phase 4: Web Interface Launch
**Objective**: Provide web access to USL members

**Features**:
- Discord OAuth login (scoped to USL members)
- MMR leaderboards
- Tracker management
- Admin panel for USL staff

## Data Migration Script Structure

```sql
-- 1. Create USL guild
INSERT INTO guilds (discord_guild_id, name, config) VALUES (...);

-- 2. Import users from sheets data
INSERT INTO users (discord_id, name, active) VALUES (...);

-- 3. Create guild memberships  
INSERT INTO user_guild_memberships (user_id, guild_id, usl_permissions) VALUES (...);

-- 4. Import MMR data
INSERT INTO player_effective_mmr (user_id, guild_id, mmr, trueskill_mu, trueskill_sigma) VALUES (...);

-- 5. Create audit trail
INSERT INTO player_historical_mmr (user_id, guild_id, mmr_after, change_reason) VALUES (...);

-- 6. Import trackers
INSERT INTO user_trackers (discord_id, url, valid) VALUES (...);
```

## Critical Considerations

### 1. Zero Downtime Requirement
- Current Discord bot must continue working during migration
- Gradual rollover rather than hard cutover
- Fallback plan to revert to sheets if issues occur

### 2. Data Integrity
- Verify all USL users are migrated correctly
- Ensure MMR values match exactly
- Preserve all historical data where possible

### 3. Permission Mapping
- Map current USL admin roles to new permission system
- Ensure proper access controls from day one

### 4. Testing Strategy
- Test migration on copy of production data
- Verify bot integration with test environment
- User acceptance testing with USL admins

## Post-Migration Verification

### Data Validation Queries
```sql
-- Verify all users migrated
SELECT COUNT(*) FROM user_guild_memberships ugm 
JOIN guilds g ON ugm.guild_id = g.id 
WHERE g.discord_guild_id = '1390537743385231451';

-- Verify MMR data integrity
SELECT COUNT(*) FROM player_effective_mmr pem
JOIN guilds g ON pem.guild_id = g.id
WHERE g.discord_guild_id = '1390537743385231451';

-- Check for orphaned records
SELECT COUNT(*) FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM user_guild_memberships ugm WHERE ugm.user_id = u.id
);
```

### Functional Testing
- [ ] Discord OAuth login works for USL members
- [ ] MMR leaderboards display correctly
- [ ] Tracker URLs validate properly
- [ ] Bot commands return correct data
- [ ] Admin functions work as expected

## Timeline
1. **Week 1**: Export data from Google Sheets
2. **Week 2**: Create and test migration scripts
3. **Week 3**: Update bot integration
4. **Week 4**: Web interface testing
5. **Week 5**: Production migration and launch

## Rollback Plan
If critical issues occur:
1. Revert Discord bot to sheets integration
2. Preserve migrated data for investigation
3. Fix issues in staging environment
4. Retry migration when ready

---

**Next Steps**: 
1. Coordinate with USL admins to export Google Sheets data
2. Create data import scripts based on exported format
3. Set up staging environment for testing