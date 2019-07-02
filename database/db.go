package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func init() {
	var err error

	// connect to cluster
	cluster := gocql.NewCluster("127.0.0.1")
	// cluster := gocql.NewCluster("database-cassandra")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "newpasscass"}
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra init done.")

	// create keyspaces
	err = Session.Query("CREATE KEYSPACE IF NOT EXISTS pdskeyspace WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	cluster.Keyspace = "pdskeyspace"

	// create table
	// admin == admin(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.admin(id UUID, fullname text, password text, email text, phone int, tokens text, PRIMARY KEY (id));").Exec()

	// dedevice == device-registeration(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.dedevice(id UUID, devicename text, registerdate text, zone text, PRIMARY KEY (id));").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// dezone == device-zone(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.dezone(id UUID, registerdate text, location text, PRIMARY KEY (id));").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// user == user(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.user(id UUID, fullname text, phone int, meta map<text,text>, PRIMARY KEY (id));").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// veblacklist == vehicle-blacklist(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.veblacklist(id UUID, devicename text, platenumber text, registerdate text, location text, roll text, reason text, PRIMARY KEY (id));").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// vemanager == vehicle-manager(model)
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.vemanager(id UUID, vehiclename text, platenumber text, registerdate text, location text, roll text, PRIMARY KEY (id));").Exec()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	err = Session.Query("CREATE TABLE IF NOT EXISTS pdskeyspace.commitlog(id UUID, date text, meta map<text,text>, PRIMARY KEY (id));").Exec()
}
