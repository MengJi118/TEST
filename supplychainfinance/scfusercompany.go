package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

type SupplyChainFinance struct {
}


type SCFUserCompany struct {
  CompanyName                  string   `json:"companyName"`            // 企业名称
  CreditGrantingCompanyId      int      `json:"creditGrantingCompanyId"`// 授信企业id
  CertificateNo                int      `json:"certificateNo"`          // 证件号码
  PassTime                     int64    `json:"passTime"`               // 实名认证时间
}


// ==============
//      Main
// ==============
func main() {
	err := shim.Start(new(SupplyChainFinance))
	if err != nil {
		fmt.Printf("Error starting SupplyChainFinance chaincode: %s", err)
	}
}


// =========================================
//       Init - Initializes chaincode
// =========================================
func (s *SupplyChainFinance) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


// ======================================================
//       Invoke - Our entry point for Invocations
// ======================================================
func (s *SupplyChainFinance) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

  if function == "CreateSCFUserCompany" {
    return s.CreateSCFUserCompany(stub, args)
  } else if function == "CreateSCFFinance" {
    return s.CreateSCFFinance(stub, args)
  } else if function == "CreateSCFLoan" {
    return s.CreateSCFLoan(stub, args)
  } else if function == "CreateSCFContract" {
    return s.CreateSCFContract(stub, args)
  } else if function == "CreateSCFLoanRepayment" {
    return s.CreateSCFLoanRepayment(stub, args)
  } else {
		return shim.Error("Function " + function + " doesn't exits, make sure function is right!")
	}
}


// ===========================================================================
//      CreateSCFUserCompany - create SCFUserCompany and write to state
// ===========================================================================
func (s *SupplyChainFinance) CreateSCFUserCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  scfUserCompany := new(SCFUserCompany)
  scfUserCompany.CreditGrantingCompanyId = -1
  scfUserCompany.CertificateNo = -1
  scfUserCompany.PassTime = -1

	err := json.Unmarshal([]byte(args[0]), &scfUserCompany)
	if err != nil {
		return shim.Error("Fail to unmarshal scfUserCompany information string " + err.Error())
	}

  if len(scfUserCompany.CompanyName) < 0 {
    return shim.Error("Please check whether argument(CompanyName) exists and is a non-empty string")
  }
  if scfUserCompany.CreditGrantingCompanyId < 0 {
    return shim.Error("Please check whether argument(CreditGrantingCompanyId) exists and is legal")
  }
  if scfUserCompany.CertificateNo < 0 {
    return shim.Error("Please check whether argument(CertificateNo) exists and is legal")
  }
  if scfUserCompany.PassTime < 0 {
    return shim.Error("Please check whether argument(PassTime) exists and is legal")
  }

	// ==== Create scfUserCompany compositekey ====
	indexName := "scfUserCompany"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(scfUserCompany.CreditGrantingCompanyId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if scfUserCompany already exists ====
	scfUserCompanyInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get scfUserCompany: " + err.Error())
	} else if scfUserCompanyInfo != nil {
		return shim.Error("The scfUserCompany already exists")
	}

	// ==== Marshal scfUserCompanyInfo to JSON ====
  scfUserCompanyAsBytes, err := json.Marshal(scfUserCompany)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save scfUserCompany to state ====
  err = stub.PutState(indexKey, scfUserCompanyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}
