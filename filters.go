package ledge

type requireContextFilter struct {
	context Context
}

func newRequireContextFilter(
	context Context,
) *requireContextFilter {
	return &requireContextFilter{
		context,
	}
}

func (r *requireContextFilter) Include(entry *Entry) bool {
	for _, context := range entry.Contexts {
		if context == r.context {
			return true
		}
	}
	return false
}

type levelFilter struct {
	level Level
}

func newLevelFilter(
	level Level,
) *levelFilter {
	return &levelFilter{
		level,
	}
}

func (l *levelFilter) Include(entry *Entry) bool {
	return l.level <= entry.Level
}
