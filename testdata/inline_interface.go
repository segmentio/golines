package fixtures

// Based on example in https://github.com/segmentio/golines/issues/40
func CreateEphemeralRecordEvent(ctx interface {
	string
	int
	string
}, id string, name string, kaid string, districtID string, status string, anotherArg string, aThirdArg string) (*string, error) {
	return nil, nil
}
