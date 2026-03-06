package runtimeprobe

type gcNode struct {
	next  *gcNode
	value int
}

func GCBarrierPath() *gcNode {
	left := &gcNode{value: 1}
	right := &gcNode{value: 2}

	left.next = right
	if left.next != nil && left.next.value == 2 {
		return left
	}

	return nil
}
