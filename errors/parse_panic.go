package errors

import (
	"strconv"
	"strings"
)


type Panic struct{message string}
func (p Panic) Error() string {
    return p.message
}

func ParsePanic(text string) (*Error, error) {
	lines := strings.Split(text, "\n")

    state := "start"

    var message string
    var stack []StackFrame

    for i := 0; i < len(lines); i += 1 {
        line := lines[i]

        if state == "start" {
            if strings.HasPrefix(line, "panic: ") {
                message = strings.TrimPrefix(line, "panic: ")
                state = "seek"
            } else {
				return nil, Errorf("bugsnag.panicParser: Invalid line (no prefix): %s", line)
            }

        } else if state == "seek" {
            if strings.HasPrefix(line, "goroutine ") && strings.HasSuffix(line, "[running]:") {
                state = "parsing"
            }

        } else if state == "parsing" {
            if line == "" || strings.HasPrefix(line, "created by ") {
                state = "done"
                break
            }

			i += 1

            if i >= len(lines) {
				return nil, Errorf("bugsnag.panicParser: Invalid line (unpaired): %s", line)
            }

            frame, err := parsePanicFrame(line, lines[i])
            if err != nil{
				return nil, err
            }

            stack = append(stack, *frame)
        }
    }

    if state == "done" || state == "parsing" {
		return &Error{Err: Panic{message}, frames: stack}, nil
	} else {
		return nil, Errorf("could not parse panic: %v", text)
	}
}

// The lines we're passing look like this:
//
//     main.(*foo).destruct(0xc208067e98)
//             /0/go/src/github.com/bugsnag/bugsnag-go/pan/main.go:22 +0x151
func parsePanicFrame(name string, line string) (*StackFrame, error) {
    idx := strings.LastIndex(name, "(")
    if idx == -1 {
		return nil, Errorf("bugsnag.panicParser: Invalid line (no call): %s", name)
    }
    name = name[:idx]
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
		return nil, Errorf("bugsnag.panicParser: Invalid line (no tab): %s", line)
    }

    idx = strings.LastIndex(line, ":")
    if idx == -1 {
        return nil, Errorf("bugsnag.panicParser: Invalid line (no line number): %s", line)
    }
    file := line[1:idx]

	number := line[idx+1:]
	if idx = strings.Index(number, " +"); idx > -1 {
		number = number[:idx]
	}

    lno, err := strconv.ParseInt(number, 10, 32)
    if err != nil {
        return nil, Errorf("bugsnag.panicParser: Invalid line (bad line number): %s", line)
    }

    return &StackFrame{
        File: file,
        LineNumber: int(lno),
        Package: pkg,
        Name: name,
    }, nil
}
