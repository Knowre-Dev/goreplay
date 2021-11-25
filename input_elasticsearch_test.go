package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestElasticSearchConfigRange(t *testing.T) {
	config := InputElasticSearchConfig{
		Address:   []string{"127.0.0.1", "127.0.0.2"},
		Index:     "cwl-raw-2021.07.03",
		FromDate:  time.Date(2021, 8, 3, 0, 0, 0, 0, time.Local),
		ToDate:    time.Date(2021, 8, 4, 0, 0, 0, 0, time.Local),
		Includes:  nil,
		Transport: nil,
	}

	assert.Equal(t, 60*24, config.Range())

}

func TestES(t *testing.T) {
	const layout = "cwl-raw-2006.01.02"
	//date := time.Now().AddDate(0,-1,1)
	date := time.Now().AddDate(0, 0, -1)
	indexName := date.Format(layout)
	transport := &http.Transport{
		ResponseHeaderTimeout: time.Second * 5,
		MaxIdleConns:          10,
	}

	y, m, d := date.Date()
	config := InputElasticSearchConfig{
		Address:   []string{"https://vpc-sl-logstrg-orange-prd-q76s3uteh4ooxa3r4brwce2yau.ap-northeast-2.es.amazonaws.com"},
		Index:     indexName,
		FromDate:  time.Date(y, m, d, 0, 0, 0, 0, time.Local),
		ToDate:    time.Date(y, m, d, 0, 1, 0, 0, time.Local),
		Includes:  []string{},
		Transport: transport,
	}

	input := NewElasticsearchInput("", &config)
	for {
		msg, err := input.PluginRead()
		if err != nil {
			assert.Equal(t, err, ErrorStopped)
			break
		}
		fmt.Println(string(msg.Meta))
		fmt.Println(string(msg.Data))
		fmt.Println(strings.Repeat("-", 80))
	}

	_ = input

	log.Println("TestES end")

}

