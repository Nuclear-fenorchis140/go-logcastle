package writer

import (
	"io"
	"sync"
	"time"
)

// BufferedWriter buffers log entries and writes them in batches for performance
type BufferedWriter struct {
	output        io.Writer
	buffer        [][]byte
	bufferSize    int
	flushInterval time.Duration
	ticker        *time.Ticker
	done          chan struct{}
	mu            sync.Mutex
	wg            sync.WaitGroup
}

// NewBufferedWriter creates a new buffered writer
func NewBufferedWriter(output io.Writer, bufferSize int, flushInterval time.Duration) *BufferedWriter {
	return &BufferedWriter{
		output:        output,
		buffer:        make([][]byte, 0, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		ticker:        time.NewTicker(flushInterval),
		done:          make(chan struct{}),
	}
}

// Start begins the auto-flush goroutine
func (w *BufferedWriter) Start() {
	w.wg.Add(1)
	go w.autoFlush()
}

// Write adds data to the buffer
func (w *BufferedWriter) Write(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Copy data to avoid mutation
	buf := make([]byte, len(data))
	copy(buf, data)
	w.buffer = append(w.buffer, buf)

	// Flush if buffer is full
	if len(w.buffer) >= w.bufferSize {
		return w.flush()
	}

	return nil
}

// Flush writes all buffered data
func (w *BufferedWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flush()
}

func (w *BufferedWriter) flush() error {
	if len(w.buffer) == 0 {
		return nil
	}

	// Write all buffered entries
	var firstErr error
	for _, entry := range w.buffer {
		if _, err := w.output.Write(entry); err != nil && firstErr == nil {
			firstErr = err
			// Continue trying to write remaining entries
		}
	}

	// Clear buffer even if there were errors
	w.buffer = w.buffer[:0]
	return firstErr
}

func (w *BufferedWriter) autoFlush() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ticker.C:
			if err := w.Flush(); err != nil {
				// Best effort - errors are logged by caller
			}
		case <-w.done:
			return
		}
	}
}

// Stop stops the auto-flush and flushes remaining data
func (w *BufferedWriter) Stop() error {
	// Stop ticker
	w.ticker.Stop()

	// Signal done
	select {
	case <-w.done:
		// Already closed
	default:
		close(w.done)
	}

	// Wait for auto-flush goroutine
	w.wg.Wait()

	// Final flush
	return w.Flush()
}
