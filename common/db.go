package common

import "github.com/gocql/gocql"

const (
	gophrDBProtoVersion = 4
	gophrDBKeyspaceName = "gophr"
)

func OpenDBConnection(address string, isDev bool) (*gocql.Session, error) {
	cluster := gocql.NewCluster(address)
	cluster.ProtoVersion = gophrDBProtoVersion
	cluster.Keyspace = gophrDBKeyspaceName
	if isDev {
		cluster.Consistency = gocql.One
	} else {
		cluster.Consistency = gocql.Quorum
	}
	return cluster.CreateSession()
}

func CloseDBConnection(session *gocql.Session) {
	session.Close()
}
