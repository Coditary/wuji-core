package driver

import "fmt"

// RAGTask identifies a RAG operation.
type RAGTask string

const (
	RAGTaskIndex  RAGTask = "index"
	RAGTaskQuery  RAGTask = "query"
	RAGTaskAnswer RAGTask = "answer"
	RAGTaskList   RAGTask = "list"
	RAGTaskInfo   RAGTask = "info"
	RAGTaskStats  RAGTask = "stats"
	RAGTaskDelete RAGTask = "delete"
	RAGTaskPurge  RAGTask = "purge"
	RAGTaskExport RAGTask = "export"
	RAGTaskImport RAGTask = "import"
	RAGTaskRename RAGTask = "rename"
)

func AllRAGTasks() []RAGTask {
	return []RAGTask{
		RAGTaskIndex, RAGTaskQuery, RAGTaskAnswer,
		RAGTaskList, RAGTaskInfo, RAGTaskStats, RAGTaskDelete, RAGTaskPurge,
		RAGTaskExport, RAGTaskImport, RAGTaskRename,
	}
}

func ParseRAGTask(name string) (RAGTask, error) {
	for _, t := range AllRAGTasks() {
		if string(t) == name {
			return t, nil
		}
	}
	return "", fmt.Errorf("unknown rag task %q", name)
}
