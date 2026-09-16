package ragstore

import "context"

// Backend persists and retrieves RAG collections (vector store abstraction).
type Backend interface {
	Exists(ctx context.Context, root, collection string) bool
	Load(ctx context.Context, root, collection string) (*Collection, error)
	Save(ctx context.Context, root string, m Manifest, chunks []Chunk, vectors []float32) error
	Delete(ctx context.Context, root, collection string) error
	ListCollections(ctx context.Context, root string) ([]Manifest, error)
	ExportCollection(ctx context.Context, root, collection, destPath string) error
	ImportCollection(ctx context.Context, root, collection, srcPath string) error
	RenameCollection(ctx context.Context, root, collection, newName string) error
}

type backendKey struct{}

// WithBackend attaches a store backend to the context for RAG engine calls.
func WithBackend(ctx context.Context, b Backend) context.Context {
	if b == nil {
		return ctx
	}
	return context.WithValue(ctx, backendKey{}, b)
}

// BackendFrom returns the backend from context, if any.
func BackendFrom(ctx context.Context) (Backend, bool) {
	b, ok := ctx.Value(backendKey{}).(Backend)
	return b, ok && b != nil
}

// ResolveBackend returns the context backend or the default filesystem store.
func ResolveBackend(ctx context.Context) Backend {
	if b, ok := BackendFrom(ctx); ok {
		return b
	}
	return DefaultBackend
}

// DefaultBackend is the built-in filesystem vector store under .wuji/rag/.
var DefaultBackend Backend = FilesystemBackend{}
