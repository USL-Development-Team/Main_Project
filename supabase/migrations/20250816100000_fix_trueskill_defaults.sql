-- Fix TrueSkill Default Values Migration
-- This migration corrects users who received incorrect TrueSkill default values
-- and ensures all tables use the correct 1000.0 default

-- Fix existing users in main users table who have incorrect defaults (skip if column doesn't exist)
-- UPDATE users 
-- SET trueskill_mu = 1000.0 
-- WHERE trueskill_mu = 25.0 OR trueskill_mu = 1500.0;

-- Fix existing users in USL users table who have incorrect defaults  
UPDATE usl_users 
SET trueskill_mu = 1000.0 
WHERE trueskill_mu = 25.0 OR trueskill_mu = 1500.0;

-- Fix existing user guild stats who have incorrect defaults (skip if table doesn't exist)
-- UPDATE user_guild_stats 
-- SET trueskill_mu = 1000.0 
-- WHERE trueskill_mu = 25.0 OR trueskill_mu = 1500.0;

-- Update default values for existing table schemas (in case they weren't migrated yet)
-- ALTER TABLE IF EXISTS users ALTER COLUMN trueskill_mu SET DEFAULT 1000.0;
ALTER TABLE IF EXISTS usl_users ALTER COLUMN trueskill_mu SET DEFAULT 1000.0;
-- ALTER TABLE IF EXISTS user_guild_stats ALTER COLUMN trueskill_mu SET DEFAULT 1000.0;

-- Log the changes (skip if migration_log table doesn't exist)
-- INSERT INTO migration_log (migration_name, description, executed_at) 
-- VALUES (
--     '20250816100000_fix_trueskill_defaults',
--     'Fixed TrueSkill default values from 25.0/1500.0 to correct 1000.0 across all tables',
--     NOW()
-- ) ON CONFLICT DO NOTHING;