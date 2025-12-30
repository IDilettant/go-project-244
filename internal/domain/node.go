package domain

import "sort"

type Node map[string]any

// KeysSorted returns node keys sorted ascending (deterministic traversal).
func (n Node) KeysSorted() []string {
	keys := make([]string, 0, len(n))
	for k := range n {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// UnionKeysSorted returns sorted union of keys of two nodes.
func (n Node) UnionKeysSorted(other Node) []string {
	unique := make(map[string]struct{}, len(n)+len(other))

	for k := range n {
		unique[k] = struct{}{}
	}
	for k := range other {
		unique[k] = struct{}{}
	}

	keys := make([]string, 0, len(unique))
	for k := range unique {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// IsEmpty is a small convenience helper.
func (n Node) IsEmpty() bool { return len(n) == 0 }

// GetNode tries to read a nested node by key.
func (n Node) GetNode(key string) (Node, bool) {
	v, ok := n[key]
	if !ok {
		return nil, false
	}
	return AsNode(v)
}

// AsNode converts value to Node if it is a nested object.
func AsNode(v any) (Node, bool) {
	switch t := v.(type) {
	case Node:
		return t, true
	case map[string]any:
		return Node(t), true
	default:
		return nil, false
	}
}

