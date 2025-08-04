package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bugsnag/bugsnag-go/v2"
	e "github.com/bugsnag/bugsnag-go/v2/errors"
)

var panicHeaders = [][]byte{
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

func parsePanic(text string) (*uncaughtPanic, *bugsnag.MetaData, error) {
	lines := strings.Split(text, "\n")
	prefixes := []string{"panic:", "fatal error:"}

	state := "start"

	var message string
	var typeName string
	var stack []e.StackFrame
	var metadata *bugsnag.MetaData

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if state == "start" {
			for _, prefix := range prefixes {
				if strings.HasPrefix(line, prefix) {
					message = strings.TrimSpace(strings.TrimPrefix(line, prefix))
					typeName = prefix[:len(prefix)-1]
					// If this was a signal, look for more signal details on the next line
					if strings.Contains(message, "signal") {
						state = "signal"
						break
					}
					state = "seek"
					break
				}
			}
			if state == "start" {
				return nil, nil, fmt.Errorf("panic-monitor: Invalid line (no prefix): %s", line)
			}

		} else if state == "signal" {
			// Capture signal details if present and store in metadata
			var re = regexp.MustCompile(`\[signal (?P<signal>[A-Z]+?): (?P<desc>.*?) code=(?P<code>[0-9xa-f]+?) addr=(?P<addr>[0-9xa-f]+?) pc=(?P<pc>[0-9xa-f]+?)\]`)

			fields := re.FindStringSubmatch(line)
			if fields != nil {
				// Note: we only add metadata here ATM. If more Metadata will be added elsewhere, then we need better initialization
				typeName = fields[1]
				message = fields[2]
				metadata = &bugsnag.MetaData{
					"signal": {
						"code": fields[3],
						"addr": fields[4],
						"pc":   fields[5],
					},
				}
			}

			// Whether or not we found signal details, start seeking
			state = "seek"

		} else if state == "seek" {
			if strings.HasPrefix(line, "goroutine ") &&
				strings.HasSuffix(line, "[running]:") || strings.HasSuffix(line, "[syscall]:") {
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
				return nil, nil, fmt.Errorf("panic-monitor: Invalid line (unpaired): '%s'", line)
			}

			frame, err := parsePanicFrame(line, lines[i], createdBy)
			if err != nil {
				return nil, nil, err
			}

			stack = append(stack, *frame)
			if createdBy {
				state = "done"
				break
			}
		}
	}

	if state == "done" || state == "parsing" {
		return &uncaughtPanic{typeName, message, stack}, metadata, nil
	}
	return nil, nil, fmt.Errorf("panic-monitor: could not parse panic: %v", text)
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
	// delineate generated files, using a separate package name to exclude from
	// project packages (and in-project detection) by default
	if isGeneratedFile(file) {
		pkg = "<generated>." + pkg
	}

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

func isGeneratedFile(file string) bool {
	// generated file patterns are documented in
	// https://go.dev/src/cmd/cgo/doc.go
	return strings.HasPrefix(file, "_cgo_") ||
		strings.HasSuffix(file, "cgo1.go") ||
		strings.HasSuffix(file, "cgo2.c")
}
