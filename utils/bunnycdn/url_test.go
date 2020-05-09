package bunnycdn

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateURL(t *testing.T) {
	gen := Generator{
		Hostname: "test-cdn.example.com",
		APIKey: "SECREY-KEY",
	}
	fmt.Println(gen.GenerateURL("/index.html", time.Minute * 60))
}