package golangspectester

import (
	"flag"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var runCmd = flag.String("run-cmd", "go run {path}", "command to run and check all results with spec files")

func TestEngine(t *testing.T) {
	t.Log("using run command:", *runCmd)
	e := NewEngine("./_test", func(path, expected string, isError bool, _ io.Reader) bool {
		return t.Run(path, func(t *testing.T) {
			require := require.New(t)
			t.Parallel()
			if isError {
				t.Skip("error check not supported yet")
			}

			execCmd := strings.ReplaceAll(*runCmd, "{path}", path)
			sc := strings.Split(execCmd, " ")
			if len(sc) < 1 {
				t.Log("wrong run command:", execCmd)
				t.FailNow()
			}

			out, err := exec.Command(sc[0], sc[1:]...).CombinedOutput()
			require.NoError(err, string(out))
			require.Equal(expected, strings.Trim(string(out), " "))
		})
	})

	require.NoError(t, e.Start())
}
