package browsersteps

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/DATA-DOG/godog"
	"github.com/tebeka/selenium"
)

func iWaitFor(amount int, unit string) error {
	u := time.Second
	fmt.Printf("Waiting for %d %s", amount, unit)
	time.Sleep(u * time.Duration(amount))
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I wait for (\d+) (milliseconds|millisecond|seconds|second)$`, iWaitFor)

	// selenium.SetDebug(true)
	capabilities := selenium.Capabilities{"browserName": "chrome"}
	capEnv := os.Getenv("SELENIUM_CAPABILITIES")
	if capEnv != "" {
		err := json.Unmarshal([]byte(capEnv), &capabilities)
		if err != nil {
			log.Panic(err)
		}
	}
	bs, _ := NewBrowserSteps(s,
		capabilities,
		os.Getenv("SELENIUM_URL"))

	var server *httptest.Server
	s.BeforeSuite(func() {
		server = httptest.NewServer(http.FileServer(http.Dir("./public")))
		u, _ := url.Parse(server.URL)
		bs.SetBaseURL(u)
	})

	s.AfterSuite(func() {
		if server != nil {
			server.Close()
			server = nil
		}
	})
}

func TestMain(m *testing.M) {
	status := godog.Run("browsersteps", FeatureContext)
	os.Exit(status)
}
