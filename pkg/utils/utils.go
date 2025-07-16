package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func IsAllowedExtension(fileURL string, allowedExtensions []string) bool {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return false
	}

	ext := strings.ToLower(filepath.Ext(parsedURL.Path))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
