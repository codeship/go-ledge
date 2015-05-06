package ledge

import "time"

var (
	systemTimerInstance = &systemTimer{}
)

type systemTimer struct{}

func (s *systemTimer) Now() time.Time {
	return time.Now().UTC()
}
