package utilities

import guuid "github.com/google/uuid"

func RandomUUID() string {
	return guuid.NewString()
}
