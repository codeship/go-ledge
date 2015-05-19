package ledge

type entryWriter struct {
	baseEntry *Entry
	logger    *logger
}

func newEntryWriter(
	baseEntry *Entry,
	logger *logger,
) *entryWriter {
	return &entryWriter{
		baseEntry,
		logger,
	}
}

func (e *entryWriter) Write(p []byte) (int, error) {
	// TODO(pedge): what if level is Level_FATAL or Level_PANIC?
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	if _, err := e.logger.write(e.logger.getEntry(e.baseEntry, p)); err != nil {
		return 0, err
	}
	return len(p), nil
}
