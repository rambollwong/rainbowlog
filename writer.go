package rainbowlog

import (
	"github.com/RambollWong/rainbowlog/level"
	"io"
	"sync"
)

// LevelWriter defines as interface a writer may implement in order
// to receive level information with payload.
type LevelWriter interface {
	io.Writer
	WriteLevel(level level.Level, bz []byte) (n int, err error)
}

type levelWriterAdapter struct {
	io.Writer
}

func (lw levelWriterAdapter) WriteLevel(l level.Level, bz []byte) (n int, err error) {
	return lw.Write(bz)
}

func LevelWriterAdapter(writer io.Writer) LevelWriter {
	if lw, ok := writer.(LevelWriter); ok {
		return lw
	}
	return &levelWriterAdapter{writer}
}

type syncWriter struct {
	mu          sync.Mutex
	levelWriter LevelWriter
}

// SyncWriter wraps writer so that each call to Write is synchronized with a mutex.
// This syncer can be used to wrap the call to writer's Write method if it is
// not thread safe. Note that you do not need this wrapper for os.File Write
// operations on POSIX and Windows systems as they are already thread-safe.
func SyncWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LevelWriter); ok {
		return &syncWriter{levelWriter: lw}
	}
	return &syncWriter{levelWriter: levelWriterAdapter{w}}
}

// Write implements the io.Writer interface.
func (s *syncWriter) Write(bz []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.levelWriter.Write(bz)
}

// WriteLevel implements the LevelWriter interface.
func (s *syncWriter) WriteLevel(l level.Level, bz []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.levelWriter.WriteLevel(l, bz)
}

type multiLevelWriter struct {
	writers []LevelWriter
}

func (t multiLevelWriter) Write(bz []byte) (n int, err error) {
	for _, writer := range t.writers {
		if _n, _err := writer.Write(bz); err == nil {
			n = _n
			if _err != nil {
				err = _err
				return
			}
			if _n != len(bz) {
				err = io.ErrShortWrite
				return
			}
		}
	}
	return
}

func (t multiLevelWriter) WriteLevel(lv level.Level, bz []byte) (n int, err error) {
	for _, writer := range t.writers {
		if _n, _err := writer.WriteLevel(lv, bz); err == nil {
			n = _n
			if _err != nil {
				err = _err
				return
			}
			if _n != len(bz) {
				err = io.ErrShortWrite
				return
			}
		}
	}
	return
}

// MultiLevelWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command. If some writers
// implement LevelWriter, their WriteLevel method will be used instead of Write.
func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	levelWriters := make([]LevelWriter, 0, len(writers))
	for _, w := range writers {
		levelWriters = append(levelWriters, LevelWriterAdapter(w))
	}
	return multiLevelWriter{levelWriters}
}
