# 커스텀사항
# Knowre ElasticSearch에서 데이터를 가져오는 부분 추가
 - elasticsearch 의 input plugin 을 개발
 - elasticsearch 의 DSL Query를 요구하는 부분은 하드코딩되어 있음
 - 다음 옵션으로 실행하면 됨
    ```
    ./goreplay -input-elasticsearch-address https://vpc-sl-logstrg-orange-prd-q76s3uteh4ooxa3r4brwce2yau.ap-northeast-2.es.amazonaws.com -input-elasticsearch-index cwl-raw-2021.08.01 --input-elasticsearch-fromDate 2021-08-01T05:00 -input-elasticsearch-toDate 2021-08-01T06:00 -output-http http://localhost:9100 -http-original-host -input-elasticsearch-match /ecs/krdky-stable
    ```
  - 추가한 입력 파라미터 값 
   - input-elasticsearch-address : 엘라스틱서치 접속주소
   - input-elasticsearch-index : 엘라스틱서치 인덱스이름
   - input-elasticsearch-fromDate : 검색을 시작할 시간 2021-08-01T00:00
   - input-elasticsearch-toDate : 검색을 종료할 시간 2021-08-02T00:05
   - input-elasticsearch-match : 제품별 로그 검색을 위해 @log_group 필드에서 검색할 값 


# 주의사항
 - gor 로 파일을 작성하면 http 요청처럼 보이지만 실제로는 \r\n이 헤더에 붙어나온다. 재전송 시에는 http 라이브러리를 사용하여 사라짐
 
#파일로 요청
./goreplay -input-file ./message.out -middleware "./custom"  -output-http-track-response -output-http http://local.knowreapp.com

#8월9일 한국시간 오전 09시부터 오전 10까지 질의하려면 UTC 0-1시로 세팅해야 한다.
./goreplay -input-elasticsearch-address https://vpc-sl-logstrg-orange-prd-q76s3uteh4ooxa3r4brwce2yau.ap-northeast-2.es.amazonaws.com -input-elasticsearch-index cwl-raw-2021.08.10 --input-elasticsearch-fromDate 2021-08-10T00:00 -input-elasticsearch-toDate 2021-08-10T01:00 -middleware "./custom"  -input-elasticsearch-match /ecs/krdky-stable -output-http-track-response -output-http http://local.knowreapp.com

