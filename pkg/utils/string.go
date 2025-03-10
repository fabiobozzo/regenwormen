package utils

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var emojiRegex = regexp.MustCompile(`\s[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F700}-\x{1F77F}]|[\x{1F780}-\x{1F7FF}]|[\x{1F800}-\x{1F8FF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA00}-\x{1FA6F}]|[\x{1FA70}-\x{1FAFF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]|[\x{1F1E6}-\x{1F1FF}]|[\x{1F191}-\x{1F251}]|[\x{1F004}]|[\x{1F0CF}]|[\x{1F18E}]|[\x{1F191}-\x{1F251}]|[\x{1F004}]|[\x{1F0CF}]|[\x{1F18E}]|[\x{1F201}-\x{1F202}]|[\x{1F21A}]|[\x{1F22F}]|[\x{1F232}-\x{1F23A}]|[\x{1F250}-\x{1F251}]|[\x{1F300}-\x{1F320}]|[\x{1F32D}-\x{1F335}]|[\x{1F337}-\x{1F37C}]|[\x{1F37E}-\x{1F393}]|[\x{1F3A0}-\x{1F3CA}]|[\x{1F3CF}-\x{1F3D3}]|[\x{1F3E0}-\x{1F3F0}]|[\x{1F3F4}]|[\x{1F3F8}-\x{1F43E}]|[\x{1F440}]|[\x{1F442}-\x{1F4FC}]|[\x{1F4FF}-\x{1F53D}]|[\x{1F54B}-\x{1F54E}]|[\x{1F550}-\x{1F567}]|[\x{1F57A}]|[\x{1F595}-\x{1F596}]|[\x{1F5A4}]|[\x{1F5FB}-\x{1F64F}]|[\x{1F680}-\x{1F6C5}]|[\x{1F6CC}]|[\x{1F6D0}]|[\x{1F6D1}-\x{1F6D2}]|[\x{1F6EB}-\x{1F6EC}]|[\x{1F6F4}-\x{1F6F9}]|[\x{1F910}-\x{1F93A}]|[\x{1F93C}-\x{1F93E}]|[\x{1F940}-\x{1F945}]|[\x{1F947}-\x{1F970}]|[\x{1F973}-\x{1F976}]|[\x{1F97A}]|[\x{1F97C}-\x{1F9A2}]|[\x{1F9B0}-\x{1F9B9}]|[\x{1F9C0}-\x{1F9C2}]|[\x{1F9D0}-\x{1F9FF}]|[\x{1FA70}-\x{1FA73}]|[\x{1FA78}-\x{1FA7A}]|[\x{1FA80}-\x{1FA82}]|[\x{1FA90}-\x{1FA95}]|[\x{1F004}]|[\x{1F0CF}]|[\x{1F18E}]|[\x{1F201}-\x{1F202}]|[\x{1F21A}]|[\x{1F22F}]|[\x{1F232}-\x{1F23A}]|[\x{1F250}-\x{1F251}]|[\x{1F300}-\x{1F320}]|[\x{1F32D}-\x{1F335}]|[\x{1F337}-\x{1F37C}]|[\x{1F37E}-\x{1F393}]|[\x{1F3A0}-\x{1F3CA}]|[\x{1F3CF}-\x{1F3D3}]|[\x{1F3E0}-\x{1F3F0}]|[\x{1F3F4}]|[\x{1F3F8}-\x{1F43E}]|[\x{1F440}]|[\x{1F442}-\x{1F4FC}]|[\x{1F4FF}-\x{1F53D}]|[\x{1F54B}-\x{1F54E}]|[\x{1F550}-\x{1F567}]|[\x{1F57A}]|[\x{1F595}-\x{1F596}]|[\x{1F5A4}]|[\x{1F5FB}-\x{1F64F}]|[\x{1F680}-\x{1F6C5}]|[\x{1F6CC}]|[\x{1F6D0}]|[\x{1F6D1}-\x{1F6D2}]|[\x{1F6EB}-\x{1F6EC}]|[\x{1F6F4}-\x{1F6F9}]|[\x{1F910}-\x{1F93A}]|[\x{1F93C}-\x{1F93E}]|[\x{1F940}-\x{1F945}]|[\x{1F947}-\x{1F970}]|[\x{1F973}-\x{1F976}]|[\x{1F97A}]|[\x{1F97C}-\x{1F9A2}]|[\x{1F9B0}-\x{1F9B9}]|[\x{1F9C0}-\x{1F9C2}]|[\x{1F9D0}-\x{1F9FF}]|[\x{1FA70}-\x{1FA73}]|[\x{1FA78}-\x{1FA7A}]|[\x{1FA80}-\x{1FA82}]|[\x{1FA90}-\x{1FA95}]`)

func MustEncloseCharAtIndex(s string, index int) (es string) {
	if index < 0 || index >= len(s) {
		log.Fatalf("failed to EncloseCharAtIndex: index out of range")

		return
	}
	// Get the parts of the string before, at, and after the specified index
	before := s[:index]
	char := s[index : index+1]
	after := s[index+1:]
	// Construct the new string
	es = fmt.Sprintf("%s(%s)%s", before, char, after)

	return
}

func RemoveEmojis(s string) string {
	return strings.TrimSpace(emojiRegex.ReplaceAllString(s, ""))
}

func Title(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
