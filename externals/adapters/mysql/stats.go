package mysql

import "database/sql"

// GetStats returns the database connection stats
func (a *Adapter) GetStats() sql.DBStats {
	return a.Pool.Stats()
}
