package drivers

import (
	"strconv"
	"strings"
)

func FindInAttribute(attributes []string, pattern string) string {
	for _, attribute := range attributes {
		if i := strings.Index(attribute, pattern); i >= 0 {
			return attribute
		}
	}
	return ""
}

func StrToInt(val string) (int, error) {
	number := ""
	for _, v := range val {
		if v >= '0' && v <= '9' {
			number += string(v)
		}
	}
	numb, err := strconv.Atoi(number)
	if err != nil {
		return 0, err
	}
	return numb, nil
}

func RemoveDuplicates(links []string) []string {
	linkMap := map[string]bool{}
	for _, link := range links {
		linkMap[link] = true
	}
	var result []string
	for link := range linkMap {
		result = append(result, link)
	}
	return result
}

func PrepareLink(link string) string {
	if clearLink, _, ok := strings.Cut(link, "?"); ok {
		return clearLink
	}
	return ""
}

func ProductFromLink(link string) string {
	if s := strings.Split(link, "/"); len(s) > 2 {
		return "/" + s[len(s)-2] + "/" + s[len(s)-1]
	}
	return ""
}
