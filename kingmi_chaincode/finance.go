package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type Finance struct {
  FinancingId                  int      `json:"financingId"`                // 融资单号
  Amount                       int64    `json:"amount"`                     // 融资金额
  DueDate                      int64    `json:"dueDate"`                    // 融资到期日
  IfProtocolInterest           int      `json:"ifProtocolInterest"`         // 是否协议付息
  ProtocolInterestProportion   int      `json:"protocolInterestProportion"` // 协议付息比例
}


type FinanceQuery struct {
  FinancingId                  int      `json:"financingId"`
  DueDate                      int64    `json:"dueDate"`
  IfProtocolInterest           int      `json:"ifProtocolInterest"`
}


// ========================================
//      Financing - Kingmi financing
// ========================================
func (m *KingmiChaincode) Financing(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  finance := new(Finance)
  finance.FinancingId = -1
  finance.Amount = -1

	err := json.Unmarshal([]byte(args[0]), &finance)
	if err != nil {
		return shim.Error("Fail to unmarshal financing information string " + err.Error())
	}

  if finance.FinancingId < 0 {
    return shim.Error("Please check whether argument(FinancingId) exits and is legal")
  }
	if finance.Amount <= 0 {
		return shim.Error("Please check whether argument(Amount) exits and it must be greater than 0")
	}

	// ==== Create financing compositekey ====
	indexName := "financing"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(finance.FinancingId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if financeInfo already exists ====
	financeInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get financeInfo: " + err.Error())
	} else if financeInfo != nil {
		return shim.Error("The financeInfo already exists")
	}

	// ==== Marshal financeInfo to JSON ====
  financeAsBytes, err := json.Marshal(finance)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save financeInfo to state ====
  err = stub.PutState(indexKey, financeAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// ===========================================================================
//       QueryFinance - Query financing information from chaincode state
// ===========================================================================
func (m *KingmiChaincode) QueryFinance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  fmt.Printf("args:\n%s\n\n", args)

  financeQuery := new(FinanceQuery)
  financeQuery.FinancingId = -1
  financeQuery.DueDate = -1
  financeQuery.IfProtocolInterest = -1

  err := json.Unmarshal([]byte(args[0]), &financeQuery)
	if err != nil {
		return shim.Error("Fail to unmarshal financing queryString " + err.Error())
	}

  if financeQuery.FinancingId < 0 {
		return shim.Error("Please check whether argument(FinancingId) exits and is legal")
	}
  if financeQuery.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exits and is legal")
	}
  if financeQuery.IfProtocolInterest != 0 && financeQuery.IfProtocolInterest != 1 {
    return shim.Error("Please check whether argument(IfProtocolInterest) exits and it must only be 0 represents non-protocalInterest and 1 represents protocalInterest")
  }

	queryString := fmt.Sprintf("{\"selector\":" + args[0] + "}}")
	queryResults, err := getResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}
