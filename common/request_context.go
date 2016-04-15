package common

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/nu7hatch/gouuid"
)

type RequestContext struct {
	DB        *gocql.Session
	RequestID string
}

func NewRequestContext(dbSession *gocql.Session) RequestContext {
	u4, err := uuid.NewV4()
	if err != nil {
		// NB: its virtually impossible to get test coverage on this line since this
		// never happens
		panic(fmt.Sprintf("Failed to generate a UUID: %v", err))
	}

	return RequestContext{
		DB:        dbSession,
		RequestID: u4.String(),
	}
}
