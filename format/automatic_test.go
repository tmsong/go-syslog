package format

import (
	"bufio"
	"bytes"
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *FormatSuite) TestAutomatic_NULFramingDetected(c *C) {
	f := Automatic{}

	messages := []string{
		"<14>May  8 12:00:00 host app: hello",
		"<14>May  8 12:00:01 host app: world",
	}
	buf := new(bytes.Buffer)
	for _, m := range messages {
		fmt.Fprintf(buf, "%s\x00", m)
	}

	scanner := bufio.NewScanner(buf)
	scanner.Split(f.GetSplitFunc())

	i := 0
	for scanner.Scan() {
		c.Assert(scanner.Text(), Equals, messages[i])
		i++
	}
	c.Assert(i, Equals, len(messages))
}

func (s *FormatSuite) TestAutomatic_NULFramingRFC5424(c *C) {
	f := Automatic{}

	msg := "<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut=\"3\"] An application event log entry"
	buf := bytes.NewReader(append([]byte(msg), 0x00))

	scanner := bufio.NewScanner(buf)
	scanner.Split(f.GetSplitFunc())

	c.Assert(scanner.Scan(), Equals, true)
	c.Assert(scanner.Text(), Equals, msg)
	c.Assert(scanner.Scan(), Equals, false)
}
