package markdown

import (
	"bufio"
	"bytes"

	"go.octolab.org/unsafe"
)

const (
	start = "---"
	end   = start
)

const (
	_ = iota
	started
	finished
)

type Splitter struct{}

func (s *Splitter) Split(raw []byte) (props, content []byte, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Split(bufio.ScanLines)

	var (
		head  = bytes.NewBuffer(nil)
		body  = bytes.NewBuffer(nil)
		state int
	)

	for scanner.Scan() {
		line := scanner.Bytes()
		if state == 0 {
			if len(line) == 0 {
				continue
			}
			if string(line) == start {
				state = started
				continue
			}
			state = finished
		}
		if state == started {
			if string(line) == end {
				state = finished
				continue
			}
			unsafe.DoSilent(head.Write(line))
			unsafe.Ignore(head.WriteByte('\n'))
			continue
		}
		unsafe.DoSilent(body.Write(line))
		unsafe.Ignore(body.WriteByte('\n'))
	}

	return head.Bytes(), body.Bytes(), scanner.Err()
}
