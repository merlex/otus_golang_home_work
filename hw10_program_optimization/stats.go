package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
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
	user := &User{}
	result := make(DomainStat)
	domain = "." + domain
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		*user = User{}
		if err := json.Unmarshal(scanner.Bytes(), user); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}
		matched := strings.Contains(user.Email, domain)
		if matched {
			indFindStr := strings.Index(user.Email, "@")
			str := strings.ToLower(user.Email[indFindStr+1 : len(user.Email)])
			result[str]++
		}
	}

	return result, scanner.Err()
}
