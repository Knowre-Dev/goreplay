package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

//CookieMap key=value; key2=value2; ... 의 값을 맵으로 바꾸어줌
type CookieMap map[string]string

type stringWriter struct {
	w io.Writer
}

func (w stringWriter) WriteString(s string) (n int, err error) {
	return w.w.Write([]byte(s))
}

func (c CookieMap) Parse(s string) {
	for k := range c {
		delete(c, k)
	}

	for _, pair := range strings.Split(s, ";") {
		z := strings.Split(pair, "=")
		if len(z) > 1 {
			c[z[0]] = z[1]
		}
	}
}

func (c CookieMap) Write(w io.Writer) error {
	return c.writeSubset(w)
}
func (c CookieMap) writeSubset(w io.Writer) error {
	ws, ok := w.(io.StringWriter)
	if !ok {
		ws = stringWriter{w}
	}

	for k, v := range c {
		ws.WriteString(k + "=" + strings.TrimSpace(v) + ";")
	}
	return nil
}

func (c CookieMap) String() string {
	buff := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buff)
	c.Write(writer)
	writer.Flush()

	return buff.String()
}
