package main;

import (
	"fmt"
    "os/exec"
    "testing"
)

import (
    "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
    if _, err := exec.LookPath("wp"); err != nil {
        panic("Failed to find executable to run system tests")
    }

	cmd := exec.Command("wp")

	str, err := cmd.CombinedOutput()
    if err != nil {
    	fmt.Println(err)
        panic("Failed to run bare command")
    }

    assert.Equal(t, string(str), "I did a thing\n")
}
