package model

import "time"

type Package struct {
	Container_id string
	Live_until   time.Time
}
