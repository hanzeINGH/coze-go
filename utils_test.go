package coze

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateRandomString(t *testing.T) {
	str1, err := generateRandomString(10)
	assert.Nil(t, err)
	str2, err := generateRandomString(10)
	assert.Nil(t, err)
	assert.NotEqual(t, str1, str2)
}

func Test_MustToJson(t *testing.T) {
	jsonStr := mustToJson(map[string]string{"test": "test"})
	assert.Equal(t, jsonStr, `{"test":"test"}`)

}
