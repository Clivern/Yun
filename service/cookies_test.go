// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitDefaultCookieOptions(t *testing.T) {
	opts := DefaultCookieOptions()

	assert.Equal(t, "/", opts.Path)
	assert.True(t, opts.HTTPOnly)
	assert.False(t, opts.Secure)
	assert.Equal(t, http.SameSiteLaxMode, opts.SameSite)
}

func TestUnitSecureCookieOptions(t *testing.T) {
	opts := SecureCookieOptions()

	assert.Equal(t, "/", opts.Path)
	assert.True(t, opts.HTTPOnly)
	assert.True(t, opts.Secure)
	assert.Equal(t, http.SameSiteStrictMode, opts.SameSite)
}

func TestUnitSetCookie(t *testing.T) {
	w := httptest.NewRecorder()

	opts := DefaultCookieOptions()
	opts.MaxAge = 3600

	SetCookie(w, "test_cookie", "test_value", opts)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	assert.Equal(t, "test_value", cookie.Value)
	assert.Equal(t, 3600, cookie.MaxAge)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
}

func TestUnitSetCookieWithNilOptions(t *testing.T) {
	w := httptest.NewRecorder()

	SetCookie(w, "test_cookie", "test_value", nil)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
}

func TestUnitSetCookieWithExpires(t *testing.T) {
	w := httptest.NewRecorder()

	opts := DefaultCookieOptions()
	opts.MaxAge = 3600

	SetCookie(w, "test_cookie", "test_value", opts)

	cookies := w.Result().Cookies()
	cookie := cookies[0]

	assert.False(t, cookie.Expires.IsZero())

	expectedExpires := time.Now().UTC().Add(time.Duration(opts.MaxAge) * time.Second)
	timeDiff := cookie.Expires.Sub(expectedExpires)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}

	assert.LessOrEqual(t, timeDiff, 5*time.Second)
}

func TestUnitSetSecureCookie(t *testing.T) {
	w := httptest.NewRecorder()

	opts := SecureCookieOptions()
	opts.MaxAge = 86400

	SetCookie(w, "secure_cookie", "secure_value", opts)

	cookies := w.Result().Cookies()
	cookie := cookies[0]

	assert.True(t, cookie.Secure)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
}

func TestUnitGetCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_cookie",
		Value: "test_value",
	})

	value := GetCookie(req, "test_cookie")
	assert.Equal(t, "test_value", value)
}

func TestUnitGetCookieNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	value := GetCookie(req, "nonexistent_cookie")
	assert.Empty(t, value)
}

func TestUnitHasCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_cookie",
		Value: "test_value",
	})

	assert.True(t, HasCookie(req, "test_cookie"))
	assert.False(t, HasCookie(req, "nonexistent_cookie"))
}

func TestUnitDeleteCookie(t *testing.T) {
	w := httptest.NewRecorder()

	DeleteCookie(w, "test_cookie")

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.Equal(t, "/", cookie.Path)
	assert.Empty(t, cookie.Value)
}

func TestUnitCookieOptionsCustomization(t *testing.T) {
	w := httptest.NewRecorder()

	opts := &CookieOptions{
		MaxAge:   7200,
		Path:     "/api",
		Domain:   "example.com",
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	SetCookie(w, "custom_cookie", "custom_value", opts)

	cookies := w.Result().Cookies()
	cookie := cookies[0]

	assert.Equal(t, 7200, cookie.MaxAge)
	assert.Equal(t, "/api", cookie.Path)
	assert.Equal(t, "example.com", cookie.Domain)
	assert.True(t, cookie.Secure)
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, http.SameSiteNoneMode, cookie.SameSite)
}

func TestUnitSessionCookieScenario(t *testing.T) {
	w := httptest.NewRecorder()

	opts := SecureCookieOptions()
	opts.MaxAge = 3600
	SetCookie(w, "session_token", "abc123def456", opts)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	sessionCookie := cookies[0]

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(sessionCookie)

	token := GetCookie(req, "session_token")
	assert.Equal(t, "abc123def456", token)

	assert.True(t, HasCookie(req, "session_token"))

	w2 := httptest.NewRecorder()
	DeleteCookie(w2, "session_token")

	deleteCookies := w2.Result().Cookies()
	assert.Equal(t, -1, deleteCookies[0].MaxAge)
}
