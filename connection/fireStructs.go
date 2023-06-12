package connection

// FireResponse store result of the fire function obout if you hit or missed the shot
type FireResponse struct {
	Result string `json:"result"`
}

// FireCoord header for fire function which store infrormation about coordination of the shot
type FireCoord struct {
	Coord string `json:"coord"`
}
