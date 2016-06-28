package nimbus

import (
	"bufio"
	"io"
	"sync"
)

type IRCReader struct {
	reader *bufio.Reader
	mu sync.Mutex
}

func NewIRCReader(reader io.Reader) *IRCReader {
	return &IRCReader {
		reader: bufio.NewReader(reader),
	}
}

func (r *IRCReader) ReadRaw() (string, error) {
	r.mu.Lock()
	line, err := r.reader.ReadString('\n')
	r.mu.Unlock()

	if err != nil {
		return "", err
	}
	
	return line, nil
}

func (r *IRCReader) Read() (*Message, error) {
	raw, err := r.ReadRaw()
	if err != nil {
		return nil, err
	}

	return ParseMessage(raw)
}

type IRCWriter struct {
	writer io.Writer
	mu sync.Mutex
}

func NewIRCWriter(writer io.Writer) *IRCWriter {
	return &IRCWriter {
		writer: writer,
	}
}

func (w *IRCWriter) Write(packet []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.writer.Write(packet)
	if err != nil {
		return err
	}

	w.writer.Write([]byte(ENDLINE))

	return nil
}
