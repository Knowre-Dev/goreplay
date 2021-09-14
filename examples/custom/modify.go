/*
로그인 이후 쿠키가 바뀔때마다 적용하게 해주는 미들웨어
엘라스틱서치의 옛 데이터로는 토큰의 문제가 있을 수 있음
1. 시간을 과거로 세팅하여 토큰을 유효하게 만든다.
 - aws의 시간을 과거로 세팅할 수 없었음
2. 엘라스틱서치에 로깅된 리퀘스트와 토큰으로는 정상적 처리가 되지 않는 문제로 김현승 샘에게 로그인 시 토큰의 인증을 하지 않도록 수정을 요청
 - 요청 후 테스트 하여 다른 API에서 토큰을 새로발급하여 사용하는 로직이 있음
 - 메인문제출제 API를 호출하면 레슨상태에 맞추어 토큰이 새로 발급되는데, 이 때 문제가 될 수 있음.
3. 엘라스틱서치에서 리퀘스트를 날린 후 온 응답의 토큰을 분석하여 ID를 추출, Map ID=token 으로 저장 후 다음 리퀘스트 처리 시 ID를 알아내어 토큰을
   바꿔치기함
 - 미들웨어 개발
./goreplay -input-elasticsearch-address https://vpc-sl-logstrg-orange-prd-q76s3uteh4ooxa3r4brwce2yau.ap-northeast-2.es.amazonaws.com -input-elasticsearch-index cwl-raw-2021.09.01 --input-elasticsearch-fromDate 2021-08-31T15:00 -input-elasticsearch-toDate 2021-09-01T15:00 -middleware "./custom"  -input-elasticsearch-match /ecs/krdky-stable -output-http-track-response -output-http http://local.knowreapp.com

*/

package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/buger/goreplay/proto"
	"github.com/buger/jsonparser"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"strconv"
)

const (
	COOKIE       = "cookie"
	SID          = "connect.sid"
	XAccessToken = "x-access-token"
)

var (
	XAccessTokens = []string{"data", "accessToken"}
)

// requestID -> originalToken
var originalTokens *TTLMap

// originalToken -> replayedToken
var xaccessToken *TTLMap

func main() {
	run(bufio.NewScanner(os.Stdin))
}

func run(scanner *bufio.Scanner) {
	originalTokens = NewTTLMap(3600)
	xaccessToken = NewTTLMap(3600)

	//scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		encoded := scanner.Bytes()
		//log.Println(string(encoded[:1534]))
		buf := make([]byte, len(encoded)/2)
		hex.Decode(buf, encoded)

		process(buf)
	}
}

func process(buf []byte) {
	// First byte indicate payload type, possible values:
	//  1 - Request
	//  2 - Response
	//  3 - ReplayedResponse

	if buf == nil || len(buf) < 1 {
		return
	}

	payloadType := buf[0]
	metaSize := bytes.IndexByte(buf, '\n') + 1
	if metaSize < 0 {
		return
	}

	metaHeader := buf[:metaSize-1]
	meta := bytes.Split(metaHeader, []byte(" "))
	if len(meta) < 0 {
		return
	}

	reqID := string(meta[1])
	payload := buf[metaSize:]

	body := proto.Body(payload)

	switch payloadType {
	case '1': // Request
		//ES에서 가져온 request
		//refresh x-access-token
		{
			oldXToken := proto.Header(payload, []byte(XAccessToken))
			if len(oldXToken) > 0 {
				//Debug(string(oldXToken))
			}

			account, err := extractUserID(oldXToken)
			if err != nil {
				log.Println(err)
			} else {
				if xToken, ok := xaccessToken.Get(account); ok {
					payload = proto.SetHeader(payload, []byte(XAccessToken), []byte(xToken))
					buf = append(buf[:metaSize], payload...)

				}
			}
		}

		//refresh cookie
		{
			rawCookies := proto.Header(payload, []byte(COOKIE))
			cMap := CookieMap{}
			cMap.Parse(string(rawCookies))
			sid := cMap[SID]

			account, err := extractUserID([]byte(sid))
			if err == nil {
				if cookie, ok := originalTokens.Get(account); ok {
					cMap.Parse(string(cookie))
					//Debug(cMap.String())
					payload = proto.SetHeader(payload, []byte(COOKIE), []byte(cMap.String()))
					buf = append(buf[:metaSize], payload...)
				}
			}
		}

		// Emitting data back
		os.Stdout.Write(encode(buf))
	case '2': // Original response
		//if _, ok := originalTokens[reqID]; ok {
		//	// Token is inside response body
		//	secureToken := proto.Body(payload)
		//	originalTokens[reqID] = secureToken
		//	Debug("Remember origial token:", string(secureToken))
		//}
	case '3': // Replayed response
		status, err := strconv.Atoi(string(proto.Status(payload)))
		if err != nil {

		} else {
			if status > 400 {
				_ = reqID
				//status code가 에러인 애는 리스폰스를 저장한다.
				//Debug(reqID)
				//
				//f, fErr := os.OpenFile(fmt.Sprintf("./err_%s.data", reqID), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				//
				//if fErr != nil {
				//	panic(fErr)
				//}
				//f.Write(payload)
				//f.Close()
			}
		}

		//set-cookie 헤더가 리스폰스에 오면 다음리퀘스트에 바꿔줄 준비
		cookie := proto.Header(payload, []byte("Set-Cookie"))
		c := CookieMap{}
		if len(cookie) > 0 {
			c.Parse(string(cookie))
			if sid, ok := c[SID]; ok {
				userID, err := extractUserID([]byte(sid))
				if err != nil {
					return
				}
				//Debug(userID, " - ", string(cookie))
				originalTokens.Put(userID, string(cookie))
			}
		}

		//response의 body에서 accessToken 추출
		xaccessTokenValue, xaccount, xerr := extractUserIDFromBody(body, extractUserID, XAccessTokens...)
		if xerr == nil {
			//Debug(xaccount, " - ", xaccessTokenValue)
			xaccessToken.Put(xaccount, xaccessTokenValue)
		}

	}
}

func extractUserID(token []byte) (string, error) {
	var account string
	tokenStr := string(token)
	if t, _ := jwt.Parse(tokenStr, nil); t != nil {
		m := t.Claims.(jwt.MapClaims)
		if m != nil {
			if m["userID"] != nil {
				account = fmt.Sprintf("%d", int(m["userID"].(float64)))
				return account, nil
			}

			if m["user_id"] != nil {
				account = fmt.Sprintf("%d", int(m["user_id"].(float64)))
				return account, nil
			}

			return "", errors.New("not found ID")
		}
	}
	return account, errors.New("not found")
}

type Function func([]byte) (string, error)

func extractUserIDFromBody(payload []byte, fn Function, keys ...string) (string, string, error) {
	var value string
	var account string
	var err error

	if value, err = jsonparser.GetString(payload, keys...); err == nil {
		if fn != nil {
			account, err = fn([]byte(value))
		}
	}
	return value, account, err
}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
}

func Debug(args ...interface{}) {
	if os.Getenv("GOR_TEST") == "" { // if we are not testing
		fmt.Fprint(os.Stderr, "[DEBUG][TOKEN-MOD] ")
		fmt.Fprintln(os.Stderr, args...)
	}
}
