package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/buger/goreplay/proto"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Empty 문자열이 {} 이거나 길이가 0이면 참을 반환한다.
func Empty(s string) bool {
	if s == "{}" || len(s) == 0 {
		return true
	}

	return false
}
func UnmarshalServerLogCookie(data []byte) (ServerLogCookie, error) {
	var r ServerLogCookie
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ServerLogCookie) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ServerLogCookie struct {
	Cookie Cookie `json:"cookie"`
}

type Cookie struct {
	ConnectSid string `json:"connect.sid"`
}

func (c Cookie) String() string {
	t := reflect.TypeOf(Cookie{})
	key, _ := t.FieldByName("ConnectSid")
	if len(c.ConnectSid) == 0 {
		return ""
	}
	return fmt.Sprintf("%s=%s", key.Tag.Get("json"), c.ConnectSid)
}

func UnmarshalElasticsearchDocument(data []byte) (ElasticsearchDocument, error) {
	var r ElasticsearchDocument
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ElasticsearchDocument) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//ES의 document구조 파싱을 위함
type ElasticsearchDocument struct {
	ID     string      `json:"_id"`
	Index  string      `json:"_index"`
	Score  interface{} `json:"_score"`
	Source Source      `json:"_source"`
	Type   string      `json:"_type"`
	Sort   []int64     `json:"sort"`
}

type Source struct {
	ID                 string       `json:"@id"`
	LogGroup           string       `json:"@log_group"`
	LogStream          string       `json:"@log_stream"`
	Message            string       `json:"@message"`
	Owner              string       `json:"@owner"`
	Timestamp          string       `json:"@timestamp"`
	Version            string       `json:"@version"`
	KafkaConsumerGroup string       `json:"kafka.consumer_group"`
	KafkaKey           string       `json:"kafka.key"`
	KafkaOffset        string       `json:"kafka.offset"`
	KafkaPartition     string       `json:"kafka.partition"`
	KafkaTimestamp     string       `json:"kafka.timestamp"`
	KafkaTopic         string       `json:"kafka.topic"`
	KnowreDaekyo       KnowreDaekyo `json:"knowre-daekyo"` //@message가 파싱된게 들어있다.
	LogType            string       `json:"logType"`
}

type KnowreDaekyo struct {
	ServerLog ServerLog `json:"serverLog"`
}

type ServerLog struct {
	AccessToken   string      `json:"accessToken,omitempty"`
	AppFlavor     string      `json:"appFlavor,omitempty"`
	AmazonTraceID string      `json:"amazonTraceId"`
	Body          string      `json:"body"`
	Cookie        string      `json:"cookie"`
	Error         interface{} `json:"error"`
	IP            string      `json:"ip"`
	Method        string      `json:"method"`
	Parameters    string      `json:"parameters"`
	Performance   float64     `json:"performance"`
	Req           interface{} `json:"req"`
	Result        bool        `json:"result"`
	Router        string      `json:"router"`
	Session       string      `json:"session"`
	SessionID     string      `json:"session_id"`
	Trace         string      `json:"trace"`
	Token         string      `json:"token,omitempty"`
	URL           string      `json:"url"`
	UserAgent     string      `json:"userAgent"`
	UserType      string      `json:"userType"`
	UserID        int64       `json:"user_id"`
}

type ElasticsearchMessage struct {
	ReqURL     string            `json:"Req_URL"`
	ReqType    string            `json:"Req_Type"`
	ReqID      string            `json:"Req_ID"`
	ReqTs      string            `json:"Req_Ts"`
	ReqMethod  string            `json:"Req_Method"`
	ReqBody    string            `json:"Req_Body,omitempty"`
	ReqHeaders map[string]string `json:"Req_Headers,omitempty"`
}

type InputElasticSearchConfig struct {
	Address   MultiOption //엘라스틱서치 주소
	Index     string      //인덱스 이름
	FromDate  time.Time   //시작시간
	ToDate    time.Time   //종료시간
	Includes  MultiOption //들어가 있을 컬럼
	Match     string      // match_phrase이 해당
	Transport *http.Transport
}

//Range FromDate와 ToDate의 간격을 분으로 변환하여 알려줌
func (c InputElasticSearchConfig) Range() int {
	//ToDate - FromDate 를 int로 나타내준다.
	diff := c.ToDate.Sub(c.FromDate)
	if diff > 0 {
		return int(diff / time.Minute)
	}
	return 0
}

type ElasticsearchInput struct {
	config   *InputElasticSearchConfig
	messages chan *ElasticsearchMessage
	quit     chan struct{}
}

