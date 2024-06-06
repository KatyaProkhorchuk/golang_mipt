package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// listDirs lists directories in given directory.
func listDirs(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, filepath.Join(dir, f.Name()))
		}
	}

	return dirs, nil
}

func doTestSubmission(t *testing.T, studentRepo, privateRepo, problem string) error {
	annotate := func(prefix string, f **os.File) func() {
		pr, pw, err := os.Pipe()
		require.NoError(t, err)

		oldF := *f
		*f = pw

		go func() {
			s := bufio.NewScanner(pr)
			for s.Scan() {
				_, _ = io.WriteString(oldF, fmt.Sprintf("%s%s\n", prefix, s.Text()))
			}
		}()

		return func() {
			pw.Close()
			*f = oldF
		}
	}

	t.Logf("=== testing started ===")
	defer annotate(">>> STDOUT >>>", &os.Stdout)()
	defer annotate(">>> STDERR >>>", &os.Stderr)()
	defer t.Logf("=== testing finished ===")

	return testSubmission(studentRepo, privateRepo, problem)
}

func Test_testSubmission_correct(t *testing.T) {
	testDirs, err := listDirs("../testdata/submissions/correct")
	require.NoError(t, err)

	for _, dir := range testDirs {
		absDir, err := filepath.Abs(dir)
		require.NoError(t, err)
		problem := path.Base(absDir)
		t.Run(problem, func(t *testing.T) {
			studentRepo := path.Join(absDir, "student")
			privateRepo := path.Join(absDir, "private")

			require.NoError(t, doTestSubmission(t, studentRepo, privateRepo, problem))
		})
	}
}

func Test_testSubmission_incorrect(t *testing.T) {
	testDirs, err := listDirs("../testdata/submissions/incorrect")
	require.NoError(t, err)

	for _, dir := range testDirs {
		absDir, err := filepath.Abs(dir)
		require.NoError(t, err)

		problem := path.Base(absDir)
		t.Run(problem, func(t *testing.T) {
			studentRepo := path.Join(absDir, "student")
			privateRepo := path.Join(absDir, "private")

			err := doTestSubmission(t, studentRepo, privateRepo, problem)
			require.Error(t, err)

			if problem == "brokentest" {
				var testFailedErr *TestFailedError
				require.True(t, errors.As(err, &testFailedErr))
			}
		})
	}
}
