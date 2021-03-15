package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	/*var json = jsoniter.Config{
		EscapeHTML:                    	false,
		OnlyTaggedField: 				true,
		ObjectFieldMustBeSimpleString: 	true,
		ValidateJsonRawMessage: 		true,
	}.Froze()*/

	i := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email := jsoniter.Get(scanner.Bytes(), "Email").ToString()
		result[i].Email = email
		i++
	}
	if err = scanner.Err(); err != nil {
		err = fmt.Errorf("invalid input: %w", err)
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, user := range u {
		if strings.HasSuffix(user.Email, domain) {
			key := strings.Split(user.Email, "@")
			if len(key) > 1 {
				result[strings.ToLower(key[1])]++
			}
		}
	}
	return result, nil
}
