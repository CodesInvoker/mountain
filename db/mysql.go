package db

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/CodesInvoker/mountain/config"
	gorp "gopkg.in/gorp.v1"
)

var ProjectDBMap, WikiDBMap, CloudDBMap *gorp.DbMap

func InitDB() error {
	var err error
	ProjectDBMap, err = OpenDB(
		config.StringOrPanic("project_db_user"),
		config.StringOrPanic("project_db_password"),
		config.StringOrPanic("project_db_ip"),
		config.StringOrPanic("project_db_name"),
	)
	if err != nil {
		return err
	}

	return nil
}

func InitCloudDB() error {
	var err error
	CloudDBMap, err = OpenCloudDB(
		config.StringOrPanic("cloud_db_user"),
		config.StringOrPanic("cloud_db_password"),
		config.StringOrPanic("cloud_db_ip"),
		config.StringOrPanic("cloud_db_name"),
	)
	if err != nil {
		return err
	}
	return nil
}

func InitCloudDBSSL() error {
	var err error
	CloudDBMap, err = OpenCloudDBSSL()
	if err != nil {
		return err
	}
	return nil
}

func OpenDB(user, password, ip, name string) (*gorp.DbMap, error) {
	dataSourceName := fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, password, ip, name)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "latin1"}}
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "Sql:", log.Lmicroseconds))
	dbmap.TraceOff()
	return dbmap, nil
}

func OpenCloudDB(user, passoword, ip, name string) (*gorp.DbMap, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Asia%sShanghai&parseTime=true", user, passoword, ip, name, `%2F`)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	dbmap := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "utf8mb4",
		},
	}
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "Sql:", log.Lmicroseconds))
	dbmap.TraceOff()
	return dbmap, nil
}

func OpenCloudDBSSL() (*gorp.DbMap, error) {
	user := config.StringOrPanic("cloud_db_user")
	pw := config.StringOrPanic("cloud_db_password")
	ip := config.StringOrPanic("cloud_db_ip")
	name := config.StringOrPanic("cloud_db_name")
	sslCa := config.StringOrPanic("cloud_ssl_ca")
	sslCert := config.StringOrPanic("cloud_ssl_cert")
	sslKey := config.StringOrPanic("cloud_ssl_key")
	serverName := config.StringOrPanic("cloud_server_name")

	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(sslCa)
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}

	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(sslCert, sslKey)
	if err != nil {
		log.Fatal(err)
	}
	clientCert = append(clientCert, certs)

	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:      rootCertPool,
		Certificates: clientCert,
		ServerName:   serverName,
	})

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Asia%sShanghai&parseTime=true&tls=custom", user, pw, ip, name, `%2F`)
	db, err := sql.Open("mysql", dataSourceName)
	dbmap := &gorp.DbMap{
		Db: db,
		Dialect: gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "utf8mb4",
		},
	}
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "Sql:", log.Lmicroseconds))
	dbmap.TraceOff()
	return dbmap, nil
}
