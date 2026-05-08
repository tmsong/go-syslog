package format

import (
	"bufio"
	"bytes"
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *FormatSuite) TestNUL_GetSplitFuncSingleMessage(c *C) {
	f := NUL{}

	buf := bytes.NewReader([]byte("<14>May  8 12:00:00 host app: hello\x00"))
	scanner := bufio.NewScanner(buf)
	scanner.Split(f.GetSplitFunc())

	c.Assert(scanner.Scan(), Equals, true)
	c.Assert(scanner.Text(), Equals, "<14>May  8 12:00:00 host app: hello")
	c.Assert(scanner.Scan(), Equals, false)
}

func (s *FormatSuite) TestNUL_GetSplitFuncMultipleMessages(c *C) {
	f := NUL{}

	messages := []string{
		"<14>May  8 12:00:00 host app: hello",
		"<14>May  8 12:00:01 host app: world",
		"<14>May  8 12:00:02 host app: foo",
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

func (s *FormatSuite) TestNUL_GetSplitFuncNoTerminator(c *C) {
	f := NUL{}

	msg := "<14>May  8 12:00:00 host app: no nul at end"
	buf := bytes.NewReader([]byte(msg))
	scanner := bufio.NewScanner(buf)
	scanner.Split(f.GetSplitFunc())

	c.Assert(scanner.Scan(), Equals, true)
	c.Assert(scanner.Text(), Equals, msg)
}

func (s *FormatSuite) TestNUL_GetParserRFC3164(c *C) {
	f := NUL{}
	line := []byte("<14>May  8 12:00:00 host app: hello")
	parser := f.GetParser(line)
	c.Assert(parser, NotNil)
	c.Assert(parser.Parse(), IsNil)
}
