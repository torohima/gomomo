package momoencoder

import (
	"log"
	"regexp"
)

func removeSymbol(str1 string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	str1 = re.ReplaceAllString(str1, "")
	return str1
}
