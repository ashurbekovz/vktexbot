package latex2img

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SyntaxError struct {
	Message string
	Line    int
	Context string
}

func (e *SyntaxError) IsUnknownError() bool {
	return e.Message == ""
}

func (e *SyntaxError) KnownLine() bool {
	return e.Line != 0
}

func (e *SyntaxError) KnownContext() bool {
	return e.Context == ""
}

func (e *SyntaxError) Error() string {
	if e.IsUnknownError() {
		return "unknown latex compilation error"
	}

	message := fmt.Sprintf("error: %s", e.Message)

	if e.KnownLine() {
		message += fmt.Sprintf(", on line: %d", e.Line)
	}

	return message
}

func (e *SyntaxError) UserError() string {
	if e.IsUnknownError() {
		return "Где-то в сообщении есть ошибка, но, к сожалению, автоматически понять по выводу LaTeX-компилятора причину не удалось."
	}

	var builder strings.Builder

	builder.WriteString("Компиляция LaTeX завершилась с ошибкой.\n")
	builder.WriteString(fmt.Sprintf("Текст ошибки: %s", e.Message))

	if e.KnownLine() {
		builder.WriteString(fmt.Sprintf("\nСтрока: %d", e.Line-13))
	}

	if e.KnownContext() {
		builder.WriteString(fmt.Sprintf("\nКонтекст: %s", e.Context))
	}

	return builder.String()
}

type UnknownCompilationError struct {
	Message string
}

func (e *UnknownCompilationError) Error() string {
	return fmt.Sprintf("unknown latex compilation error: %s", e.Message)
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
func parseLatexError(output []byte) error {
	err := parseSyntaxError(output)
	if err != nil {
		return err
	}

	return &UnknownCompilationError{Message: string(output)}
}

func parseSyntaxError(output []byte) error {
	var currentError *SyntaxError
	var contextBuffer []string

	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()

		if currentError == nil {
			if strings.HasPrefix(line, "! ") {
				currentError = &SyntaxError{
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

	if currentError == nil {
		return nil
	}
	return currentError
}
