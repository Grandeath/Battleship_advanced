package ship

import (
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
	Row    int
	Column int
}

type QuadTree struct {
	Coord      ShipCoord
	NorthChild *QuadTree
	WestChild  *QuadTree
	EastChild  *QuadTree
	SouthChild *QuadTree
}

func NewQuadTree(coord ShipCoord, enemyBoardState [10][10]gui.State, current currentNode, searchState gui.State) *QuadTree {
	newNode := &QuadTree{Coord: coord}
	enemyBoardState[coord.Column][coord.Row] = gui.Empty
	newNode.FindNode(&enemyBoardState, current, searchState)
	return newNode
}

func (q *QuadTree) FindNode(enemyBoardState *[10][10]gui.State, current currentNode, searchState gui.State) {

	if q.Coord.Row+1 < 10 && current != South {
		if enemyBoardState[q.Coord.Column][q.Coord.Row+1] == searchState {
			enemyBoardState[q.Coord.Column][q.Coord.Row+1] = gui.Empty
			q.NorthChild = NewQuadTree(ShipCoord{Column: q.Coord.Column, Row: q.Coord.Row + 1}, *enemyBoardState, North, searchState)
		}
	}
	if q.Coord.Row-1 >= 0 && current != North {
		if enemyBoardState[q.Coord.Column][q.Coord.Row-1] == searchState {
			enemyBoardState[q.Coord.Column][q.Coord.Row-1] = gui.Empty
			q.SouthChild = NewQuadTree(ShipCoord{Column: q.Coord.Column, Row: q.Coord.Row - 1}, *enemyBoardState, South, searchState)
		}
	}
	if q.Coord.Column+1 < 10 && current != West {
		if enemyBoardState[q.Coord.Column+1][q.Coord.Row] == searchState {
			enemyBoardState[q.Coord.Column+1][q.Coord.Row] = gui.Empty
			q.EastChild = NewQuadTree(ShipCoord{Column: q.Coord.Column + 1, Row: q.Coord.Row}, *enemyBoardState, East, searchState)
		}
	}
	if q.Coord.Column-1 >= 0 && current != East {
		if enemyBoardState[q.Coord.Column-1][q.Coord.Row] == searchState {
			enemyBoardState[q.Coord.Column-1][q.Coord.Row] = gui.Empty
			q.WestChild = NewQuadTree(ShipCoord{Column: q.Coord.Column - 1, Row: q.Coord.Row}, *enemyBoardState, West, searchState)
		}
	}

}

func (q *QuadTree) GetAllCoords() []ShipCoord {
	visited := make(map[ShipCoord]bool)
	return q.getAllCoordsUnique(visited)
}

func (q *QuadTree) getAllCoordsUnique(visited map[ShipCoord]bool) []ShipCoord {
	var coords []ShipCoord

	// Check if the current coordinate has been visited before
	if !visited[q.Coord] {
		visited[q.Coord] = true
		coords = append(coords, q.Coord)
	}

	if q.NorthChild != nil {
		childCoords := q.NorthChild.getAllCoordsUnique(visited)
		coords = append(coords, childCoords...)
	}
	if q.WestChild != nil {
		childCoords := q.WestChild.getAllCoordsUnique(visited)
		coords = append(coords, childCoords...)
	}
	if q.EastChild != nil {
		childCoords := q.EastChild.getAllCoordsUnique(visited)
		coords = append(coords, childCoords...)
	}
	if q.SouthChild != nil {
		childCoords := q.SouthChild.getAllCoordsUnique(visited)
		coords = append(coords, childCoords...)
	}

	return coords
}
