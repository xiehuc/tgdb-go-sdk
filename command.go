//
// command.go
// Copyright (C) 2018 toraxie <toraxie@tencent.com>
//
// Distributed under terms of the Tencent license.
//

package tgdb

import ()

type Command interface {
	Name() string
}
