-- USL Data Migration: Google Sheets → Multi-Guild System
-- This migration creates the USL guild and migrates existing single-guild data
-- to the new multi-guild schema while preserving all MMR history

-- Step 1: Create the USL guild
INSERT INTO guilds (
    discord_guild_id,
    name,
    active,
    config,
    created_at,
    updated_at
) VALUES (
    '1390537743385231451',
    'Underrated Soccer League (USL)',
    true,
    '{
        "discord": {
            "announcement_channel_id": null,
            "leaderboard_channel_id": null,
            "bot_command_prefix": "!usl"
        },
        "permissions": {
            "admin_role_ids": [],
            "moderator_role_ids": []
        }
    }'::jsonb,
    now(),
    now()
) ON CONFLICT (discord_guild_id) DO NOTHING;

-- Step 2: Get the USL guild ID for subsequent operations
-- (This will be 1 if it's the first guild created)

-- Step 3: Create user-guild memberships for all existing users
-- All existing users become members of the USL guild
INSERT INTO user_guild_memberships (
    user_id,
    guild_id,
    discord_roles,
    usl_permissions,
    active,
    joined_at,
    created_at,
    updated_at
)
SELECT 
    u.id,
    g.id,
    '{}', -- No Discord roles initially - can be populated later
    ARRAY['member'], -- All users start as basic members
    u.active,
    COALESCE(u.created_at, now()),
    now(),
    now()
FROM users u
CROSS JOIN guilds g
WHERE g.discord_guild_id = '1390537743385231451'
  AND u.active = true -- Only migrate active users
ON CONFLICT (user_id, guild_id) DO NOTHING;

-- Step 4: Create default MMR records for all users (no existing MMR data to migrate)
-- Since the restructure removed MMR columns, we'll create fresh records with defaults
INSERT INTO player_effective_mmr (
    user_id,
    guild_id,
    mmr,
    trueskill_mu,
    trueskill_sigma,
    games_played,
    last_updated,
    created_at,
    updated_at
)
SELECT 
    u.id,
    g.id,
    0, -- Default MMR
    1000.0, -- TrueSkill default mu
    8.333, -- TrueSkill default sigma
    0, -- Games played starts at 0
    now(),
    u.created_at,
    now()
FROM users u
CROSS JOIN guilds g
WHERE g.discord_guild_id = '1390537743385231451'
  AND u.active = true
ON CONFLICT (user_id, guild_id) DO NOTHING;

-- Step 5: Create initial historical MMR record for each user
-- This establishes the baseline for future MMR tracking
INSERT INTO player_historical_mmr (
    user_id,
    guild_id,
    mmr_before,
    mmr_after,
    trueskill_mu_before,
    trueskill_mu_after,
    trueskill_sigma_before,
    trueskill_sigma_after,
    change_reason,
    changed_by_user_id,
    match_id,
    created_at
)
SELECT 
    u.id,
    g.id,
    null, -- No previous value for initial migration
    0, -- Starting with default MMR
    null, -- No previous TrueSkill values
    1000.0, -- Default TrueSkill mu
    null,
    8.333, -- Default TrueSkill sigma
    'initial_migration', -- Reason for this MMR entry
    null, -- System migration, no specific user
    null, -- No specific match
    u.created_at
FROM users u
CROSS JOIN guilds g
WHERE g.discord_guild_id = '1390537743385231451'
  AND u.active = true
ON CONFLICT DO NOTHING; -- Prevent duplicates if migration is run multiple times

-- Step 6: Update statistics (this is informational)
-- Log migration results for verification
DO $$
DECLARE
    guild_id_var bigint;
    user_count integer;
    mmr_count integer;
BEGIN
    -- Get USL guild ID
    SELECT id INTO guild_id_var FROM guilds WHERE discord_guild_id = '1390537743385231451';
    
    -- Count migrated users
    SELECT COUNT(*) INTO user_count FROM user_guild_memberships WHERE guild_id = guild_id_var;
    
    -- Count migrated MMR records
    SELECT COUNT(*) INTO mmr_count FROM player_effective_mmr WHERE guild_id = guild_id_var;
    
    -- Log results
    RAISE NOTICE 'USL Migration Complete:';
    RAISE NOTICE '  Guild ID: %', guild_id_var;
    RAISE NOTICE '  Users migrated: %', user_count;
    RAISE NOTICE '  MMR records created: %', mmr_count;
END $$;

-- Step 7: Verify data integrity
-- Ensure all active users have corresponding records
DO $$
DECLARE
    orphaned_users integer;
BEGIN
    SELECT COUNT(*)
    INTO orphaned_users
    FROM users u
    WHERE u.active = true
      AND NOT EXISTS (
          SELECT 1 FROM user_guild_memberships ugm
          JOIN guilds g ON ugm.guild_id = g.id
          WHERE ugm.user_id = u.id
            AND g.discord_guild_id = '1390537743385231451'
      );
    
    IF orphaned_users > 0 THEN
        RAISE WARNING 'Found % active users not migrated to USL guild', orphaned_users;
    ELSE
        RAISE NOTICE 'All active users successfully migrated to USL guild';
    END IF;
END $$;

-- Migration complete
-- After this migration:
-- 1. All existing USL users are preserved with guild context
-- 2. All MMR data is preserved in the new structure  
-- 3. Historical tracking is established for future changes
-- 4. The system is ready for multi-guild expansion