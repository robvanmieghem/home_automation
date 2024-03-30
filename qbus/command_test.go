package qbus

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGroupsCommand(t *testing.T) {
	c := NewGetGroupsCommand()
	jsonEncoded, err := json.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, `{"Type":10,"Value":null}`, string(jsonEncoded))
}

func TestErrors(t *testing.T) {
	//No error
	jsonResponse := `{"Type":1,"Value":null}`
	r := &Response{}
	assert.NoError(t, json.Unmarshal([]byte(jsonResponse), r))
	assert.Nil(t, r.Error())
	//Test if returned error responsecodes are converted to the correct errors
	type testcase struct {
		ErrorCode     int
		ExpectedError error
	}
	testCases := []testcase{
		{1, ErrBusy},
		{2, ErrSessionTimeout},
		{3, ErrTooManyConnections},
		{4, ErrFailed},
		{5, ErrSessionStartFailure},
		{6, ErrUnknownCommand},
		{7, ErrNoEQWebConfig},
		{8, ErrSystemManagerConnected},
		{255, ErrUndefined},
	}

	for _, c := range testCases {
		jsonResponse := fmt.Sprintf(`{"Type":0,"Value":{"Error":%d}}`, c.ErrorCode)
		r := &Response{}
		assert.NoError(t, json.Unmarshal([]byte(jsonResponse), r))
		assert.Equal(t, c.ExpectedError, r.Error())
	}
}
