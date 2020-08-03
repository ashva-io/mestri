package mestri

import "fmt"

const (
	host     = "ec2-34-193-117-204.compute-1.amazonaws.com"
	port     = 5432
	user     = "oszrwdkweikqbw"
	password = "2b606c4bd60baf639c557547c1fb0a38f414d774af7c7e00ae12fd6eb064a18d"
	dbname   = "d27aq3mo3jlkv6"
)

// PsqlInfo : PsqlInfo is const of connection string to connect to DB
var PsqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s", host, port, user, password, dbname)

//var PsqlInfo = "postgres://yzydeggkwlzlzu:abb532e1351fcef2db086340610d29ecd385e6af18c212010f4bde63290c4d9f@ec2-54-234-44-238.compute-1.amazonaws.com:5432/dabc8bfcspdah9"
