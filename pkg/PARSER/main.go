package PARSER

func FetchDomainEntries() ([]*DomainEntry, error) {
	return readEntries()
}
