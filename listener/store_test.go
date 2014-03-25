package listener

// dummyAppStore is a fake implementation of AppStore.
type dummyAppStore struct {
	listAppNames func() ([]string, error)
	getApp       func(string) (*Application, error)
	getUserIds   func(string) ([]string, error)
}

func (s *dummyAppStore) ListAppNames() ([]string, error) {
	return s.listAppNames()
}

func (s *dummyAppStore) GetApp(name string) (*Application, error) {
	return s.getApp(name)
}

func (s *dummyAppStore) ListAppUserIds(name string) ([]string, error) {
	return s.getUserIds(name)
}
