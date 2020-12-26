package bluepet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalPages(t *testing.T) {
	var (
		num int
		err error
	)

	num, err = GetTotalPages("a")
	assert.Error(t, err)

	num, err = GetTotalPages("13")
	assert.NoError(t, err)
	assert.Equal(t, 2, num)

	num, err = GetTotalPages("14")
	assert.NoError(t, err)
	assert.Equal(t, 2, num)

	num, err = GetTotalPages("15")
	assert.NoError(t, err)
	assert.Equal(t, 3, num)
}

func TestWriteInCSV(t *testing.T) {
	var err error
	err = WritePetitionsInCSV(nil)
	assert.Error(t, err)
	err = WritePetitionsInCSV([]Petition{
		{
			JSONID: "123",
		},
	})
	assert.NoError(t, err)
}
