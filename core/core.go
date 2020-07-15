package core

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetNick generates a n digit random code
func GetNick(size int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetHash gets the hash of string
func GetHash(name string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// SafeFilename gives a filename if a file with same name is present
func SafeFilename(name string) string {
	if _, err := os.Stat(name); err != nil {
		return name
	}
	ext := filepath.Ext(name)
	basename := strings.TrimSuffix(name, ext)
	ix := 1
	for {
		name := fmt.Sprintf("%s_%d%s", basename, ix, ext)
		if _, err := os.Stat(name); err != nil {
			return name
		}
		ix++
	}
}
