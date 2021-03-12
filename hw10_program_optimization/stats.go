package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"log"
	_ "net/http/pprof"
	"strings"

	"github.com/valyala/fastjson"
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

/*type User struct {
	ID       int
	Name     []byte
	Username []byte
	Email    []byte
	Phone    []byte
	Password []byte
	Address  []byte
}*/

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	//fmt.Printf("len: %v \n", u)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

//type users [100_000]User

func getUsers(r io.Reader) (result []User, err error) {
	//i := 0
	var p fastjson.Parser
	//user := new(User)
	var user User
	//var tmp []User
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		v, err := p.Parse(line)
		if err != nil {
			log.Fatalf("cannot parse json: %s", err)
		}
		user.ID = v.GetInt("ID")
		user.Name = string(v.GetStringBytes("Name"))
		user.Username = string(v.GetStringBytes("Username"))
		user.Email = string(v.GetStringBytes("Email"))
		user.Phone = string(v.GetStringBytes("Phone"))
		user.Password = string(v.GetStringBytes("Password"))
		user.Address = string(v.GetStringBytes("Address"))
		/*user.Name 		= v.GetStringBytes("Name")
		user.Username 	= v.GetStringBytes("Username")
		user.Email    	= v.GetStringBytes("Email")
		user.Phone    	= v.GetStringBytes("Phone")
		user.Password 	= v.GetStringBytes("Password")
		user.Address  	= v.GetStringBytes("Address")*/
		result = append(result, user)
		//result[i] = user
		//i++
	}
	//fmt.Printf("len: %v \n", len(result))
	if err := scanner.Err(); err != nil {
		fmt.Println("reading standard input:", err)
	}

	return result, nil
}

func countDomains(u []User, domain string) (DomainStat, error) {
	result := make(DomainStat)
	//fmt.Printf("len: %s \n", u[0].Name)
	for _, user := range u {
		//fmt.Printf("%s = %s\n", user.Email, []byte(domain))
		if len(user.Email) == 0 {
			continue
		}
		//fmt.Printf("%s = %s\n", user.Email, []byte(domain))
		if strings.Contains(user.Email, domain) {
			//if bytes.Contains(user.Email, []byte(domain)) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[key]
			num++
			result[key] = num
		}
	}
	return result, nil
}