func NewElasticsearchInput(address string, config *InputElasticSearchConfig) *ElasticsearchInput {

	if config.Transport == nil {
		config.Transport = &http.Transport{
			ResponseHeaderTimeout: time.Second * 5,
			MaxIdleConns:          10,
		}
	}

	e := &ElasticsearchInput{
		config:   config,
		messages: make(chan *ElasticsearchMessage),
		quit:     make(chan struct{}),
	}

	go func(config *InputElasticSearchConfig) {
		es(config, e.messages)

	}(config)

	return e
}

func (e *ElasticsearchInput) PluginRead() (*Message, error) {
	var message *ElasticsearchMessage
	var err error
	var msg Message
	select {
	case <-e.quit:
		return nil, ErrorStopped
	case message = <-e.messages:
	}

	msg.Data, err = message.Dump()
	if err != nil {
		log.Fatal(1, "[ELASTICSEARCH] failed to decode: ", err)
		return nil, err
	}

	if isOriginPayload(msg.Data) {
		msg.Meta, msg.Data = payloadMetaWithBody(msg.Data)
	}

	return &msg, nil
}

func (e *ElasticsearchInput) Close() error {
	close(e.quit)
	return nil
}

func (e *ElasticsearchInput) String() string {
	return "ElasticsearchInput: " + strings.Join(e.config.Address, ",")
}

func es(c *InputElasticSearchConfig, messages chan *ElasticsearchMessage) {
	const layout = "2006-01-02T15:04:05.000Z"

	cfg := elasticsearch7.Config{
		Addresses: c.Address,
		Transport: c.Transport,
	}

	es, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	loop := c.Range()
	log.Printf("loop range %d\n", loop)
	var size = 1000
	var scrollID string
	var batchNum = 0
	var documents = 0
	var lastTime int64 = -1

	for i := 0; i < loop; i++ {
		var subDocuments = 0
		batchNum = 0
		var buf bytes.Buffer
		gte := c.FromDate.Add(time.Duration(i) * time.Minute)
		lt := gte.Add(time.Duration(59)*time.Second + 999*time.Millisecond)
		log.Println(gte, "  ", lt)
		//time.Sleep(time.Second*30)

		body := map[string]interface{}{
			//"_source": map[string]interface{}{
			//	"includes": c.Includes,
			//},
			"sort": []interface{}{
				map[string]interface{}{
					"@timestamp": "asc",
				},
			},
			"query": map[string]interface{}{

				"bool": map[string]interface{}{
					"must": []interface{}{
						map[string]interface{}{
							"match_phrase": map[string]interface{}{
								"@log_group": c.Match,
							},
						},
						map[string]interface{}{
							"match_phrase": map[string]interface{}{
								"logType": "formattedLog",
							},
						},
					},
					"filter": map[string]interface{}{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"time_zone": "+09:00",
								"gte":       gte.Format(layout),
								"lt":        lt.Format(layout),
							},
						},
					},
				},
			},
		}

		dslQuery, _ := json.Marshal(body)
		log.Println(string(dslQuery))

		if err = json.NewEncoder(&buf).Encode(body); err != nil {
			log.Fatalf("Error encoding query : %s", err)
		}

		res, err = es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(c.Index),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
			es.Search.WithSize(size),
			es.Search.WithScroll(time.Minute),
		)

		if err != nil {
			log.Fatalf("Error getting response: %s", err)
			res.Body.Close()
		}
		if res.IsError() {
			b, _ := json.Marshal(body)
			log.Println(string(b))
			log.Fatalf("Error response: %s", res)
		}

		var resJson map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&resJson); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		scrollID = resJson["_scroll_id"].(string)
		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d hits; took: %dms   scroll_id:%s",
			res.Status(),
			int(resJson["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(resJson["took"].(float64)),
			scrollID,
		)
		hits := resJson["hits"].(map[string]interface{})["hits"].([]interface{})
		documents += len(hits)
		subDocuments += len(hits)
		if len(hits) > 0 {
			//log.Println("Batch   ", batchNum)
			//log.Println("ScrollID", scrollID)

			log.Printf("0-Batch %d, ScrollID %s  message_len %d\n", batchNum, scrollID, len(messages))

			var e []byte
			var ems *ElasticsearchMessage
			for _, hit := range resJson["hits"].(map[string]interface{})["hits"].([]interface{}) {
				e, err = json.Marshal(hit)
				if err != nil {
					log.Fatalf("json marshal err : %s\n", err)
				}
				doc, uErr := UnmarshalElasticsearchDocument(e)
				if uErr != nil {
					log.Fatalf("document decode err : %s\n", uErr)
				}
				ems, err = NewElasticsearchMessage(doc)
				if err != nil {
					log.Fatal(err)
				}

				limiter(ems, &lastTime)

				messages <- ems

			}
			//log.Println(strings.Repeat("-", 80))
		}

		for {
			batchNum++
			res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			if res.IsError() {
				log.Fatalf("Error response: %s", res)
			}

			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				res.Body.Close()
				log.Fatalf("Error parsing the response body: %s", err)
			}

			// Extract the scrollID from response
			scrollID = r["_scroll_id"].(string)

			// Extract the search results
			hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
			documents += len(hits)
			subDocuments += len(hits)

			// Break out of the loop when there are no results
			if len(hits) < 1 {
				log.Println("Finished scrolling. SubDocuments ", subDocuments)
				res.Body.Close()
				break
			} else {
				log.Printf("1-Batch %d, ScrollID %s  message_len %d\n", batchNum, scrollID, len(messages))

				var e []byte
				var ems *ElasticsearchMessage
				for _, hit := range resJson["hits"].(map[string]interface{})["hits"].([]interface{}) {
					//log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
					//messages <- hit
					e, err = json.Marshal(hit)
					if err != nil {
						log.Fatalf("json marshal err : %s\n", err)
					}

					//f, ferr := ioutil.TempFile("./tmp", "gor")
					//if ferr != nil {
					//	panic(ferr)
					//}
					//
					//bytes.NewReader(e).WriteTo(f)
					//f.Close()

					doc, uErr := UnmarshalElasticsearchDocument(e)
					if uErr != nil {
						log.Fatalf("document decode err : %s\n", uErr)
					}
					ems, err = NewElasticsearchMessage(doc)
					if err != nil {
						log.Fatal(err)
					}
					limiter(ems, &lastTime)

					messages <- ems
				}
				log.Println(strings.Repeat("-", 80))
			}
			res.Body.Close()
		}
		res.Body.Close()

	}

	log.Println("Total : ", documents)
}

