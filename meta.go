package rainbowlog

type metaKeys struct {
	keys          []string
	consoleColors map[string][]int
}

func emptyMetaKes() *metaKeys {
	return &metaKeys{keys: []string{}, consoleColors: map[string][]int{}}
}

func (m *metaKeys) SetKeys(keys []string) {
	if m == nil {
		return
	}
	m.keys = keys
}

func (m *metaKeys) SetKeyColors(key string, colors []int) {
	m.consoleColors[key] = colors
}

func (m *metaKeys) Keys() []string {
	if m == nil {
		return nil
	}
	return m.keys
}

func (m *metaKeys) ConsoleColors(key string) []int {
	if m == nil {
		return nil
	}
	cs, ok := m.consoleColors[key]
	if !ok || len(cs) == 0 {
		return nil
	}
	return cs
}

func (m *metaKeys) Clone() *metaKeys {
	if m == nil {
		return nil
	}
	res := &metaKeys{}
	res.SetKeys(m.keys)
	for s, ints := range m.consoleColors {
		res.SetKeyColors(s, ints)
	}
	return res
}
