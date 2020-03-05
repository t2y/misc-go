package main

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
)

func createMock(t *testing.T) *MockS3API {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return NewMockS3API(ctrl)
}

func TestGetObject(t *testing.T) {
	var data = []struct {
		key  string
		body io.ReadCloser
	}{
		{
			"single-line.txt",
			ioutil.NopCloser(strings.NewReader("dummy\n")),
		},

		{
			"multiple-line.txt",
			ioutil.NopCloser(strings.NewReader("one\ntwo\nthree\n")),
		},
	}

	m := createMock(t)
	for _, tt := range data {
		output := &s3.GetObjectOutput{
			Body: tt.body,
		}

		// set return an object when GetObject() would called
		m.EXPECT().
			GetObject(gomock.Any()).
			Return(output, nil)

		getObject(m, "bucket", tt.key)
	}
}
