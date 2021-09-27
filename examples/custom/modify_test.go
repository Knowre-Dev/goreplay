package main

import (
	"github.com/buger/jsonparser"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestJWT(t *testing.T) {
	tokenString := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJjb25uZWN0ZWRBdCI6IjIwMjEtMDktMDFUMDA6MDE6NDUuOTA5WiIsInVzZXJJRCI6NDQ3NDU1LCJ1c2VyTmFtZSI6IiIsImV4dGVybmFsQWNjb3VudCI6bnVsbCwiY3VycmljdWx1bV9pZCI6OCwiY3VycmljdWx1bV90eXBlIjoiVkFDQVRJT04iLCJkaWZmaWN1bHR5IjoyNSwicGF5bWVudENoZWNrIjpmYWxzZSwidG9rZW5fYWNjb3VudCI6ImVsZTQydjkwIiwicHJvZHVjdFR5cGUiOiJFTEVNIiwic3RhbXAiOiJzaGExJGYwYWFhNzM1JDEkNWE5ZjA3ZTlmZDlhMzdhMTViZWEzODg1ZGNkNWUzZGMxMDFkYTI5NCIsImNsaWVudFZlcnNpb24iOjEwMjAwMDMsImFwaV92ZXJzaW9uIjoidjIiLCJ1c2VyVHlwZSI6IktOT1dSRV9URVNUIiwidGljayI6MiwiYWN0aXZpdHkiOnsidHJpYWxJRCI6IjQ0NzQ1NUkwVDI2OTQ1ODU3MzAzIiwidW5pdElEIjoyNDM3NCwiaXNTdGVwIjpmYWxzZX0sInRyeUNvdW50IjowLCJpYXQiOjE2MzA0NTU2MzksImV4cCI6MTYzMDQ2MTAzOX0._O26cYHPMx1qm7cm1PF1CPT501jUMoxcZBA23wAX1WE`
	token, _ := jwt.Parse(tokenString, nil)

	if token != nil {
		m := token.Claims.(jwt.MapClaims)
		for k, v := range m {
			log.Printf("[%s] [%s]", k, v)
		}
	}

}
func TestExtractUserIDFromBody(t *testing.T) {

	payload := `{"success":true,"session":true,"error":null,"errorCode":null,"data":{"isFirst":true,"curriculum":{"id":8,"type":"VACATION"},"nickname":"ele42v90","curriculumId":8,"curriculumType":"VACATION","accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0NDc0NTUsImFjY291bnQiOiJlbGU0MnY5MCIsImV4dGVybmFsX2FjY291bnQiOiJlbGU0MnY5MCIsImNvbnN1bWVyS2V5IjoiZGFla3lvIiwiaXNzdWVyIjoibG9jYWwua25vd3JlYXBwLmNvbSIsImlzc3VlX2RhdGUiOiIyMDIxLTA5LTAxVDA0OjEyOjE3Ljg3N1oiLCJzZXNzaW9uX2lkIjoiNDQ3NDU1VDA0Njk1Mzc4NzdEIiwiaWF0IjoxNjMwNDY5NTM3fQ.lsN34fo7TBJrQBmgkrv5oQetn_MVY6vyI4hxvxfPeO4","isNewVersion":true,"version":5030101,"url":"https://knowre-daekyo-prod.s3.amazonaws.com/apk/5030101/20160118/app-daekyo-release-5.3.1-RC1.239-2015-1210-110350.apk","alertMsg":"......... ...... !!! TEST MESSAGE","forceUpdate":false},"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJjb25uZWN0ZWRBdCI6IjIwMjEtMDktMDFUMDQ6MTI6MTcuODQyWiIsInVzZXJJRCI6NDQ3NDU1LCJ1c2VyTmFtZSI6IiIsImV4dGVybmFsQWNjb3VudCI6bnVsbCwiY3VycmljdWx1bV9pZCI6OCwiY3VycmljdWx1bV90eXBlIjoiVkFDQVRJT04iLCJkaWZmaWN1bHR5IjoyNSwicGF5bWVudENoZWNrIjpmYWxzZSwidG9rZW5fYWNjb3VudCI6ImVsZTQydjkwIiwicHJvZHVjdFR5cGUiOiJFTEVNIiwic3RhbXAiOiJzaGExJGE2OTRjZGUxJDEkYjIzZGM2MTZkZjY3M2IwYmVkZTRmMWRlYjNhNmNhNjM5MDkyYjAzMCIsImNsaWVudFZlcnNpb24iOjEwMjAwMDMsImFwaV92ZXJzaW9uIjoidjIiLCJ1c2VyVHlwZSI6IktOT1dSRV9URVNUIiwidGljayI6MSwiaWF0IjoxNjMwNDY5NTM3LCJleHAiOjE2MzA0NzQ5Mzd9.dbakyQrw-TdXZK_C2WHoIGIiWfye0YE4P61iCZgAX7g"}`
	value, account, err := extractUserIDFromBody([]byte(payload), extractUserID, "token")

	assert.NoError(t, err)
	assert.Equal(t, account, "447455")
	log.Println(value)

	{
		value, account, err := extractUserIDFromBody([]byte(payload), extractUserID, XAccessTokens...)
		assert.NoError(t, err)
		assert.Equal(t, account, "447455")
		log.Println(value)

	}

}

func TestExtractUserID(t *testing.T) {
	tokenStr := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJjb25uZWN0ZWRBdCI6IjIwMjEtMDktMDFUMDQ6MTI6MTcuODQyWiIsInVzZXJJRCI6NDQ3NDU1LCJ1c2VyTmFtZSI6IiIsImV4dGVybmFsQWNjb3VudCI6bnVsbCwiY3VycmljdWx1bV9pZCI6OCwiY3VycmljdWx1bV90eXBlIjoiVkFDQVRJT04iLCJkaWZmaWN1bHR5IjoyNSwicGF5bWVudENoZWNrIjpmYWxzZSwidG9rZW5fYWNjb3VudCI6ImVsZTQydjkwIiwicHJvZHVjdFR5cGUiOiJFTEVNIiwic3RhbXAiOiJzaGExJGE2OTRjZGUxJDEkYjIzZGM2MTZkZjY3M2IwYmVkZTRmMWRlYjNhNmNhNjM5MDkyYjAzMCIsImNsaWVudFZlcnNpb24iOjEwMjAwMDMsImFwaV92ZXJzaW9uIjoidjIiLCJ1c2VyVHlwZSI6IktOT1dSRV9URVNUIiwidGljayI6MSwiaWF0IjoxNjMwNDY5NTM3LCJleHAiOjE2MzA0NzQ5Mzd9.dbakyQrw-TdXZK_C2WHoIGIiWfye0YE4P61iCZgAX7g`
	account, err := extractUserID([]byte(tokenStr))
	assert.NoError(t, err)
	assert.Equal(t, account, "447455")
}

func TestExtractProblemID(t *testing.T) {
	jsonData := `{
		"success": true,
		"session": false,
		"error": null,
		"errorCode": null,
		"data": {
		"compositeId": "KNRLESS39003",
			"difficulty": 25,
			"problems": [
	{
	"prob": 437901,
	"ptrn": 66833,
	"unit": 24424,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I0T24977530723"
	},
	{
	"prob": 437904,
	"ptrn": 66838,
	"unit": 24425,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I1T24977530723"
	},
	{
	"prob": 437920,
	"ptrn": 66838,
	"unit": 24427,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I2T24977530723"
	},
	{
	"prob": 437935,
	"ptrn": 66843,
	"unit": 24428,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I3T24977530723"
	},
	{
	"prob": 437960,
	"ptrn": 66846,
	"unit": 24430,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I4T24977530723"
	},
	{
	"prob": 437968,
	"ptrn": 66846,
	"unit": 24431,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I5T24977530723"
	},
	{
	"prob": 438113,
	"ptrn": 66846,
	"unit": 24434,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I6T24977530723"
	},
	{
	"prob": 438152,
	"ptrn": 66853,
	"unit": 24437,
	"status": "",
	"publishType": "D",
	"publishCount": 1,
	"trialID": "447455I7T24977530723"
	}
	],
	"curriculum": {
	"chapterSeq": 1,
	"chapterName": "분수의 덧셈과 뺄셈",
	"lessonSeq": 2,
	"lessonName": "대분수의 덧셈",
	"objectives": [
	"분수의 덧셈 원리와 형식을 이해하고 계산할 수 있다."
	],
	"repLessonSeq": 2
	},
	"lessonInfo": {
	"lessonType": "LESSON",
	"curriculumId": 8,
	"curriculumType": "VACATION",
	"compositeId": "KNRLESS39003",
	"chapterSeq": 1,
	"lessonSeq": 2,
	"lessonId": 66832
	},
	"externalKey": "KNRLESS39003",
	"continue": 0,
	"time_left": 2400000,
	"progress": 0,
	"totalDuration": 2400000,
	"startedAt": null,
	"expectedScore": -1,
	"recent_trial_id": "447455I0T24977530723"
	},
	"system": null,
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfdG9rZW5FeHBpcmVUaW1lIjo1NDAwLCJjb25uZWN0ZWRBdCI6IjIwMjEtMDktMjdUMDA6NDY6MDguMTgwWiIsInVzZXJJRCI6NDQ3NDU1LCJ1c2VyTmFtZSI6IiIsImV4dGVybmFsQWNjb3VudCI6bnVsbCwiY3VycmljdWx1bV9pZCI6OCwiY3VycmljdWx1bV90eXBlIjoiVkFDQVRJT04iLCJkaWZmaWN1bHR5IjoyNSwicGF5bWVudENoZWNrIjpmYWxzZSwidG9rZW5fYWNjb3VudCI6ImVsZTQydjkwIiwicHJvZHVjdFR5cGUiOiJFTEVNIiwic3RhbXAiOiJzaGExJDU2Y2Q5Njg0JDEkYTUyZTBhNjZhZDkwOTY1YjUxMjc2OWFmZjQ0YjcxYTA1Y2Y0MDJlNCIsImNsaWVudFZlcnNpb24iOjEwMjAwMDMsImFwaV92ZXJzaW9uIjoidjIiLCJ1c2VyVHlwZSI6IktOT1dSRV9URVNUIiwidGljayI6MiwibGVzc29uSW5mbyI6eyJsZXNzb25UeXBlIjoiTEVTU09OIiwiY3VycmljdWx1bUlkIjo4LCJjdXJyaWN1bHVtVHlwZSI6IlZBQ0FUSU9OIiwiY29tcG9zaXRlSWQiOiJLTlJMRVNTMzkwMDMiLCJjaGFwdGVyU2VxIjoxLCJsZXNzb25TZXEiOjIsImxlc3NvbklkIjo2NjgzMn0sImxlc3NvbkV4dGVybmFsS2V5IjoiS05STEVTUzM5MDAzIiwiYWN0aXZpdHkiOm51bGwsInRyeUNvdW50IjpudWxsLCJyZXRyeSI6bnVsbCwiaWF0IjoxNjMyNzAzNTc4LCJleHAiOjE2MzI3MDg5Nzh9.a9s8DJV9nVjOIFpmd4gypH8Q-sS7CyX0IQ2SpC163OM"
}`
	answers := []int64{437901, 437904, 437920, 437935, 437960, 437968, 438113, 438152}
	var problems []int64
	jsonparser.ArrayEach([]byte(jsonData), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		problemID, e := jsonparser.GetInt(value, "prob")
		if e != nil {
			log.Fatal(e)
		}
		problems = append(problems, problemID)

	}, "data", "problems")
	assert.EqualValues(t, answers, problems)

}
