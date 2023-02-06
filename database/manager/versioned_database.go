// Copyright (C) 2022-2023, Dijets Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package manager

import (
	"github.com/lasthyphen/dijetsnodego/database"
	"github.com/lasthyphen/dijetsnodego/utils"
	"github.com/lasthyphen/dijetsnodego/version"
)

var _ utils.Sortable[*VersionedDatabase] = (*VersionedDatabase)(nil)

type VersionedDatabase struct {
	Database database.Database
	Version  *version.Semantic
}

// Close the underlying database
func (db *VersionedDatabase) Close() error {
	return db.Database.Close()
}

// Note this sorts in descending order (newest version --> oldest version)
func (db *VersionedDatabase) Less(other *VersionedDatabase) bool {
	return db.Version.Compare(other.Version) > 0
}
