package utils

import "os/exec"

func Clrscr() {
	cmd := exec.Command("clear")
	cmd.Run()
}
