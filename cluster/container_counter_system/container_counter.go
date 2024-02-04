package container_counter_system

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/broadcast"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/model"
	"github.com/Trip1eLift/container-counter/cluster/container_counter_system/queue"
)

const baseTime = 30 * time.Second

type Manager struct {
	self_id       string
	mu_sleep      sync.RWMutex
	live_until    time.Time
	mu_containers sync.RWMutex
	containers    map[string]time.Time
	bc            *broadcast.Client
}

var manager Manager

func init() {
	manager.self_id = os.Getenv("CONTAINER_ID")
	manager.live_until = time.Now().Add(-1 * time.Second)
	manager.containers = make(map[string]time.Time)
	manager.bc = broadcast.New(context.Background())
	manager.bc.Subscribe(context.Background())

	// clean up expired containers every 1s
	go func() {
		for {
			time.Sleep(1 * time.Second)
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
			// fmt.Printf("pack recieved %v\n", pack)

			// if dead() {
			// 	// When I am dead, I don't count.
			// 	continue
			// }
			// I keep counting when I am dead (reduce delay problem)

			if pack.Container_id == manager.self_id {
				// I do not retrieve my own id from redis.
				continue
			}

			manager.mu_containers.Lock()
			_, old_container := manager.containers[pack.Container_id] // If this is a new container, the noob doesn't know I am alive, so I need to re-enroll myself.
			manager.containers[pack.Container_id] = pack.Live_until
			manager.mu_containers.Unlock()

			if !old_container {
				publish_enrollment()
				// delay is unnecessary
				// go func() {
				// 	time.Sleep(1 * time.Second)
				// 	publish_enrollment()
				// 	// delay publish to noobs in case noobs are not ready to recieve
				// }()
			}
		}
	}()
}

// For every 30s, publish self to redis with 35s ttl
func OnTraffic() {
	// I am alive
	if dead() == false {
		return
	}

	now := time.Now()
	//aliveTime := baseTime + utils.RandTimeS(5) // 25s + 0-5s
	aliveTime := baseTime // rand alive time is unnecessary

	manager.mu_sleep.Lock()
	manager.live_until = now.Add(aliveTime) // now + aliveTime (spawn + 30s)
	manager.mu_sleep.Unlock()

	manager.mu_containers.Lock()
	manager.containers[manager.self_id] = now.Add(aliveTime + 5*time.Second) // now + aliveTime + 5s (spawn + 35s)
	manager.mu_containers.Unlock()

	publish_enrollment()
}

func dead() bool {
	manager.mu_sleep.RLock()
	dead := time.Now().After(manager.live_until)
	manager.mu_sleep.RUnlock()
	return dead
}

func publish_enrollment() {
	if dead() {
		return
	}

	// enroll self to other containers
	manager.mu_sleep.RLock()
	pack := model.Package{
		Container_id: manager.self_id,
		Live_until:   manager.live_until.Add(5 * time.Second), // spawn + 35s
	}
	manager.mu_sleep.RUnlock()
	manager.bc.Publish(context.Background(), pack)
}

func GetCount() int {
	keys := getMap()

	manager.mu_containers.RLock()
	count := len(manager.containers)
	manager.mu_containers.RUnlock()

	fmt.Printf("container map: %v count: %d\n", keys, count)
	return count
}

func getMap() []string {
	keys := []string{}
	manager.mu_containers.RLock()
	for k := range manager.containers {
		keys = append(keys, k)
	}
	manager.mu_containers.RUnlock()
	return keys
}
