// doskey np=notepad++.exe $*

package shell

const (
	CMD = "cmd"
)

func (a *Alias) cmd() *Alias {
	if a.Type == Command {
		a.template = "macrofile:write(\"{{ .Name }}={{ escapeString .Value }}\", \"\\n\")"
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

func cmdAliasPre() string {
	return `
local filename  = os.tmpname()
local macrofile = io.open(filename, "w+")
`
}

func cmdAliasPost() string {
	return `
macrofile:close()
local _ = io.popen(string.format("doskey /macrofile=%s", filename)):close()
os.remove(filename)
`
}
