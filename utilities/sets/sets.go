package sets


type Set[T string|int|uint] struct {
	items map[T]bool
}

// Builds a new `Set` from a variadic list of `items`.
func NewSet[T string|int|uint](items ...T) Set[T] {
	setItems := make(map[T]bool)
	for _, item := range items {
		setItems[item] = true
	}
	return Set[T]{setItems}
}

// Builds a new empty `Set`.
func NewEmptySet[T string|int|uint]() Set[T] {
	setItems := make(map[T]bool)
	return Set[T]{setItems}
}

// ----- SET METHODS -----

// Returns a deep copy of this `Set`.
func (s *Set[T]) Clone() Set[T] {
	return NewSet(s.GetItems()...)
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Add(item T) {
	s.items[item] = true
}

func (s *Set[T]) Delete(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Has(item T) bool {
	exists := s.items[item]
	return exists
}

func (s *Set[T]) GetItems() []T {
	itemsArr := make([]T, len(s.items))
	idx := 0
	for item := range s.items {
		itemsArr[idx] = item
		idx++
	}
	return itemsArr
}

func (this *Set[T]) Union(other Set[T]) Set[T] {
	newSet := NewEmptySet[T]()
	for item := range this.items {
		newSet.Add(item)
	}
	for item := range other.items {
		newSet.Add(item)
	}
	return newSet
}

func (this *Set[T]) Intersection(other Set[T]) Set[T] {
	newSet := NewEmptySet[T]()
	for item := range this.items {
		if (other.Has(item)) {
			newSet.Add(item)
		}
	}
	return newSet
}