# Privilege 特權

## 1. App Overview

- **App Name:** Privilege (特權)
- **Platform:** iOS only
- **One-liner:** A two-tier dating app where paid users hold the power to
  discover and invite free users, while free users experience the anticipation
  of being chosen.
- **Core Problem:** Existing dating apps give paying users more features but not
  a fundamentally different social experience. Privilege creates a genuine power
  asymmetry — paid users feel superior and in control, free users feel the
  tension and excitement of being selected.

---

## 2. Target User

- **Free users:** 18–30, mobile-native generation, familiar with internet and
  app culture
- **Paid (Privilege) users:** 25–45, willing to pay for a status-driven, curated
  experience
- **Technical comfort:** High — users expect smooth, modern, gesture-driven UI
- **Core pain point:** Current dating apps blur the line between tiers.
  Privilege makes the difference visceral and social, not just functional

---

## 3. Privilege Tier System

Two tiers only:

|                    | Free User                                | Privilege User (Paid)              |
| ------------------ | ---------------------------------------- | ---------------------------------- |
| Feed               | Sees only Privilege users who liked them | Sees all Free users nearby         |
| Discovery          | Passive — waits to be chosen             | Active — browses and initiates     |
| Photos             | Up to 3                                  | Up to 6                            |
| UI Style           | Minimal, muted, greyscale-leaning        | Rich, gold or dark luxury themed   |
| Matching           | Accepts or rejects invitations           | Sends invitations by swiping right |
| Same-tier matching | ❌ Not allowed                           | ❌ Not allowed                     |

- Cross-tier only: a Privilege user swipes right → Free user receives an
  invitation → Free user accepts or rejects → Match formed
- Two Free users cannot match. Two Privilege users cannot match.

---

## 4. Core Features — MVP

**In scope:**

- Login screen
- Signup flow with paywall prompt
- Swiping / invitation feed (tier-specific behavior)
- Match state
- Chatroom list (matched users)
- Individual chatroom
- My Profile screen
- Edit Profile screen
- Basic info and photo upload
- Backend infrastructure (mock/local API for development)

**Out of scope for v1:**

- Cloud Backend infrastructure (API for development)
- Filtering / search
- Settings page
- Payment integration
- Threads / Instagram OAuth
- Profile verification
- Admin tools

---

## 5. Screens & Navigation

### Screen List

1. Onboarding / Privilege Info Screen
2. Login Screen
3. Signup Screen
4. Paywall Screen (Privilege upsell — shown once during signup)
5. Swipe Feed (Privilege user) — card stack of nearby Free users
6. Invitation Feed (Free user) — sealed letter stack UI; open to reveal who
   liked you
7. Chatroom List Screen (matched conversations)
8. Individual Chatroom Screen
9. My Profile Screen
10. Edit Profile Screen
11. General Setting Screen

### Navigation Structure

- **Tab Bar** (main app, post-login):
  - Feed (swipe or invitation depending on tier)
  - Matches / Chatrooms
  - My Profile
- **Navigation Stack** within each tab for drill-down (e.g. chatroom list →
  individual chat, profile → edit profile)
- **Modal** for paywall screen during signup

### Onboarding Flow

```
App Launch
  → Onboarding Info Screen (explains tier difference)
  → Login / Signup choice
      → Signup → basic info (incl. birthdate + ToS agreement) → age gate (block if under 18) → Paywall prompt (become Privilege?) → Feed
      → Login → Feed
```

---

## 6. User Profile Fields

| Field               | Type         | Notes                                                                             |
| ------------------- | ------------ | --------------------------------------------------------------------------------- |
| Name                | Text         | Display name                                                                      |
| Age                 | Number       | Derived from birthdate; birthdate stored server-side and immutable after signup   |
| Bio                 | Text         | Short free text                                                                   |
| Photos              | Images       | 3 max (free), 6 max (paid)                                                        |
| Job                 | Text         | Optional                                                                          |
| Height              | Number       | cm                                                                                |
| Weight              | Number       | kg                                                                                |
| Sex / Gender        | Select       | Male, Female, Non-binary, Genderqueer, Prefer not to say, Other                   |
| Relationship Status | Select       | Single, In a Relationship, Married, Divorced, Open Relationship, It's Complicated |
| Lifestyle Tags      | Multi-select | Pets, Workout Frequency, Habits (drinking, smoking), Religion, Politics           |