func TestDump(t *testing.T) {
	esJSON := `{"_id":"7iJuLnsBy4VzcorOmQtT","_index":"cwl-raw-2021.08.10","_score":null,"_source":{"@id":"36318360321655520334835791342972977475959172674164162595","@log_group":"/ecs/krdky-stable","@log_stream":"ecs/krdky-stable/2e4c2b0032fd4fa28d0918662da0caf5","@message":"{\"logType\":\"formattedLog\",\"knowre-daekyo\":{\"serverLog\":{\"url\":\"/api/v3/result/lesson/KNRREVI47783?errorCode=true\u0026timestamp=1628571601312\",\"method\":\"GET\",\"ip\":\"106.247.213.40\",\"router\":\"daekyo-stable.knowreapp.com:29\",\"body\":\"{}\",\"userAgent\":\"okhttp/3.12.1\",\"accessToken\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2ODE2NzgsImFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImV4dGVybmFsX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImNvbnN1bWVyS2V5Ijoia25vd3JlIiwiaXNzdWVyIjoiZGFla3lvLXN0YWJsZS5rbm93cmVhcHAuY29tIiwiaXNzdWVfZGF0ZSI6IjIwMjEtMDgtMTBUMDQ6MzU6MTIuNTY5WiIsInNlc3Npb25faWQiOiI2ODE2NzhUODU3MDExMjU2OUQzYTIwZTAiLCJpYXQiOjE2Mjg1NzAxMTJ9.37N3jmvvv8uMdzRAYz_PYRIFnb-bWiL5djfuGmZJHMA\",\"token\":\"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJncm93dGhJbmZvIjoie1wia3VubnJcIjpcIjAwNTUzNTg4NjNcIixcImFjY2Vzc1Rva2VuXCI6XCJLbGRKZHE5OTQ5SUdsNU1CVXlmMXlnNThhdTd0SE5ZNzVEMDJ2TEVUYmw1VU9qaWJhMXc5WUpBazVLQjlhKzcvXCIsXCJsb2dpblNlcVwiOlwiNjA4MTMxMDNcIixcInByb2R1Y3RJZFwiOlwiUERNXCIsXCJjdXJyaWN1bHVtSWRcIjo0LFwiY3VycmljdWx1bVR5cGVcIjpcIlZBQ0FUSU9OXCJ9IiwiY29ubmVjdGVkQXQiOiIyMDIxLTA4LTEwVDA0OjM1OjEyLjUwMVoiLCJ1c2VySUQiOjY4MTY3OCwidXNlck5hbWUiOiIiLCJleHRlcm5hbEFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImN1cnJpY3VsdW1faWQiOjQsImN1cnJpY3VsdW1fdHlwZSI6IlZBQ0FUSU9OIiwiZGlmZmljdWx0eSI6NTAsInRva2VuX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsInByb2R1Y3RUeXBlIjoiRUxFTSIsInN0YW1wIjoic2hhMSQzNjZlMTJkZSQxJDgyNWYyMGVmM2NjZGJkMzg1OGJkZThjNWUyNjEyZTkzNjBkNjY5ZmYiLCJjbGllbnRWZXJzaW9uIjoxMDUwMDk5LCJ1c2VyRFBJIjoyNDAsImFwaV92ZXJzaW9uIjoidjMiLCJpc0V4dGVybmFsTGVhZ3VlVXNlciI6dHJ1ZSwidXNlclR5cGUiOiJEQUVLWU9fTEMiLCJ0aWNrIjo0NywibGVzc29uSW5mbyI6eyJsZXNzb25UeXBlIjoiUkVWSUVXIiwiY3VycmljdWx1bUlkIjo0LCJjdXJyaWN1bHVtVHlwZSI6IlZBQ0FUSU9OIiwiY29tcG9zaXRlSWQiOiJLTlJSRVZJNDc3ODMiLCJjaGFwdGVyU2VxIjoxLCJsZXNzb25TZXEiOjIsImxlc3NvbklkIjo4NjYwNn0sImxlc3NvbkV4dGVybmFsS2V5IjoiS05SUkVWSTQ3NzgzIiwiYWN0aXZpdHkiOnsidHJpYWxJRCI6IjY4MTY3OEk5VDU1MDMzMDE3NTgyIiwidW5pdElEIjoyODM2NCwiaXNTdGVwIjpmYWxzZX0sInJldHJ5IjpudWxsLCJ0cnlDb3VudCI6MCwiaWF0IjoxNjI4NTcxNTU4LCJleHAiOjE2Mjg1NzY5NTh9.z9_Ay1fr_UGBw5YquSePZrOgAa8ZYoJBngK4Fu4DWUY\",\"appFlavor\":\"summitScore\",\"cookie\":\"{\\\"cookie\\\":{\\\"connect.sid\\\":\\\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJncm93dGhJbmZvIjoie1wia3VubnJcIjpcIjAwNTUzNTg4NjNcIixcImFjY2Vzc1Rva2VuXCI6XCJLbGRKZHE5OTQ5SUdsNU1CVXlmMXlnNThhdTd0SE5ZNzVEMDJ2TEVUYmw1VU9qaWJhMXc5WUpBazVLQjlhKzcvXCIsXCJsb2dpblNlcVwiOlwiNjA4MTMxMDNcIixcInByb2R1Y3RJZFwiOlwiUERNXCIsXCJjdXJyaWN1bHVtSWRcIjo0LFwiY3VycmljdWx1bVR5cGVcIjpcIlZBQ0FUSU9OXCJ9IiwiY29ubmVjdGVkQXQiOiIyMDIxLTA4LTEwVDA0OjM1OjEyLjUwMVoiLCJ1c2VySUQiOjY4MTY3OCwidXNlck5hbWUiOiIiLCJleHRlcm5hbEFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImN1cnJpY3VsdW1faWQiOjQsImN1cnJpY3VsdW1fdHlwZSI6IlZBQ0FUSU9OIiwiZGlmZmljdWx0eSI6NTAsInRva2VuX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsInByb2R1Y3RUeXBlIjoiRUxFTSIsInN0YW1wIjoic2hhMSQzNjZlMTJkZSQxJDgyNWYyMGVmM2NjZGJkMzg1OGJkZThjNWUyNjEyZTkzNjBkNjY5ZmYiLCJjbGllbnRWZXJzaW9uIjoxMDUwMDk5LCJ1c2VyRFBJIjoyNDAsImFwaV92ZXJzaW9uIjoidjMiLCJpc0V4dGVybmFsTGVhZ3VlVXNlciI6dHJ1ZSwidXNlclR5cGUiOiJEQUVLWU9fTEMiLCJ0aWNrIjo0NywibGVzc29uSW5mbyI6eyJsZXNzb25UeXBlIjoiUkVWSUVXIiwiY3VycmljdWx1bUlkIjo0LCJjdXJyaWN1bHVtVHlwZSI6IlZBQ0FUSU9OIiwiY29tcG9zaXRlSWQiOiJLTlJSRVZJNDc3ODMiLCJjaGFwdGVyU2VxIjoxLCJsZXNzb25TZXEiOjIsImxlc3NvbklkIjo4NjYwNn0sImxlc3NvbkV4dGVybmFsS2V5IjoiS05SUkVWSTQ3NzgzIiwiYWN0aXZpdHkiOnsidHJpYWxJRCI6IjY4MTY3OEk5VDU1MDMzMDE3NTgyIiwidW5pdElEIjoyODM2NCwiaXNTdGVwIjpmYWxzZX0sInJldHJ5IjpudWxsLCJ0cnlDb3VudCI6MCwiaWF0IjoxNjI4NTcxNTU4LCJleHAiOjE2Mjg1NzY5NTh9.z9_Ay1fr_UGBw5YquSePZrOgAa8ZYoJBngK4Fu4DWUY\\\"}}\",\"session\":\"{\\\"_domain\\\":{\\\"domain\\\":\\\"daekyo-stable.knowreapp.com\\\"},\\\"_secret\\\":\\\"knowre_prod\\\",\\\"_cookieKey\\\":\\\"connect.sid\\\",\\\"_tokenExpireTime\\\":5400,\\\"growthInfo\\\":\\\"{\\\\\\\"kunnr\\\\\\\":\\\\\\\"0055358863\\\\\\\",\\\\\\\"accessToken\\\\\\\":\\\\\\\"KldJdq9949IGl5MBUyf1yg58au7tHNY75D02vLETbl5UOjiba1w9YJAk5KB9a+7/\\\\\\\",\\\\\\\"loginSeq\\\\\\\":\\\\\\\"60813103\\\\\\\",\\\\\\\"productId\\\\\\\":\\\\\\\"PDM\\\\\\\",\\\\\\\"curriculumId\\\\\\\":4,\\\\\\\"curriculumType\\\\\\\":\\\\\\\"VACATION\\\\\\\"}\\\",\\\"connectedAt\\\":\\\"2021-08-10T04:35:12.501Z\\\",\\\"userID\\\":681678,\\\"userName\\\":\\\"\\\",\\\"externalAccount\\\":\\\"000S-0055358863_P\\\",\\\"curriculum_id\\\":4,\\\"curriculum_type\\\":\\\"VACATION\\\",\\\"difficulty\\\":50,\\\"token_account\\\":\\\"000S-0055358863_P\\\",\\\"productType\\\":\\\"ELEM\\\",\\\"stamp\\\":\\\"sha1$366e12de$1$825f20ef3ccdbd3858bde8c5e2612e9360d669ff\\\",\\\"clientVersion\\\":1050099,\\\"userDPI\\\":240,\\\"api_version\\\":\\\"v3\\\",\\\"isExternalLeagueUser\\\":true,\\\"userType\\\":\\\"DAEKYO_LC\\\",\\\"tick\\\":47,\\\"lessonInfo\\\":{\\\"lessonType\\\":\\\"REVIEW\\\",\\\"curriculumId\\\":4,\\\"curriculumType\\\":\\\"VACATION\\\",\\\"compositeId\\\":\\\"KNRREVI47783\\\",\\\"chapterSeq\\\":1,\\\"lessonSeq\\\":2,\\\"lessonId\\\":86606},\\\"lessonExternalKey\\\":\\\"KNRREVI47783\\\",\\\"activity\\\":{\\\"trialID\\\":\\\"681678I9T55033017582\\\",\\\"unitID\\\":28364,\\\"isStep\\\":false},\\\"retry\\\":null,\\\"tryCount\\\":0,\\\"iat\\\":1628571558,\\\"exp\\\":1628576958,\\\"_tokenHash\\\":\\\"366a361cbc9384239d5e0f0daf46121f\\\"}\",\"trace\":\"{\\\"elapsedTime\\\":{\\\"server-log-start\\\":{\\\"time\\\":[0,56052],\\\"end\\\":true,\\\"elapsed\\\":0.056052},\\\"jwtSessionOut\\\":{\\\"time\\\":[0,56948],\\\"end\\\":true,\\\"elapsed\\\":0.056948},\\\"responseOut\\\":{\\\"time\\\":[0,247856],\\\"end\\\":true,\\\"elapsed\\\":0.247856},\\\"totalElapsedTime\\\":8.255396,\\\"state\\\":\\\"normal\\\"}}\",\"performance\":8.255396,\"error\":null,\"result\":true,\"req\":null,\"session_id\":\"681678T8570112569D3a20e0\",\"amazonTraceId\":\"Root=1-611207d1-7e85e30452f1f5f02a8540e9\",\"user_id\":681678,\"userType\":\"DAEKYO_LC\",\"parameters\":\"{\\\"query\\\":{\\\"errorCode\\\":\\\"true\\\",\\\"timestamp\\\":\\\"1628571601312\\\"},\\\"body\\\":{}}\"}}}","@owner":"468720534852","@timestamp":"2021-08-10T05:00:01.457Z","@version":"1","kafka.consumer_group":"logstash-cwl-raw","kafka.key":"%{[@metadata][kafka][key]}","kafka.offset":"774285035","kafka.partition":"0","kafka.timestamp":"1628571604340","kafka.topic":"formattedLog","knowre-daekyo":{"serverLog":{"accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2ODE2NzgsImFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImV4dGVybmFsX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImNvbnN1bWVyS2V5Ijoia25vd3JlIiwiaXNzdWVyIjoiZGFla3lvLXN0YWJsZS5rbm93cmVhcHAuY29tIiwiaXNzdWVfZGF0ZSI6IjIwMjEtMDgtMTBUMDQ6MzU6MTIuNTY5WiIsInNlc3Npb25faWQiOiI2ODE2NzhUODU3MDExMjU2OUQzYTIwZTAiLCJpYXQiOjE2Mjg1NzAxMTJ9.37N3jmvvv8uMdzRAYz_PYRIFnb-bWiL5djfuGmZJHMA","amazonTraceId":"Root=1-611207d1-7e85e30452f1f5f02a8540e9","appFlavor":"summitScore","body":"{}","cookie":"{\"cookie\":{\"connect.sid\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJncm93dGhJbmZvIjoie1wia3VubnJcIjpcIjAwNTUzNTg4NjNcIixcImFjY2Vzc1Rva2VuXCI6XCJLbGRKZHE5OTQ5SUdsNU1CVXlmMXlnNThhdTd0SE5ZNzVEMDJ2TEVUYmw1VU9qaWJhMXc5WUpBazVLQjlhKzcvXCIsXCJsb2dpblNlcVwiOlwiNjA4MTMxMDNcIixcInByb2R1Y3RJZFwiOlwiUERNXCIsXCJjdXJyaWN1bHVtSWRcIjo0LFwiY3VycmljdWx1bVR5cGVcIjpcIlZBQ0FUSU9OXCJ9IiwiY29ubmVjdGVkQXQiOiIyMDIxLTA4LTEwVDA0OjM1OjEyLjUwMVoiLCJ1c2VySUQiOjY4MTY3OCwidXNlck5hbWUiOiIiLCJleHRlcm5hbEFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImN1cnJpY3VsdW1faWQiOjQsImN1cnJpY3VsdW1fdHlwZSI6IlZBQ0FUSU9OIiwiZGlmZmljdWx0eSI6NTAsInRva2VuX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsInByb2R1Y3RUeXBlIjoiRUxFTSIsInN0YW1wIjoic2hhMSQzNjZlMTJkZSQxJDgyNWYyMGVmM2NjZGJkMzg1OGJkZThjNWUyNjEyZTkzNjBkNjY5ZmYiLCJjbGllbnRWZXJzaW9uIjoxMDUwMDk5LCJ1c2VyRFBJIjoyNDAsImFwaV92ZXJzaW9uIjoidjMiLCJpc0V4dGVybmFsTGVhZ3VlVXNlciI6dHJ1ZSwidXNlclR5cGUiOiJEQUVLWU9fTEMiLCJ0aWNrIjo0NywibGVzc29uSW5mbyI6eyJsZXNzb25UeXBlIjoiUkVWSUVXIiwiY3VycmljdWx1bUlkIjo0LCJjdXJyaWN1bHVtVHlwZSI6IlZBQ0FUSU9OIiwiY29tcG9zaXRlSWQiOiJLTlJSRVZJNDc3ODMiLCJjaGFwdGVyU2VxIjoxLCJsZXNzb25TZXEiOjIsImxlc3NvbklkIjo4NjYwNn0sImxlc3NvbkV4dGVybmFsS2V5IjoiS05SUkVWSTQ3NzgzIiwiYWN0aXZpdHkiOnsidHJpYWxJRCI6IjY4MTY3OEk5VDU1MDMzMDE3NTgyIiwidW5pdElEIjoyODM2NCwiaXNTdGVwIjpmYWxzZX0sInJldHJ5IjpudWxsLCJ0cnlDb3VudCI6MCwiaWF0IjoxNjI4NTcxNTU4LCJleHAiOjE2Mjg1NzY5NTh9.z9_Ay1fr_UGBw5YquSePZrOgAa8ZYoJBngK4Fu4DWUY\"}}","error":null,"ip":"106.247.213.40","method":"GET","parameters":"{\"query\":{\"errorCode\":\"true\",\"timestamp\":\"1628571601312\"},\"body\":{}}","performance":8.255396,"req":null,"result":true,"router":"daekyo-stable.knowreapp.com:29","session":"{\"_domain\":{\"domain\":\"daekyo-stable.knowreapp.com\"},\"_secret\":\"knowre_prod\",\"_cookieKey\":\"connect.sid\",\"_tokenExpireTime\":5400,\"growthInfo\":\"{\\\"kunnr\\\":\\\"0055358863\\\",\\\"accessToken\\\":\\\"KldJdq9949IGl5MBUyf1yg58au7tHNY75D02vLETbl5UOjiba1w9YJAk5KB9a+7/\\\",\\\"loginSeq\\\":\\\"60813103\\\",\\\"productId\\\":\\\"PDM\\\",\\\"curriculumId\\\":4,\\\"curriculumType\\\":\\\"VACATION\\\"}\",\"connectedAt\":\"2021-08-10T04:35:12.501Z\",\"userID\":681678,\"userName\":\"\",\"externalAccount\":\"000S-0055358863_P\",\"curriculum_id\":4,\"curriculum_type\":\"VACATION\",\"difficulty\":50,\"token_account\":\"000S-0055358863_P\",\"productType\":\"ELEM\",\"stamp\":\"sha1$366e12de$1$825f20ef3ccdbd3858bde8c5e2612e9360d669ff\",\"clientVersion\":1050099,\"userDPI\":240,\"api_version\":\"v3\",\"isExternalLeagueUser\":true,\"userType\":\"DAEKYO_LC\",\"tick\":47,\"lessonInfo\":{\"lessonType\":\"REVIEW\",\"curriculumId\":4,\"curriculumType\":\"VACATION\",\"compositeId\":\"KNRREVI47783\",\"chapterSeq\":1,\"lessonSeq\":2,\"lessonId\":86606},\"lessonExternalKey\":\"KNRREVI47783\",\"activity\":{\"trialID\":\"681678I9T55033017582\",\"unitID\":28364,\"isStep\":false},\"retry\":null,\"tryCount\":0,\"iat\":1628571558,\"exp\":1628576958,\"_tokenHash\":\"366a361cbc9384239d5e0f0daf46121f\"}","session_id":"681678T8570112569D3a20e0","token":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJncm93dGhJbmZvIjoie1wia3VubnJcIjpcIjAwNTUzNTg4NjNcIixcImFjY2Vzc1Rva2VuXCI6XCJLbGRKZHE5OTQ5SUdsNU1CVXlmMXlnNThhdTd0SE5ZNzVEMDJ2TEVUYmw1VU9qaWJhMXc5WUpBazVLQjlhKzcvXCIsXCJsb2dpblNlcVwiOlwiNjA4MTMxMDNcIixcInByb2R1Y3RJZFwiOlwiUERNXCIsXCJjdXJyaWN1bHVtSWRcIjo0LFwiY3VycmljdWx1bVR5cGVcIjpcIlZBQ0FUSU9OXCJ9IiwiY29ubmVjdGVkQXQiOiIyMDIxLTA4LTEwVDA0OjM1OjEyLjUwMVoiLCJ1c2VySUQiOjY4MTY3OCwidXNlck5hbWUiOiIiLCJleHRlcm5hbEFjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsImN1cnJpY3VsdW1faWQiOjQsImN1cnJpY3VsdW1fdHlwZSI6IlZBQ0FUSU9OIiwiZGlmZmljdWx0eSI6NTAsInRva2VuX2FjY291bnQiOiIwMDBTLTAwNTUzNTg4NjNfUCIsInByb2R1Y3RUeXBlIjoiRUxFTSIsInN0YW1wIjoic2hhMSQzNjZlMTJkZSQxJDgyNWYyMGVmM2NjZGJkMzg1OGJkZThjNWUyNjEyZTkzNjBkNjY5ZmYiLCJjbGllbnRWZXJzaW9uIjoxMDUwMDk5LCJ1c2VyRFBJIjoyNDAsImFwaV92ZXJzaW9uIjoidjMiLCJpc0V4dGVybmFsTGVhZ3VlVXNlciI6dHJ1ZSwidXNlclR5cGUiOiJEQUVLWU9fTEMiLCJ0aWNrIjo0NywibGVzc29uSW5mbyI6eyJsZXNzb25UeXBlIjoiUkVWSUVXIiwiY3VycmljdWx1bUlkIjo0LCJjdXJyaWN1bHVtVHlwZSI6IlZBQ0FUSU9OIiwiY29tcG9zaXRlSWQiOiJLTlJSRVZJNDc3ODMiLCJjaGFwdGVyU2VxIjoxLCJsZXNzb25TZXEiOjIsImxlc3NvbklkIjo4NjYwNn0sImxlc3NvbkV4dGVybmFsS2V5IjoiS05SUkVWSTQ3NzgzIiwiYWN0aXZpdHkiOnsidHJpYWxJRCI6IjY4MTY3OEk5VDU1MDMzMDE3NTgyIiwidW5pdElEIjoyODM2NCwiaXNTdGVwIjpmYWxzZX0sInJldHJ5IjpudWxsLCJ0cnlDb3VudCI6MCwiaWF0IjoxNjI4NTcxNTU4LCJleHAiOjE2Mjg1NzY5NTh9.z9_Ay1fr_UGBw5YquSePZrOgAa8ZYoJBngK4Fu4DWUY","trace":"{\"elapsedTime\":{\"server-log-start\":{\"time\":[0,56052],\"end\":true,\"elapsed\":0.056052},\"jwtSessionOut\":{\"time\":[0,56948],\"end\":true,\"elapsed\":0.056948},\"responseOut\":{\"time\":[0,247856],\"end\":true,\"elapsed\":0.247856},\"totalElapsedTime\":8.255396,\"state\":\"normal\"}}","url":"/api/v3/result/lesson/KNRREVI47783?errorCode=true\u0026timestamp=1628571601312","userAgent":"okhttp/3.12.1","userType":"DAEKYO_LC","user_id":681678}},"logType":"formattedLog"},"_type":"_doc","sort":[1628571601457]}`

	doc, err := UnmarshalElasticsearchDocument([]byte(esJSON))
	assert.NoError(t, err, "error message %s")
	ems, emsErr := NewElasticsearchMessage(doc)
	assert.NoError(t, emsErr, "error message %s")
	b, bErr := ems.Dump()
	assert.NoError(t, bErr, "error message %s")

	var msg Message
	msg.Meta, msg.Data = payloadMetaWithBody(b)

	assert.True(t, strings.HasSuffix(string(msg.Meta), "\n"))

	headerStart := bytes.Index(msg.Data, []byte("\n")) + 1
	reader := bufio.NewReader(strings.NewReader(string(msg.Data[headerStart:])))
	tp := textproto.NewReader(reader)
	mimeHeader, mimeErr := tp.ReadMIMEHeader()
	if mimeErr != nil {
		log.Fatal(mimeErr)
	}

	httpHeader := http.Header(mimeHeader)
	fmt.Println(httpHeader)

	//	assert.True()

}

func TestDecodeBody(t *testing.T) {
	raws := [][]byte{
		[]byte(`{\"input\":\"{\\\"account\\\":\\\"jimy@knowre.com\\\", \\\"password\\\":\\\"111111\\\", \\\"productType\\\": \\\"AIMS\\\"}\"}`),
		[]byte(`{\"account\":\"jimy@knowre.com\", \"password\":\"111111\", \"productType\": \"AIMS\"}`),
		[]byte(`{"altToken":"a4de2a0924a0f13509da5ab5791886e4aad97f3349b9db"}`),
	}

	for _, raw := range raws {
		var msg string
		var err error
		body := string(raw)
		if IsJSON(body) {
			continue
		}

		body = "\"" + body + "\""
		msg, err = strconv.Unquote(body)

		if err != nil {
			log.Fatal(err, " ", string(raw))
		}

		log.Println(msg)
	}
	//
	//for _, raw := range raws {
	//	body := "\"" + string(raw) + "\""
	//	s, err := strconv.Unquote(body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Println(s)
	//}

}
