package bunnycdn

import (
	"crypto/md5"
	"encoding/base64"
	"strconv"
	"time"
)

// Generator stores the CDN's hostname and API key
type Generator struct {
	APIKey   string
	Hostname string
}

// GenerateURL generates an expiring url to access the CDN
// expirationDuration can be in any unit as long as it's type is time.Duration
// path should begin with a forward slash like this: "/content/image.jpg"
func GenerateURL(generator Generator, path string, expirationDuration time.Duration) string {
	expirationDate := time.Now().Add(expirationDuration)
	expirationUnix := strconv.FormatInt(expirationDate.Unix(), 10)
	tokenContent := []byte(generator.APIKey + path + expirationUnix)
	hash := md5.Sum(tokenContent)
	tokenB64 := base64.RawURLEncoding.EncodeToString(hash[:]) // pass hash as a slice

	url := generator.Hostname + path + "?token=" + tokenB64 + "&expires" + expirationUnix
	return url
}
