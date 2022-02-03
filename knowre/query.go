package knowre

import (
	"encoding/json"
	"time"
)

func UnmarshalESQuery(data []byte) (ESQuery, error) {
	var r ESQuery
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ESQuery) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ESQuery struct {
	Query Query  `json:"query"`
	Sort  []Sort `json:"sort"`
}

type Query struct {
	Bool Bool `json:"bool"`
}

type Bool struct {
	Filter Filter `json:"filter"`
	Must   []Must `json:"must"`
}

type Filter struct {
	Range Range `json:"range"`
}

type Range struct {
	Timestamp Timestamp `json:"@timestamp"`
}

type Timestamp struct {
	Gte      string `json:"gte"`
	Lt       string `json:"lt"`
	TimeZone string `json:"time_zone"`
}

type Must struct {
	MatchPhrase MatchPhrase `json:"match_phrase"`
}

type MatchPhrase struct {
	LogGroup                    *string `json:"@log_group,omitempty"`
	LogType                     *string `json:"logType,omitempty"`
	KnowreDaekyoServerLogUserID *int64  `json:"knowre-daekyo.serverLog.user_id,omitempty"`
	KnowreDaekyoServerLogResult *bool   `json:"knowre-daekyo.serverLog.result,omitempty"`
}

type Sort struct {
	Timestamp string `json:"@timestamp"`
}

func MakeQuery(fromDate time.Time, match string, userID int, i int, duration int) (string, *ESQuery, error) {
	const layout = "2006-01-02T15:04:05.000Z"
	gte := fromDate.Add(time.Duration(i*duration) * time.Minute)
	lt := gte.Add(time.Duration(duration)*time.Minute - 1*time.Millisecond)

	logType := "formattedLog"
	logResult := true

	must := []Must{
		{MatchPhrase{
			LogGroup: &match,
		}},
		{MatchPhrase{
			LogType: &logType,
		}},
		{MatchPhrase{
			KnowreDaekyoServerLogResult: &logResult,
		}},
	}

	if userID > 0 {
		uid := int64(userID)
		must = append(must, Must{MatchPhrase{KnowreDaekyoServerLogUserID: &uid}})
	}

	boolVar := Bool{
		Filter: Filter{
			Range: Range{Timestamp{
				Gte:      gte.Format(layout),
				Lt:       lt.Format(layout),
				TimeZone: "+09:00",
			}},
		},
		Must: must,
	}

	query := &ESQuery{
		Query: Query{
			Bool: boolVar,
		},
		Sort: []Sort{{Timestamp: "asc"}},
	}

	q, err := query.Marshal()
	return string(q), query, err
}
