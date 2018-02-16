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
	entry, _ := client.GetEntry("6036126")
	if entry == nil {
		t.Errorf("Failed to get entry")
	}
}
func TestGetNotExistingEntry(t *testing.T) {
	_, err := client.GetEntry("29963685")
	if err.(*ErrorResponse).ErrorObject.Code != 61 {
		t.Errorf("Somehow managed to get removed entry")
	}
}
func GetConversations(t *testing.T) {
	if client.userKey == "" {
		t.Errorf("You have to be logged to perform this action")
		return
	}
	_, err := client.GetConversationList()
	if err != nil {
		t.Errorf("Failed to get entries. Error: %v", err)
	}
}
func ObserveUser(t *testing.T) {
	if client.userKey == "" {
		t.Errorf("You have to be logged to perform this action")
		return
	}
	_, err := client.Observe("sokytsinolop")
	if err != nil {
		t.Errorf("Failed to Observe user Error:\n%v", err)
	}
}
func UnObserveUser(t *testing.T) {
	if client.userKey == "" {
		t.Errorf("You have to be logged to perform this action")
		return
	}
	_, err := client.Unobserve("sokytsinolop")
	if err != nil {
		t.Errorf("Failed to Unobserve user Error:\n%v", err)
	}
}
func TestAuthorizedEndpoints(t *testing.T) {
	t.Run("authorizationneeded", func(t *testing.T) {
		t.Run("GetConversations", GetConversations)
		t.Run("Observe", ObserveUser)
		t.Run("Unobserve", UnObserveUser)
	})
}