---

## 7. Matching & Interaction Logic

**Privilege user flow:**

1. Opens Feed → sees card stack of nearby Free users
2. Swipes right → sends invitation to Free user
3. Swipes left → passes, Free user never knows
4. On mutual acceptance → Match created → Chatroom unlocked

**Feed ordering (MVP):** Most recently active Free users first, filtered to
unseen only. Left-swipe is permanent — passed Free users never reappear. Backend
maintains a lightweight append-only "passed" log per Privilege user.

**Feed ordering (post-MVP):** Ranking score algorithm (TBD).

**Free user flow:**

1. Opens Feed → sees stack of sealed letter-style cards
2. Taps a letter → opens to reveal the Privilege user who liked them
3. Accepts → Match created → Chatroom unlocked
4. Ignores / rejects → invitation dismissed, no notification to Privilege user

**Chat:**

- Both sides can send the first message after matching
- Real-time messaging required — no missed messages, smooth delivery
- **Protocol:** WebSockets (Socket.io) for real-time layer; all messages
  persisted to DB as source of truth
- **Offline delivery:** Message saved to DB → APNs push notification triggers
  app to fetch on next open
- **Read receipts:** None for MVP
- **History:** Load last 100 messages on open, paginate older on scroll; no hard
  retention limit for MVP

**Invitation lifecycle notes (to be decided in later stage):**

- Invitation expiry duration and any urgency signals shown to Free users: TBD
- When a user account is deleted, all associated invitations (sent or received)
  are immediately marked void

---

## 8. Authentication

**Supported login methods (MVP):**

- Email + password
- Sign in with Apple
- Sign in with Google

**Later stage:**

- Threads
- Instagram

**Rules:**

- No guest mode — login required to use the app
- Single user per session

**Token strategy:**

- JWT access token (15-min expiry) + refresh token (30-day expiry, rotating)
- JWT payload encodes current tier (`free` / `privilege`) — tier changes
  propagate to client within one refresh cycle
- Client refreshes on every app foreground; new JWT reflects current backend
  tier (handles subscription lapse automatically)
- Both tokens stored in iOS Keychain (never UserDefaults)

---

## 9. Location

- City-level precision (not exact GPS coordinates)
- Default discovery radius: 50 km
- Radius adjustable by user (in scope post-MVP via settings)
- Used by Privilege users to surface nearby Free users in feed
- **Implementation:** iOS CoreLocation, foreground permission only (no
  background tracking)
- Location captured on every app foreground, reverse-geocoded to city name;
  lat/lng stored as geohash on backend for efficient radius queries
- Traveling users see their current location automatically — no manual city
  selection for MVP

---

## 10. Notifications

Triggered events:

