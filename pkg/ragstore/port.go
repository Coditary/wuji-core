package ragstore

import "context"

// ExportCollection copies a collection directory to destPath.
func ExportCollection(ctx context.Context, root, collection, destPath string) error {
	return DefaultBackend.ExportCollection(ctx, root, collection, destPath)
}

// ImportCollection copies a collection snapshot into the store.
func ImportCollection(ctx context.Context, root, collection, srcPath string) error {
	return DefaultBackend.ImportCollection(ctx, root, collection, srcPath)
}

// RenameCollection moves a collection directory to a new name.
func RenameCollection(ctx context.Context, root, collection, newName string) error {
	return DefaultBackend.RenameCollection(ctx, root, collection, newName)
}
