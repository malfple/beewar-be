package mapmanager

import "gitlab.com/beewar/beewar-be/internal/access"

// CreateEmptyMap creates an empty map of fixed size and name
func CreateEmptyMap(userID uint64) uint64 {
	mapID, err := access.CreateEmptyMap(0, 10, 10, "Untitled", userID)
	if err != nil {
		panic("unexpected error when creating empty map")
	}
	return mapID
}
