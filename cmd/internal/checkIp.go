package internal

import (
	"fmt"
	"sync"
	"time"
)

const MAXIMUM_REQUESTS = 3
const BLOCKED_TIME = 24 * time.Hour

type IpList struct {
	mu     sync.Mutex
	IpAddr map[string]*IpAddr
}

type IpAddr struct {
	Address         string
	Count           int
	TimeCount       []time.Time
	Blocked         bool
	Expire          time.Time
	LastInteraction time.Time
}

func CreateIpList() *IpList {
	return &IpList{
		IpAddr: make(map[string]*IpAddr),
	}
}

func (i *IpList) AddAndCheckIP(IpAddress string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	addr, exists := i.IpAddr[IpAddress]
	if !exists {
		i.IpAddr[IpAddress] = &IpAddr{
			Count: 1,
			TimeCount: []time.Time{
				time.Now(),
			},
			LastInteraction: time.Now(),
		}
		return nil
	}

	if addr.Blocked {
		if time.Now().Before(addr.Expire) {
			return fmt.Errorf("Endereço bloqueado por rate limits")
		}
		addr.Blocked = false
		addr.Count = 1
		addr.LastInteraction = time.Now()
		return nil
	}

	addr.Count++
	addr.TimeCount = append(addr.TimeCount, time.Now())
	if addr.Count > MAXIMUM_REQUESTS {
		// Se a diferença entre o primeiro e ultimo registro for menor que 5 segundos bloqueia
		if addr.TimeCount[len(addr.TimeCount)-1].Sub(addr.TimeCount[0]).Seconds() < 5 {
			addr.Blocked = true
			addr.Expire = time.Now().Add(BLOCKED_TIME)
			addr.LastInteraction = time.Now()
			return fmt.Errorf("Endereço bloqueado por time limits")
		}
	}

	return nil
}
