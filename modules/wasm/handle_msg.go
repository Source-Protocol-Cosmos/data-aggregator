package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Source-Protocol-Cosmos/juno/v3/types"
	"github.com/Source-Protocol-Cosmos/juno/v3/types/cw20"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleMsg allows to handle the different utils related to the gov module
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *types.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}
	switch cosmosMsg := msg.(type) {
	case *wasmtypes.MsgStoreCode:
		return m.handleMsgStoreCode(tx, index, cosmosMsg)
	case *wasmtypes.MsgInstantiateContract:
		return m.handleMsgInstantiateContract(tx, index, cosmosMsg)
	case *wasmtypes.MsgExecuteContract:
		return m.handleMsgExecuteContract(tx, index, cosmosMsg)
	case *wasmtypes.MsgMigrateContract:
		return m.handleMsgMigrateContract(tx, index, cosmosMsg)
	case *wasmtypes.MsgClearAdmin:
		return m.handleMsgClearAdmin(tx, index, cosmosMsg)
	case *wasmtypes.MsgUpdateAdmin:
		return m.handleMsgUpdateAdmin(tx, index, cosmosMsg)
	}
	fmt.Println("nil")
    fmt.Println(reflect.TypeOf(msg))

	return nil
}

func (m *Module) handleMsgStoreCode(tx *types.Tx, index int, msg *wasmtypes.MsgStoreCode) error {
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeStoreCode)
	if err != nil {
		return err
	}

	codeIDVal, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyCodeID)
	if err != nil {
		return err
	}

	codeID, _ := strconv.Atoi(codeIDVal)
	response, err := m.client.Code(context.Background(), &wasmtypes.QueryCodeRequest{
		CodeId: uint64(codeID),
	})
	if err != nil {
		return err
	}

	hash := response.CodeInfoResponse.DataHash.String()
	size := len(response.Data)
	code := types.NewCode(uint64(codeID), msg.Sender, hash, uint64(size), tx.Timestamp, tx.Height)

	return m.db.SaveCode(code)
}

func (m *Module) handleMsgInstantiateContract(tx *types.Tx, index int, msg *wasmtypes.MsgInstantiateContract) error {
	contracts, err := GetAllContracts(tx, index, wasmtypes.EventTypeInstantiate)
	if err != nil {
		return err
	}

	if len(contracts) == 0 {
		return fmt.Errorf("No contract address found")
	}

	createdAt := &wasmtypes.AbsoluteTxPosition{
		BlockHeight: uint64(tx.Height),
		TxIndex:     uint64(index),
	}

	fee := tx.GetFee()
	feeAmount := int64(0)
	if fee.Len() == 1 {
		feeAmount = fee[0].Amount.Int64()
	}

	ctx := context.Background()
	for i, contractAddress := range contracts {
		response, err := m.client.ContractInfo(ctx, &wasmtypes.QueryContractInfoRequest{
			Address: contractAddress,
		})
		if err != nil {
			return err
		}

		creator, _ := sdk.AccAddressFromBech32(response.Creator)
		admin, _ := sdk.AccAddressFromBech32(response.Admin)
		contractInfo := wasmtypes.NewContractInfo(response.CodeID, creator, admin, response.Label, createdAt)
		contract := types.NewContract(&contractInfo, contractAddress, tx.Timestamp)

		if i == 0 {
			err = m.db.SaveContract(contract, tx.GasUsed, feeAmount)
		} else {
			err = m.db.SaveContract(contract, 0, 0)
		}

		if err != nil {
			return err
		}

		// Store code data
		data, err := m.db.GetCodeData(response.CodeID)
		if err != nil {
			return err
		}

		// Check if cw20 token
		var tokenInfo *cw20.TokenInfo
		if data.Version == nil || *data.CW20 {
			cw20Response, err := m.client.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
				Address:   contractAddress,
				QueryData: []byte(`{"token_info":{}}`),
			})
			if err == nil {
				var token cw20.TokenInfo
				err = json.Unmarshal(cw20Response.Data, &token)
				if err == nil {
					tokenInfo = &token
				}
			}
		}

		if data.Version == nil {
			res, err := m.client.RawContractState(ctx, &wasmtypes.QueryRawContractStateRequest{
				Address:   contractAddress,
				QueryData: []byte("contract_info"),
			})

			version := "none"
			if err == nil && res.Data != nil {
				version = string(res.Data)
			}

			isIBC := response.IBCPortID != ""
			isCW20 := tokenInfo != nil
			newData := types.NewCodeData(response.CodeID, version, isIBC, isCW20)
			err = m.db.SetCodeData(newData)
			if err != nil {
				return err
			}
		}

		// Save cw20 token
		if isValidCw20Token(tokenInfo) {
			token := types.NewToken(contractAddress, tokenInfo.Name, tokenInfo.Symbol, tokenInfo.Decimals, tokenInfo.TotalSupply)
			err = m.db.SaveToken(token)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Module) handleMsgExecuteContract(tx *types.Tx, index int, msg *wasmtypes.MsgExecuteContract) error {
	contracts, err := GetAllContracts(tx, index, wasmtypes.EventTypeExecute)
	fmt.Println("handleMsgExecuteContract 1")
	if err != nil {
		return err
	}
	fmt.Println("handleMsgExecuteContract 2")

	if len(contracts) == 0 {
		return fmt.Errorf("No contract address found")
	}
	fmt.Println("handleMsgExecuteContract 3")

	fee := tx.GetFee()
	fmt.Println("handleMsgExecuteContract 4")
	feeAmount := int64(0)
	fmt.Println("handleMsgExecuteContract 5")
	if fee.Len() == 1 {
		feeAmount = fee[0].Amount.Int64()
	}
	fmt.Println("handleMsgExecuteContract 6")

	for i, contract := range contracts {
		fmt.Println("handleMsgExecuteContract 7"+fmt.Sprintf("%d",i))
		if i == 0 {
			err = m.db.UpdateContractStats(contract, 1, tx.GasUsed, feeAmount)
		} else {
			err = m.db.UpdateContractStats(contract, 1, 0, 0)
		}

		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			return err
		}
	}
	fmt.Println("handleMsgExecuteContract 8")

	return nil
}

func (m *Module) handleMsgMigrateContract(tx *types.Tx, index int, msg *wasmtypes.MsgMigrateContract) error {

	return m.db.SaveContractCodeID(msg.Contract, msg.CodeID)
}

func (m *Module) handleMsgClearAdmin(tx *types.Tx, index int, msg *wasmtypes.MsgClearAdmin) error {

	return m.db.UpdateContractAdmin(msg.Contract, "")
}

func (m *Module) handleMsgUpdateAdmin(tx *types.Tx, index int, msg *wasmtypes.MsgUpdateAdmin) error {

	return m.db.UpdateContractAdmin(msg.Contract, msg.NewAdmin)
}

func GetAllContracts(tx *types.Tx, index int, eventType string) ([]string, error) {
	contracts := []string{}
	event, err := tx.FindEventByType(index, eventType)
	if err != nil {
		return contracts, err
	}

	for _, attr := range event.Attributes {
		if attr.Key == wasmtypes.AttributeKeyContractAddr {
			contracts = append(contracts, attr.Value)
		}
	}

	return contracts, nil
}

func isValidCw20Token(token *cw20.TokenInfo) bool {
	return token != nil && token.Name != "" && token.Symbol != ""
}
