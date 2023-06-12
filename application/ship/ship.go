// contain algoritm which find all ship nodes
package ship

import (
	gui "github.com/grupawp/warships-gui/v2"
)

// currentNode indicator side of current node
type currentNode uint8

// All possible heading direction: East, North, West, South and Center
const (
	East   currentNode = iota
	North  currentNode = iota
	West   currentNode = iota
	South  currentNode = iota
	Center currentNode = iota
)

// ShipCoord contain Row and Column of the ship
type ShipCoord struct {
	Row    int
	Column int
}

// QuadTree contain Coord of current node and pointer to all children
type QuadTree struct {
	Coord      ShipCoord
	NorthChild *QuadTree
	WestChild  *QuadTree
	EastChild  *QuadTree
	SouthChild *QuadTree
}

// NewQuadTree take coordinates, boardState, direction of script and searching state and creates pointer to QuadTree with finding all children of the node
func NewQuadTree(coord ShipCoord, enemyBoardState [10][10]gui.State, current currentNode, searchState gui.State) *QuadTree {
	newNode := &QuadTree{Coord: coord}
	enemyBoardState[coord.Column][coord.Row] = gui.Empty
	newNode.FindNode(&enemyBoardState, current, searchState)
	return newNode
}

// FindNode find all the children of the node
func (q *QuadTree) FindNode(enemyBoardState *[10][10]gui.State, current currentNode, searchState gui.State) {
	// check if heading North is correct and if previous node is not in the South
	if q.Coord.Row+1 < 10 && current != South {
		// Check check if North node is searched state
		if enemyBoardState[q.Coord.Column][q.Coord.Row+1] == searchState {
			// Delete this node to search again
			enemyBoardState[q.Coord.Column][q.Coord.Row+1] = gui.Empty
			// Create new node
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

// GetAllCoords make a slice of all unique coordinates
func (q *QuadTree) GetAllCoords() []ShipCoord {
	visited := make(map[ShipCoord]bool)
	return q.getAllCoordsUnique(visited)
}

// getAllCoordsUnique search for all unique coordinates
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
