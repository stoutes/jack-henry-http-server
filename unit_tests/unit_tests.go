package unit_tests

import (
	"net/http"
	"testing"
)

// go test -v
// go test -run TestGetDownload/statusok -v
// go test -run TestGetDownload/statusnotfound -v

const succeed = "\u2713"
const failed = "\u2717"

// TestGetDownload validates the http Get function can download content.
func TestGetDownload(t *testing.T) {
	statusCode := 200
	url := "http://localhost:6666/getWeatherReport/666/666/666"
	t.Log("Given the need to test downloading content.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen checking %q for status code %d", testID, url, statusCode)
		{
			resp, err := http.Get(url)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to make the Get call : %v", failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to make the Get call.", succeed, testID)

			defer resp.Body.Close()

			if resp.StatusCode == statusCode {
				t.Logf("\t%s\tTest %d:\tShould receive a %d status code.", succeed, testID, statusCode)
			} else {
				t.Errorf("\t%s\tTest %d:\tShould receive a %d status code : %d", failed, testID, statusCode, resp.StatusCode)
			}
		}
	}
}
