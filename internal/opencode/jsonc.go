package opencode

import (
	"bytes"
)

func sanitizeJSONC(input []byte) []byte {
	return removeTrailingCommas(stripJSONC(input))
}

func stripJSONC(input []byte) []byte {
	out := &bytes.Buffer{}
	inString := false
	escaped := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
				out.WriteByte(ch)
			}
			continue
		}

		if inBlockComment {
			if ch == '\n' {
				out.WriteByte(ch)
			}
			if ch == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if inString {
			out.WriteByte(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out.WriteByte(ch)
			continue
		}

		if ch == '/' && i+1 < len(input) {
			next := input[i+1]
			if next == '/' {
				inLineComment = true
				i++
				continue
			}
			if next == '*' {
				inBlockComment = true
				i++
				continue
			}
		}

		out.WriteByte(ch)
	}

	return out.Bytes()
}

func removeTrailingCommas(input []byte) []byte {
	out := &bytes.Buffer{}
	inString := false
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if inString {
			out.WriteByte(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			out.WriteByte(ch)
			continue
		}

		if ch == ',' {
			j := i + 1
			for j < len(input) {
				next := input[j]
				if next == ' ' || next == '\t' || next == '\n' || next == '\r' {
					j++
					continue
				}
				if next == '}' || next == ']' {
					i = j - 1
					ch = 0
				}
				break
			}
			if ch == 0 {
				continue
			}
		}

		out.WriteByte(ch)
	}

	return out.Bytes()
}
