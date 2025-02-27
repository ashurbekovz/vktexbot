package latex2img

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type LatexCompilationError struct {
	Message string 
	Line    int  
	Context string
}

func (e *LatexCompilationError) IsUnknownError() bool {
    return e.Message == ""
}

func (e *LatexCompilationError) KnownLine() bool {
    return e.Line != 0
}

func (e *LatexCompilationError) Error() string {
    if e.IsUnknownError() {
        return "unknown latex compilation error"
    } 
    
    message := fmt.Sprintf("error: %s", e.Message)

    if e.KnownLine() {
        message += fmt.Sprintf(", on line: %d", e.Line)
    }

    return message
}

const (
    maxContextLines = 3
)

var (
	lineNumRegex = regexp.MustCompile(`l\.(\d+)`)
)

// parseLatexError return first error founded in the logs.
// Latex has no documentation on how errors are organized,
// so this parsing is written purely empirically.
func parseLatexError(log []byte) error {
	var currentError *LatexCompilationError
	var contextBuffer []string

	scanner := bufio.NewScanner(bytes.NewReader(log))
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if currentError == nil {
            if strings.HasPrefix(line, "! ") {
                currentError = &LatexCompilationError{
                    Message: strings.TrimPrefix(line, "! "),
                }
                contextBuffer = nil
            }

            continue
		}

        if matches := lineNumRegex.FindStringSubmatch(line); currentError.Line == 0 && len(matches) > 1 {
            if num, err := strconv.Atoi(matches[1]); err == nil {
                currentError.Line = num
            }
        }

        if strings.TrimSpace(line) != "" {
            contextBuffer = append(contextBuffer, line)
            if len(contextBuffer) > maxContextLines {
                contextBuffer = contextBuffer[1:]
            }
        }

        if strings.TrimSpace(line) == "" {
            currentError.Context = strings.Join(contextBuffer, "\n")
            break
        }
	}

	return currentError
}
