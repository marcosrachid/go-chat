package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetenvDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if s, _ = reader.ReadString('\n'); true {
		s = strings.TrimSpace(s)
	}
	return s
}
