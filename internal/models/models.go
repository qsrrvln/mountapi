package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Mountain struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"unique;not null" json:"name"`
	ElevationM int       `gorm:"not null" json:"elevation_m"`
	Province   string    `json:"province"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Routes     []Route   `gorm:"foreignKey:MountainID" json:"routes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Route struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	MountainID     uuid.UUID `gorm:"type:uuid;not null" json:"mountain_id"`
	Name           string    `gorm:"not null" json:"name"`
	StartPointName string    `json:"start_point_name"`
	Difficulty     string    `json:"difficulty"`
	LengthM        int       `json:"length_m"`
	ElevationGainM int       `json:"elevation_gain_m"`
	IsOfficial     bool      `gorm:"default:false" json:"is_official"`
	Posts          []Post    `gorm:"foreignKey:RouteID" json:"posts,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Post struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RouteID            uuid.UUID `gorm:"type:uuid;not null" json:"route_id"`
	Name               string    `gorm:"not null" json:"name"`
	OrderIndex         int       `gorm:"not null" json:"order_index"`
	AltitudeM          int       `json:"altitude_m"`
	DistanceFromStartM int       `json:"distance_from_start_m"`
	SegmentToNextM     int       `json:"segment_to_next_m"`
	DistanceToSummitM  int       `gorm:"not null" json:"distance_to_summit_m"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Mountain{}, &Route{}, &Post{})
}
