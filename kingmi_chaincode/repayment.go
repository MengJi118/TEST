package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type RepaymentHistory struct {
  FinancingId             int       `json:"financingId"`    // 融资单id
  DueRepayment            int64     `json:"dueRepayment"`   // 应还金额
  DueInterest             int64     `json:"dueInterest"`    // 应还利息
  DueDate                 int64     `json:"dueDate"`        // 融资到期日
  RepaymentAmount         int64     `json:"repaymentAmount"`// 实际还款金额
  RepaymentDate           int64     `json:"repaymentDate"`  // 实际还款日期
  OverduePayment          int64     `json:"overduePayment"` // 逾期罚息
  OverdueDays             int       `json:"overdueDays"`    // 逾期天数
}


type RepaymentHistoryQuery struct {
  FinancingId             int       `json:"financingId"`
  DueDate                 int64     `json:"dueDate"`
  RepaymentDate           int64     `json:"repaymentDate"`
}


// ==========================================================
//      WriteRepaymentHistory - Write repayment history
// ==========================================================
func (m *KingmiChaincode) WriteRepaymentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  repaymentHistory := new(RepaymentHistory)
  repaymentHistory.FinancingId = -1
  repaymentHistory.DueRepayment = -1
  repaymentHistory.DueInterest = -1
  repaymentHistory.DueDate = -1
  repaymentHistory.RepaymentAmount = -1
  repaymentHistory.RepaymentDate = -1

	err := json.Unmarshal([]byte(args[0]), &repaymentHistory)
	if err != nil {
		return shim.Error("Fail to unmarshal repaymentHistory information string " + err.Error())
	}

  if repaymentHistory.FinancingId < 0 {
    return shim.Error("Please check whether argument(FinancingId) exits and is legal")
  }
	if repaymentHistory.DueRepayment <= 0 {
		return shim.Error("Please check whether argument(DueRepayment) exits and it must be greater than 0")
	}
  if repaymentHistory.DueInterest <= 0 {
    return shim.Error("Please check whether argument(DueInterest) exits and it must be greater than 0")
  }
  if repaymentHistory.DueDate <= 0 {
    return shim.Error("Please check whether argument(DueDate) exits and is legal")
  }
  if repaymentHistory.RepaymentAmount <= 0 {
    return shim.Error("Please check whether argument(RepaymentAmount) exits and it must be greater than 0")
  }
  if repaymentHistory.RepaymentDate <= 0 {
		return shim.Error("Please check whether argument(RepaymentDate) exits and is legal")
	}

	// ==== Create repayment compositekey ====
	indexName := "repaymentHistory"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(repaymentHistory.FinancingId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if repaymentHistory already exists ====
	repaymentHistoryInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get repaymentHistoryInfo: " + err.Error())
	} else if repaymentHistoryInfo != nil {
		return shim.Error("The repaymentHistoryInfo already exists")
	}

	// ==== Marshal repaymentHistoryInfo to JSON ====
  repaymentHistoryAsBytes, err := json.Marshal(repaymentHistory)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save repaymentHistory to state ====
  err = stub.PutState(indexKey, repaymentHistoryAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// ===============================================================================
//       QueryRepaymentHistory - Query repaymentHistory from chaincode state
// ===============================================================================
func (m *KingmiChaincode) QueryRepaymentHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  fmt.Printf("args:\n%s\n\n", args)

  repaymentHistoryQuery := new(RepaymentHistoryQuery)
  repaymentHistoryQuery.FinancingId = -1
  repaymentHistoryQuery.DueDate = -1
  repaymentHistoryQuery.RepaymentDate = -1

  err := json.Unmarshal([]byte(args[0]), &repaymentHistoryQuery)
	if err != nil {
		return shim.Error("Fail to unmarshal repaymentHistory queryString " + err.Error())
	}

  if repaymentHistoryQuery.FinancingId < 0 {
    return shim.Error("Please check whether argument(FinancingId) exits and is legal")
  }
  if repaymentHistoryQuery.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exits and is legal")
	}
  if repaymentHistoryQuery.RepaymentDate <= 0 {
    return shim.Error("Please check whether argument(RepaymentDate) exits and is legal")
  }

	queryString := fmt.Sprintf("{\"selector\":" + args[0] + "}")
	queryResults, err := getResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}