func limiter(ems *ElasticsearchMessage, lastTime *int64) {
	timestamp, _ := strconv.ParseInt(ems.ReqTs, 10, 64)
	if *lastTime != -1 {
		diff := timestamp - *lastTime
		*lastTime = timestamp
		_ = diff

		//배속을 조절할 수 있음
		//if i.speedFactor != 1 {
		//	diff = int64(float64(diff) / i.speedFactor)
		//}

		time.Sleep(time.Duration(diff))
	} else {
		*lastTime = timestamp
	}
}

func NewElasticsearchMessage(doc ElasticsearchDocument) (*ElasticsearchMessage, error) {
	const layout = "2006-01-02T15:04:05.000Z"

	serverLog := doc.Source.KnowreDaekyo.ServerLog
	timestamp := doc.Source.Timestamp
	requestID := doc.ID
	url := serverLog.URL
	host := strings.Split(serverLog.Router, ":")[0]
	method := serverLog.Method
	body := serverLog.Body
	//auth := serverLog.Token
	accessToken := serverLog.AccessToken
	appFlavor := serverLog.AppFlavor
	cookie := serverLog.Cookie

	if len(serverLog.Cookie) > 0 {
		c, _ := UnmarshalServerLogCookie([]byte(serverLog.Cookie))
		cookie = c.Cookie.String()
	}

	//cehck empty string
	if Empty(body) {
		body = ""
	}

	var headers = map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Host":         host,
	}
	if !Empty(cookie) {
		headers["Cookie"] = cookie
	}
	if !Empty(body) {
		headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	}
	//if len(auth) > 0 {
	//	headers["Authorization"] = auth
	//}
	if !Empty(accessToken) {
		headers["x-access-token"] = accessToken
	}
	if !Empty(appFlavor) {
		headers["x-app-flavor"] = appFlavor
	}

	logTimestamp, terr := time.Parse(layout, timestamp)
	if terr != nil {
		return nil, terr
	}
	ems := &ElasticsearchMessage{
		ReqURL:     url,
		ReqType:    "1",
		ReqID:      requestID,
		ReqTs:      fmt.Sprintf("%d", logTimestamp.UnixNano()),
		ReqMethod:  method,
		ReqBody:    body,
		ReqHeaders: headers,
	}
	return ems, nil
}

// Dump returns the given request in its HTTP/1.x wire
// representation.
func (m ElasticsearchMessage) Dump() ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("%s %s %s\n", m.ReqType, m.ReqID, m.ReqTs))
	b.WriteString(fmt.Sprintf("%s %s HTTP/1.1", m.ReqMethod, m.ReqURL))
	b.Write(proto.CRLF)
	for key, value := range m.ReqHeaders {
		b.WriteString(fmt.Sprintf("%s: %s", key, value))
		b.Write(proto.CRLF)
	}

	b.Write(proto.CRLF)
	b.WriteString(m.ReqBody)

	return b.Bytes(), nil
}
