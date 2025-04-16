package entities

import "time"

type ControllerConfig struct {
	ConfigName   string
	LastUpdate   time.Time
	ConfigBranch string
}
