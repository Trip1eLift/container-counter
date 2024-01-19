package container_counter_system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/broadcast"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/queue"
	"github.com/google/uuid"
)

type Manager struct {
	self_id       string
	mu_sleep      sync.RWMutex
	sleep_until   time.Time
	mu_containers sync.RWMutex
	containers    map[string]time.Time
	bc            *broadcast.Client
}

var manager Manager

func init() {
	manager.self_id = uuid.New().String()
	manager.sleep_until = time.Now()
	manager.containers = make(map[string]time.Time)
	manager.bc = broadcast.New(context.Background())
	manager.bc.Subscribe(context.Background())

	// clean up expired containers every 5s
	go func() {
		for {
			time.Sleep(5 * time.Second)
			manager.mu_containers.Lock()
			now := time.Now()
			for container_id, live_until := range manager.containers {
				if now.After(live_until) {
					delete(manager.containers, container_id)
				}
			}
			manager.mu_containers.Unlock()
		}
	}()

	// subscribe to queue
	go func() {
		for {
			pack := queue.Pop_package() // only return when new pack arrives
			fmt.Printf("pack recieved %v\n", pack)
			manager.mu_containers.Lock()
			manager.containers[pack.Container_id] = pack.Live_until
			manager.mu_containers.Unlock()
		}
	}()
}

// For every 30s, publish self to redis with 35s ttl
func OnTraffic() {
	manager.mu_sleep.RLock()
	if time.Now().Before(manager.sleep_until) {
		manager.mu_sleep.RUnlock()
		return
	} else {
		manager.mu_sleep.RUnlock()
	}

	manager.mu_sleep.Lock()
	manager.sleep_until = time.Now().Add(30 * time.Second) // now + 30s
	manager.mu_sleep.Unlock()

	pack := model.Package{
		Container_id: manager.self_id,
		Live_until:   time.Now().Add(35 * time.Second), // now + 30s
	}
	manager.bc.Publish(context.Background(), pack)
}

func GetCount() int {
	manager.mu_containers.RLock()
	count := len(manager.containers)
	manager.mu_containers.RUnlock()
	return count
}
