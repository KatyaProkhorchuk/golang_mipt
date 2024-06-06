//go:build !solution

package fileleak

import (
	"os"
	"io/fs"
)
type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}
func ReadMap(nameDir string, data map[string]int, dir []fs.DirEntry) {
	for _, file := range dir {
		key, _ := os.Readlink(nameDir + "/" + file.Name())
		if _, ok := data[key]; ok {
			data[key]++
		} else {
			data[key] = 0
		}
	}
}
func VerifyNone(t testingT) {
	nameDir := "/proc/self/fd"
	dir, err := os.ReadDir(nameDir)
	if err != nil {
		return
	}
	beforeRun := make(map[string]int)
	ReadMap(nameDir, beforeRun, dir)

	t.Cleanup(func() {
		dir, err := os.ReadDir(nameDir)
		if err != nil {
			t.Errorf("don't read directory %s", nameDir)
		}
		afterRun := make(map[string]int)
		ReadMap(nameDir, afterRun, dir)
		if len(afterRun) > len(beforeRun) {
			t.Errorf("delect leak %s", nameDir)

		}
		for newKey, newVal := range afterRun {
			if oldVal, ok := beforeRun[newKey]; ok {
				if oldVal < newVal {
					t.Errorf("delect leak %s", newKey)
				}

			} else {
				t.Errorf("delect leak %s", newKey)
			}
		}
	})
	

}