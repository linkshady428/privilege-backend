-- Seed data for local development
-- All passwords are: password123

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Wipe existing seed rows so this is idempotent
DELETE FROM users WHERE email LIKE '%@seed.test';

-- Privilege users
INSERT INTO users (email, password_hash, name, birthdate, sex, tier, bio, agree_tos) VALUES
  ('alice@seed.test', crypt('password123', gen_salt('bf', 10)), 'Alice Chen',  '1992-03-15', 'female', 'privilege', 'CEO by day, surfer by dawn.',           true),
  ('bob@seed.test',   crypt('password123', gen_salt('bf', 10)), 'Bob Kim',     '1988-07-22', 'male',   'privilege', 'Architect. Dog dad. Espresso snob.',     true);

-- Free users
INSERT INTO users (email, password_hash, name, birthdate, sex, tier, bio, agree_tos) VALUES
  ('carol@seed.test', crypt('password123', gen_salt('bf', 10)), 'Carol Wu',    '1997-11-08', 'female', 'free', 'Film photographer. Bookshop hoarder.',  true),
  ('dan@seed.test',   crypt('password123', gen_salt('bf', 10)), 'Dan Osei',    '1994-05-30', 'male',   'free', 'Trail runner. Amateur chef.',            true),
  ('emma@seed.test',  crypt('password123', gen_salt('bf', 10)), 'Emma Torres', '1999-01-19', 'female', 'free', 'Ceramics nerd. Coffee shop regular.',    true);
