package driver

import (
	"fmt"
	"strconv"

	"github.com/coditary/wuji-core/pkg/data"
)

func RAGQueryToRecordSet(resp *RAGQueryResponse) *data.RecordSet {
	return RAGQueryToRecordSetOpts(resp, true)
}

func RAGQueryToRecordSetOpts(resp *RAGQueryResponse, includeScores bool) *data.RecordSet {
	if resp == nil {
		return &data.RecordSet{}
	}
	rs := &data.RecordSet{
		MetaData: data.Meta{
			Capability: "rag",
			Task:       string(RAGTaskQuery),
			Source:     map[string]string{"collection": resp.Collection, "query": resp.Query},
		},
		Records: hitsToRecords(resp.Hits, includeScores),
	}
	return rs
}

func RAGIndexToRecordSet(resp *RAGIndexResponse) *data.RecordSet {
	if resp == nil {
		return &data.RecordSet{}
	}
	fields := map[string]data.Value{
		"collection":    data.NewString(resp.Collection),
		"chunks_added":  data.NewInt(int64(resp.ChunksAdded)),
		"total_chunks":  data.NewInt(int64(resp.TotalChunks)),
		"bytes_indexed": data.NewInt(resp.BytesIndexed),
		"message":       data.NewString(resp.Message),
		"dry_run":       data.NewBool(resp.DryRun),
	}
	return &data.RecordSet{
		MetaData: data.Meta{Capability: "rag", Task: string(RAGTaskIndex)},
		Records: []data.Record{{
			ID: "index_result", Role: "result", Fields: fields,
		}},
	}
}

func RAGAnswerToRecordSet(resp *RAGAnswerResponse) *data.RecordSet {
	var rs *data.RecordSet
	if resp.NoContext {
		rs = &data.RecordSet{
			MetaData: data.Meta{
				Capability: "rag",
				Task:       string(RAGTaskAnswer),
				Source:     map[string]string{"collection": resp.Collection, "query": resp.Query},
			},
		}
	} else {
		rs = RAGQueryToRecordSet(&RAGQueryResponse{
			Query: resp.Query, Collection: resp.Collection, Hits: resp.Hits,
		})
		rs.MetaData.Task = string(RAGTaskAnswer)
	}
	if resp.Answer != "" {
		answerRec := data.Record{
			ID: "answer", Role: "answer",
			Fields: map[string]data.Value{"text": data.NewString(resp.Answer)},
		}
		rs.Records = append([]data.Record{answerRec}, rs.Records...)
	}
	if resp.TokensUsed > 0 {
		rs.MetaData.Extra = map[string]string{"tokens_used": strconv.Itoa(resp.TokensUsed)}
	}
	return rs
}

func RAGManageToRecordSet(resp *RAGManageResponse) *data.RecordSet {
	if resp == nil {
		return &data.RecordSet{}
	}
	rs := &data.RecordSet{
		MetaData: data.Meta{Capability: "rag", Task: string(resp.Task)},
	}
	if resp.Message != "" {
		rs.Records = append(rs.Records, data.Record{
			ID: "message", Role: "result", Fields: map[string]data.Value{"message": data.NewString(resp.Message)},
		})
	}
	for i, c := range resp.Collections {
		rs.Records = append(rs.Records, collectionRecord(c, i))
	}
	if resp.Collection.Collection != "" {
		rs.Records = append(rs.Records, collectionRecord(resp.Collection, 0))
	}
	return rs
}

func collectionRecord(c RAGCollectionInfo, seq int) data.Record {
	fields := map[string]data.Value{
		"collection":  data.NewString(c.Collection),
		"embed_model": data.NewString(c.EmbedModel),
		"dims":        data.NewInt(int64(c.Dims)),
		"chunk_count": data.NewInt(int64(c.ChunkCount)),
	}
	if c.ChunkSize > 0 {
		fields["chunk_size"] = data.NewInt(int64(c.ChunkSize))
	}
	if c.Overlap > 0 {
		fields["overlap"] = data.NewInt(int64(c.Overlap))
	}
	if c.UpdatedAt != "" {
		fields["updated_at"] = data.NewString(c.UpdatedAt)
	}
	return data.Record{
		ID: c.Collection, Role: "collection", Seq: seq, Fields: fields,
	}
}

// RAGQueryFromShape extracts query hits from a WDD record set.
func RAGQueryFromShape(shape data.DataShape) (*RAGQueryResponse, error) {
	rs, ok := shape.(*data.RecordSet)
	if !ok || rs == nil {
		return nil, fmt.Errorf("expected recordset from rag query, got %T", shape)
	}
	resp := &RAGQueryResponse{
		Query:      rs.MetaData.Source["query"],
		Collection: rs.MetaData.Source["collection"],
	}
	for _, rec := range rs.Records {
		if rec.Role != "chunk" {
			continue
		}
		hit := RAGHit{ID: rec.ID}
		if v, ok := rec.Fields["score"]; ok {
			hit.Score = float32(v.F)
		}
		if v, ok := rec.Fields["text"]; ok {
			hit.Text = v.S
		}
		if v, ok := rec.Fields["source"]; ok {
			hit.Source = v.S
		}
		if v, ok := rec.Fields["chunk_index"]; ok {
			hit.ChunkIdx = int(v.I)
		}
		resp.Hits = append(resp.Hits, hit)
	}
	return resp, nil
}

func hitsToRecords(hits []RAGHit, includeScores bool) []data.Record {
	recs := make([]data.Record, len(hits))
	for i, hit := range hits {
		fields := map[string]data.Value{
			"text":   data.NewString(hit.Text),
			"source": data.NewString(hit.Source),
		}
		if includeScores {
			fields["score"] = data.NewFloat(float64(hit.Score))
		}
		if hit.ChunkIdx > 0 {
			fields["chunk_index"] = data.NewInt(int64(hit.ChunkIdx))
		}
		recs[i] = data.Record{
			ID: hit.ID, Role: "chunk", Seq: i, Fields: fields,
		}
	}
	return recs
}
