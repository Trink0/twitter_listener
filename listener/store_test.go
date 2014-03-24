package listener

type dummyAppStore struct {
  listAppNames func()([]string, error)
  getApp func(string)(*Application, error)
}

func (s *dummyAppStore) ListAppNames() ([]string, error) {
  return s.listAppNames()
}

func (s *dummyAppStore) GetApp(name string) (*Application, error) {
  return s.getApp(name)
}