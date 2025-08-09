package rainbowlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rambollwong/rainbowlog/level"
)

// LevelWriter is an interface that wraps the standard io.Writer interface
// and adds a WriteLevel method to allow writing data with log level information.
type LevelWriter interface {
	io.Writer
	WriteLevel(level level.Level, bz []byte) (n int, err error)
}

// levelWriterAdapter is a wrapper that adapts a standard io.Writer
// to implement the LevelWriter interface by ignoring the level information.
type levelWriterAdapter struct {
	io.Writer
}

// WriteLevel implements the LevelWriter interface by calling the underlying
// Write method and ignoring the level parameter.
func (lw levelWriterAdapter) WriteLevel(_ level.Level, bz []byte) (n int, err error) {
	return lw.Write(bz)
}

// LevelWriterAdapter converts an io.Writer to a LevelWriter.
// If the writer already implements LevelWriter, it is returned as-is.
// Otherwise, it is wrapped in a levelWriterAdapter.
func LevelWriterAdapter(writer io.Writer) LevelWriter {
	if lw, ok := writer.(LevelWriter); ok {
		return lw
	}
	return &levelWriterAdapter{writer}
}

// syncWriter is a thread-safe wrapper for LevelWriter.
// It uses a mutex to synchronize access to the underlying writer.
type syncWriter struct {
	mu          sync.Mutex
	levelWriter LevelWriter
}

// SyncWriter wraps a writer to make it thread-safe.
// Each call to Write is synchronized with a mutex.
// This wrapper is useful for writers that are not thread-safe.
// Note: os.File operations on POSIX and Windows systems are already thread-safe
// and do not require this wrapper.
func SyncWriter(w io.Writer) io.Writer {
	if lw, ok := w.(LevelWriter); ok {
		return &syncWriter{levelWriter: lw}
	}
	return &syncWriter{levelWriter: levelWriterAdapter{w}}
}

// Write implements the io.Writer interface in a thread-safe manner.
func (s *syncWriter) Write(bz []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.levelWriter.Write(bz)
}

// WriteLevel implements the LevelWriter interface in a thread-safe manner.
func (s *syncWriter) WriteLevel(l level.Level, bz []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.levelWriter.WriteLevel(l, bz)
}

// multiLevelWriter writes to multiple LevelWriters simultaneously.
type multiLevelWriter struct {
	writers []LevelWriter
}

// Write implements the io.Writer interface by writing to all underlying writers.
// It returns an error if any of the writers fails or performs a short write.
func (t multiLevelWriter) Write(bz []byte) (n int, err error) {
	for _, writer := range t.writers {
		n, err = writer.Write(bz)
		if err != nil {
			return
		}
		if n != len(bz) {
			err = io.ErrShortWrite
			return
		}
	}
	return
}

// WriteLevel implements the LevelWriter interface by writing to all underlying writers.
// It returns an error if any of the writers fails or performs a short write.
func (t multiLevelWriter) WriteLevel(lv level.Level, bz []byte) (n int, err error) {
	for _, writer := range t.writers {
		n, err = writer.WriteLevel(lv, bz)
		if err != nil {
			return
		}
		if n != len(bz) {
			err = io.ErrShortWrite
			return
		}
	}
	return
}

// MultiLevelWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
// Each provided writer is wrapped with LevelWriterAdapter to ensure
// it implements the LevelWriter interface.
func MultiLevelWriter(writers ...io.Writer) LevelWriter {
	levelWriters := make([]LevelWriter, 0, len(writers))
	for _, w := range writers {
		levelWriters = append(levelWriters, LevelWriterAdapter(w))
	}
	return multiLevelWriter{levelWriters}
}

// BufferedWriter is a buffered writer that uses a double buffering mechanism to improve write performance.
// It automatically switches to another buffer when the current buffer is full and asynchronously flushes
// the filled buffer data to the underlying writer.
type BufferedWriter struct {
	mu         sync.Mutex     // Buffer switch lock to ensure atomicity of buffer switching
	wMu        sync.Mutex     // Underlying writer lock to ensure exclusive access to the underlying writer
	wg         sync.WaitGroup // Used to wait for all asynchronous flush operations to complete
	w          io.Writer      // Underlying writer, the actual data write target
	bufSize    int            // Defined buffer size
	bufA       *bytes.Buffer  // Buffer A, one of the double buffers
	bufB       *bytes.Buffer  // Buffer B, one of the double buffers
	bufCurrent *bytes.Buffer  // Pointer to the currently used buffer
	closed     bool           // Flag indicating whether the writer is closed
}

