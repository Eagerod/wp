package main;

import (
    "os/exec"
    "testing"
)

import (
    "github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
    _, err := exec.LookPath("wp")
    assert.NoError(t, err)

	cmd := exec.Command("wp", "--help")

	_, err = cmd.CombinedOutput()
    assert.NoError(t, err)
}
