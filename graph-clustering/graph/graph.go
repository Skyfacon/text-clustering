package graph

// Graph 无向图？如果是有向图呢？其实这里涉及到一个问题，就是 A搜索得到B的相似度和B搜索A得到的相似度不匹配的问题
type Graph struct {
}

func (g *Graph) newGraph() *Graph {

}

type Edge struct {
	This *Node
	That *Node
}

func (e *Edge) one() *Node {
	return e.This
}

func (e *Edge) another(a *Node) *Node {
	if e.This == a {
		return e.That
	} else if e.That == a {
		return e.This
	}
	return nil
}

type Node struct {
	NodeId int
}

func (n *Node) Equal(other *Node) bool {
	if n.NodeId != other.NodeId {
		return false
	}
	return true
}