- Someone (Privilege user) liked you → notify Free user
- You have a new match
- New message received
- Subscription expiring soon
- (All subject to user's notification permission settings)

---

## 11. Design & Visual Identity

**Free user experience:**

- Minimal, clean, muted palette
- Greyscale-leaning UI elements
- Functional but intentionally understated — the "waiting" experience

**Privilege user experience:**

- Rich visual treatment — gold accents OR dark luxury theme (user's choice of
  variant)
- Elevated card designs, premium feel
- Communicates power, control, and exclusivity

**General direction:**

- Modern, gesture-driven, smooth animations
- Reference: best-in-class dating apps (Hinge, Bumble) but with stronger visual
  tier differentiation
- Easy to learn, minimal cognitive load

**Modes:** Support both Light and Dark mode

---

## 12. Monetization

- **Model:** Freemium with subscription
- **Free tier:** Full access as a Free user (passive, invitation-only feed)
- **Privilege tier:** Monthly/annual subscription unlocks active discovery and
  full profile features
- **Paywall trigger:** Once during signup flow (modal), and accessible later via
  profile/upgrade screen
- **In-app purchases:** None for MVP

**Upgrade moment behavior:**

- On successful upgrade, backend issues a new JWT with `privilege` tier
  immediately
- iOS app detects tier change from refreshed token and reloads feed + UI
  in-place (no restart required)
- Pending unread invitations received as Free user are discarded — no longer
  applicable
- Photo limit updates to 6 immediately
- UI theme switches to gold/luxury immediately — upgrade feels instant and
  magical

**Subscription lapse behavior:**

- On expiry, user reverts to Free tier immediately — new JWT on next refresh
  encodes `free`
- Existing matched chats are fully preserved (both sides consented to the match)
- All pending unaccepted outgoing invitations are voided silently
- User re-enters passive Free feed (sealed letter view) — no memory of past
  swipes carried over
- No grace period — lapse is immediate

---

## 13. Technical Constraints

- **Minimum iOS version:** iOS 16+ (covers ~80%+ of target demographic as
  of 2024)
- **Xcode version:** 14.2 (macOS Monterey constraint — local dev only)
- **App Store build:** Requires Xcode 16 via GitHub Actions CI/CD pipeline
  (separate from local dev)
- **Offline support:** None — app requires active internet connection
- **Backend (dev/test):** Mock API via local Docker (Node.js or similar) for
  local development
- **Performance-critical:** Chatroom must be real-time, zero message loss,
  smooth UX
- **Data storage:** Remote-first — all user data on backend; local storage only
  for auth token

---

## 14. Account Deletion

Required by App Store guideline 5.1.1 — in-app deletion must be offered.

- **Trigger:** Delete account option in profile/settings screen
- **Immediate effect:** Account deactivated — hidden from all feeds, all sent
  invitations voided, chats locked for both sides (partner sees "this user is no
  longer available" placeholder, chat is not silently removed)
- **Grace period:** 30-day soft-delete window — user can recover their account
  within 30 days
- **Hard delete:** After 30 days, all personal data is permanently deleted (GDPR
  compliance)

---

## 15. Safety & Moderation

**MVP (stub implementation):**

- Block button: accessible from chatroom and profile view — API call logs the
  action server-side; no enforcement logic yet
- Report button: accessible from chatroom and profile view — reason selector
  (Spam, Inappropriate content, Feels unsafe) — logs to a reports table; no
  automated action
- ToS acceptance checkbox during signup for legal ground

**Post-MVP:**

- Actual block enforcement (mutual hide from feeds, chat locked)
- Report review queue and admin action
- Photo moderation pipeline

---

## 16. Tech Stack

### iOS

| Concern                    | Choice                                                      |
| -------------------------- | ----------------------------------------------------------- |
| Platform                   | iOS 16+ only                                                |
| UI Framework               | UIKit                                                       |
| Architecture               | MVVM with manual dependency injection via `AppDependencies` |
| HTTP Networking            | Alamofire                                                   |
| WebSocket (chat)           | URLSessionWebSocketTask (native)                            |
| Keychain                   | KeychainAccess                                              |
| Testing                    | XCTest                                                      |
| Local Xcode (dev)          | Xcode 14.2 (macOS Monterey constraint)                      |
| CI Xcode (App Store build) | Xcode 16 via GitHub Actions                                 |

### Backend

| Concern        | Choice                      |
| -------------- | --------------------------- |
| Language       | Go                          |
| Web Framework  | Echo                        |
| Real-time chat | Echo native WebSockets      |
| API Contract   | OpenAPI 3.x via swaggo/swag |
| Database       | PostgreSQL + PostGIS        |
| Query layer    | sqlc                        |
| Migrations     | golang-migrate              |

### Infrastructure & Services

| Concern                 | Choice                                                            |
| ----------------------- | ----------------------------------------------------------------- |
| Hosting                 | Railway                                                           |
| Auth                    | Firebase Auth                                                     |
| Push notifications      | Firebase Cloud Messaging (FCM); `firebase-admin-go` on backend    |
| Image upload & delivery | Cloudflare Images                                                 |
| CI/CD                   | GitHub Actions (iOS App Store builds + backend deploy to Railway) |
| Error tracking          | Sentry (Go backend + iOS)                                         |
| Analytics               | Firebase Analytics                                                |

---

## 17. Infrastructure & Deployment

This section describes the full workflow from a blank machine to a running cloud
environment. The guiding principle: **scripts do the work, not humans**. Every
step that must be repeated is automated; only one-time bootstrap steps are
manual.

---

### Key Design Decisions

**Migrations replace `init.sql`** There is no separate `init.sql`. The first
`golang-migrate` file is the init script — it creates all tables, indexes, and
extensions (`CREATE EXTENSION IF NOT EXISTS postgis;`). Running `migrate up` on
a fresh database applies everything in order. Schema is version-controlled like
code.

**Migrations run on app startup** The Go binary embeds `golang-migrate` and
calls `migrate up` before Echo starts accepting traffic. This means every
Railway deploy automatically migrates the production database — no manual
commands, no separate Railway job.

**PostGIS via migration, not template** Use the standard Railway Postgres
template, then enable PostGIS in migration `0001`. This is more portable than
relying on a marketplace variant and keeps the extension declaration in version
control.

---

### Step 1 — One-Time Railway Project Setup (manual, done once)

These steps cannot be scripted because they require an account and UI
authorization.

1. Create a Railway project at [railway.com](https://railway.com)
2. Add a **PostgreSQL** service from the Railway template marketplace
3. Copy the `DATABASE_URL` from Railway's Variables tab for the Postgres service
4. Create a **Backend** service (empty — GitHub Actions will deploy into it)
5. Generate a **Railway API token**: Account Settings → Tokens → New Token
6. Note the backend service's **Service ID** from its Settings tab

That's it for the dashboard. Everything from here is scripted.

---

### Step 2 — GitHub Secrets (one-time, per repo)

Set these in GitHub → Repo Settings → Secrets and variables → Actions:

| Secret               | Value                                                            |
| -------------------- | ---------------------------------------------------------------- |
| `RAILWAY_TOKEN`      | The API token from Step 1                                        |
| `RAILWAY_SERVICE_ID` | The backend service ID from Step 1                               |
| `DATABASE_URL`       | Production Postgres URL from Railway (for prod migrations in CI) |

---

### Step 3 — Repository File Structure

```
privilege-backend/
├── cmd/
│   └── server/
│       └── main.go                     # Echo setup, starts server
├── internal/
│   ├── config/
│   │   └── config.go                   # Reads env vars (PORT, DATABASE_URL, JWT_SECRET, SKIP_AUTH)
│   ├── handler/
│   │   ├── health.go                   # GET /health
│   │   ├── auth.go                     # register, login, refresh, logout
│   │   ├── user.go                     # me CRUD, photo upload, location update
│   │   ├── feed.go                     # tier-aware feed, pass (left swipe)
│   │   ├── invitation.go               # send / accept / reject invitation
│   │   ├── match.go                    # list matches
│   │   ├── chat.go                     # message history, REST send
│   │   ├── ws.go                       # WebSocket upgrade for real-time chat
│   │   └── safety.go                   # block, report
│   ├── middleware/
│   │   └── auth.go                     # JWT validation; RequireTier gate
│   └── router/
│       └── router.go                   # All route registrations
├── migrations/
│   ├── 000001_init.up.sql              # PostGIS, all tables, indexes
│   ├── 000001_init.down.sql
│   └── ...                             # future migrations via make migrate-new
├── Dockerfile                          # production image (multi-stage)
├── docker-compose.yml                  # local dev: PostGIS + backend
├── railway.toml                        # Railway build/deploy config
├── Makefile                            # make run / up / down / migrate-up / migrate-new
└── .github/
    └── workflows/
        └── deploy.yml                  # CI/CD pipeline (stub — configure after Railway setup)
```

---

### Step 4 — Key Config Files

**`railway.toml`** — declarative Railway config (committed to repo):

```toml
[build]
builder = "dockerfile"

[deploy]
startCommand = "./server"
healthcheckPath = "/health"
healthcheckTimeout = 30
restartPolicyType = "always"
```

**`Makefile`** — local dev shortcuts:

```makefile
up:
    docker compose up --build

down:
    docker compose down -v

migrate-up:
    migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
    migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-new:
    migrate create -ext sql -dir migrations -seq $(name)
```

**`docker-compose.yml`** — local environment (PostGIS + backend):

```yaml
services:
  db:
    image: postgis/postgis:16-3.4
    environment:
      POSTGRES_DB: privilege
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      retries: 5

  server:
    build: .
    depends_on:
      db:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://postgres:postgres@db:5432/privilege?sslmode=disable
    ports:
      - "8080:8080"
```

When you run `make up`, Docker builds the backend image, starts Postgres, waits
for it to be healthy, then starts the backend — which runs `migrate up` and
begins listening. No manual steps.

---

### Step 5 — GitHub Actions CI/CD Pipeline

**`.github/workflows/deploy.yml`** — triggered on every push to `main`:

```yaml
name: CI/CD

on:
  push:
    branches: [main]

jobs:
  test-and-deploy:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgis/postgis:16-3.4
        env:
          POSTGRES_DB: privilege_test
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready --health-interval 10s --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Install golang-migrate
        run:
          go install -tags 'postgres'
          github.com/golang-migrate/migrate/v4/cmd/migrate@latest

      # 1. Validate migrations against a fresh test database
      - name: Run migrations on test DB
        run: |
          migrate -path migrations \
            -database "postgres://postgres:postgres@localhost:5432/privilege_test?sslmode=disable" \
            up

      # 2. Run all Go tests (against the migrated test DB)
      - name: Run tests
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/privilege_test?sslmode=disable
        run: go test ./...

      # 3. Run migrations on production DB (safe: migrate is idempotent)
      - name: Migrate production database
        env:
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        run: |
          migrate -path migrations -database "$DATABASE_URL" up

      # 4. Deploy backend to Railway
      - name: Deploy to Railway
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
        run: |
          npm install -g @railway/cli
          railway up --service ${{ secrets.RAILWAY_SERVICE_ID }} --detach
```

**Pipeline guarantee**: tests always run against a migrated schema. Production
migration runs before the new binary is deployed. If either step fails, the
deploy does not proceed.

---

### Step 6 — Changing PostgreSQL Settings After Initial Setup

| What you want to change                              | How                                                                                                                                                                 |
| ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Add a new extension                                  | Write a new migration: `CREATE EXTENSION IF NOT EXISTS <ext>;`                                                                                                      |
| Enable PostGIS (already in `000001`)                 | No action needed after first deploy                                                                                                                                 |
| Tune `work_mem`, `max_connections`                   | Run `ALTER SYSTEM SET work_mem = '256MB';` via Railway's query console, then restart the Postgres service from the Railway dashboard. One-time.                     |
| Increase shared memory (for PostGIS heavy workloads) | Set `RAILWAY_SHM_SIZE_BYTES=134217728` (128 MB) as an env var on the Postgres service in Railway dashboard. Restart service.                                        |
| Upgrade Postgres major version                       | Railway does not support in-place major version upgrades. Provision a new Postgres service, `pg_dump` the old one, restore, update `DATABASE_URL` secret, redeploy. |

---

### Step 7 — Zero-to-Cloud Checklist

```
[ ] Step 1: Create Railway project, add Postgres service, create backend service
[ ] Step 2: Add RAILWAY_TOKEN, RAILWAY_SERVICE_ID, DATABASE_URL to GitHub Secrets
[ ] Step 3: railway.toml, Dockerfile, docker-compose.yml committed to repo
[ ] Step 4: .github/workflows/deploy.yml committed to repo
[ ] Step 5: Push to main — GitHub Actions runs the full pipeline
[ ] Step 6: Verify Railway backend service shows healthy (check /health endpoint)
[ ] Step 7: Verify Postgres has all tables (Railway → Postgres → Query)
```

After Step 5, every subsequent `git push main` is fully automated: test →
migrate → deploy.

---

### Reference

- **simplebank** (TechSchool) — Go + sqlc + golang-migrate + PostgreSQL + GitHub
  Actions:
  [github.com/techschool/simplebank](https://github.com/techschool/simplebank).
  Closest public example to this stack.
- Railway config-as-code docs:
  [docs.railway.com/guides/config-as-code](https://docs.railway.com/guides/config-as-code)
- golang-migrate library (embedded in Go):
  [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)

---

## 15. Out of Scope (Future Versions)

- Filtering and search
- Settings page
- Payment integration (Stripe or equivalent)
- Threads / Instagram OAuth
- Profile verification
- Boost / spotlight features
- Video profiles
- Voice messages
- Admin / moderation dashboard
