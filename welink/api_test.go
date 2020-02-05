package welink

import (
	"os"
	"sort"
	"testing"

	"go.uber.org/zap"

	"github.com/fhyx/welink-api-go/log"
)

var (
	api *API
)

func TestMain(m *testing.M) {
	_logger, _ := zap.NewDevelopment()
	defer _logger.Sync() // flushes buffer, if any
	sugar := _logger.Sugar()
	log.SetLogger(sugar)

	api = NewAPI()
	os.Exit(m.Run())
}

// TestAPIDepartment test api // WELINK_CORP_ID= WELINK_CORP_SECRET=
func TestAPIDepartment(t *testing.T) {

	data, err := api.ListDepartment(0, true)
	if err != nil {
		t.Fatal(err)
	}

	sort.Sort(data)

	for _, dept := range data {
		t.Logf("dept %v", dept)
	}

}

func TestUser(t *testing.T) {

	user, err := api.GetUser(os.Getenv("WELINK_TEST_UID"), "uid")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user %v", user)
}
