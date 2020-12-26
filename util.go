package bluepet

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatusCode(response *http.Response) {
	if response.StatusCode != 200 {
		log.Fatal("Status Code:", response.StatusCode)
	}
}

// WritePetitionsInCSV CSV 형식으로 청원 정보를 저장하기
// 위한 함수입니다. 저장하는 내용은 JSONID, 제목, 참여인원수
// 세 가지입니다.
func WritePetitionsInCSV(petitions []Petition) {
	file, err := os.Create("petitions.csv")
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"JSONID", "title", "agree"}
	err = w.Write(headers)
	checkError(err)

	for _, petition := range petitions {
		petitionCSV := []string{petition.JSONID, petition.Title, petition.Agreement}
		err = w.Write(petitionCSV)
		checkError(err)
	}
}
