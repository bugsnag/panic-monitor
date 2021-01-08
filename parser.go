package main

import (
	"fmt"
	"strconv"
	"strings"
	e "github.com/bugsnag/bugsnag-go/v2/errors"
)

var panicHeaders = [][]byte {
	[]byte("panic:"),
	[]byte("fatal error:"),
}

type uncaughtPanic struct {
	typeName string
	message  string
	frames   []e.StackFrame
}

func (p uncaughtPanic) Error() string {
	return p.message
}

func (p uncaughtPanic) TypeName() string {
	return p.typeName
}

func (p uncaughtPanic) StackFrames() []e.StackFrame {
	return p.frames
}

func parsePanic(text string) (*uncaughtPanic, error) {
	lines := strings.Split(text, "\n")
	prefixes := []string{"panic:", "fatal error:"}

	state := "start"

	var message string
	var typeName string
	var stack []e.StackFrame

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if state == "start" {
			for _, prefix := range prefixes {
				if strings.HasPrefix(line, prefix) {
					message = strings.TrimSpace(strings.TrimPrefix(line, prefix))
					typeName = prefix[:len(prefix) - 1]
					state = "seek"
					break
				}
			}
			if state == "start" {
				return nil, fmt.Errorf("panic-monitor: Invalid line (no prefix): %s", line)
			}

		} else if state == "seek" {
			if strings.HasPrefix(line, "goroutine ") && strings.HasSuffix(line, "[running]:") {
				state = "parsing"
			}

		} else if state == "parsing" {
			if line == "" || strings.HasPrefix(line, "...") {
				state = "done"
				break
			}
			createdBy := false
			if strings.HasPrefix(line, "created by ") {
				line = strings.TrimPrefix(line, "created by ")
				createdBy = true
			}

			i++

			if i >= len(lines) {
				return nil, fmt.Errorf("panic-monitor: Invalid line (unpaired): '%s'", line)
			}

			frame, err := parsePanicFrame(line, lines[i], createdBy)
			if err != nil {
				return nil, err
			}

			stack = append(stack, *frame)
			if createdBy {
				state = "done"
				break
			}
		}
	}

	if state == "done" || state == "parsing" {
		return &uncaughtPanic{typeName, message, stack}, nil
	}
	return nil, fmt.Errorf("panic-monitor: could not parse panic: %v", text)
}

// The lines we're passing look like this:
//
//     main.(*foo).destruct(0xc208067e98)
//             /0/go/src/github.com/bugsnag/bugsnag-go/pan/main.go:22 +0x151
func parsePanicFrame(name string, line string, createdBy bool) (*e.StackFrame, error) {
	idx := strings.LastIndex(name, "(")
	if idx == -1 && !createdBy {
		return nil, fmt.Errorf("panic-monitor: Invalid line (no call): %s", name)
	}
	if idx != -1 {
		name = name[:idx]
	}
	pkg := ""

	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)

	if !strings.HasPrefix(line, "\t") {
		return nil, fmt.Errorf("panic-monitor: Invalid line (no tab): %s", line)
	}

	idx = strings.LastIndex(line, ":")
	if idx == -1 {
		return nil, fmt.Errorf("panic-monitor: Invalid line (no line number): %s", line)
	}
	file := line[1:idx]

	number := line[idx+1:]
	if idx = strings.Index(number, " +"); idx > -1 {
		number = number[:idx]
	}

	lno, err := strconv.ParseInt(number, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("panic-monitor: Invalid line (bad line number): %s", line)
	}

	return &e.StackFrame{
		File:       file,
		LineNumber: int(lno),
		Package:    pkg,
		Name:       name,
	}, nil
}
