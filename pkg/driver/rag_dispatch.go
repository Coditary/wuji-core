package driver

import (
	"context"
	"fmt"

	"github.com/coditary/wuji-core/pkg/capability"
	"github.com/coditary/wuji-core/pkg/data"
)

// RunRAGTask dispatches a RAG request to driver interfaces.
func RunRAGTask(ctx context.Context, d Driver, req RAGRequest) (data.DataShape, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.Task != RAGTaskAnswer {
		if prod, err := As[RAGProducer](d, capability.RAG); err == nil {
			return prod.ProduceRAG(ctx, req)
		}
	}
	switch req.Task {
	case RAGTaskIndex:
		indexer, err := As[RAGIndexer](d, capability.RAG)
		if err != nil {
			return nil, err
		}
		resp, err := indexer.Index(ctx, req)
		if err != nil {
			return nil, err
		}
		return RAGIndexToRecordSet(resp), nil
	case RAGTaskQuery:
		querier, err := As[RAGQuerier](d, capability.RAG)
		if err != nil {
			return nil, err
		}
		resp, err := querier.Query(ctx, req)
		if err != nil {
			return nil, err
		}
		return RAGQueryToRecordSetOpts(resp, req.IncludeScores), nil
	case RAGTaskAnswer:
		return nil, fmt.Errorf("answer task must be handled by core.RunRAGAnswer")
	case RAGTaskList, RAGTaskInfo, RAGTaskStats, RAGTaskDelete, RAGTaskPurge, RAGTaskExport, RAGTaskImport, RAGTaskRename:
		mgr, err := As[RAGCollectionManager](d, capability.RAG)
		if err != nil {
			return nil, err
		}
		resp, err := runRAGManage(ctx, mgr, req)
		if err != nil {
			return nil, err
		}
		return RAGManageToRecordSet(resp), nil
	default:
		return nil, fmt.Errorf("unsupported rag task %q", req.Task)
	}
}

func runRAGManage(ctx context.Context, mgr RAGCollectionManager, req RAGRequest) (*RAGManageResponse, error) {
	switch req.Task {
	case RAGTaskList:
		items, err := mgr.ListCollections(ctx, req)
		if err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Collections: items, Message: "collections listed"}, nil
	case RAGTaskInfo:
		info, err := mgr.CollectionInfo(ctx, req)
		if err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Collection: info, Message: "collection info"}, nil
	case RAGTaskStats:
		info, err := mgr.CollectionInfo(ctx, req)
		if err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Collection: info, Message: "collection stats"}, nil
	case RAGTaskExport:
		if err := mgr.ExportCollection(ctx, req); err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Message: fmt.Sprintf("exported collection %q to %s", req.Collection, req.ExportPath)}, nil
	case RAGTaskImport:
		if err := mgr.ImportCollection(ctx, req); err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Message: fmt.Sprintf("imported collection %q from %s", req.Collection, req.ImportPath)}, nil
	case RAGTaskRename:
		if err := mgr.RenameCollection(ctx, req); err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Message: fmt.Sprintf("renamed collection %q to %q", req.Collection, req.RenameTo)}, nil
	case RAGTaskDelete:
		if err := mgr.DeleteCollection(ctx, req); err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Message: fmt.Sprintf("deleted collection %q", req.Collection)}, nil
	case RAGTaskPurge:
		if err := mgr.PurgeSource(ctx, req); err != nil {
			return nil, err
		}
		return &RAGManageResponse{Task: req.Task, Message: fmt.Sprintf("purged source %q", req.PurgeSource)}, nil
	default:
		return nil, fmt.Errorf("unsupported manage task %q", req.Task)
	}
}
