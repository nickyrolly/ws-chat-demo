package domain

const (
	HostPostgreSQL = "postgres:5432"
	HostNSQlookupd = "nsqlookupd:4160"
	HostNSQd       = "nsqd:4150"
	HostNSQadmin   = "nsqadmin:4171"
)

var (
	ServiceMap = map[string]string{
		"PostgreSQL": HostPostgreSQL,
		"NSQlookupd": HostNSQlookupd,
		"NSQd":       HostNSQd,
		"NSQadmin":   HostNSQadmin,
	}
)
