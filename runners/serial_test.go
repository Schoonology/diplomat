package runners_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/mocks"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/runners"
)

func TestRunSerial(t *testing.T) {
	assert := assert.New(t)

	client := mocks.Client{}
	differ := mocks.Differ{}

	subject := runners.Serial{
		Client: &client,
		Differ: &differ,
	}
	spec := parsers.Spec{
		Tests: []parsers.Test{
			parsers.Test{
				Request:  http.NewRequest("METHOD", "path"),
				Response: http.NewResponse(200, "STATUS TEXT"),
			},
		},
	}

	client.
		On("Do", http.NewRequest("METHOD", "path")).
		Return(http.NewResponse(200, "STATUS TEXT"), nil)
	differ.
		On("Diff", http.NewResponse(200, "STATUS TEXT"), http.NewResponse(200, "STATUS TEXT")).
		Return("some diff", nil)

	result, err := subject.Run(&spec)
	assert.Nil(err)
	assert.NotNil(result)

	client.AssertExpectations(t)
	differ.AssertExpectations(t)
}