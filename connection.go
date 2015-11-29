package nimbus

import (
	"bufio"
	"io"
	"sync"
)

type IRCReader struct {
	Reader *bufio.Reader
	Mu sync.Mutex
}

func NewIRCReader(reader io.Reader) *IRCReader {
	return &IRCReader {
		Reader: bufio.NewReader(reader),
	}
}

func (r *IRCReader) ReadRaw() (string, error) {
	r.Mu.Lock()
	line, err := r.Reader.ReadString('\n')
	r.Mu.Unlock()

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
	Writer io.Writer
	Mu sync.Mutex
}

func NewIRCWriter(writer io.Writer) *IRCWriter {
	return &IRCWriter {
		Writer: writer,
	}
}

func (w *IRCWriter) Write(packet []byte) error {
	w.Mu.Lock()

	_, err := w.Writer.Write(packet)
	if err != nil {
		w.Mu.Unlock()
		return err
	}

	w.Mu.Unlock()
	return nil
}