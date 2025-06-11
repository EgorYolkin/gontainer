package writer

import (
	"gontainer/internal/kgsp/spvalue"
)

func (w *Writer) Write(v spvalue.Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
