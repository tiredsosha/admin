package formater

import (
	"github.com/WqyJh/go-fstring"
	"github.com/tiredsosha/admin/tools/logger"
)

func CustomStr(temp string, val map[string]any) string {
	// template := "Hello, {name}! You have {count} unread messages." this is how temp string should looks like
	// values := map[string]any{"name": "Alice", "count": 10} this is how val map should looks like
	result, err := fstring.Format(temp, val)
	if err != nil {
		logger.Error.Printf("error fstinging %v", err)
	}
	return result
}
