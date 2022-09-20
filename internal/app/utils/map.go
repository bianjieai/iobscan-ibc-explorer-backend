package utils

func ContainsKey(m map[string]string, key string) bool {
	for k, _ := range m {
		if k == key {
			return true
		}
	}
	return false
}

func ContainsValue(m map[string]string, value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}

func MapKeys(m map[string]string) []string {
	var res []string
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}

func MapValues(m map[string]string) []string {
	var res []string
	for _, v := range m {
		res = append(res, v)
	}
	return res
}
