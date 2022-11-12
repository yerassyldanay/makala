package convx

import "strconv"

func NewFeedVersion(version string) (int32, error) {
	i, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		return 0, err
	}

	return int32(i), err
}
