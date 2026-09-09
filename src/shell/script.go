package shell

type Scripts []*Script

type Script struct {
	Value Template `yaml:"value"`
	If    If       `yaml:"if"`
}

func (s *Script) String() string {
	script := s.Value.Parse()
	return string(script)
}

func (s Scripts) Render() {
	renderEntries(s, false, func(entry *Script) bool {
		return entry.If.Ignore()
	}, func(entry *Script) string {
		return entry.String()
	})
}
