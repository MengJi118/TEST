package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type SCFLoan struct {
  CreditGrantingCompanyId   int      `json:"creditGrantingCompanyId"`// 授信企业id
  CoreCompanyId             int      `json:"coreCompanyId"`          // 核心企业id
  FinancingNo               int      `json:"financingNo"`            // 融资编号
  FinancingAmount           int      `json:"financingAmount"`        // 融资金额
  Term                      int      `json:"term"`                   // 融资期限
  CredentialNo              int      `json:"CredentialNo"`           // 凭证编号（多个）
  PayeeName                 string   `json:"payeeName"`              // 收款方名称
  LoanAmount                int      `json:"loanAmount"`             // 实际放款金额
  LoanDate                  int64    `json:"loanDate"`               // 放款日期
  DueDate                   int64    `json:"dueDate"`                // 到期日
  Interest                  int      `json:"interest"`               // 利率
  ApplyDate                 int64    `json:"applyDate"`              // 申请日期
}


// ===================================================================
//      CreateSCFFinance - create scfFinance and write to state
// ===================================================================
func (s *SupplyChainFinance) CreateSCFFinance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  scfLoan := new(SCFLoan)
  scfLoan.CreditGrantingCompanyId = -1
  scfLoan.CoreCompanyId = -1
  scfLoan.FinancingNo = -1
  scfLoan.FinancingAmount = -1
  scfLoan.Term = -1
  scfLoan.CredentialNo = -1
  scfLoan.Interest = -1
  scfLoan.ApplyDate = -1

  err := json.Unmarshal([]byte(args[0]), &scfLoan)
	if err != nil {
		return shim.Error("Fail to unmarshal scfLoan information string " + err.Error())
	}

  if scfLoan.CreditGrantingCompanyId < 0 {
    return shim.Error("Please check whether argument(CreditGrantingCompanyId) exists and is legal")
  }
  if scfLoan.CoreCompanyId < 0 {
    return shim.Error("Please check whether argument(CoreCompanyId) exists and is legal")
  }
  if scfLoan.FinancingNo < 0 {
    return shim.Error("Please check whether argument(FinancingNo) exists and is legal")
  }
  if scfLoan.FinancingAmount < 0 {
    return shim.Error("Please check whether argument(FinancingAmount) exists and is legal")
  }
  if scfLoan.Term < 0 {
    return shim.Error("Please check whether argument(Term) exists and is legal")
  }
  if scfLoan.CredentialNo < 0 {
    return shim.Error("Please check whether argument(CredentialNo) exists and is legal")
  }
  if scfLoan.Interest < 0 {
    return shim.Error("Please check whether argument(Interest) exists and is legal")
  }
  if scfLoan.ApplyDate < 0 {
    return shim.Error("Please check whether argument(ApplyDate) exists and is legal")
  }

  // ==== Create scfLoan compositekey ====
	indexName := "scfLoan"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(scfLoan.CreditGrantingCompanyId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

  // ==== Check if scfLoan already exists ====
	scfLoanInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get scfLoan: " + err.Error())
	} else if scfLoanInfo != nil {
		return shim.Error("The scfLoan already exists")
	}

  // ==== Marshal scfLoanInfo to JSON ====
  scfLoanAsBytes, err := json.Marshal(scfLoan)
	if err != nil {
		return shim.Error(err.Error())
	}

  // ==== Save scfLoan to state ====
  err = stub.PutState(indexKey, scfLoanAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// =============================================================
//      CreateSCFLoan - create scfLoan and write to state
// =============================================================
func (s *SupplyChainFinance) CreateSCFLoan(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  scfLoan := new(SCFLoan)
  scfLoan.CreditGrantingCompanyId = -1
  scfLoan.FinancingNo = -1
  scfLoan.LoanAmount = -1
  scfLoan.LoanDate = -1
  scfLoan.DueDate = -1

  err := json.Unmarshal([]byte(args[0]), &scfLoan)
	if err != nil {
		return shim.Error("Fail to unmarshal scfLoan information string " + err.Error())
	}

  if scfLoan.CreditGrantingCompanyId < 0 {
    return shim.Error("Please check whether argument(CreditGrantingCompanyId) exists and is legal")
  }
  if scfLoan.FinancingNo < 0 {
    return shim.Error("Please check whether argument(FinancingNo) exists and is legal")
  }
  if len(scfLoan.PayeeName) < 0 {
    return shim.Error("Please check whether argument(PayeeName) exists and is a non-empty string")
  }
  if scfLoan.LoanAmount < 0 {
    return shim.Error("Please check whether argument(LoanAmount) exists and is legal")
  }
  if scfLoan.LoanDate < 0 {
    return shim.Error("Please check whether argument(LoanDate) exists and is legal")
  }
  if scfLoan.DueDate < 0 {
    return shim.Error("Please check whether argument(DueDate) exists and is legal")
  }

  // ==== Create scfLoan compositekey ====
	indexName := "scfLoan"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(scfLoan.CreditGrantingCompanyId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

  // ==== Check if scfLoan already exists ====
	scfLoanInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get scfLoan: " + err.Error())
	} else if scfLoanInfo != nil {
		return shim.Error("The scfLoan already exists")
	}

  // ==== Marshal scfLoanInfo to JSON ====
  scfLoanAsBytes, err := json.Marshal(scfLoan)
	if err != nil {
		return shim.Error(err.Error())
	}

  // ==== Save scfLoan to state ====
  err = stub.PutState(indexKey, scfLoanAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}
