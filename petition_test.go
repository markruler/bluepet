package bluepet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPetitions(t *testing.T) {
	petitions, err := GetPetitions(0, 1, 2)
	// for _, petition := range petitions {
	// 	t.Log(petition)
	// }
	assert.NoError(t, err)
	assert.Len(t, petitions, 350)
}
