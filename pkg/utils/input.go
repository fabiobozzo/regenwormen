package utils

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func MustReadString(reader *bufio.Reader, message string) string {
	if message != "" {
		fmt.Print(message)
	}

	rs, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("failed to read input:", err)
	}

	return rs
}

func MustReadInt(reader *bufio.Reader, message string) (int, error) {
	s := MustReadString(reader, message)

	i, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, err
	}

	return i, nil
}

func MustReadBool(reader *bufio.Reader, message, trueAnswer, falseAnswer string) (answer, hasAnswer bool) {
	s := strings.TrimSpace(MustReadString(reader, fmt.Sprintf("%s (%s/%s) ", message, trueAnswer, falseAnswer)))

	if len(s) == 0 {
		return
	}

	if strings.Contains(strings.ToLower(trueAnswer), strings.ToLower(s)) {
		answer, hasAnswer = true, true
	} else if strings.Contains(strings.ToLower(falseAnswer), strings.ToLower(s)) {
		answer, hasAnswer = false, true
	}

	return
}
