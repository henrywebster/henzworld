package statuscafe

import (
	"testing"
)

func TestGetStatus(t *testing.T) {
	expectedContent := "example status"
	expectedFace := "😌"
	expectedTimeAgo := "1 day ago"

	response := &Response{
		Content: expectedContent,
		Face:    expectedFace,
		TimeAgo: expectedTimeAgo,
	}

	status := response.GetStatus()

	if status == nil {
		t.Error("Expected value, got nil")
		return
	}
	if status.Content != expectedContent {
		t.Errorf("Expected Content to be %s, got: %s", expectedContent, status.Content)
		return
	}
	if status.Face != expectedFace {
		t.Errorf("Expected Face to be %s, got: %s", expectedFace, status.Face)
		return
	}
	if status.TimeAgo != expectedTimeAgo {
		t.Errorf("Expected TimeAgo to be %s, got: %s", expectedTimeAgo, status.TimeAgo)
		return
	}
}

func TestGetStatusEscaped(t *testing.T) {
	expectedContent := ">:)"

	response := &Response{
		Content: "&gt;:)",
	}

	status := response.GetStatus()

	if status == nil {
		t.Error("Expected value, got nil")
		return
	}
	if status.Content != expectedContent {
		t.Errorf("Expected Content to be %s, got: %s", expectedContent, status.Content)
		return
	}
}
