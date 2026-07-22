package interfaces

import "context"

type Player map[string]any

type GameSession map[string]any

type LaunchGameRequest struct {
	Player  Player
	Session GameSession
}

type LaunchGameResult struct {
	GameSource string
	LaunchType string
}

// GamingIntegration is the minimum capability required by every gaming partner.
type GamingIntegration interface {
	LaunchGame(ctx context.Context, request LaunchGameRequest) (LaunchGameResult, error)
}

// PlayerCreator is an optional capability for partners that require player creation.
type PlayerCreator interface {
	CreatePlayer(ctx context.Context, player Player) error
}

type ControlEventRequest struct {
	Name    string
	Player  Player
	Session GameSession
	Inputs  map[string]any
}

// ControlEventHandler is an optional, extensible partner capability.
type ControlEventHandler interface {
	ControlEvent(ctx context.Context, request ControlEventRequest) (map[string]any, error)
}
