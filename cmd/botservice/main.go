package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/tcotav/travelbot/bot"
	loc "github.com/tcotav/travelbot/location"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// route rules
// random size route chosen, ranging from 2 to 5 locations + home
func getRandomRoute(rmap *map[string]loc.Location) []loc.Location {
	route := make([]loc.Location, 0)
	route = append(route, (*rmap)["Hom"])

	// get list of all the keys
	keys := make([]string, 0)
	for k := range *rmap {
		if k != "Hom" {
			keys = append(keys, k)
		}
	}

	// get how many keys we need
	numKeys := rand.Intn(2) + 3
	for i := 0; i < int(numKeys); i++ {
		// get a random key
		idx := rand.Intn(len(keys))
		randKey := keys[idx]
		// add it to the route
		route = append(route, (*rmap)[randKey])
		// remove it from the list of keys
		keys = append(keys[:idx], keys[idx+1:]...)
	}

	return route

}

//go:embed files/routemap.json
var jsonData []byte

func main() {
	// handle command line flags for:
	// - ship name
	// - ship speed
	// - move sleep time

	// get the ship name
	shipName := flag.String("ship", "TravelBot", "The name of the ship")

	// get the ship speed
	shipSpeed := flag.String("speed", "2.0", "The speed of the ship")

	// get the sleep time
	sleepTime := flag.String("sleep", "2", "The time to sleep between moves")

	flag.Parse()

	routeMap := make(map[string]loc.Location)

	err := json.Unmarshal(jsonData, &routeMap)
	if err != nil {
		log.Fatalf("unmarshall failed, %v", err)
	}
	stopMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "spaceship_total_stops_made",
		Help: "The total number of stops made by the ship",
	}, []string{"ship_name"})

	distanceMetric := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "spaceship_total_distance_traveled",
		Help: "The total distance traversed by the ship",
	}, []string{"ship_name"})

	floatSpeed, err := strconv.ParseFloat(*shipSpeed, 64)
	if err != nil {
		log.Fatalf("float convert of %s failed, %s", *shipSpeed, err.Error())
	}

	floatSleep, err := strconv.Atoi(*sleepTime)
	if err != nil {
		log.Fatalf("Int convert of %s failed, %s", *sleepTime, err.Error())
	}

	b := bot.NewTravelBot(
		*shipName,
		floatSpeed,
		loc.Point{X: 0, Y: 0},
		bot.NewRouteList(getRandomRoute(&routeMap)),
		routeMap,
		stopMetric,
		distanceMetric,
	)

	mainloop(b, floatSleep)

}

func mainloop(b *bot.TravelBot, sleepTime int) {
	for {
		log.Printf("Ship %s is at %v and heading to %v, %d\n",
			b.Name,
			b.Position,
			b.Destination,
			b.Route.CurrentLoc,
		)
		b.Move()
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
