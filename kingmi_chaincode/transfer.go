package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type Transfer struct {
  KingmiId             int     `json:"kingmiId"`      // 原始金米单号
  Transferee           int     `json:"transferee"`    // 被转让方企业
  Amount               int64   `json:"amount"`        // 转让金额
  TransferDate         int64   `json:"transferDate"`  // 转让日期
  DueDate              int64   `json:"dueDate"`       // 转让单到期日
  Note                 string  `json:"note"`          // 备注
}


type TransferQuery struct {
  KingmiId             int     `json:"kingmiId"`
  Transferee           int     `json:"transferee"`
  TransferDate         int64   `json:"transferDate"`
  DueDate              int64   `json:"dueDate"`
}


// =============================================
//      Transfer - Kingmi account transfer
// =============================================
func (m *KingmiChaincode) Transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  transfer := new(Transfer)
  transfer.KingmiId = -1
  transfer.Transferee = -1
  transfer.Amount = -1
  transfer.TransferDate = -1
  transfer.DueDate = -1

	err := json.Unmarshal([]byte(args[0]), &transfer)
	if err != nil {
		return shim.Error("Fail to unmarshal transfer information string " + err.Error())
	}

  if transfer.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exits and is legal")
  }
  if transfer.Transferee < 0 {
    return shim.Error("Please check whether argument(Transferee) exits and is legal")
  }
	if transfer.Amount <= 0 {
		return shim.Error("Please check whether argument(Amount) exits and it must be greater than 0")
	}
	if transfer.TransferDate <= 0 {
		return shim.Error("Please check whether argument(TransferDate) exits and is legal")
	}
	if transfer.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exits and is legal")
	}
  if transfer.TransferDate >= transfer.DueDate {
    return shim.Error("DueDate must be later than TransferDate")
  }

	// ==== Create transfer compositekey ====
	indexName := "transfer"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(transfer.KingmiId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if transferInfo already exists ====
	transferInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get transferInfo: " + err.Error())
	} else if transferInfo != nil {
		return shim.Error("The transferInfo already exists")
	}

	// ==== Marshal transferInfo to JSON ====
  transferAsBytes, err := json.Marshal(transfer)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save application to state ====
  err = stub.PutState(indexKey, transferAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// ==============================================================================
//       QueryTransfer - Query the transfer from chaincode state exactly
// ==============================================================================
func (m *KingmiChaincode) QueryTransfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  fmt.Printf("args:\n%s\n\n", args)

  transferQuery := new(TransferQuery)
  transferQuery.KingmiId = -1
  transferQuery.Transferee = -1
  transferQuery.TransferDate = -1
  transferQuery.DueDate = -1

  err := json.Unmarshal([]byte(args[0]), &transferQuery)
	if err != nil {
		return shim.Error("Fail to unmarshal transfer queryString " + err.Error())
	}

  if transferQuery.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exits and is legal")
  }
  if transferQuery.Transferee < 0 {
    return shim.Error("Please check whether argument(Transferee) exits and is legal")
  }
  if transferQuery.TransferDate <= 0 {
		return shim.Error("Please check whether argument(TransferDate) exits and is legal")
	}
  if transferQuery.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exits and is legal")
	}
  if transferQuery.TransferDate >= transferQuery.DueDate {
    return shim.Error("DueDate must be later than TransferDate")
  }

	queryString := fmt.Sprintf("{\"selector\":" + args[0] + "}")
	queryResults, err := getResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}
