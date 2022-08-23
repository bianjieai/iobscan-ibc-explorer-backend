package redis

import (
	"fmt"
	"testing"
	"time"

	v8 "github.com/go-redis/redis/v8"
)

type people struct {
	Name   string
	Age    int
	Friend []string
}

var client *Client

func TestMain(m *testing.M) {
	client = New("127.0.0.1:6379", "", "", "single", 0)
	m.Run()
}

func TestRedisFunc(t *testing.T) {
	key := "vis"
	value := "sss"

	if err := client.Set(key, value, time.Hour*2); err != nil {
		t.Fatal(err)
	}
	v, err := client.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if v != value {
		t.Fatal(err)
	}

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}

	key = "myhash"
	if _, err := client.HSet(key, "k1", "v1", "k2", "v2"); err != nil {
		t.Fatal(err)
	}

	if _, err = client.HSet(key, []string{"k3", "v3", "k4", "v4"}); err != nil {
		t.Fatal(err)
	}

	if _, err = client.HSet(key, map[string]interface{}{"k5": "v5", "k6": "v6"}); err != nil {
		t.Fatal(err)
	}

	v, err = client.HGet(key, "k1")
	if err != nil {
		t.Fatal(err)
	}
	if v != "v1" {
		t.Fatal(v)
	}

	get, err := client.HGet("sss", "sss")
	if err != nil && err != v8.Nil {
		t.Fatal(err)
	}
	fmt.Println(get)

	all, err := client.HGetAll(key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(all)

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}

	key = "myset"
	if _, err := client.SAdd(key, "s1", "s2", "s3"); err != nil {
		t.Fatal(err)
	}
	members, err := client.SMembers(key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(members)

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}
}

func TestRedisMarshalFunc(t *testing.T) {
	key := "ss"
	value := people{
		Name:   "kamir",
		Age:    13,
		Friend: []string{"catlin", "lee"},
	}

	if err := client.MarshalSet(key, value, 0); err != nil {
		t.Fatal(err)
	}
	var v people
	if err := client.UnmarshalGet(key, &v); err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}

	key = "myset"
	value2 := people{
		Name:   "jax",
		Age:    13,
		Friend: []string{"catlin", "lee"},
	}
	if _, err := client.MarshalSAdd(key, value, value2); err != nil {
		t.Fatal(err)
	}
	var members []people
	if err := client.UnmarshalSMembers(key, &members); err != nil {
		t.Fatal(err)
	}
	for _, v := range members {
		fmt.Println(v)
	}

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}

	key = "myhash"
	vm := map[string]people{"p1": value, "p2": value2}
	if _, err := client.MarshalHSet(key, vm); err != nil {
		t.Fatal(err)
	}

	var p1 people
	if err := client.UnmarshalHGet(key, "p1", &p1); err != nil {
		t.Fatal(err)
	}
	fmt.Println(p1)

	var pMap map[string]people
	if err := client.UnmarshalHGetAll(key, &pMap); err != nil {
		t.Fatal(err)
	}
	fmt.Println(pMap)

	if _, err := client.Del(key); err != nil {
		t.Fatal(err)
	}

}
