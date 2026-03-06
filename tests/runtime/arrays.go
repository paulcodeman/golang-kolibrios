package runtimeprobe

func ArrayValueOps() bool {
	left := [4]byte{'k', 'o', 's', '!'}
	right := left
	match := [4]byte{'k', 'o', 's', '!'}

	right[3] = '?'

	return left == match && right != match && len(left) == 4 && left[1] == 'o'
}
