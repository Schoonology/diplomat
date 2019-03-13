package runners_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/mocks"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/runners"
)

func streamTests(tests []parsers.Test) chan parsers.Test {
	testChannel := make(chan parsers.Test)

	go func() {
		for _, test := range tests {
			testChannel <- test
		}
	}()

	return testChannel
}

func TestRunSerial(t *testing.T) {
	assert := assert.New(t)

	client := mocks.Client{}
	differ := mocks.Differ{}

	subject := runners.Serial{
		Client: &client,
		Differ: &differ,
	}
	tests := streamTests([]parsers.Test{
		parsers.Test{
			Request:  http.NewRequest("METHOD", "path"),
			Response: http.NewResponse(200, "STATUS TEXT"),
		},
	})
	errors := make(chan error)

	client.
		On("Do", http.NewRequest("METHOD", "path")).
		Return(http.NewResponse(200, "STATUS TEXT"), nil)
	differ.
		On("Diff", http.NewResponse(200, "STATUS TEXT"), http.NewResponse(200, "STATUS TEXT")).
		Return("some diff", nil)

	results := subject.Run(tests, errors)

	result := <-results

	client.AssertExpectations(t)
	differ.AssertExpectations(t)

	assert.NotNil(result)
}
