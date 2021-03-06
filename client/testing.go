package client

import (
	"errors"

	"github.com/compose/transporter/message"
	"github.com/compose/transporter/message/ops"
)

// Mock can be used for mocking tests that need no actual client or Session.
type Mock struct {
	Closed bool
}

// Connect satisfies the Client interface.
func (c *Mock) Connect() (Session, error) {
	return &MockSession{}, nil
}

// Close satisfies the Closer interface.
func (c *Mock) Close() { c.Closed = true }

// MockErr can be used for mocking tests that need no actual client or Session.
type MockErr struct {
}

// Connect satisfies the Client interface.
func (c *MockErr) Connect() (Session, error) {
	return nil, errors.New("connect failed")
}

// MockSession can be used for mocking tests the do not need to use anything in the Session.
type MockSession struct {
}

// Close satisfies the Closer interface.
func (s *MockSession) Close() {}

// MockReader can be used to send a random number of messages
type MockReader struct {
	MsgCount int
}

func (r *MockReader) Read(filterFn NsFilterFunc) MessageChanFunc {
	return func(s Session, done chan struct{}) (chan message.Msg, error) {
		out := make(chan message.Msg)
		go func() {
			defer close(out)
			for i := 0; i < r.MsgCount; i++ {
				out <- message.From(ops.Insert, "test", map[string]interface{}{"id": i})
			}
		}()

		return out, nil
	}
}

// MockWriter can be used to count the number of messages sent to Write.
type MockWriter struct {
	MsgCount int
}

// Writer satisfies the Writer interface.
func (w *MockWriter) Write(msg message.Msg) func(Session) (message.Msg, error) {
	return func(s Session) (message.Msg, error) {
		w.MsgCount++
		return msg, nil
	}
}
