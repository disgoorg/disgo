package discord

import "iter"

// componentIter returns an iterator over the given components and their children.
func componentIter(components []LayoutComponent) iter.Seq[Component] {
	return func(yield func(Component) bool) {
		for _, c := range components {
			if !yield(c) {
				return
			}
			ic, ok := c.(ComponentIter)
			if !ok {
				continue
			}

			for cc := range ic.SubComponents() {
				if !yield(cc) {
					return
				}
			}
		}
	}
}
