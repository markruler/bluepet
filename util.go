package bluepet

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"
)

// GetTotalPages 국민청원 총 페이지 수를 구합니다.
func GetTotalPages(totalPages string) (totalNumber int, err error) {
	totalNumber, err = strconv.Atoi(totalPages)
	if err != nil {
		return -1, err
	}
	result := totalNumber / petitionNumberByPage
	if totalNumber%petitionNumberByPage > 0 {
		return result + 1, nil
	}
	return result, nil
}

// WritePetitionsInCSV CSV 형식으로 청원 정보를 저장하기
// 위한 함수입니다. 저장하는 내용은 JSONID, 제목, 참여인원수
// 세 가지입니다.
func WritePetitionsInCSV(petitions []Petition) error {
	if petitions == nil {
		return errors.New("'petitions' is empty")
	}

	file, err := os.Create(strconv.FormatInt(time.Now().Unix(), 10) + "-petitions.csv")
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	defer w.Flush()
	headers := []string{"JSONID", "title", "agree"}
	if err = w.Write(headers); err != nil {
		return err
	}

	for _, petition := range petitions {
		petitionCSV := []string{petition.JSONID, petition.Title, petition.Agreement}
		if err = w.Write(petitionCSV); err != nil {
			return err
		}
	}
	return nil
}
