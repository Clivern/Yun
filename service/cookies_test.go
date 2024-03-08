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

func TestDefaultCookieOptions(t *testing.T) {
	opts := DefaultCookieOptions()

	assert.Equal(t, "/", opts.Path)
	assert.True(t, opts.HTTPOnly)
	assert.False(t, opts.Secure)
	assert.Equal(t, http.SameSiteLaxMode, opts.SameSite)
}

func TestSecureCookieOptions(t *testing.T) {
	opts := SecureCookieOptions()

	assert.Equal(t, "/", opts.Path)
	assert.True(t, opts.HTTPOnly)
	assert.True(t, opts.Secure)
	assert.Equal(t, http.SameSiteStrictMode, opts.SameSite)
}

func TestSetCookie(t *testing.T) {
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

func TestSetCookieWithNilOptions(t *testing.T) {
	w := httptest.NewRecorder()

	SetCookie(w, "test_cookie", "test_value", nil)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	// Should use default options
	assert.Equal(t, "/", cookie.Path)
	assert.True(t, cookie.HttpOnly)
}

func TestSetCookieWithExpires(t *testing.T) {
	w := httptest.NewRecorder()

	opts := DefaultCookieOptions()
	opts.MaxAge = 3600

	SetCookie(w, "test_cookie", "test_value", opts)

	cookies := w.Result().Cookies()
	cookie := cookies[0]

	// Check that Expires is set when MaxAge is positive
	assert.False(t, cookie.Expires.IsZero())

	// Expires should be approximately Now + MaxAge seconds
	expectedExpires := time.Now().UTC().Add(time.Duration(opts.MaxAge) * time.Second)
	timeDiff := cookie.Expires.Sub(expectedExpires)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}

	// Allow 5 second tolerance for test execution time
	assert.LessOrEqual(t, timeDiff, 5*time.Second)
}

func TestSetSecureCookie(t *testing.T) {
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

func TestGetCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_cookie",
		Value: "test_value",
	})

	value := GetCookie(req, "test_cookie")
	assert.Equal(t, "test_value", value)
}

func TestGetCookieNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	value := GetCookie(req, "nonexistent_cookie")
	assert.Empty(t, value)
}

func TestHasCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "test_cookie",
		Value: "test_value",
	})

	assert.True(t, HasCookie(req, "test_cookie"))
	assert.False(t, HasCookie(req, "nonexistent_cookie"))
}

func TestDeleteCookie(t *testing.T) {
	w := httptest.NewRecorder()

	DeleteCookie(w, "test_cookie")

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.Empty(t, cookie.Value)
}

func TestDeleteCookieWithOptions(t *testing.T) {
	w := httptest.NewRecorder()

	opts := &CookieOptions{
		Path:   "/api",
		Domain: "example.com",
		Secure: true,
	}

	DeleteCookieWithOptions(w, "test_cookie", opts)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	cookie := cookies[0]
	assert.Equal(t, "test_cookie", cookie.Name)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.Equal(t, "/api", cookie.Path)
	assert.Equal(t, "example.com", cookie.Domain)
	assert.True(t, cookie.Secure)
	// Check that Expires is set to Unix epoch
	assert.True(t, cookie.Expires.Equal(time.Unix(0, 0)))
}

func TestDeleteCookieWithNilOptions(t *testing.T) {
	w := httptest.NewRecorder()

	DeleteCookieWithOptions(w, "test_cookie", nil)

	cookies := w.Result().Cookies()
	cookie := cookies[0]

	assert.Equal(t, -1, cookie.MaxAge)
	// Should use default path
	assert.Equal(t, "/", cookie.Path)
}

func TestGetAllCookies(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "cookie1", Value: "value1"})
	req.AddCookie(&http.Cookie{Name: "cookie2", Value: "value2"})
	req.AddCookie(&http.Cookie{Name: "cookie3", Value: "value3"})

	cookies := GetAllCookies(req)

	assert.Len(t, cookies, 3)

	// Verify all cookies are present
	cookieMap := make(map[string]string)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = cookie.Value
	}

	assert.Equal(t, "value1", cookieMap["cookie1"])
	assert.Equal(t, "value2", cookieMap["cookie2"])
	assert.Equal(t, "value3", cookieMap["cookie3"])
}

func TestGetAllCookiesEmpty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	cookies := GetAllCookies(req)

	assert.Empty(t, cookies)
}

func TestCookieOptionsCustomization(t *testing.T) {
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

func TestSessionCookieScenario(t *testing.T) {
	// Simulate a typical session cookie workflow
	w := httptest.NewRecorder()

	// Set a session cookie
	opts := SecureCookieOptions()
	opts.MaxAge = 3600 // 1 hour
	SetCookie(w, "session_token", "abc123def456", opts)

	// Verify cookie was set
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)

	sessionCookie := cookies[0]

	// Create a new request with the session cookie
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(sessionCookie)

	// Retrieve the session token
	token := GetCookie(req, "session_token")
	assert.Equal(t, "abc123def456", token)

	// Verify cookie exists
	assert.True(t, HasCookie(req, "session_token"))

	// Delete the session cookie
	w2 := httptest.NewRecorder()
	DeleteCookie(w2, "session_token")

	deleteCookies := w2.Result().Cookies()
	assert.Equal(t, -1, deleteCookies[0].MaxAge)
}
