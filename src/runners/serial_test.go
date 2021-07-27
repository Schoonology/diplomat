package runners_test

import (
	"testing"

	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/http"
	"github.com/schoonology/diplomat/mocks"
	"github.com/schoonology/diplomat/runners"
	"github.com/stretchr/testify/assert"
)

func TestRunSerial(t *testing.T) {
	assert := assert.New(t)

	client := mocks.Client{}
	differ := mocks.Differ{}

	subject := runners.Serial{
		Client: &client,
		Differ: &differ,
	}
	test := builders.Test{
		Request:  http.NewRequest("METHOD", "path"),
		Response: http.NewResponse(200, "STATUS TEXT"),
	}

	client.
		On("Do", http.NewRequest("METHOD", "path")).
		Return(http.NewResponse(200, "STATUS TEXT"), nil)
	differ.
		On("Diff", http.NewResponse(200, "STATUS TEXT"), http.NewResponse(200, "STATUS TEXT")).
		Return("some diff", nil)

	result, err := subject.Run(test)

	assert.Nil(err)

	client.AssertExpectations(t)
	differ.AssertExpectations(t)

	assert.NotNil(result)
}
