package aws

import "encoding/json"

func createEmailPayload(email, token string) ([]byte, error) {
	var payload = make(map[string]string)

	payload["email"] = email
	payload["token"] = token

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return body, nil

}
