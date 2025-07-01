package utils

import "strings"

func ExtractNameFromEmail(email string) string {
  parts := strings.SplitN(email, "@", 2)
  if len(parts) > 1 && parts[0] != "" {
    return parts[0]
  }
  return "User"
}