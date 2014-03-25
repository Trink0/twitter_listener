package listener

// dummyStore is a fake implementation of Store.
type dummyStore struct {
	listAppNames func() ([]string, error)
	getApp       func(string) (*Application, error)
	getUserIds   func(string) ([]string, error)
}

func (s *dummyStore) ListAppNames() ([]string, error) {
	return s.listAppNames()
}

func (s *dummyStore) GetApp(name string) (*Application, error) {
	return s.getApp(name)
}

func (s *dummyStore) ListAppUserIds(name string) ([]string, error) {
	return s.getUserIds(name)
}
