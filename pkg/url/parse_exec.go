package url

import (
	"strings"

	"github.com/pkg/errors"
)

// execURLPrefix is the lowercase version of the Docker URL prefix.
const execURLPrefix = "exec:"

// isExecURL checks whether or not a URL is an exec URL.
func isExecURL(raw string) bool {
	return strings.HasPrefix(strings.ToLower(raw), execURLPrefix)
}

// parseExec parses an exec URL.
func parseExec(raw string, kind Kind) (*URL, error) {
	// Strip off the prefix.
	raw = raw[len(execURLPrefix):]

	// Split string by colon, and correctly handle quotes around the components
	quoted := false
	quoteMark := '"'
	split := strings.FieldsFunc(raw, func(r rune) bool {
		if quoted && r == quoteMark {
			quoted = false
		} else if !quoted && (r == '"' || r == '\'') {
			quoted = true
			quoteMark = r
		}
		return !quoted && r == ':'
	})

	path := ""
	command := ""

	if len(split) == 0 {
		return nil, errors.New("no command or path specified")
	} else if len(split) > 2 {
		return nil, errors.New("too many separators")
	} else {
		command = split[0]

		// Strip quotes off command
		if len(command) >= 2 && ((command[0] == '\'' && command[len(command) - 1] == '\'') || (command[0] == '"' && command[len(command) - 1] == '"')) {
			command = command[1:len(command) - 1]
		}

		// Set path, if provided
		if len(split) == 2 {
			path = split[1]
		}
	}

	// We return the command as the host
	return &URL{
		Kind:        kind,
		Protocol:    Protocol_Exec,
		Host:        command,
		Path:        path,
	}, nil
}
