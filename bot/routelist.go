package bot

import (
	"math"

	loc "github.com/tcotav/travelbot/location"
)

type RouteList struct {
	Locations  []loc.Location
	CurrentLoc int
	Traverse   int // 1 = forward, -1 = backward
}

func NewRouteList(locations []loc.Location) *RouteList {
	locations = findShortestPath(locations)
	return &RouteList{
		Locations:  locations,
		CurrentLoc: 0,
		Traverse:   1,
	}
}

// insert a new destination into the route list at a specified index
func (r *RouteList) Append(l loc.Location) {
	r.Locations = findShortestPath(append(r.Locations, l))
}

// get the next destination in the route list
// assume we return back along the same route but in reverse
func (r *RouteList) Next() loc.Location {
	nextLoc := r.CurrentLoc + r.Traverse
	// switch directions in traversal if we're at the start or end of the list
	if nextLoc > len(r.Locations)-1 || nextLoc < 0 {
		// change direction and try again
		r.Traverse *= -1
		nextLoc = r.CurrentLoc + r.Traverse
	}
	r.CurrentLoc = nextLoc
	return r.Locations[r.CurrentLoc]
}

func distance(p1, p2 loc.Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func swap(locations []loc.Location, i, j int) {
	locations[i], locations[j] = locations[j], locations[i]
}

func allPermutations(locations []loc.Location, start int, result *[]loc.Location, minDist *float64) {
	if start == len(locations) {
		dist := 0.0
		for i := 1; i < len(locations); i++ {
			dist += distance(locations[i-1].Position, locations[i].Position)
		}

		if dist < *minDist {
			*minDist = dist
			*result = make([]loc.Location, len(locations))
			copy(*result, locations)
		}
		return
	}

	for i := start; i < len(locations); i++ {
		swap(locations, start, i)
		allPermutations(locations, start+1, result, minDist)
		swap(locations, start, i)
	}
}

func findShortestPath(l []loc.Location) []loc.Location {
	if len(l) <= 1 {
		return l
	}

	minDist := math.Inf(1)
	var result []loc.Location
	allPermutations(l, 0, &result, &minDist)

	return result
}
