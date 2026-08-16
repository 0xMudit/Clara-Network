// Package env provides small helpers for reading configuration from the
// environment with defaults.
package env

import (
	"os"
	"strconv"
)

// Get returns the value of key or fallback if unset or empty.
func Get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Int returns the integer value of key or fallback if unset or invalid.
func Int(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
