package wykop

import (
	"testing"
)

const APIKEY = "aNd401dAPp"
const SECRET = "3lWf1lCxD6" //Wykop Android Application API key and secret

var client = Create(APIKEY, SECRET)

func TestInitializeClient(t *testing.T) {
	LocalClient := Create(APIKEY, SECRET)
	if LocalClient == nil {
		t.Errorf("Couldnt create wykop client instance")
	}
}
func TestGetEntry(t *testing.T) {
	entry := client.GetEntry("6036126")
	if entry == nil {
		t.Errorf("Failed to get entry")
	}
}
func TestGetNotExistingEntry(t *testing.T) {
	entry := client.GetEntry("29963685")
	if entry != nil {
		t.Errorf("Somehow managed to get removed entry")
	}
}
