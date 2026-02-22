package config

import (
	"encoding/json"
	"os"
)

const configPath = "config.json"

type Config struct {
	AvoidanceFactor float64 `json:"avoidance_factor"`
	AlignmentFactor float64 `json:"alignment_factor"`
	GatheringFactor float64 `json:"gathering_factor"`
	AvoidanceRadius float64 `json:"avoidance_radius"`
	DetectionRadius float64 `json:"detection_radius"`
	MaxSpeed        float64 `json:"max_speed"`
}

func Default() Config {
	return Config{
		AvoidanceFactor: 1.0,
		AlignmentFactor: 0.005,
		GatheringFactor: 0.001,
		AvoidanceRadius: 20,
		DetectionRadius: 100,
		MaxSpeed:        3,
	}
}

func Load() Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Default()
	}
	cfg := Default()
	json.Unmarshal(data, &cfg)
	return cfg
}

func Save(cfg Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(configPath, data, 0644)
}
