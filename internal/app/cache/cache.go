package cache

import (
	"encoding/json"
)

type Cache interface {
	AddToStore(key string, value json.RawMessage)
	Get(key string) (json.RawMessage, error)
}
