// doskey np=notepad++.exe $*

package shell

const (
	CMD = "cmd"
)

func (a *Alias) cmd() *Alias {
	if a.Type == Command {
		a.template = `local p = assert(io.popen("doskey {{ .Alias }}={{ .Value }}"))
p:close()`
	}

	return a
}

func (e *Echo) cmd() *Echo {
	e.template = `message = [[
{{ .Message }}
]]
print(message)`
	return e
}

func (e *Variable) cmd() *Variable {
	e.template = `os.setenv("{{ .Name }}", {{ formatString .Value }})`
	return e
}

func (p *PathEntry) cmd() *PathEntry {
	p.template = `os.setenv("PATH", "{{ cleanString .Value }};" .. os.getenv("PATH"))`
	return p
}
