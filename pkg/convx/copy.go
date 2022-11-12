package convx

import "encoding/json"

func Copy(fromThis, toThis interface{}) error {
	b, err := json.Marshal(fromThis)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &toThis)
}
