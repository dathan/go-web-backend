// Package register handles the registration logic
package register_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/System-Glitch/goyave/v3"
	"github.com/dathan/go-web-backend/pkg/http/routes"
	"github.com/davecgh/go-spew/spew"
)

type CustomTestSuite struct {
	goyave.TestSuite
}

func (suite *CustomTestSuite) TestRegister() {

	var tt = []struct {
		name     string
		input    map[string]interface{}
		respCode int
		ok       bool
	}{
		{"fail", map[string]interface{}{"invalid": "invalid"}, 422, false},
	}

	suite.RunServer(routes.Register, func() {
		headers := map[string]string{"Content-Type": "application/json"}
		for _, tc := range tt {
			// TODO make a test table (tt) of test cases for the input and the expecte4d response
			body, _ := json.Marshal(tc.input)
			resp, err := suite.Post("/auth/register", headers, bytes.NewReader(body))
			suite.Nil(err)
			suite.NotNil(resp)
			if resp != nil {
				defer resp.Body.Close()
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					suite.Error(err)
				}
				suite.Equal(tc.respCode, resp.StatusCode)
				if tc.respCode != 200 && tc.ok == true {
					data := make(map[string]interface{})
					json.Unmarshal(bodyBytes, &data)
					suite.T().Log("Invalid return")
					suite.T().Fatal(spew.Sdump(data))
				}
			}
		}
	})
}

//TestRegister executes a set of tests all which should pass,
//where we send the approprate and invalid requests for the type of response expected.
func TestRegister(t *testing.T) {
	goyave.RunTest(t, new(CustomTestSuite))
}
