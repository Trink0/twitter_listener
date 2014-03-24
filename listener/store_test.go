// +build test

package listener

import (
	"fmt"
)

func (s *AppStore) ListAppNames() ([]string, error) {
	return []string{"chumhum"}, nil
}

func (s *AppStore) GetApp(name string) (app *Application, err error) {
	switch name {
	case "chumhum":
		app = &Application{Name: name}
		err = nil
	default:
		err = fmt.Errorf("Fake app %q not found ", name)
	}
	return
}
