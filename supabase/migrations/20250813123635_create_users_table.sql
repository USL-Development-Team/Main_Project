-- Create users table matching the Google Sheets User schema
create table if not exists users (
  id bigint primary key generated always as identity,
  name text not null check (length(name) >= 1 and length(name) <= 255),
  discord_id text unique not null check (discord_id ~ '^\d{17,19}$'),
  active boolean not null default true,
  banned boolean not null default false,
  mmr integer not null default 0 check (mmr >= 0),
  trueskill_mu decimal(10,3) not null default 1000.000 check (trueskill_mu >= 0 and trueskill_mu <= 5000),
  trueskill_sigma decimal(10,3) not null default 8.333 check (trueskill_sigma >= 0 and trueskill_sigma <= 20),
  trueskill_last_updated timestamptz not null default now(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- Create indexes for performance
create index if not exists idx_users_discord_id on users(discord_id);
create index if not exists idx_users_active on users(active) where active = true;
create index if not exists idx_users_trueskill_last_updated on users(trueskill_last_updated);

-- Enable RLS (Row Level Security)
alter table users enable row level security;

-- Create policies (basic read access for authenticated users)
create policy "Users are viewable by authenticated users" on users
  for select using (auth.role() = 'authenticated');

create policy "Users can be inserted by authenticated users" on users
  for insert with check (auth.role() = 'authenticated');

create policy "Users can be updated by authenticated users" on users  
  for update using (auth.role() = 'authenticated');

-- Trigger to update updated_at timestamp
create or replace function handle_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;

create trigger users_updated_at
  before update on users
  for each row execute function handle_updated_at();