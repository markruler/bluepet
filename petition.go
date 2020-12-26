package bluepet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseURL              = "https://www1.president.go.kr/api/petitions/list"
	petitionNumberByPage = 7
	limitPages           = 50
)

// BlueHousePetition 청와대 국민청원 한 페이지에 대한 오브젝트입니다.
type BlueHousePetition struct {
	Status string     `json:"status"`
	Total  string     `json:"total"`
	Page   int        `json:"page"`
	Item   []Petition `json:"item"`
}

// Petition 하나의 국민청원 데이터를 나타내는 오브젝트입니다.
// source: https://www1.president.go.kr/petitions
type Petition struct {
	JSONID    string `json:"id"`
	PagingID  int    `json:"paging_id"`
	Title     string `json:"title"`
	Agreement string `json:"agreement"`
	Category  string `json:"category"`
	Created   string `json:"created"`
	Finished  string `json:"finished"`
	Provider  string `json:"provider"`
}

// GetPetitions 실시간 국민청원 데이터를 가져옵니다.
func GetPetitions(category, only, order int) ([]Petition, error) {
	channelPetition := make(chan []Petition, 7)
	defer close(channelPetition)
	petitions := make([]Petition, 0, 7)

	bluepet := RequestToBlueHouse(category, only, order, 1)
	totalPages, err := GetTotalPages(bluepet.Total)
	if err != nil {
		return nil, err
	}

	// 청와대 제한: 약 3분 간 403 Forbidden
	// 현재까지 maximum: 350 = 50(pages) * 7(petitions) * 1(request)
	totalPages = limitPages
	channelPetition <- bluepet.Item
	for index := 2; index <= totalPages; index++ {
		go func(index int) {
			channelPetition <- RequestToBlueHouse(category, only, order, index).Item
		}(index)
	}
	for index := 1; index <= totalPages; index++ {
		pagePetitions := <-channelPetition
		petitions = append(petitions, pagePetitions...)
	}
	return petitions, nil
}

// RequestToBlueHouse HTTP 요청을 통해 한 페이지에 해당하는
// 국민청원 데이터를 가져옵니다. 청와대 국민청원 오픈 API가
// 별도로 제공되지 않기 때문에 홈페이지에서 직접 참조하여 만든
// 요청 함수입니다. 아래는 제가 임의로 추측한 파라미터 정보입니다.
//
//  parameter  | value         | description
//  -----------+---------------+-------------
//  c          | 0             | 전체
//             | 35            | 정치개혁
//             | 36            | 외교/통일/국방
//             | 37            | 일자리
//             | 38            | 미래
//             | 39            | 성장동력
//             | 40            | 농산어촌
//             | 41            | 보건복지
//             | 42            | 육아/교육
//             | 43            | 안전/환경
//             | 44            | 저출산/고령화대책
//             | 45            | 행정
//             | 46            | 반려동물
//             | 47            | 교통/건축/국토
//             | 48            | 경제민주화
//             | 49            | 인권/성평등
//             | 50            | 문화/예술/체육/언론
//             | 51            | 기타
//  -----------+---------------+-------------
//  only       | 1             | 진행 중 청원
//             | 2             | 만료된 청원
//  -----------+---------------+-------------
//  order      | 1             | 최신순 보기
//             | 2             | 추천순 보기
//  -----------+---------------+-------------
//  page       | number        | 현재 페이지 번호
func RequestToBlueHouse(category, only, order, page int) (bluepet *BlueHousePetition) {
	form := url.Values{
		"c":     {strconv.Itoa(category)},
		"only":  {strconv.Itoa(only)},
		"order": {strconv.Itoa(order)},
		"page":  {strconv.Itoa(page)},
	}

	request, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || 300 <= response.StatusCode {
		panic(response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &bluepet)
	if err != nil {
		panic(err)
	}

	return bluepet
}
