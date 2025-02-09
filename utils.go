package discord

import (
	"fmt"
	"log"
	"strings"
)

func truncateContent(content string) (string, bool) {
	truncated := false
	if len(content) >= 1850 { // discord content limit is 2000
		log.Printf("Content is too long (%d), truncating", len(content))
		content = content[:1850]
		content = strings.TrimRight(content, "\\")
		// prepend TRUNCATED to the content
		content = fmt.Sprintf("TRUNCATED\n%s", content)
		truncated = true
	}
	return content, truncated
}
