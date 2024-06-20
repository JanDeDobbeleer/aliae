package shell

import (
	"fmt"

	"github.com/jandedobbeleer/aliae/src/context"
)

type If string
type Ifs []string

func (i If) Ignore() bool {
	if len(i) == 0 {
		return false
	}
	template := fmt.Sprintf(`{{ if %s }}false{{ else }}true{{ end }}`, i)

	got, err := parse(template, context.Current)
	if err != nil {
		return false
	}

	return got == "true"
}

func (ifs Ifs) Ignore() bool {
	for _, i := range ifs {
		if If(i).Ignore() {
			return true
		}
	}

	return false
}

func convertToIfs(interfaces []interface{}) Ifs {
	strs := make(Ifs, len(interfaces))
	for index, value := range interfaces {
		str, ok := value.(string)
		if !ok {
			continue
		}
		strs[index] = str
	}
	return strs
}

func checkIf(i any) bool {
	switch i := i.(type) {
	case string:
		return If(i).Ignore()
	case []string:
		return Ifs(i).Ignore()
	case If:
		return i.Ignore()
	case Ifs:
		return i.Ignore()
	case []interface{}:
		return convertToIfs(i).Ignore()
	default:
		return false
	}
}
