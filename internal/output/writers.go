package output

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

var (
	writersMu sync.RWMutex
	stdout    io.Writer = os.Stdout
	stderr    io.Writer = os.Stderr
)

func Stdout() io.Writer {
	writersMu.RLock()
	defer writersMu.RUnlock()
	return stdout
}

func Stderr() io.Writer {
	writersMu.RLock()
	defer writersMu.RUnlock()
	return stderr
}

func SetWritersForTest(out, err io.Writer) func() {
	writersMu.Lock()
	prevOut := stdout
	prevErr := stderr
	if out != nil {
		stdout = out
	}
	if err != nil {
		stderr = err
	}
	writersMu.Unlock()
	return func() {
		writersMu.Lock()
		stdout = prevOut
		stderr = prevErr
		writersMu.Unlock()
	}
}

type NDJSONWriter struct {
	enc *json.Encoder
}

func NewNDJSONWriter(w io.Writer) *NDJSONWriter {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &NDJSONWriter{enc: enc}
}

func (n *NDJSONWriter) WriteItem(item any) error {
	return n.enc.Encode(item)
}

func (n *NDJSONWriter) WritePagination(p *Pagination) error {
	return n.enc.Encode(map[string]any{"@pagination": p})
}
