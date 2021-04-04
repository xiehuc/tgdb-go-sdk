//
// session.go
// Copyright (C) 2018 toraxie <toraxie@tencent.com>
//
// Distributed under terms of the Tencent license.
//

package tgdb

import (
	"github.com/mindstand/go-bolt/connection"
)

type Session struct {
	db   string
	conn connection.IConnection
	p    *Client // parent
}

type Params = connection.QueryParams
type IResult = connection.IResult

func (s *Session) Close() error {
	return s.conn.Close()
}

func (s *Session) Cypher(q string, p Params) ([][]interface{}, IResult, error) {
	if p == nil {
		p = make(map[string]interface{})
	}
	p["cypher"] = q
	// TODO FailureMessage should not print stack
	return s.conn.Query(s.db+"->cypher", p)
}

// Run 是 底层操作
func (s *Session) Run(m Command, p Params) ([][]interface{}, IResult, error) {
	return s.conn.Query(s.db+"->"+m.Name(), p)
}
