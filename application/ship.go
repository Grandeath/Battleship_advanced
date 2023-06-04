package application

import (
	"log"

	gui "github.com/grupawp/warships-gui/v2"
)

type currentNode uint8

const (
	East   currentNode = iota
	North  currentNode = iota
	West   currentNode = iota
	South  currentNode = iota
	Center currentNode = iota
)

type ShipCoord struct {
	row    int
	column int
}

type QuadTree struct {
	Coord      ShipCoord
	NorthChild *QuadTree
	WestChild  *QuadTree
	EastChild  *QuadTree
	SouthChild *QuadTree
}

func NewQuadTree(coord ShipCoord, enemyBoardState [10][10]gui.State, current currentNode) *QuadTree {
	newNode := &QuadTree{Coord: coord}
	newNode.FindNode(enemyBoardState, current)
	return newNode
}

func (q *QuadTree) FindNode(enemyBoardState [10][10]gui.State, current currentNode) {
	log.Println("beep")

	if q.Coord.row+1 < 10 && current != South {
		if enemyBoardState[q.Coord.column][q.Coord.row+1] == gui.Hit {
			q.NorthChild = NewQuadTree(ShipCoord{column: q.Coord.column, row: q.Coord.row + 1}, enemyBoardState, North)
		} else {
			enemyBoardState[q.Coord.column][q.Coord.row+1] = gui.Empty
			if q.Coord.column+1 < 10 {
				enemyBoardState[q.Coord.column+1][q.Coord.row+1] = gui.Empty
			}
			if q.Coord.column-1 >= 0 {
				enemyBoardState[q.Coord.column-1][q.Coord.row+1] = gui.Empty
			}
		}
	}
	if q.Coord.row-1 >= 0 && current != North {
		if enemyBoardState[q.Coord.column][q.Coord.row-1] == gui.Hit {
			q.SouthChild = NewQuadTree(ShipCoord{column: q.Coord.column, row: q.Coord.row - 1}, enemyBoardState, South)
		} else {
			enemyBoardState[q.Coord.column][q.Coord.row-1] = gui.Empty
			if q.Coord.column+1 < 10 {
				enemyBoardState[q.Coord.column+1][q.Coord.row-1] = gui.Empty
			}
			if q.Coord.column-1 >= 0 {
				enemyBoardState[q.Coord.column-1][q.Coord.row-1] = gui.Empty
			}
		}
	}
	if q.Coord.column+1 < 10 && current != West {
		if enemyBoardState[q.Coord.column+1][q.Coord.row] == gui.Hit {
			q.EastChild = NewQuadTree(ShipCoord{column: q.Coord.column + 1, row: q.Coord.row}, enemyBoardState, East)
		} else {
			enemyBoardState[q.Coord.column+1][q.Coord.row] = gui.Empty
			if q.Coord.row+1 < 10 {
				enemyBoardState[q.Coord.column+1][q.Coord.row+1] = gui.Empty
			}
			if q.Coord.row-1 >= 0 {
				enemyBoardState[q.Coord.column+1][q.Coord.row-1] = gui.Empty
			}

		}
	}
	if q.Coord.column-1 >= 0 && current != East {
		if enemyBoardState[q.Coord.column-1][q.Coord.row] == gui.Hit {
			q.WestChild = NewQuadTree(ShipCoord{column: q.Coord.column - 1, row: q.Coord.row}, enemyBoardState, West)
		} else {
			enemyBoardState[q.Coord.column-1][q.Coord.row] = gui.Empty
			if q.Coord.row+1 < 10 {
				enemyBoardState[q.Coord.column-1][q.Coord.row+1] = gui.Empty
			}
			if q.Coord.row-1 >= 0 {
				enemyBoardState[q.Coord.column-1][q.Coord.row-1] = gui.Empty
			}
		}
	}

}

func (q *QuadTree) GetAllCoords() []ShipCoord {
	var coords []ShipCoord

	coords = append(coords, q.Coord)

	if q.NorthChild != nil {
		coords = append(coords, q.NorthChild.GetAllCoords()...)
	}
	if q.WestChild != nil {
		coords = append(coords, q.WestChild.GetAllCoords()...)
	}
	if q.EastChild != nil {
		coords = append(coords, q.EastChild.GetAllCoords()...)
	}
	if q.SouthChild != nil {
		coords = append(coords, q.SouthChild.GetAllCoords()...)
	}

	return coords
}
