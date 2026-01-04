package main

import (
	"log"

	"github.com/qsrrvln/mountapi/internal/config"
	"github.com/qsrrvln/mountapi/internal/database"
	"github.com/qsrrvln/mountapi/internal/models"
)

func main() {
	// 1. Load Config & Connect
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Define Data
	// Mount Guntur
	guntur := models.Mountain{
		Name:       "Mount Guntur",
		ElevationM: 2249,
		Province:   "Jawa Barat",
		Latitude:   -7.140,
		Longitude:  107.828,
	}

	// Create Mountain to get its ID
	if err := db.FirstOrCreate(&guntur, models.Mountain{Name: guntur.Name}).Error; err != nil {
		log.Fatalf("Failed to seed mountain: %v", err)
	}
	log.Printf("Seeded Mountain: %s (ID: %s)", guntur.Name, guntur.ID)

	// Route Via Citiis
	citiis := models.Route{
		MountainID:     guntur.ID,
		Name:           "Via Citiis",
		StartPointName: "Basecamp Citiis",
		Difficulty:     "Sedang",
		LengthM:        5000,
		ElevationGainM: 1400,
		IsOfficial:     true,
	}

	if err := db.Where("mountain_id = ? AND name = ?", guntur.ID, citiis.Name).FirstOrCreate(&citiis).Error; err != nil {
		log.Fatalf("Failed to seed route: %v", err)
	}
	log.Printf("Seeded Route: %s (ID: %s)", citiis.Name, citiis.ID)

	// Posts
	posts := []models.Post{
		{
			RouteID:            citiis.ID,
			Name:               "Basecamp Citiis",
			OrderIndex:         0,
			AltitudeM:          850,
			DistanceFromStartM: 0,
			SegmentToNextM:     1500,
			DistanceToSummitM:  5000,
		},
		{
			RouteID:            citiis.ID,
			Name:               "Pos 1",
			OrderIndex:         1,
			AltitudeM:          1100,
			DistanceFromStartM: 1500,
			SegmentToNextM:     1500,
			DistanceToSummitM:  3500,
		},
		{
			RouteID:            citiis.ID,
			Name:               "Pos 2",
			OrderIndex:         2,
			AltitudeM:          1500,
			DistanceFromStartM: 3000,
			SegmentToNextM:     1000,
			DistanceToSummitM:  2000,
		},
		{
			RouteID:            citiis.ID,
			Name:               "Pos 3",
			OrderIndex:         3,
			AltitudeM:          1900,
			DistanceFromStartM: 4000,
			SegmentToNextM:     1000,
			DistanceToSummitM:  1000,
		},
		{
			RouteID:            citiis.ID,
			Name:               "Puncak 1",
			OrderIndex:         4,
			AltitudeM:          2200,
			DistanceFromStartM: 5000,
			SegmentToNextM:     0,
			DistanceToSummitM:  0,
		},
	}

	for _, p := range posts {
		// Check strictly by RouteID and OrderIndex to avoid duplicates
		var existing models.Post
		if err := db.Where("route_id = ? AND order_index = ?", p.RouteID, p.OrderIndex).First(&existing).Error; err != nil {
			// Not found, create it
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Failed to create post %s: %v", p.Name, err)
			} else {
				log.Printf("Created Post: %s", p.Name)
			}
		} else {
			log.Printf("Post already exists: %s", p.Name)
		}
	}
}
