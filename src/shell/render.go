package shell

// renderEntries writes the non-empty, non-ignored scripts produced by content
// for each entry to DotFile, separated by a blank line and prefixed with a
// blank line if DotFile already has content.
//
// ignoreFirst controls whether the ignore check happens before or after
// content is computed, since content may have side effects (e.g. Link
// creating directories, Path mutating context.Current.Path) that some
// callers need to skip entirely when an entry is ignored, and others don't.
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
