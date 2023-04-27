package bot

import (
	"github.com/prometheus/client_golang/prometheus"
	loc "github.com/tcotav/travelbot/location"
)

type TravelBot struct {
	Name             string
	Speed            float64
	Position         loc.Point
	Destination      loc.Location
	Route            *RouteList
	Map              map[string]loc.Location
	IsErrored        bool
	MaintRequired    bool
	StopsMade        *prometheus.CounterVec
	DistanceTraveled *prometheus.CounterVec
}

func NewTravelBot(name string,
	speed float64,
	position loc.Point,
	route *RouteList,
	routemap map[string]loc.Location,
	stopsMade *prometheus.CounterVec,
	distanceTraveled *prometheus.CounterVec,
) *TravelBot {
	return &TravelBot{
		Name:             name,
		Speed:            speed,
		Position:         route.Next().Position,
		Destination:      route.Next(),
		Route:            route,
		Map:              routemap,
		IsErrored:        false,
		MaintRequired:    false,
		StopsMade:        stopsMade,
		DistanceTraveled: distanceTraveled,
	}
}

func (bot *TravelBot) Move() {
	totalDistance := bot.Position.DistanceTo(bot.Destination.Position)

	dx := bot.Destination.Position.X - bot.Position.X
	dy := bot.Destination.Position.Y - bot.Position.Y

	direction := loc.Point{X: dx / totalDistance, Y: dy / totalDistance}
	moveDelta := loc.Point{X: direction.X * bot.Speed, Y: direction.Y * bot.Speed}

	// If the next move would overshoot the target, set currentPosition to p2
	if bot.Position.DistanceTo(bot.Destination.Position) < bot.Speed {
		// we are going to stop now at a world

		// check if we need maintenance work done
		// if so, this flag is set and is only unset by external process
		if bot.MaintRequired {
			bot.MaintRequired = false
			// and we skip this move to do our maintenance
		} else {
			// we set our next destination
			bot.DistanceTraveled.WithLabelValues(bot.Name).Add(bot.Position.DistanceTo(bot.Destination.Position))
			bot.Position = bot.Destination.Position
			sNextDest := bot.Route.Next()
			bot.Destination = bot.Map[sNextDest.Name]
			bot.StopsMade.WithLabelValues(bot.Name).Inc()
		}
	} else {
		bot.DistanceTraveled.WithLabelValues(bot.Name).Add(bot.Speed)
		bot.Position.X += moveDelta.X
		bot.Position.Y += moveDelta.Y
	}
}
