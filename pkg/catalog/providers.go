package catalog

import (
	"sort"
	"strings"

	wujicfg "github.com/coditary/wuji-core/pkg/config"
)

// MergeIntoProviderCatalog adds models.dev catalog providers not already present in summaries.
func MergeIntoProviderCatalog(root string, summaries map[string]wujicfg.ProviderSummary) {
	idx, err := LoadIndex(root)
	if err != nil {
		return
	}
	for id, p := range idx.Providers {
		if _, exists := summaries[id]; exists {
			continue
		}
		modelIDs := make([]string, 0, len(p.Models))
		for mid := range p.Models {
			modelIDs = append(modelIDs, mid)
		}
		sort.Strings(modelIDs)
		defaultModel := ""
		if len(modelIDs) > 0 {
			defaultModel = modelIDs[0]
		}
		summaries[id] = wujicfg.ProviderSummary{
			Name:   firstNonEmpty(strings.TrimSpace(p.Name), id),
			Model:  defaultModel,
			Driver: id,
			Kind:   "catalog",
			Models: modelIDs,
			API:    strings.TrimSpace(p.API),
			NPM:    strings.TrimSpace(p.NPM),
			Env:    append([]string(nil), p.Env...),
		}
	}
}
