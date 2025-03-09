package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

var clearFuncs map[string]func()

func initClearScreen() {
	clearFuncs = make(map[string]func())
	clearFuncs["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	clearFuncs["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
	clearFuncs["darwin"] = clearFuncs["linux"]
}

func clearScreen() {
	if value, ok := clearFuncs[runtime.GOOS]; ok {
		value()
	} else {
		log.Println("[WARNING] Cannot clearFuncs terminal screen on unsupported platform: ", runtime.GOOS)
	}
}
