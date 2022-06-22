package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Basename removes directory components
// e.g., a => a, a.go => a.go, a/b/c.go => c.go
func Basename(s string) string {
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	s = s[slash+1:]
	return s
}

// PrefixName delete str after first slash
// e.g., a/b/c => a, fds/dd => fds, c.go//11 => c.go
func PrefixName(s string) string {
	slash := strings.Index(s, "/") // -1 if "/" not found
	if slash == -1 {
		return s
	}
	s = s[:slash]
	return s
}

// ParentPath delete str after last slash
// e.g., a/b/c => a/b, fds/dd => fds, c.go//11 => c.go/
func ParentPath(s string) string {
	slash := strings.LastIndex(s, "/")
	if slash == -1 {
		return s
	}
	return s[:slash]
}

// RemoveFirstLevelPath delete str before slash
// e.g., a/b/c => b/c, fds/dd => dd, c.go//11 => /11
func RemoveFirstLevelPath(s string) string {
	slash := strings.Index(s, "/")
	if slash == -1 {
		return s
	}
	return s[slash+1:]
}

func IsContain(sources []string, target string) bool {
	for _, source := range sources {
		if source == target {
			return true
		}
	}
	return false
}

func DistinctSliceStr(slice []string) []string {
	set := make(map[string]int)
	for _, v := range slice {
		set[v] = 0
	}

	res := make([]string, 0, len(set))
	for k, _ := range set {
		res = append(res, k)
	}
	return res
}

func HttpGet(url string) (bz []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode != 200")
	}

	bz, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}

func HttpPost(url string, reqBody interface{}) (bz []byte, err error) {
	reqBz := MarshalJsonIgnoreErr(reqBody)
	reader := strings.NewReader(string(reqBz))
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode != 200")
	}

	bz, err = ioutil.ReadAll(resp.Body)
	return bz, nil
}

func InArray(arr []string, e string) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}
	return false
}
