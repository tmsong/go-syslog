package format

import (
	"bufio"
	"bytes"

	"github.com/tmsong/go-syslog/internal/syslogparser/rfc3164"
	"github.com/tmsong/go-syslog/internal/syslogparser/rfc5424"
)

// NUL implements NUL-framing as used by Python's SysLogHandler over TCP.
// Each syslog message is terminated by a NUL byte (0x00).
type NUL struct{}

func (f *NUL) GetParser(line []byte) LogParser {
	switch detect(line) {
	case detectedRFC5424:
		return &parserWrapper{rfc5424.NewParser(line)}
	default:
		return &parserWrapper{rfc3164.NewParser(line)}
	}
}

func (f *NUL) GetSplitFunc() bufio.SplitFunc {
	return nulScannerSplit
}

func nulScannerSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, 0x00); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		// No NUL terminator found; return remaining data as-is
		return len(data), data, nil
	}

	return 0, nil, nil
}
