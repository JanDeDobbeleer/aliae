package shell

// renderEntries writes the non-empty, non-ignored scripts produced by content
// to DotFile. ignoreFirst controls whether ignore is checked before content,
// since content can have side effects some callers need to skip.
func renderEntries[T any](entries []T, ignoreFirst bool, ignore func(T) bool, content func(T) string) {
	if len(entries) == 0 {
		return
	}

	first := true
	for _, entry := range entries {
		var script string

		if ignoreFirst {
			if ignore(entry) {
				continue
			}

			script = content(entry)
			if len(script) == 0 {
				continue
			}
		} else {
			script = content(entry)
			if len(script) == 0 || ignore(entry) {
				continue
			}
		}

		if first && DotFile.Len() > 0 {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString("\n")
		DotFile.WriteString(script)

		first = false
	}
}
