package logcastle

import (
	"io"
	"sync"
	"time"
)

// BufferedWriter batches log writes for improved performance.
// Automatically flushes based on buffer size and time interval.
type BufferedWriter struct {
	output        io.Writer
	buffer        [][]byte
	bufferSize    int
	flushInterval time.Duration

	mu     sync.Mutex
	ticker *time.Ticker
	done   chan struct{}
	wg     sync.WaitGroup
}

// NewBufferedWriter creates a buffered writer with specified capacity and flush interval.
// bufferSize: max entries before auto-flush. flushInterval: time between flushes.
func NewBufferedWriter(output io.Writer, bufferSize int, flushInterval time.Duration) *BufferedWriter {
	return &BufferedWriter{
		output:        output,
		buffer:        make([][]byte, 0, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		done:          make(chan struct{}),
	}
}

// Start begins the auto-flush goroutine for periodic flushing.
// Called automatically by Init. Don't call manually.
func (w *BufferedWriter) Start() {
	w.ticker = time.NewTicker(w.flushInterval)
	w.wg.Add(1)
	go w.autoFlush()
}

// Stop stops the writer and flushes remaining data to output.
// Idempotent - safe to call multiple times. Returns error if flush fails.
func (w *BufferedWriter) Stop() error {
	w.mu.Lock()
	// Check if already stopped
	select {
	case <-w.done:
		// Already stopped
		w.mu.Unlock()
		return nil
	default:
		close(w.done)
	}
	w.mu.Unlock()

	w.ticker.Stop()
	w.wg.Wait()
	return w.Flush()
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

	var firstErr error
	// Write all buffered entries
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
