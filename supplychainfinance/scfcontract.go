package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type SCFContract struct {
  CreditGrantingCompanyId    int      `json:"creditGrantingCompanyId"`
  CoreCompanyId              int      `json:"coreCompanyId"`
  FinancingNo                int      `json:"financingNo"`
  ContractNo                 int      `json:"contractNo"`
  SignDate                   int64    `json:"signDate"`
  Subject                    string   `json:"subject"`
  ContractTemplateId         int      `json:"contractTemplateId"`
}


// =====================================================================
//      CreateSCFContract - create scfContract and write to state
// =====================================================================
func (s *SupplyChainFinance) CreateSCFContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  scfContract := new(SCFContract)
  scfContract.CreditGrantingCompanyId = -1
  scfContract.CoreCompanyId = -1
  scfContract.FinancingNo = -1
  scfContract.ContractNo = -1
  scfContract.SignDate = -1
  scfContract.ContractTemplateId = -1

  err := json.Unmarshal([]byte(args[0]), &scfContract)
	if err != nil {
		return shim.Error("Fail to unmarshal scfContract information string " + err.Error())
	}

  if scfContract.CreditGrantingCompanyId < 0 {
    return shim.Error("Please check whether argument(CreditGrantingCompanyId) exists and is legal")
  }
  if scfContract.CoreCompanyId < 0 {
    return shim.Error("Please check whether argument(CoreCompanyId) exists and is legal")
  }
  if scfContract.FinancingNo < 0 {
    return shim.Error("Please check whether argument(FinancingNo) exists and is legal")
  }
  if scfContract.ContractNo < 0 {
    return shim.Error("Please check whether argument(ContractNo) exists and is legal")
  }
  if scfContract.SignDate < 0 {
    return shim.Error("Please check whether argument(SignDate) exists and is legal")
  }
  if len(scfContract.Subject) < 0 {
    return shim.Error("Please check whether argument(Subject) exists and is a non-empty string")
  }
  if scfContract.ContractTemplateId < 0 {
    return shim.Error("Please check whether argument(ContractTemplateId) exists and is legal")
  }

  // ==== Create scfContract compositekey ====
	indexName := "scfContract"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(scfContract.CreditGrantingCompanyId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

  // ==== Check if scfContract already exists ====
	scfContractInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get scfContract: " + err.Error())
	} else if scfContractInfo != nil {
		return shim.Error("The scfContract already exists")
	}

  // ==== Marshal scfContractInfo to JSON ====
  scfContractAsBytes, err := json.Marshal(scfContract)
	if err != nil {
		return shim.Error(err.Error())
	}

  // ==== Save scfContract to state ====
  err = stub.PutState(indexKey, scfContractAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}
