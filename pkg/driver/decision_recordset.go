package driver

import (
	"fmt"

	"github.com/coditary/wuji-core/pkg/data"
)

// DecisionResponseToRecordSet converts a decision response into a WDD recordset.
func DecisionResponseToRecordSet(resp *DecisionResponse, meta data.Meta) (*data.RecordSet, error) {
	if resp == nil {
		return nil, fmt.Errorf("empty decision response")
	}
	rs := &data.RecordSet{
		MetaData: meta,
		Schema: []data.FieldDef{
			{Name: "question_id", Type: data.KindString},
			{Name: "type", Type: data.KindString},
			{Name: "choice", Type: data.KindString},
			{Name: "score", Type: data.KindFloat},
			{Name: "noul", Type: data.KindFloat},
			{Name: "confidence", Type: data.KindFloat},
			{Name: "probabilities", Type: data.KindObject},
		},
		Records: make([]data.Record, 0, len(resp.Answers)),
	}
	for _, ans := range resp.Answers {
		fields := map[string]data.Value{
			"question_id": data.NewString(ans.QuestionID),
			"type":        data.NewString(string(ans.Type)),
			"confidence":  data.NewFloat(ans.Confidence),
		}
		if ans.Choice != "" {
			fields["choice"] = data.NewString(ans.Choice)
		}
		if ans.Score != 0 {
			fields["score"] = data.NewFloat(ans.Score)
		}
		if ans.Noul != 0 {
			fields["noul"] = data.NewFloat(ans.Noul)
		}
		if len(ans.Probabilities) > 0 {
			obj := make(map[string]data.Value, len(ans.Probabilities))
			for k, v := range ans.Probabilities {
				obj[k] = data.NewFloat(v)
			}
			fields["probabilities"] = data.Value{Kind: data.KindObject, O: obj}
		}
		rs.Records = append(rs.Records, data.Record{
			ID:     ans.QuestionID,
			Role:   "decision",
			Fields: fields,
		})
	}
	return rs, nil
}