// NewBufferedWriter creates a new buffered writer.
// The parameter w is the underlying io.Writer, and bufSize is the size of each buffer.
// The return value is a BufferedWriter instance that implements the io.Writer interface.
func NewBufferedWriter(w io.Writer, bufSize int) *BufferedWriter {
	bw := &BufferedWriter{
		w:       w,
		bufSize: bufSize,
		bufA:    bytes.NewBuffer(make([]byte, 0, bufSize)), // Initialize buffer A
		bufB:    bytes.NewBuffer(make([]byte, 0, bufSize)), // Initialize buffer B
		closed:  false,
	}
	bw.bufCurrent = bw.bufA // Initially use buffer A
	return bw
}

// BufferedLevelWriter creates a new LevelWriter that wraps the provided io.Writer with buffering capabilities.
// It uses the BufferedWriter to buffer writes and flushes them asynchronously for improved performance.
// The bufSize parameter specifies the size of each internal buffer.
func BufferedLevelWriter(w io.Writer, bufSize int) LevelWriter {
	return LevelWriterAdapter(NewBufferedWriter(w, bufSize))
}

// Write implements the io.Writer interface, writing data to the buffer.
// When the buffer is full, it automatically switches to another buffer and asynchronously flushes the filled buffer.
func (bw *BufferedWriter) Write(bz []byte) (n int, err error) {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	if bw.closed {
		return 0, os.ErrClosed
	}

	writeSize := len(bz)
	written := 0
	for written < writeSize {
		remaining := writeSize - written
		available := bw.bufSize - bw.bufCurrent.Len()
		if remaining > available {
			// Current buffer cannot accommodate all remaining data, first fill the current buffer
			n, err := bw.bufCurrent.Write(bz[written : written+available])
			if err != nil {
				return written + n, err
			}
			written += available

			// Switch buffer and asynchronously write the filled buffer
			if err = bw.swapAndFlush(); err != nil {
				return written, err
			}
		} else {
			// Remaining data can be completely written to the current buffer
			n, err := bw.bufCurrent.Write(bz[written:])
			written += n
			if err != nil {
				return written, err
			}

			// If the current buffer is full, switch buffer and asynchronously write
			if bw.bufCurrent.Len() >= bw.bufSize {
				if err = bw.swapAndFlush(); err != nil {
					return written, err
				}
			}
			break
		}
	}

	return written, nil
}

// Flush forcibly flushes the data in the current buffer to the underlying writer.
func (bw *BufferedWriter) Flush() error {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return bw.flushCurrent()
}

// Close closes the writer, ensuring all buffer data is flushed.
// This method waits for all asynchronous flush operations to complete before returning.
func (bw *BufferedWriter) Close() error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	if bw.closed {
		return nil
	}

	// Flush the data in the current buffer
	if err := bw.flushCurrent(); err != nil {
		return err
	}

	// Wait for all asynchronous flush operations to complete
	bw.wg.Wait()

	bw.closed = true
	return nil
}

// swapAndFlush switches the currently used buffer and starts a goroutine to asynchronously flush the old buffer's data.
func (bw *BufferedWriter) swapAndFlush() error {
	// Wait for previous asynchronous flush operations to complete to ensure no concurrent access to the same buffer
	bw.wg.Wait()

	// Determine the buffer to flush and switch the current buffer
	toWrite := bw.bufCurrent
	if bw.bufCurrent == bw.bufA {
		bw.bufCurrent = bw.bufB
	} else {
		bw.bufCurrent = bw.bufA
	}

	// Start the asynchronous flush goroutine
	bw.wg.Add(1)
	go func(toWrite *bytes.Buffer) {
		defer bw.wg.Done()
		e := bw.flushBuffer(toWrite)
		if e != nil {
			// Use panic instead of returning an error because this is running in a goroutine.
			// Consider using a more graceful error handling mechanism, such as logging the error
			// and possibly notifying the user or taking corrective action.
			// For now, we will log the error and continue.
			fmt.Fprintf(os.Stderr, "failed to flush buffer: %v\n", e)
		}
	}(toWrite)

	return nil
}

// flushCurrent flushes the data in the current buffer.
// If the current buffer is empty, no operation is performed.
func (bw *BufferedWriter) flushCurrent() error {
	bw.wg.Wait()

	if bw.bufCurrent.Len() == 0 {
		return nil
	}
	return bw.flushBuffer(bw.bufCurrent)
}

// flushBuffer writes the data in the specified buffer to the underlying writer and resets the buffer.
func (bw *BufferedWriter) flushBuffer(buf *bytes.Buffer) error {
	// If the buffer is empty, no flush is needed
	if buf.Len() == 0 {
		return nil
	}

	// Use wMu to ensure exclusive access to the underlying writer
	bw.wMu.Lock()
	defer bw.wMu.Unlock()

	// Write data to the underlying writer
	_, err := bw.w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// Reset the buffer for next use
	buf.Reset()
	return nil
}
