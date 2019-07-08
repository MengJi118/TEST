package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type SCFLoanRepayment struct {
  CreditGrantingCompanyId    int     `json:"creditGrantingCompanyId"` // 授信企业id
  FinancingNo                int     `json:"financingNo"`             // 融资编号
  RepaymentAmount            int     `json:"repaymentAmount"`         // 还款金额
  RepaymentInterest          int     `json:"repaymentInterest"`       // 实际还款利息
  RepaymentPrincipal         int     `json:"repaymentPrincipal"`      // 实际还款本金
}


// ===============================================================================
//      CreateSCFLoanRepayment - create scfLoanRepayment and write to state
// ===============================================================================
func (s *SupplyChainFinance) CreateSCFLoanRepayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  scfLoanRepayment := new(SCFLoanRepayment)
  scfLoanRepayment.CreditGrantingCompanyId = -1
  scfLoanRepayment.FinancingNo = -1
  scfLoanRepayment.RepaymentAmount = -1
  scfLoanRepayment.RepaymentInterest = -1
  scfLoanRepayment.RepaymentPrincipal = -1

  err := json.Unmarshal([]byte(args[0]), &scfLoanRepayment)
	if err != nil {
		return shim.Error("Fail to unmarshal scfLoanRepayment information string " + err.Error())
	}

  if scfLoanRepayment.CreditGrantingCompanyId < 0 {
    return shim.Error("Please check whether argument(CreditGrantingCompanyId) exists and is legal")
  }
  if scfLoanRepayment.FinancingNo < 0 {
    return shim.Error("Please check whether argument(FinancingNo) exists and is legal")
  }
  if scfLoanRepayment.RepaymentAmount < 0 {
    return shim.Error("Please check whether argument(RepaymentAmount) exists and is a non-empty string")
  }
  if scfLoanRepayment.RepaymentInterest < 0 {
    return shim.Error("Please check whether argument(RepaymentInterest) exists and is legal")
  }
  if scfLoanRepayment.RepaymentPrincipal < 0 {
    return shim.Error("Please check whether argument(RepaymentPrincipal) exists and is legal")
  }

  // ==== Create scfLoanRepayment compositekey ====
	indexName := "scfLoanRepayment"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(scfLoanRepayment.CreditGrantingCompanyId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

  // ==== Check if scfLoanRepayment already exists ====
	scfLoanRepaymentInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get scfLoanRepayment: " + err.Error())
	} else if scfLoanRepaymentInfo != nil {
		return shim.Error("The scfLoanRepayment already exists")
	}

  // ==== Marshal scfLoanRepaymentInfo to JSON ====
  scfLoanRepaymentAsBytes, err := json.Marshal(scfLoanRepayment)
	if err != nil {
		return shim.Error(err.Error())
	}

  // ==== Save scfLoanRepayment to state ====
  err = stub.PutState(indexKey, scfLoanRepaymentAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}
