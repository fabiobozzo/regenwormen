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
