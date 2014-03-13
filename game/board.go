package game

import ("log")

type board [9] interface{}

func (b *board) process(move int, p interface{}) bool {

	if move >= 0 && move < len(b) && b[move] == nil {
		b[move] = p
		return true
	}

	return false
}

func (b *board) done() bool {

	var ver = make([][]int, 7)
	ver[0] = []int{1, 3, 4}
	ver[1] = []int{3}
	ver[2] = []int{3, 2}
	ver[3] = []int{1}
	ver[6] = []int{1}

	for i, v := range ver {

		for _, v2 := range v {

			if b[i] != nil &&
				b[i] == b[i+v2] &&
				b[i] == b[i+v2+v2] {

				log.Printf("Done\n")
				return true
			}
		}
	}

	return false
}

func (b *board) stale() bool {

	for _, v := range b {
		if v == nil {
			return false
		}
	}
	return true
}
