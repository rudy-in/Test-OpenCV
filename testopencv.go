package main

import (
	"container/heap"
	"fmt"
	"math"
)

// Cell represents each cell in the grid
type Cell struct {
	parentI int
	parentJ int
	f       float64
	g       float64
	h       float64
}

// Define the size of the grid
const (
	ROW = 9
	COL = 10
)

// PriorityQueueItem represents an item in the priority queue
type PriorityQueueItem struct {
	f    float64
	row  int
	col  int
	index int
}

// PriorityQueue implements heap.Interface and holds PriorityQueueItems
type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func is_valid(row, col int) bool {
	return row >= 0 && row < ROW && col >= 0 && col < COL
}

func is_unblocked(grid [][]int, row, col int) bool {
	return grid[row][col] == 1
}

func is_destination(row, col int, dest [2]int) bool {
	return row == dest[0] && col == dest[1]
}

func calculate_h_value(row, col int, dest [2]int) float64 {
	return math.Sqrt(float64((row-dest[0])*(row-dest[0]) + (col-dest[1])*(col-dest[1])))
}

func trace_path(cellDetails [][]Cell, dest [2]int) {
	fmt.Println("The Path is ")
	path := [][]int{}
	row, col := dest[0], dest[1]

	for !(cellDetails[row][col].parentI == row && cellDetails[row][col].parentJ == col) {
		path = append(path, []int{row, col})
		tempRow := cellDetails[row][col].parentI
		tempCol := cellDetails[row][col].parentJ
		row = tempRow
		col = tempCol
	}

	path = append(path, []int{row, col})
	for i := len(path) - 1; i >= 0; i-- {
		fmt.Printf("-> (%d,%d) ", path[i][0], path[i][1])
	}
	fmt.Println()
}

func a_star_search(grid [][]int, src, dest [2]int) {
	if !is_valid(src[0], src[1]) || !is_valid(dest[0], dest[1]) {
		fmt.Println("Source or destination is invalid")
		return
	}

	if !is_unblocked(grid, src[0], src[1]) || !is_unblocked(grid, dest[0], dest[1]) {
		fmt.Println("Source or the destination is blocked")
		return
	}

	if is_destination(src[0], src[1], dest) {
		fmt.Println("We are already at the destination")
		return
	}

	closedList := make([][]bool, ROW)
	for i := range closedList {
		closedList[i] = make([]bool, COL)
	}

	cellDetails := make([][]Cell, ROW)
	for i := range cellDetails {
		cellDetails[i] = make([]Cell, COL)
		for j := range cellDetails[i] {
			cellDetails[i][j] = Cell{
				parentI: -1,
				parentJ: -1,
				f:       math.Inf(1),
				g:       math.Inf(1),
				h:       0,
			}
		}
	}

	i, j := src[0], src[1]
	cellDetails[i][j].f = 0.0
	cellDetails[i][j].g = 0.0
	cellDetails[i][j].h = 0.0
	cellDetails[i][j].parentI = i
	cellDetails[i][j].parentJ = j

	openList := &PriorityQueue{}
	heap.Init(openList)
	heap.Push(openList, &PriorityQueueItem{f: 0.0, row: i, col: j})

	foundDest := false
	directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

	for openList.Len() > 0 {
		pqItem := heap.Pop(openList).(*PriorityQueueItem)
		i, j = pqItem.row, pqItem.col
		closedList[i][j] = true

		for _, dir := range directions {
			newI, newJ := i+dir[0], j+dir[1]

			if is_valid(newI, newJ) && is_unblocked(grid, newI, newJ) && !closedList[newI][newJ] {
				if is_destination(newI, newJ, dest) {
					cellDetails[newI][newJ].parentI = i
					cellDetails[newI][newJ].parentJ = j
					fmt.Println("The destination cell is found")
					trace_path(cellDetails, dest)
					foundDest = true
					return
				} else {
					gNew := cellDetails[i][j].g + 1.0
					hNew := calculate_h_value(newI, newJ, dest)
					fNew := gNew + hNew

					if cellDetails[newI][newJ].f == math.Inf(1) || cellDetails[newI][newJ].f > fNew {
						heap.Push(openList, &PriorityQueueItem{f: fNew, row: newI, col: newJ})
						cellDetails[newI][newJ] = Cell{
							parentI: i,
							parentJ: j,
							f:       fNew,
							g:       gNew,
							h:       hNew,
						}
					}
				}
			}
		}
	}

	if !foundDest {
		fmt.Println("Failed to find the destination cell")
	}
}

func main() {
	grid := [][]int{
		{1, 0, 1, 1, 1, 1, 0, 1, 1, 1},
		{1, 1, 1, 0, 1, 1, 1, 0, 1, 1},
		{1, 1, 1, 0, 1, 1, 0, 1, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 1, 1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 1, 1, 1, 1, 0, 1, 0, 0},
		{1, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 0, 1, 1, 1},
		{1, 1, 1, 0, 0, 0, 1, 0, 0, 1},
	}

	src := [2]int{8, 0}
	dest := [2]int{0, 0}

	a_star_search(grid, src, dest)
}
