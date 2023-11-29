// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"net/http"
	"time"
)

// CookieOptions holds configuration for cookie creation.
type CookieOptions struct {
	// MaxAge is the maximum age of the cookie in seconds.
	// If zero, the cookie is deleted when the browser closes (session cookie).
	// If negative, the cookie is deleted immediately.
	MaxAge int

	// Path specifies the cookie path. Defaults to "/".
	Path string

	// Domain specifies the cookie domain.
	Domain string

	// Secure indicates if the cookie should only be sent over HTTPS.
	Secure bool

	// HttpOnly makes the cookie inaccessible to JavaScript.
	HttpOnly bool

	// SameSite controls the cookie's SameSite attribute.
	// Options: SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode.
	SameSite http.SameSite
}

// DefaultCookieOptions returns sensible default options for cookies.
//
// Example:
//
//	opts := DefaultCookieOptions()
//	opts.MaxAge = 3600 // 1 hour
//	SetCookie(w, "session_token", "abc123", opts)
func DefaultCookieOptions() *CookieOptions {
	return &CookieOptions{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

// SecureCookieOptions returns secure cookie options for production use.
// These options should be used when serving over HTTPS.
//
// Example:
//
//	opts := SecureCookieOptions()
//	opts.MaxAge = 86400 // 24 hours
//	SetCookie(w, "session_token", token, opts)
func SecureCookieOptions() *CookieOptions {
	return &CookieOptions{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

// SetCookie sets a cookie with the given name, value, and options.
//
// Example:
//
//	opts := DefaultCookieOptions()
//	opts.MaxAge = 3600 // 1 hour
//	SetCookie(w, "user_preference", "dark_mode", opts)
func SetCookie(w http.ResponseWriter, name, value string, options *CookieOptions) {
	if options == nil {
		options = DefaultCookieOptions()
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}

	// If MaxAge is set, also set Expires for older browsers
	if options.MaxAge > 0 {
		cookie.Expires = time.Now().Add(time.Duration(options.MaxAge) * time.Second)
	}

	http.SetCookie(w, cookie)
}

// GetCookie retrieves a cookie value by name.
// Returns an empty string if the cookie is not found.
//
// Example:
//
//	token := GetCookie(r, "session_token")
//	if token == "" {
//		// No session cookie found
//	}
func GetCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// HasCookie checks if a cookie with the given name exists.
//
// Example:
//
//	if HasCookie(r, "session_token") {
//		// User has an active session
//	}
func HasCookie(r *http.Request, name string) bool {
	_, err := r.Cookie(name)
	return err == nil
}

// DeleteCookie deletes a cookie by setting its MaxAge to -1.
//
// Example:
//
//	DeleteCookie(w, "session_token")
func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

// DeleteCookieWithOptions deletes a cookie with specific path and domain.
// This is useful when you need to delete cookies set with custom options.
//
// Example:
//
//	opts := &CookieOptions{
//		Path:   "/api",
//		Domain: "example.com",
//	}
//	DeleteCookieWithOptions(w, "api_token", opts)
func DeleteCookieWithOptions(w http.ResponseWriter, name string, options *CookieOptions) {
	if options == nil {
		options = DefaultCookieOptions()
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	}

	http.SetCookie(w, cookie)
}

// GetAllCookies returns all cookies from the request.
//
// Example:
//
//	cookies := GetAllCookies(r)
//	for _, cookie := range cookies {
//		fmt.Printf("%s: %s\n", cookie.Name, cookie.Value)
//	}
func GetAllCookies(r *http.Request) []*http.Cookie {
	return r.Cookies()
}
