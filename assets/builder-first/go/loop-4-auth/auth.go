// Loop 4 — auth. Sessions over JWT for simplicity and revocability.
//
// Don't paste this into your main.go yet — read through, port piece by piece,
// understand what each line is doing.

package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

// ─────────────────────────────────────────────────────────────────────────────
// Password hashing — argon2id with sensible defaults.
// ─────────────────────────────────────────────────────────────────────────────

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MiB
	argonThreads = 4
	argonKeyLen  = 32
	argonSaltLen = 16
)

// HashPassword returns a string in the format "argon2id$<salt-b64>$<hash-b64>".
// TODO: implement using argon2.IDKey.
func HashPassword(plain string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	_ = argonTime
	_ = argonMemory
	_ = argonThreads
	_ = argonKeyLen
	_ = argon2.IDKey
	return "", errors.New("TODO: implement HashPassword")
}

// VerifyPassword compares a plaintext password against a stored hash.
// MUST use a constant-time comparison (subtle.ConstantTimeCompare).
//
// TODO: parse the encoded hash, recompute, compare with subtle.ConstantTimeCompare.
func VerifyPassword(plain, encoded string) bool {
	_ = strings.SplitN
	_ = base64.RawStdEncoding
	_ = subtle.ConstantTimeCompare
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// Session management.
// ─────────────────────────────────────────────────────────────────────────────

const sessionTTL = 30 * 24 * time.Hour // 30 days

func newSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// CreateSession inserts a sessions row and returns the (id, expiresAt).
// TODO: implement.
func CreateSession(ctx context.Context, db *pgxpool.Pool, userID int64) (string, time.Time, error) {
	return "", time.Time{}, errors.New("TODO: implement CreateSession")
}

// LookupSession returns the user ID for a valid, unexpired session, or
// (0, false, nil) if missing/expired. err only on infrastructure failure.
// TODO: implement.
func LookupSession(ctx context.Context, db *pgxpool.Pool, sid string) (userID int64, ok bool, err error) {
	_ = pgx.ErrNoRows
	return 0, false, errors.New("TODO: implement LookupSession")
}

// DeleteSession removes a session (logout).
func DeleteSession(ctx context.Context, db *pgxpool.Pool, sid string) error {
	return errors.New("TODO: implement DeleteSession")
}

// ─────────────────────────────────────────────────────────────────────────────
// Cookie helpers.
// ─────────────────────────────────────────────────────────────────────────────

const cookieName = "session"

func setSessionCookie(w http.ResponseWriter, sid string, expires time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // false for local http; true once you're on https
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}

func clearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // far past
		MaxAge:   -1,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Middleware.
// ─────────────────────────────────────────────────────────────────────────────

type ctxKey int

const ctxUserID ctxKey = 0

// RequireAuth wraps a handler, requiring a valid session cookie. On miss → 401.
// On hit, stashes userID in r.Context().
//
// TODO: read the cookie, LookupSession, inject ctx, call next.
func RequireAuth(db *pgxpool.Pool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
		http.Error(w, "TODO: implement RequireAuth", http.StatusInternalServerError)
	}
}

// UserIDFrom returns the authenticated user's ID, or 0 if not authenticated.
func UserIDFrom(ctx context.Context) int64 {
	uid, _ := ctx.Value(ctxUserID).(int64)
	return uid
}
