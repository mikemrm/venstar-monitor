package jsonPrinter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mikemrm/venstar-monitor"
)

// JSONWriter implements monitor.ResultsWriter
type JSONWriter struct {
	writer io.Writer
}

func (w *JSONWriter) WriteResults(results *monitor.Results) error {
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}
	fmt.Fprintln(w.writer, string(data))
	return nil
}

func NewWriter(writer ...io.Writer) (*JSONWriter, error) {
	var w io.Writer = os.Stdout
	if len(writer) != 0 {
		w = writer[0]
	}
	return &JSONWriter{
		writer: w,
	}, nil
}
