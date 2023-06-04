// connection package handle connection to rest-api server Warships Online API
package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ConnectionClient implementation of Client interface
type ConnectionClient struct {
	client         http.Client
	host           string
	token          string
	StartingHeader StartingHeader
}

// NewClient constructor of ConnectionClient struct
func NewClient(host string) ConnectionClient {
	return ConnectionClient{client: http.Client{}, host: host}
}

// SetStartingHeader set starting header
func (c *ConnectionClient) SetStartingHeader(setHeader StartingHeader) {
	c.StartingHeader = setHeader
}

// StartGame starting game with given params in starting header
// which decide if you want to play against bot or player or wait for opponent,
// if you want send your nick and description if specify them
// if you want send your own board and at the end store authentication token in ConnectionClient struct
func (c *ConnectionClient) StartGame(ctx context.Context) error {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game")
	if err != nil {
		return err
	}

	bodyJson, err := json.Marshal(c.StartingHeader)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(bodyJson)

	// Create a new POST request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, connectionString, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return err
		}
		return &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}

	fmt.Println(resp.Header.Get("X-Auth-Token"))
	c.token = resp.Header.Get("X-Auth-Token")

	if len(c.token) == 0 {
		return &TokenError{Token: c.token}
	}
	return nil

}

// GerBoard return game board from server using authentication token and decoding it to boardRespons struct
// and an error
func (c ConnectionClient) GetBoard(ctx context.Context) (BoardRespons, error) {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game/board")
	if err != nil {
		return BoardRespons{}, err
	}

	// Create a new GET request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectionString, http.NoBody)
	if err != nil {
		return BoardRespons{}, fmt.Errorf("cannot create request: %w", nil)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return BoardRespons{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return BoardRespons{}, err
		}
		return BoardRespons{}, &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return BoardRespons{}, &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}

	newBoard := BoardRespons{}
	err = json.NewDecoder(resp.Body).Decode(&newBoard)
	if err != nil {
		return BoardRespons{}, err
	}
	return newBoard, nil
}

// GetDescription return your and enemy login and description using authentication token
// in Description struct and an error
func (c ConnectionClient) GetDescription(ctx context.Context) (Description, error) {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game/desc")
	if err != nil {
		return Description{}, err
	}

	// Create a new GET request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectionString, http.NoBody)
	if err != nil {
		return Description{}, fmt.Errorf("cannot create request: %w", nil)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return Description{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return Description{}, err
		}
		return Description{}, &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return Description{}, &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}

	newStatus := Description{}
	err = json.NewDecoder(resp.Body).Decode(&newStatus)
	if err != nil {
		return Description{}, err
	}
	return newStatus, nil
}

// Fire is making a shot to a server with specified coordinates and return response about the outcome
// by decoding message to FireResponse struct and returning an error when function faild
func (c ConnectionClient) Fire(ctx context.Context, coordinates string) (FireResponse, error) {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game/fire")
	if err != nil {
		return FireResponse{}, err
	}

	bodyJson, err := json.Marshal(FireCoord{coordinates})
	if err != nil {
		return FireResponse{}, err
	}

	body := bytes.NewBuffer(bodyJson)

	// Create a new POST request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, connectionString, body)
	if err != nil {
		return FireResponse{}, fmt.Errorf("cannot create request: %w", nil)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return FireResponse{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return FireResponse{}, err
		}

		return FireResponse{}, &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return FireResponse{}, &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}

	newFire := FireResponse{}
	err = json.NewDecoder(resp.Body).Decode(&newFire)
	if err != nil {
		return FireResponse{}, err
	}
	return newFire, nil
}

// GetStatus retrieves the current game status from the server using the authentication token.
// It returns the GameStatus struct containing information about the game status and an error if the request fails.
func (c ConnectionClient) GetStatus(ctx context.Context) (GameStatus, error) {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game")
	if err != nil {
		return GameStatus{}, fmt.Errorf("cannot create request: %w", nil)
	}

	// Create a new GET request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectionString, http.NoBody)
	if err != nil {
		return GameStatus{}, fmt.Errorf("cannot create request: %w", nil)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return GameStatus{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return GameStatus{}, &RequestError{StatusCode: resp.StatusCode, Err: ""}
		}
		return GameStatus{}, &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return GameStatus{}, &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}
	newStatus := GameStatus{}
	err = json.NewDecoder(resp.Body).Decode(&newStatus)
	if err != nil {
		return GameStatus{}, err
	}
	return newStatus, nil
}

// GetPlayerList retrieves the list of players waiting for game from the server using the authentication token.
// It returns the PlayerList struct containing the list of players and an error if the request fails.
func (c ConnectionClient) GetPlayerList(ctx context.Context) (PlayerList, error) {
	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	connectionString, err := url.JoinPath(c.host, "/api/game/list")
	if err != nil {
		return PlayerList{}, err
	}

	// Create a new GET request with the connection string
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectionString, http.NoBody)
	if err != nil {
		return PlayerList{}, fmt.Errorf("cannot create request: %w", nil)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return PlayerList{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == 403 {
		errorBody := ErrorMessage{}
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return PlayerList{}, &RequestError{StatusCode: resp.StatusCode, Err: ""}
		}
		return PlayerList{}, &RequestError{StatusCode: resp.StatusCode, Err: errorBody.Message}
	}
	if resp.StatusCode != http.StatusOK {
		statusText := http.StatusText(resp.StatusCode)
		return PlayerList{}, &RequestError{StatusCode: resp.StatusCode, Err: statusText}
	}

	playerList := PlayerList{}
	err = json.NewDecoder(resp.Body).Decode(&playerList)
	if err != nil {
		return PlayerList{}, err
	}
	return playerList, nil
}
