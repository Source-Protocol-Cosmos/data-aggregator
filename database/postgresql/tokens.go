package postgresql

import (
	"github.com/Source-Protocol-Cosmos/juno/v3/types"
)

// SaveToken allows to save token info into the database.
func (db *Database) SaveToken(token types.Token) error {
	stmt := `
INSERT INTO tokens (address, name, symbol, decimals, supply) 
VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Sql.Exec(stmt, token.Contract, token.Name, token.Symbol, token.Decimals, token.InitalSupply)
	return err
}
