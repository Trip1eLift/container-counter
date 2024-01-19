package queue

import (
	"sync"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
)

type package_queue struct {
	boxes     []model.Package
	mu        sync.RWMutex
	mu_level2 sync.RWMutex
}

var Packages package_queue

func init() {
	Packages.boxes = make([]model.Package, 0)
}

func Push_package(pack model.Package) {
	Packages.mu.Lock()
	Packages.boxes = append(Packages.boxes, pack)
	Packages.mu.Unlock()
}

func Pop_package() model.Package {
	// prevent Pop_package being called multiple times at once
	// Without level2 lock, when 2 Pop_package are waiting and a package arrive,
	// they can rush to stage 2 together and break the code.
	Packages.mu_level2.Lock()

	// Stage 1: Waiting room
	for {
		Packages.mu.RLock()
		if len(Packages.boxes) > 0 {
			Packages.mu.RUnlock()
			break
		} else {
			Packages.mu.RUnlock()
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Stage 2: Pop queue's head
	Packages.mu.Lock()
	pack := Packages.boxes[0]
	Packages.boxes = Packages.boxes[1:]
	Packages.mu.Unlock()
	Packages.mu_level2.Unlock()

	return pack
}
