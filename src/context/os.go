package context

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	DARWIN  = "darwin"
)

func PathDelimiter() string {
	if Current == nil {
		return ":"
	}

	switch Current.OS {
	case WINDOWS:
		return ";"
	default:
		return ":"
	}
}

func PathSeparator() string {
	if Current == nil {
		return "/"
	}

	switch Current.OS {
	case WINDOWS:
		return "\\"
	default:
		return "/"
	}
}
