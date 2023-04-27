package bot

import (
	"testing"

	loc "github.com/tcotav/travelbot/location"
)

var routeMap map[string]loc.Location

func TestNewRouteList(t *testing.T) {
	l := []loc.Location{routeMap["Hom"], routeMap["Erothol"], routeMap["Tyros"], routeMap["Eliatri"]}
	rl := NewRouteList(l)
	if rl.CurrentLoc != 0 {
		t.Errorf("Expected CurrentLoc to be 0, got %d", rl.CurrentLoc)
	}
	if rl.Locations[1].Name != "Tyros" {
		t.Errorf("Expected second location to be Tyros, got %s", rl.Locations[1].Name)
	}

}

func TestMain(m *testing.M) {
	routeMap = make(map[string]loc.Location)
	routeMap["Hom"] = loc.Location{
		Name: "Hom",
		Position: loc.Point{
			X: 0,
			Y: 0,
		},
	}
	routeMap["Erothol"] = loc.Location{
		Name: "Erothol",
		Position: loc.Point{
			X: 10.0,
			Y: 20.1,
		},
	}
	routeMap["Muk"] = loc.Location{
		Name: "Muk",
		Position: loc.Point{
			X: 22.0,
			Y: 14.1,
		},
	}
	routeMap["Tyros"] = loc.Location{
		Name: "Tyros",
		Position: loc.Point{
			X: 5.0,
			Y: 4.1,
		},
	}
	routeMap["Eliatri"] = loc.Location{
		Name: "Eliatri",
		Position: loc.Point{
			X: 19.0,
			Y: 21.1,
		},
	}
}
