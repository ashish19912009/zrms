package permissions

import "sync"

type Permission struct {
	Resource string
	Action   string
}

var (
	mu             sync.RWMutex
	methodPermsMap = make(map[string]Permission)
)

func Register(fullMethod, resource, action string) {
	mu.Lock()
	defer mu.Unlock()
	methodPermsMap[fullMethod] = Permission{
		Resource: resource,
		Action:   action,
	}
}

func Get(fullMethod string) (Permission, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := methodPermsMap[fullMethod]
	return p, ok
}
