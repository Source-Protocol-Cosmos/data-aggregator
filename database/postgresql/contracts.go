package postgresql

import (
	"github.com/Source-Protocol-Cosmos/juno/v3/types"
)

// SaveContract allows to save the given contract into the database.
func (db *Database) SaveContract(contract types.Contract, gas, fees int64) error {
	stmt := `
INSERT INTO contracts (code_id, address, creator, admin, label, creation_time, height, gas, fees) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := db.Sql.Exec(stmt, contract.CodeID, contract.Address, contract.Creator, contract.Admin, contract.Label, contract.CreatedTime, contract.Created.BlockHeight, gas, fees)
	return err
}

// UpdateContractStats update stats by contract call.
func (db *Database) UpdateContractStats(contract string, tx, gas, fees int64) error {
	stmt := `
UPDATE contracts SET tx=tx+$2, gas=gas+$3, fees=fees+$4 WHERE address = $1`
	_, err := db.Sql.Exec(stmt, contract, tx, gas, fees)
	return err
}

// SaveContractCodeID save new contract CodeID.
func (db *Database) SaveContractCodeID(contract string, codeID uint64) error {
	stmt := `
UPDATE contracts SET code_id = $2 WHERE address = $1`
	_, err := db.Sql.Exec(stmt, contract, codeID)
	return err
}

// UpdateContractAdmin update contract admin.
func (db *Database) UpdateContractAdmin(contract string, admin string) error {
	stmt := `
UPDATE contracts SET admin = $2 WHERE address = $1`
	_, err := db.Sql.Exec(stmt, contract, admin)
	return err
}
