# WIN — Loop 4 done

## Code

- [ ] Signup, login, logout, /me — all implemented and tested
- [ ] Passwords stored with argon2id (or bcrypt as a defensible second choice); never plaintext, never reversible hashes
- [ ] Session IDs are 256 bits of crypto/rand, base64url-encoded
- [ ] Cookies set with: `HttpOnly`, `Secure` (in non-dev), `SameSite=Lax`, scoped path, expiry
- [ ] `/links/*` requires authentication; users see only their own links
- [ ] Session rotation on login *and* logout

## Verification

- [ ] CSRF attack page (BREAK 1) fails against your service with SameSite=Lax + Content-Type check
- [ ] No password or token appears in logs from any auth endpoint
- [ ] Session-fixation test passes: pre-login session ID is not valid post-login

## Understanding

1. **Why argon2id over bcrypt over PBKDF2 over SHA-256?**
   *Outline: argon2id is memory-hard (resists GPU/ASIC brute-force); bcrypt is the second-best widely-supported option (still resistant but less so); PBKDF2 is acceptable when constrained to FIPS-approved primitives but slower per unit of resistance. SHA-256 (even salted) is unsuitable — it's GPU-friendly and a billion-hash/sec rig is commodity. Salted, hashed, and slow are the three required properties.*

2. **Why is `SameSite=Lax` enough for most state-changing routes?**
   *Outline: it tells browsers not to send the cookie on cross-site POSTs/PUTs/DELETEs (the vectors for CSRF). Top-level navigations (clicking a link to your site) still get the cookie, which is what users expect. `Strict` is safer but breaks "log in via email link" flows; `Lax` is the right default. None means no CSRF protection at all.*

3. **You set a 30-day session TTL. A user logs in on a public computer and forgets to log out. What's the recovery path, and what's the design implication?**
   *Outline: server-side session storage means logout is real (DELETE FROM sessions WHERE id = ...). User can also log out from "all sessions" via a UI that lists their active sessions. With JWT (stateless), there's no equivalent without a revocation list — which is essentially server-side state, defeating the JWT advantage. This is one of the strongest arguments for sessions over JWT for first-party web.*

## Reflection

What surprised you? Common ones:
- Cookie flags are easy to get wrong; the defaults are *all bad*
- argon2id is slow on purpose — that's the security feature
- "Just JWT it" is wrong way more often than people admit

## What's next

Loop 5 — async work via a queue. You'll move "send notification" off the request path, and meet at-least-once delivery, idempotent consumers, and dead-letter queues.
