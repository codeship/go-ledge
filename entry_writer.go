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
	// TODO(pedge): what if level is LevelFatal or LevelPanic?
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	return e.logger.write(e.logger.getEntry(e.baseEntry, p))
}
