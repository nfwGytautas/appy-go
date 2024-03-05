package driver_websocket

var (
	cNewline = []byte{'\n'}
	cSpace   = []byte{' '}
)

// Timeout for client to send ready in seconds
const cClientReadyTimeout = 10
