package local

// LegacyDriverIDs maps old driver IDs to local-rag.
var LegacyDriverIDs = map[string]string{
	"local": DriverID,
	"rag":   DriverID,
}

// NormalizeDriverID resolves legacy aliases to local-rag.
func NormalizeDriverID(id string) string {
	if mapped, ok := LegacyDriverIDs[id]; ok {
		return mapped
	}
	return id
}
