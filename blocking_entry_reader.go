package ledge

import "fmt"

type blockingEntryReader struct {
	entryReader EntryReader
}

func newBlockingEntryReader(
	entryReader EntryReader,
) *blockingEntryReader {
	return &blockingEntryReader{
		entryReader,
	}
}

func (b *blockingEntryReader) Entries() ([]*Entry, error) {
	var entries []*Entry
	var errs []error
	entryC := b.entryReader.Channel()
	for {
		entryResponse, ok := <-entryC
		if !ok {
			break
		}
		if entryResponse.Entry != nil {
			entries = append(entries, entryResponse.Entry)
		}
		if entryResponse.Error != nil {
			errs = append(errs, entryResponse.Error)
		}
	}
	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("%v", errs)
	}
	return entries, err
}
