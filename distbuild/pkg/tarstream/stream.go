package tarstream

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// Send рекурсивно обходит директорию и сериализует её содержимое в поток w.
func Send(dir string, w io.Writer) error {
	tw := tar.NewWriter(w)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		switch {
		case info.IsDir():
			return tw.WriteHeader(&tar.Header{
				Name:     rel,
				Typeflag: tar.TypeDir,
			})

		default:
			h := &tar.Header{
				Typeflag: tar.TypeReg,
				Name:     rel,
				Size:     info.Size(),
				Mode:     int64(info.Mode()),
			}

			if err := tw.WriteHeader(h); err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(tw, f)
			return err
		}
	})

	if err != nil {
		return err
	}

	return tw.Close()
}

// Receive читает поток r и материализует содержимое потока внутри dir.
func Receive(dir string, r io.Reader) error {
	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		absPath := filepath.Join(dir, h.Name)

		if h.Typeflag == tar.TypeDir {
			if err := os.Mkdir(absPath, 0777); err != nil {
				return err
			}
		} else {
			writeFile := func() error {
				f, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY, os.FileMode(h.Mode))
				if err != nil {
					return err
				}
				defer f.Close()

				_, err = io.Copy(f, tr)
				return err
			}

			if err := writeFile(); err != nil {
				return err
			}
		}
	}
}
