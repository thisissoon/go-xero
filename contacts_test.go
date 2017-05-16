package xero

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactIterator_url(t *testing.T) {
	type testcase struct {
		tname       string
		page        int
		expectedURL string
	}
	tt := []testcase{
		testcase{
			tname:       "page 1",
			page:        1,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=1",
		},
		testcase{
			tname:       "page 2",
			page:        2,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=2",
		},
		testcase{
			tname:       "page 3",
			page:        3,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=3",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			i := ContactIterator{tc.page, &Client{}}
			assert.Equal(t, tc.expectedURL, i.url())
		})
	}
}
