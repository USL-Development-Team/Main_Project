-- Create user_trackers table matching the Google Sheets UserTracker schema
create table if not exists user_trackers (
  id bigint primary key generated always as identity,
  discord_id text not null check (discord_id ~ '^\d{17,19}$'),
  url text not null check (length(url) >= 1 and length(url) <= 1000),
  ones_current_season_peak integer not null default 0 check (ones_current_season_peak >= 0),
  ones_previous_season_peak integer not null default 0 check (ones_previous_season_peak >= 0),
  ones_all_time_peak integer not null default 0 check (ones_all_time_peak >= 0),
  ones_current_season_games integer not null default 0 check (ones_current_season_games >= 0),
  ones_previous_season_games integer not null default 0 check (ones_previous_season_games >= 0),
  twos_current_season_peak integer not null default 0 check (twos_current_season_peak >= 0),
  twos_previous_season_peak integer not null default 0 check (twos_previous_season_peak >= 0),
  twos_all_time_peak integer not null default 0 check (twos_all_time_peak >= 0),
  twos_current_season_games integer not null default 0 check (twos_current_season_games >= 0),
  twos_previous_season_games integer not null default 0 check (twos_previous_season_games >= 0),
  threes_current_season_peak integer not null default 0 check (threes_current_season_peak >= 0),
  threes_previous_season_peak integer not null default 0 check (threes_previous_season_peak >= 0),
  threes_all_time_peak integer not null default 0 check (threes_all_time_peak >= 0),
  threes_current_season_games integer not null default 0 check (threes_current_season_games >= 0),
  threes_previous_season_games integer not null default 0 check (threes_previous_season_games >= 0),
  last_updated timestamptz not null default now(),
  valid boolean not null default true,
  calculated_mmr integer not null default 0 check (calculated_mmr >= 0),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  
  -- Foreign key constraint
  constraint fk_user_trackers_discord_id foreign key (discord_id) references users(discord_id) on delete cascade,
  
  -- Unique constraint for discord_id + url combination
  constraint user_trackers_discord_id_url_unique unique (discord_id, url)
);

-- Create indexes for performance
create index if not exists idx_user_trackers_discord_id on user_trackers(discord_id);
create index if not exists idx_user_trackers_valid on user_trackers(valid) where valid = true;
create index if not exists idx_user_trackers_last_updated on user_trackers(last_updated);
create index if not exists idx_user_trackers_calculated_mmr on user_trackers(calculated_mmr);

-- Enable RLS (Row Level Security)
alter table user_trackers enable row level security;

-- Create policies (basic read access for authenticated users)
create policy "User trackers are viewable by authenticated users" on user_trackers
  for select using (auth.role() = 'authenticated');

create policy "User trackers can be inserted by authenticated users" on user_trackers
  for insert with check (auth.role() = 'authenticated');

create policy "User trackers can be updated by authenticated users" on user_trackers
  for update using (auth.role() = 'authenticated');

-- Trigger to update updated_at timestamp
create trigger user_trackers_updated_at
  before update on user_trackers
  for each row execute function handle_updated_at();