package db

import (
	"fmt"
	"testing"

	"github.com/CodesInvoker/mountain/config"
	"gopkg.in/gorp.v1"
)

func TestDB(t *testing.T) {
	err := config.LoadConfigs("")
	if err != nil {
		fmt.Println("load config error:", err)
		return
	}
	// err = InitDB()
	// if err != nil {
	// 	fmt.Println("init db error:", err)
	// 	return
	// }
	// users, err := ListUser(ProjectDBMap)
	// if err != nil {
	// 	fmt.Println("list user error:", err)
	// 	return
	// }
	// for _, u := range users {
	// 	fmt.Println("user=", *u)
	// }
	err = InitCloudDBSSL()
	if err != nil {
		fmt.Println("init db error:", err)
		return
	}
	users, err := ListCloudUser(CloudDBMap)
	if err != nil {
		fmt.Println("list user error:", err)
		return
	}
	for _, u := range users {
		fmt.Println("user=", *u)
	}
}

type User struct {
	UUID string `db:"uuid"`
	Name string `db:"name"`
}

type CloudUser struct {
	Email string `db:"email"`
	Name  string `db:"name"`
}

func ListUser(src gorp.SqlExecutor) ([]*User, error) {
	sql := "SELECT uuid, name FROM user;"
	users := make([]*User, 0)
	_, err := src.Select(&users, sql)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func ListCloudUser(src gorp.SqlExecutor) ([]*CloudUser, error) {
	sql := "SELECT email, name FROM user;"
	users := make([]*CloudUser, 0)
	_, err := src.Select(&users, sql)
	if err != nil {
		return nil, err
	}
	return users, nil
}
