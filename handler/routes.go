package handler

import (
	"context"

	"github.com/jtarchie/sqlettuce/db"
	"github.com/jtarchie/sqlettuce/router"
)

//nolint:funlen
func NewRoutes(
	ctx context.Context,
	client *db.Client,
) router.Command {
	commands := router.Command{
		"APPEND": appendRouter(ctx, client),
		"CONFIG": router.Command{
			"GET": router.Command{
				"save":       router.StaticResponseRouter(router.EmptyStringResponse),
				"appendonly": router.StaticResponseRouter("+no\r\n"),
			},
		},
		"COMMAND": router.Command{
			"DOCS": router.StaticResponseRouter(router.EmptyStringResponse),
		},
		"DECR":        decrRouter(ctx, client),
		"DECRBY":      decrByRouter(ctx, client),
		"DEL":         delRouter(ctx, client),
		"ECHO":        echoRouter(),
		"FLUSHALL":    flushAllRouter(ctx, client),
		"GET":         getRouter(ctx, client),
		"GETDEL":      getDelRouter(ctx, client),
		"GETRANGE":    getRangeRouter(ctx, client),
		"INCR":        incrRouter(ctx, client),
		"INCRBY":      incrByRouter(ctx, client),
		"INCRBYFLOAT": incrByFloatRouter(ctx, client),
		"LRANGE":      lrangeRouter(ctx, client),
		"MGET":        mgetRouter(ctx, client),
		"MSET":        msetRouter(ctx, client),
		"PING":        router.StaticResponseRouter("+PONG\r\n"),
		"RPUSH":       rpushRouter(ctx, client),
		"RPUSHX":      rpushXRouter(ctx, client),
		"SET":         setRouter(ctx, client),
		"STRLEN":      strlenRouter(ctx, client),

		// deprecated commands, let's not support them
		"RPOPLPUSH":  router.StaticResponseRouter("-Deprecated command, please use LMOVE with the RIGHT and LEFT\r\n"),
		"BRPOPLPUSH": router.StaticResponseRouter("-Deprecated command, please use LMOVE with the RIGHT and LEFT\r\n"),

		"GETSET": router.StaticResponseRouter("-Deprecated command, please use SET with the GET argument\r\n"),
		"PSETEX": router.StaticResponseRouter("-Deprecated command, please use SET with the PX argument\r\n"),
		"SETEX":  router.StaticResponseRouter("-Deprecated command, please use SET with the EX argument\r\n"),
		"SETNX":  router.StaticResponseRouter("-Deprecated command, please use SET with the NX argument\r\n"),
		"SUBSTR": router.StaticResponseRouter("-Deprecated command, please use GETRANGE\r\n"),
	}

	commands["FLUSHDB"] = commands["FLUSHALL"]
	commands["UNLINK"] = commands["DEL"]

	return commands
}
