# BREAK — Loop 4

Three breaks. Run all three.

---

## BREAK 1 — CSRF on a state-changing route

**Foundations filled here:** T8.3 (XSS / CSRF), T1.7 (sessions).

### The setup

You ship `POST /links` requiring a session cookie. A logged-in user visits attacker.example with an HTML page like:

```html
<form action="https://yourservice.example/links" method="POST" enctype="text/plain">
  <input name='{"url":"https://evil.example","title":"x","x":"' value='"}'>
</form>
<script>document.forms[0].submit()</script>
```

Without CSRF protection, the browser **sends the user's session cookie with the cross-site POST**. Attacker injects a link into the user's account.

### Why this works

Cookies are sent on cross-origin requests by default. The `text/plain` content type avoids the CORS preflight that `application/json` would trigger.

### The fix (in order of robustness)

1. **`SameSite=Lax` on the session cookie** — blocks cookies on cross-site POSTs. **Lax**, not None. Gets you ~90% there.
2. **Strict Content-Type check** — only accept `application/json`. Rejects the `text/plain` trick. (You did this in Loop 1; double-check it's still there.)
3. **CSRF token (defense in depth)** — for form-style submissions, double-submit cookie token.

Verify: try the attack page locally against your service with each defense in turn, confirm it fails.

---

## BREAK 2 — Token / password leakage in logs

**Foundations filled here:** T8.4 (secrets management).

### The setup

You added a request log middleware that prints something like:

```go
log.Printf("%s %s %s", r.Method, r.URL.Path, r.Header.Get("Authorization"))
```

Or:
```go
body, _ := io.ReadAll(r.Body)
log.Printf("body: %s", body)
```

Hit `POST /auth/login` with body `{"email":"a@b.com","password":"hunter2"}`. Now `grep hunter2 logs/*` finds it.

### Why this happens by accident

It's never written deliberately. Common patterns:

- A debug `log.Printf("req: %+v", r)` left in
- A 500-handler that dumps the request body
- A test helper that's printf'd in a debug session and never reverted
- A third-party logging middleware that's helpful by default

### The fix

1. **Never log request bodies on auth routes.** Per-route middleware decides what's safe.
2. **Scrub Authorization headers, cookies, and known-sensitive fields** in your generic logger. A small allowlist of safe fields is more reliable than a denylist.
3. **Run `grep -i password logs/*`** in CI. If it ever finds anything, fail the build.

This is also where the practice of *short-lived* tokens earns its keep — even a leaked token expires soon.

---

## BREAK 3 — Session fixation (and rotation on login)

**Foundations filled here:** T1.7 (sessions), T8.6 (authorization).

### The setup

User visits a public page; you set a session cookie tied to the *anonymous* user (e.g., for a shopping cart). User logs in. You associate the *existing* session ID with the now-authenticated user.

Attacker discovered this and:
1. Hit your public page; got session ID `abc123`.
2. Sent the user a link with `Set-Cookie: session=abc123` (e.g., via a domain they control or a sub-domain).
3. User logged in. Their session ID is now `abc123` — and the attacker still has it.

### The fix

**Rotate the session ID on every privilege change**, especially login:
1. User logs in.
2. *Delete* the old session row.
3. *Create* a fresh session ID.
4. Set the new cookie.

Same on logout: delete the session server-side, clear the cookie. Don't reuse session IDs.

### Test

Write an integration test:
1. GET a page (cookie set: `S1`)
2. POST `/auth/login` with valid credentials
3. Assert: response sets a *new* session cookie value, and `S1` no longer authenticates.

---

## The takeaway

Auth correctness is **binary**. "Almost right" gets you on the front page of TechCrunch.

The three patterns:
- **SameSite=Lax + strict Content-Type** beats CSRF
- **Audit your logs** for sensitive data; scrub at the middleware level
- **Rotate sessions on login/logout** to defeat fixation

Add these to your review queue. They'll come up again — security checklists, code reviews, every audit.
