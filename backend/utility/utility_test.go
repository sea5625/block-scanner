package utility

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestIsAlphanumericString(t *testing.T) {

	assert.Equal(t, true, IsAlphanumericString("asdf"))
	assert.Equal(t, false, IsAlphanumericString("asdf!!"))
}

func TestIsIdenticalSlice(t *testing.T) {
	stringSlice1 := []string{"a", "b", "c"}
	stringSlice2 := []string{"b", "a", "c"}
	stringSlice3 := []string{"a", "b", "d"}

	assert.Equal(t, true, IsIdenticalSlice(stringSlice1, stringSlice2))
	assert.Equal(t, false, IsIdenticalSlice(stringSlice1, stringSlice3))
}

func TestIsExistValueInList(t *testing.T) {
	intList := []int{1, 2, 3, 4, 5}
	stringList := []string{"a", "b", "c"}

	assert.Equal(t, true, IsExistValueInList(3, intList))
	assert.Equal(t, false, IsExistValueInList(7, intList))
	assert.Equal(t, true, IsExistValueInList("a", stringList))
	assert.Equal(t, false, IsExistValueInList("e", stringList))
}

func TestConvertInterfaceToStringSlice(t *testing.T) {
	var interfaceSlice1 interface{} = []string{"string 1", "string 2"}

	stringSlice1 := ConvertInterfaceToStringSlice(interfaceSlice1)

	assert.Equal(t, "[]string", GetTypeOf(stringSlice1))
}

func TestGetTypeOf(t *testing.T) {
	int1 := 1
	intSlice1 := []int{1, 2, 3}
	string1 := "a"
	stringSlice1 := []string{"a", "b", "c"}

	assert.Equal(t, "int", GetTypeOf(int1))
	assert.Equal(t, "[]int", GetTypeOf(intSlice1))
	assert.Equal(t, "string", GetTypeOf(string1))
	assert.Equal(t, "[]string", GetTypeOf(stringSlice1))
}

func TestMatchTypeOf(t *testing.T) {
	int1 := 1
	intSlice1 := []int{1, 2, 3}
	string1 := "a"
	stringSlice1 := []string{"a", "b", "c"}

	assert.Equal(t, true, IsMatchTypeOf(int1, "int"))
	assert.Equal(t, true, IsMatchTypeOf(intSlice1, "[]int"))
	assert.Equal(t, true, IsMatchTypeOf(string1, "string"))
	assert.Equal(t, true, IsMatchTypeOf(stringSlice1, "[]string"))

	assert.Equal(t, false, IsMatchTypeOf(int1, "float"))
	assert.Equal(t, false, IsMatchTypeOf(intSlice1, "float"))
	assert.Equal(t, false, IsMatchTypeOf(string1, "float"))
	assert.Equal(t, false, IsMatchTypeOf(stringSlice1, "float"))
}

func TestGetHostAndPortInURL(t *testing.T) {
	host := "192.168.12.54:8080"
	url := "https://" + host

	hostAndPort, _ := GetHostAndPortInURL(url)

	assert.Equal(t, host, hostAndPort)
}
