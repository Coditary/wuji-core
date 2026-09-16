package catalog

import (
	"path/filepath"
)

// Dir returns the catalog cache directory under the Wuji project root.
func Dir(root string) string {
	return filepath.Join(root, ".wuji", "catalog")
}

// IndexPath returns the slim index file path.
func IndexPath(root string) string {
	return filepath.Join(Dir(root), "index.json")
}

// MetaPath returns the sync metadata file path.
func MetaPath(root string) string {
	return filepath.Join(Dir(root), "meta.json")
}
