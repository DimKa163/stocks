package types

import "github.com/beevik/guid"

type Path struct {
	nodes  []guid.Guid
	length int
}

func NewPath(capacity int) *Path {
	return &Path{nodes: make([]guid.Guid, capacity)}
}

func NewPathFromSlice(slice []guid.Guid) *Path {
	return &Path{nodes: slice, length: len(slice)}
}

func (p *Path) Destination() guid.Guid {
	if len(p.nodes) == 0 {
		return guid.Guid{}
	}
	return p.nodes[0]
}

func (p *Path) Foreach(fn func(node guid.Guid) (bool, error)) error {
	for i := 0; i < p.length; i++ {
		cont, err := fn(p.nodes[i])
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}

func (p *Path) ForeachReverse(fn func(node guid.Guid) (bool, error)) error {
	for i := p.length - 1; i > 0; i-- {
		cont, err := fn(p.nodes[i])
		if err != nil {
			return err
		}
		if !cont {
			break
		}
	}
	return nil
}

func (p *Path) AddNode(n guid.Guid) {
	if len(p.nodes) > p.length {
		p.nodes[p.length] = n
		p.length += 1
	} else {
		p.nodes = append(p.nodes, n)
		p.length += 1
	}
}

func (p *Path) Len() int {
	return p.length
}
