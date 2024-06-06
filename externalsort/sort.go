//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type lineReader struct {
	reader *bufio.Reader
}

func (lr *lineReader) ReadLine() (string, error) {
	line, err := lr.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	line = strings.TrimRight(line, "\n")
	return line, err
}

type lineWriter struct {
	writer *bufio.Writer
}

func (wl *lineWriter) Write(l string) error {
	_, err := wl.writer.WriteString(l + "\n")
	if err != nil {
		return err
	}
	return wl.writer.Flush()
}

func NewReader(r io.Reader) LineReader {
	return &lineReader{reader: bufio.NewReader(r)}
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{writer: bufio.NewWriter(w)}
}

type lineItem struct {
	line   string
	reader LineReader
}

type minHeap []*lineItem

func (h minHeap) Len() int {
	return len(h)
}

func (h minHeap) Less(i, j int) bool {
	return strings.Compare(h[i].line, h[j].line) < 0
}

func (h minHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*lineItem))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
func Merge(w LineWriter, readers ...LineReader) error {
	minHeap := &minHeap{}
	for _, reader := range readers {
		line, err := reader.ReadLine()
		if err != nil && err != io.EOF || (err == io.EOF && line == "") {
			continue
		}
		heap.Push(minHeap, &lineItem{line, reader})
	}
	heap.Init(minHeap)
	for minHeap.Len() > 0 {
		item := heap.Pop(minHeap).(*lineItem)
		if err := w.Write(item.line); err != nil {
			return err
		}
		line, err := item.reader.ReadLine()
		if err == nil || (err == io.EOF && line != "") {
			heap.Push(minHeap, &lineItem{line, item.reader})
			heap.Fix(minHeap, minHeap.Len()-1)
		}
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	var readers []LineReader

	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			f.Close()
			return err
		}
		lr := NewReader(f)
		var lines []string

		for {
			line, errFile := lr.ReadLine()

			if line != "" && errFile == io.EOF {
				lines = append(lines, line)
			}
			if errFile != nil {
				break
			}
			lines = append(lines, line)
		}
		sort.Strings(lines)
		f.Close()
		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		lw := NewWriter(f)
		for _, line := range lines {
			err := lw.Write(line)
			if err != nil {
				return err
			}
		}
		f.Close()

	}
	lw := NewWriter(w)
	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			f.Close()
			return err
		}
		lr := NewReader(f)
		readers = append(readers, lr)
	}
	return Merge(lw, readers...)
}
