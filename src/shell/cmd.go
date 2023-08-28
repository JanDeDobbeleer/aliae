// doskey np=notepad++.exe $*

package shell

const (
	CMD = "cmd"
)

func (a *Alias) cmd() *Alias {
	if a.Type == Command {
		a.template = `local p = assert(io.popen("doskey {{ .Name }}={{ escapeString .Value }}"))
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

func (e *Env) cmd() *Env {
	e.template = `os.setenv("{{ .Name }}", {{ formatString .Value }})`
	return e
}

func (p *Path) cmd() *Path {
	p.template = `os.setenv("PATH", "{{ escapeString .Value }};" .. os.getenv("PATH"))`
	return p
}
