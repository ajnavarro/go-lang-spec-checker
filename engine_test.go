package golangspectester

import (
	"flag"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var runCmd = flag.String("run-cmd", "go run {path}", "command to run and check all results with spec files")
var unsupported = flag.String("unsupported", "", "list of unsupported functionalities. Tests based on these features will be skipped.")

func TestEngine(t *testing.T) {
	t.Log("using run command:", *runCmd)
	t.Log("unsupported features:", *unsupported)

	ufs := strings.Split(*unsupported, ",")

	e := NewEngine("./_test", func(path, expected string, isError bool, _ io.Reader) bool {
		dir, file := filepath.Split(path)
		base := filepath.Base(dir)

		t.Run(file, func(t *testing.T) {
			for _, uf := range ufs {
				if strings.Trim(uf, " ") == base {
					t.Skipf("unsupported feature: %q", base)
				}
			}

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

		return true
	})

	require.NoError(t, e.Start())
}
