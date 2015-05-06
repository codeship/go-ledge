package ledge

import "github.com/satori/go.uuid"

var (
	uuidAllocatorInstance = &uuidAllocator{}
)

type uuidAllocator struct{}

func (u *uuidAllocator) Allocate() string {
	return uuid.NewV4().String()
}
