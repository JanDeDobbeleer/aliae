package shell

type Scripts []*Script

type Script struct {
	Value Template    `yaml:"value"`
	If    interface{} `yaml:"if"`
}

func (s *Script) String() string {
	script := s.Value.Parse()
	return string(script)
}

func (s Scripts) Render() {
	if len(s) == 0 {
		return
	}

	first := true
	for _, script := range s {
		scriptBlock := script.String()
		if len(scriptBlock) == 0 || checkIf(script.If) {
			continue
		}

		if first && DotFile.Len() > 0 {
			DotFile.WriteString("\n\n")
		}

		if !first {
			DotFile.WriteString("\n")
		}

		DotFile.WriteString(scriptBlock)

		first = false
	}
}
