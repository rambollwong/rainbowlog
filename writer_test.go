package rainbowlog

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/rambollwong/rainbowlog/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferedWriter(t *testing.T) {
	t.Run("BasicWrite", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bw := NewBufferedWriter(buf, 64)

		n, err := bw.Write([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, 5, bw.bufCurrent.Len())

		// Data should not be flushed yet
		assert.Equal(t, "", buf.String())

		// Force flush
		err = bw.Flush()
		assert.NoError(t, err)
		assert.Equal(t, "hello", buf.String())
	})

	t.Run("AutoFlushOnBufferFull", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bw := NewBufferedWriter(buf, 5)

		// Write data that exactly fills the buffer
		n, err := bw.Write([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)

		// Data should be automatically flushed
		assert.Equal(t, "hello", buf.String())
	})

	t.Run("WriteAcrossBufferBoundary", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bw := NewBufferedWriter(buf, 4)

		// Write data larger than buffer size
		n, err := bw.Write([]byte("hello world"))
		assert.NoError(t, err)
		assert.Equal(t, 11, n)

		// Data should be flushed
		err = bw.Flush()
		assert.NoError(t, err)
		assert.Equal(t, "hello world", buf.String())
	})

	t.Run("CloseFlushesBuffer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bw := NewBufferedWriter(buf, 64)

		_, err := bw.Write([]byte("hello"))
		assert.NoError(t, err)

		// Data should not be flushed yet
		assert.Equal(t, "", buf.String())

		// Close should flush
		err = bw.Close()
		assert.NoError(t, err)
		assert.Equal(t, "hello", buf.String())

		// Subsequent writes should fail
		_, err = bw.Write([]byte("world"))
		assert.Error(t, err)
	})

	t.Run("ConcurrentWrites", func(t *testing.T) {
		buf := &bytes.Buffer{}
		bw := NewBufferedWriter(buf, 256)

		// Perform concurrent writes
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := bw.Write([]byte("message"))
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		// Flush and check result
		err := bw.Flush()
		assert.NoError(t, err)

		result := buf.String()
		assert.Len(t, result, 70) // 10 * "message" = 70 characters
		assert.Contains(t, result, "message")
	})
}

func TestBufferedLevelWriter(t *testing.T) {
	t.Run("WriteLevel", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := BufferedLevelWriter(buf, 64)

		n, err := writer.WriteLevel(level.Info, []byte("info message"))
		assert.NoError(t, err)
		assert.Equal(t, 12, n)

		n, err = writer.WriteLevel(level.Error, []byte("error message"))
		assert.NoError(t, err)
		assert.Equal(t, 13, n)

		// For BufferedLevelWriter, we can't directly access the underlying BufferedWriter
		// so we'll just flush through the interface
		flusher, ok := writer.(interface{ Flush() error })
		if ok {
			err = flusher.Flush()
			assert.NoError(t, err)
		}

		// Check that data was written
		result := buf.String()
		assert.Contains(t, result, "info message")
		assert.Contains(t, result, "error message")
	})
}

func TestLevelWriterAdapter(t *testing.T) {
	t.Run("AdaptStandardWriter", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := LevelWriterAdapter(buf)

		// Should be able to write with level
		n, err := writer.WriteLevel(level.Debug, []byte("debug message"))
		assert.NoError(t, err)
		assert.Equal(t, 13, n)

		// Check that data was written (level is ignored)
		assert.Equal(t, "debug message", buf.String())
	})

	t.Run("PassthroughExistingLevelWriter", func(t *testing.T) {
		buf := &bytes.Buffer{}
		original := LevelWriterAdapter(buf)
		adapted := LevelWriterAdapter(original)

		// Should return the same instance
		assert.Equal(t, original, adapted)
	})
}

func TestSyncWriter(t *testing.T) {
	t.Run("ThreadSafeWrites", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := SyncWriter(buf)

		// Perform concurrent writes
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := writer.Write([]byte("test"))
				assert.NoError(t, err)
			}()
		}
		wg.Wait()

		// Check result (order may vary)
		result := buf.String()
		assert.Len(t, result, 40) // 10 * "test" = 40 characters
	})

	t.Run("ThreadSafeLevelWrites", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writer := SyncWriter(buf)
		levelWriter, ok := writer.(LevelWriter)
		require.True(t, ok)

		// Perform concurrent level writes
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				_, err := levelWriter.WriteLevel(level.Info, []byte("test"))
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		// Check result (order may vary)
		result := buf.String()
		assert.Len(t, result, 40) // 10 * "test" = 40 characters
	})
}

func TestMultiLevelWriter(t *testing.T) {
	t.Run("WriteToMultipleWriters", func(t *testing.T) {
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		writer := MultiLevelWriter(buf1, buf2)

		n, err := writer.Write([]byte("hello"))
		assert.NoError(t, err)
		assert.Equal(t, 5, n)

		assert.Equal(t, "hello", buf1.String())
		assert.Equal(t, "hello", buf2.String())
	})

	t.Run("WriteLevelToMultipleWriters", func(t *testing.T) {
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		writer := MultiLevelWriter(buf1, buf2)

		n, err := writer.WriteLevel(level.Warn, []byte("warning"))
		assert.NoError(t, err)
		assert.Equal(t, 7, n)

		assert.Equal(t, "warning", buf1.String())
		assert.Equal(t, "warning", buf2.String())
	})

	t.Run("HandleWriterErrors", func(t *testing.T) {
		buf := &bytes.Buffer{}
		errorWriter := &errorWriter{}

		writer := MultiLevelWriter(buf, errorWriter)

		// First write should succeed (to buf)
		n, err := writer.Write([]byte("hello"))
		assert.Error(t, err)
		assert.Equal(t, 0, n) // ErrorWriter returns 0, causing ErrShortWrite

		// Both writers should have received the data (before the error)
		assert.Equal(t, "hello", buf.String())
	})
}

// errorWriter is a test helper that always returns an error on write
type errorWriter struct{}

func (w *errorWriter) Write(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
