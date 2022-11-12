package inmemory

import "time"

type UpdatedAtKey time.Time

func (t UpdatedAtKey) GetKey() string {
	return "key_sorted_set_updated_at"
}

var LastUpdatedAtKey UpdatedAtKey
