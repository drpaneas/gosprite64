package residency

type ResidencyClass uint8

const (
	ResidencyPermanent ResidencyClass = iota
	ResidencyScene
	ResidencyStreamed
)

type Store struct {
	classes map[string]ResidencyClass
}

func NewStore() *Store {
	return &Store{
		classes: make(map[string]ResidencyClass),
	}
}

func (s *Store) Pin(name string, class ResidencyClass) {
	if s == nil {
		return
	}
	s.classes[name] = class
}

func (s *Store) ClassOf(name string) ResidencyClass {
	if s == nil {
		return ResidencyStreamed
	}
	class, ok := s.classes[name]
	if !ok {
		return ResidencyStreamed
	}
	return class
}
