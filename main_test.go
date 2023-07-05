package main

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestIsTwitterLink(t *testing.T) {
	tests := []struct {
		content      string
		expectedLink string
		expectedBool bool
	}{
		{
			content:      "Check out this tweet: https://twitter.com/username/status/12345",
			expectedLink: "Check out this tweet: https://vxtwitter.com/username/status/12345",
			expectedBool: true,
		},
		{
			content:      "No Twitter link here",
			expectedLink: "",
			expectedBool: false,
		},
		{
			content:      "Another tweet: https://twitter.com/another_user/status/67890",
			expectedLink: "Another tweet: https://vxtwitter.com/another_user/status/67890",
			expectedBool: true,
		},
		{
			content: "https://twitter.com/another_user",
			expectedLink: "",
			expectedBool: false,
		},
	}

	for _, test := range tests {
		log := logrus.New().WithField("test_case", test.content)
		link, ok := isTwitterLink(test.content, log)

		if link != test.expectedLink {
			t.Errorf("Expected replaced link: %s, got: %s", test.expectedLink, link)
		}

		if ok != test.expectedBool {
			t.Errorf("Expected bool value: %t, got: %t", test.expectedBool, ok)
		}
	}
}