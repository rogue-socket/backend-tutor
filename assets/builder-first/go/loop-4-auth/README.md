# Loop 4 — Auth

**Tier mapping:** T1.7 (sessions), T1.8 (JWT) optionally, T8 (security)
**Time:** 120–180 minutes
**Theme:** *shipping anything to "real users" means knowing who they are.*
**Prereqs:** Loop 2 done (you need persistence for users + sessions).

## What you're building

User accounts and authentication. **Sessions over JWT** is the recommended default for first-party web apps (simpler, revocable, the right tool for this shape). Add JWT later if you have a use case (mobile, cross-domain SSO).

Endpoints:
```
POST   /auth/signup     {email, password}        → 201, sets session cookie
POST   /auth/login      {email, password}        → 200, sets session cookie
POST   /auth/logout                              → 204, clears + invalidates session
GET    /auth/me                                  → 200 with the current user, or 401
```

Plus: middleware that wraps `/links/*` to require authentication. Each user sees only their own links.

## Schema (new migration `0005_users_and_sessions.sql`)

```sql
CREATE TABLE users (
    id              BIGSERIAL    PRIMARY KEY,
    email           TEXT         NOT NULL UNIQUE,
    password_hash   TEXT         NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id          TEXT         PRIMARY KEY,        -- random 256-bit, base64url
    user_id     BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX sessions_user_id_idx  ON sessions(user_id);
CREATE INDEX sessions_expires_idx  ON sessions(expires_at);

-- Add owner_id to links so each user owns their links.
ALTER TABLE links ADD COLUMN owner_id BIGINT REFERENCES users(id) ON DELETE CASCADE;
-- Backfill (Loop 3 pattern) if you have existing links: assign to a "system" user.
```

## Tasks

1. **Add the migration.** Apply.
2. **Implement signup.**
   - Validate email format and password length (≥12 chars; reject anything in a small banlist of common passwords)
   - Hash password with **argon2id** (preferred) or **bcrypt** (acceptable). Do **not** use MD5, SHA-1, or unsalted SHA-256.
   - Insert user; create a session; set the cookie.
3. **Implement login.** Same shape as signup minus the user creation.
4. **Implement the auth middleware.** Reads the session cookie; looks up the session row; checks `expires_at`; injects the user ID into the request context.
5. **Implement logout.** Deletes the session row; clears the cookie.
6. **Implement `/auth/me`.** Returns the user record (without `password_hash`).
7. **Wire `/links/*` through the auth middleware.** All link operations are now scoped to the current user.
8. **Run the BREAKs.** See `BREAK.md`.

## Cookie flags (non-negotiable)

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    sessionID,
    Path:     "/",
    HttpOnly: true,                // not readable from JS — no XSS exfil
    Secure:   true,                // HTTPS only (set false for local dev only)
    SameSite: http.SameSiteLaxMode, // CSRF defense; "Strict" if you don't need cross-site links
    Expires:  expiresAt,
})
```

If any of these are missing or wrong, the audit fails. They're tested in `BREAK.md`.

## Stretch

- **Argon2id parameter tuning.** The defaults from `golang.org/x/crypto/argon2` are reasonable. Measure the hash time on your hardware; tune `time` and `memory` parameters so a single hash takes ~100ms (resistant to brute-force at scale, fast enough for login).
- **Password reset flow.** A token in a separate `password_reset_tokens` table, single-use, 1-hour expiry, sent by email (just log it for now).
- **Rate limiting on `/auth/login`.** Per-IP; 5 attempts per 5 minutes; 429 on excess. (Loop 1.6 territory; foreshadows real rate limiting in Loop 5+.)
