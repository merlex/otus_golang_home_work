package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"ENV1": EnvValue{"val1", false},
		"ENV2": EnvValue{"", true},
		"ENV3": EnvValue{"", false},
		"ENV4": EnvValue{"abcd\x00dcba", false},
	}
	retCode := RunCmd([]string{"./testdata/echo.sh", "arg1=1", "arg2=2"}, env)
	require.Equal(t, retCode, 0)

	code := RunCmd([]string{"bash", "-c", "exit 123"}, nil)
	require.Equal(t, code, 123)

	retCode = RunCmd([]string{""}, env)
	require.Equal(t, retCode, 1)

	retCode = RunCmd(nil, env)
	require.Equal(t, retCode, 1)
}
