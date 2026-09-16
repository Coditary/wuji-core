package catalog

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
)

// BuildIndex converts a models.dev catalog.json payload into a slim index.
func BuildIndex(data []byte, source string, syncedAt time.Time) (*Index, *Meta, error) {
	raw, err := decodeRawCatalog(data)
	if err != nil {
		return nil, nil, err
	}
	if syncedAt.IsZero() {
		syncedAt = time.Now().UTC()
	}
	if strings.TrimSpace(source) == "" {
		source = DefaultSource
	}

	idx := &Index{
		Version:   IndexVersion,
		Source:    source,
		SyncedAt:  syncedAt.UTC().Format(time.RFC3339),
		Providers: make(map[string]ProviderEntry, len(raw.Providers)),
		Models:    make(map[string]ModelEntry, len(raw.Models)),
		Labs:      map[string]LabEntry{},
	}

	providerModels := 0
	for pid, p := range raw.Providers {
		id := strings.TrimSpace(p.ID)
		if id == "" {
			id = pid
		}
		models := map[string]string{}
		for mid, m := range p.Models {
			name := strings.TrimSpace(m.Name)
			if name == "" {
				name = mid
			}
			models[mid] = name
			providerModels++
		}
		idx.Providers[id] = ProviderEntry{
			ID:     id,
			Name:   firstNonEmpty(strings.TrimSpace(p.Name), id),
			API:    strings.TrimSpace(p.API),
			NPM:    strings.TrimSpace(p.NPM),
			Doc:    strings.TrimSpace(p.Doc),
			Env:    append([]string(nil), p.Env...),
			Models: models,
		}
	}

	for mid, m := range raw.Models {
		id := strings.TrimSpace(m.ID)
		if id == "" {
			id = mid
		}
		lab := labFromModelID(id)
		entry := ModelEntry{
			ID:        id,
			Name:      firstNonEmpty(strings.TrimSpace(m.Name), id),
			Lab:       lab,
			Family:    strings.TrimSpace(m.Family),
			Reasoning: m.Reasoning,
			ToolCall:  m.ToolCall,
			Context:   m.Limit.Context,
			Output:    m.Limit.Output,
			Open:      m.OpenWeights,
		}
		idx.Models[id] = entry

		labEntry := idx.Labs[lab]
		if labEntry.ID == "" {
			labEntry = LabEntry{ID: lab, Name: formatLabName(lab)}
		}
		labEntry.Models = append(labEntry.Models, id)
		idx.Labs[lab] = labEntry
	}

	for lab, entry := range idx.Labs {
		sort.Strings(entry.Models)
		idx.Labs[lab] = entry
	}

	meta := &Meta{
		Source:         source,
		SyncedAt:       idx.SyncedAt,
		ProviderCount:  len(idx.Providers),
		ModelCount:     len(idx.Models),
		LabCount:       len(idx.Labs),
		ProviderModels: providerModels,
	}
	return idx, meta, nil
}

// BuildIndexFromRaw is a helper for tests.
func BuildIndexFromRaw(raw *rawCatalog, source string, syncedAt time.Time) (*Index, *Meta, error) {
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal test catalog: %w", err)
	}
	return BuildIndex(data, source, syncedAt)
}

func labFromModelID(id string) string {
	if i := strings.Index(id, "/"); i > 0 {
		return id[:i]
	}
	return id
}

func formatLabName(lab string) string {
	if lab == "" {
		return ""
	}
	parts := strings.Split(lab, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
