package main


import (
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

type KingmiChaincode struct {
}


type Application struct {
  KingmiId                 int      `json:"kingmiId"`              // 金米单号
  OpenEnterpriseId         int      `json:"openEnterpriseId"`      // 开立企业id
  DownEnterpriseId         int      `json:"downEnterpriseId"`      // 供应商企业id
  FactoryingEnterpriseId   int      `json:"factoryingEnterpriseId"`// 保理公司id
  CredentialType           int      `json:"credentialType"`        // 凭证类型
  Amount                   int64    `json:"amount"`                // 金米金额
  CreatedDate              int64    `json:"createdDate"`           // 金米开立日期
  DueDate                  int64    `json:"dueDate"`               // 金米到期日
  IfGuarantee              int      `json:"ifGuarantee"`           // 是否担保
  IfTransfer               int      `json:"ifTransfer"`            // 是否转让
  IfDelay                  int      `json:"ifDelay"`               // 是否延期
  Description              string   `json:"description"`           // 金米信息描述
  Note                     string   `json:"note"`                  // 备注
}


type ApplicationQuery struct {
  KingmiId                 int      `json:"kingmiId"`
  OpenEnterpriseId         int      `json:"openEnterpriseId"`
  DownEnterpriseId         int      `json:"downEnterpriseId"`
  FactoryingEnterpriseId   int      `json:"factoryingEnterpriseId"`
  CredentialType           int      `json:"credentialType"`
  CreatedDate              int64    `json:"createdDate"`
  DueDate                  int64    `json:"dueDate"`
  IfGuarantee              int      `json:"ifGuarantee"`
  IfTransfer               int      `json:"ifTransfer"`
  IfDelay                  int      `json:"ifDelay"`
}


type Settlement struct {
  KingmiId                 int      `json:"kingmiId"`      // 金米单号
  DueDate                  int64    `json:"dueDate"`       // 金米到期日
  SoaDate                  int64    `json:"soaDate"`       // 实际结算日期
  Amount                   int64    `json:"amount"`        // 金米金额
  DueInterest              int64    `json:"dueInterest"`   // 利息
  OverdueDays              int      `json:"overdueDays"`   // 逾期天数
  OverduePenalty           int64    `json:"overduePenalty"`// 逾期罚息
}


type SettlementQuery struct {
  KingmiId                 int      `json:"kingmiId"`
  DueDate                  int64    `json:"dueDate"`
  SoaDate                  int64    `json:"soaDate"`
}


const (
  FUNCTION_NAME_ERROR = iota + 3000
  ARGUMENT_ERROR
)


// ==============
//      Main
// ==============
func main() {
	err := shim.Start(new(KingmiChaincode))
	if err != nil {
		fmt.Printf("Error starting Kingmi chaincode: %s", err)
	}
}


// =========================================
//       Init - Initializes chaincode
// =========================================
func (m *KingmiChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}


// ======================================================
//       Invoke - Our entry point for Invocations
// ======================================================
func (m *KingmiChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

  if function == "Apply" {
    return m.Apply(stub, args)
  } else if function == "QueryApplication" {
    return m.QueryApplication(stub, args)
  } else if function == "Settle" {
    return m.Settle(stub, args)
  } else if function == "QuerySettlement" {
    return m.QuerySettlement(stub, args)
  } else if function == "Transfer" {
    return m.Transfer(stub, args)
  } else if function == "QueryTransfer" {
    return m.QueryTransfer(stub, args)
  } else if function == "Financing" {
    return m.Financing(stub, args)
  } else if function == "QueryFinance" {
    return m.QueryFinance(stub, args)
  } else if function == "WriteRepaymentHistory" {
    return m.WriteRepaymentHistory(stub, args)
  } else if function == "QueryRepaymentHistory" {
    return m.QueryRepaymentHistory(stub, args)
  } else {
		return shim.Error("Function " + function + " doesn't exits, make sure function is right!")
	}
}


// ==========================================
//      Apply - Apply kingmi account
// ==========================================
func (m *KingmiChaincode) Apply(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  application := new(Application)
  application.KingmiId = -1
  application.OpenEnterpriseId = -1
  application.DownEnterpriseId = -1
  application.FactoryingEnterpriseId = -1
  application.CredentialType = -1
  application.Amount = -1
  application.CreatedDate = -1
  application.DueDate = -1
  application.IfGuarantee = -1
  application.IfTransfer = -1
  application.IfDelay = -1

	err := json.Unmarshal([]byte(args[0]), &application)
	if err != nil {
		return shim.Error("Fail to unmarshal apply information string " + err.Error())
	}

  if application.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exists and is legal")
  }
  if application.OpenEnterpriseId < 0 {
    return shim.Error("Please check whether argument(OpenEnterpriseId) exists and is legal")
  }
  if application.DownEnterpriseId < 0 {
    return shim.Error("Please check whether argument(DownEnterpriseId) exists and is legal")
  }
	if application.FactoryingEnterpriseId < 0 {
		return shim.Error("Please check whether argument(FactoryingEnterpriseId) exists and is legal")
	}
	if application.CredentialType != 1 && application.CredentialType != 2 {
		return shim.Error("Please check whether argument(CredentialType) exists and it must only be 1 represents accounts receivable or 2 represents accounts payable")
	}
	if application.Amount <= 0 {
		return shim.Error("Please check whether argument(Amount) exists and it must be greater than 0")
	}
  if application.CreatedDate <= 0 {
		return shim.Error("Please check whether argument(CreatedDate) exists and is legal")
	}
  if application.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exists and is legal")
	}
  if application.CreatedDate >= application.DueDate {
    return shim.Error("DueDate must be later than CreatedDate")
  }
	if application.IfGuarantee != 0 && application.IfGuarantee != 1 {
		return shim.Error("Please check whether argument(IfGuarantee) exists and it must only be 0 represents unguaranteed and 1 represents guaranteed")
	}
	if application.IfTransfer != 0 && application.IfTransfer != 1 {
		return shim.Error("Please check whether argument(IfTransfer) exists and it must only be 0 represents non-transferable and 1 represents transferable")
	}
  if application.IfDelay != 0 && application.IfDelay != 1 {
		return shim.Error("Please check whether argument(IfDelay) exists and it must only be 0 represents non-postponable and 1 represents postponable")
	}

	// ==== Create application compositekey ====
	indexName := "apply"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(application.KingmiId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if application already exists ====
	applicationInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get application: " + err.Error())
	} else if applicationInfo != nil {
		return shim.Error("The application already exists")
	}

	// ==== Marshal applicationInfo to JSON ====
  applicationAsBytes, err := json.Marshal(application)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save application to state ====
  err = stub.PutState(indexKey, applicationAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// =======================================================================
//       QueryApplication - Query applications from chaincode state
// =======================================================================
func (m *KingmiChaincode) QueryApplication(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  fmt.Printf("args:\n%s\n\n", args)

  applicationQuery := new(ApplicationQuery)
  applicationQuery.KingmiId = -1
  applicationQuery.OpenEnterpriseId = -1
  applicationQuery.DownEnterpriseId = -1
  applicationQuery.FactoryingEnterpriseId = -1
  applicationQuery.CredentialType = -1
  applicationQuery.CreatedDate = -1
  applicationQuery.DueDate = -1
  applicationQuery.IfGuarantee = -1
  applicationQuery.IfTransfer = -1
  applicationQuery.IfDelay = -1

  err := json.Unmarshal([]byte(args[0]), &applicationQuery)
	if err != nil {
		return shim.Error("Fail to unmarshal application queryString " + err.Error())
	}

  if applicationQuery.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exits and is legal")
  }
  if applicationQuery.OpenEnterpriseId < 0 {
    return shim.Error("Please check whether argument(OpenEnterpriseId) exits and is legal")
  }
  if applicationQuery.DownEnterpriseId < 0 {
    return shim.Error("Please check whether argument(DownEnterpriseId) exits and is legal")
  }
  if applicationQuery.FactoryingEnterpriseId < 0 {
		return shim.Error("Please check whether argument(FactoryingEnterpriseId) exits and is legal")
	}
  if applicationQuery.CredentialType != 1 && applicationQuery.CredentialType != 2 {
		return shim.Error("Please check whether argument(CredentialType) exits and it must only be 1 represents accounts receivable or 2 represents accounts payable")
	}
  if applicationQuery.CreatedDate <= 0 {
		return shim.Error("Please check whether argument(CreatedDate) exists and is legal")
	}
  if applicationQuery.DueDate <= 0 {
		return shim.Error("Please check whether argument(DueDate) exists and is legal")
	}
  if applicationQuery.CreatedDate >= applicationQuery.DueDate {
    return shim.Error("DueDate must be later than CreatedDate")
  }
  if applicationQuery.IfGuarantee != 0 && applicationQuery.IfGuarantee != 1 {
		return shim.Error("Please check whether argument(IfGuarantee) exists and it must only be 0 represents unguaranteed and 1 represents guaranteed")
	}
	if applicationQuery.IfTransfer != 0 && applicationQuery.IfTransfer != 1 {
		return shim.Error("Please check whether argument(IfTransfer) exists and it must only be 0 represents non-transferable and 1 represents transferable")
	}
  if applicationQuery.IfDelay != 0 && applicationQuery.IfDelay != 1 {
		return shim.Error("Please check whether argument(IfDelay) exists and it must be only 0 represents non-postponable and 1 represents postponable")
	}

	queryString := fmt.Sprintf("{\"selector\":" + args[0] + "}")
	queryResults, err := getResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}


// =============================================
//      Settle - Kingmi account settle
// =============================================
func (m *KingmiChaincode) Settle(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
  if len(args) < 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  fmt.Printf("args:\n%s\n\n", args)

  settlement := new(Settlement)
  settlement.KingmiId = -1
  settlement.DueDate = -1
  settlement.SoaDate = -1
  settlement.Amount = -1
  settlement.DueInterest = -1

	err := json.Unmarshal([]byte(args[0]), &settlement)
	if err != nil {
		return shim.Error("Fail to unmarshal settle information string " + err.Error())
	}

  if settlement.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exits and is legal")
  }
  if settlement.DueDate <= 0 {
    return shim.Error("Please check whether argument(DueDate) exits and is legal")
  }
	if settlement.SoaDate <= 0 {
		return shim.Error("Please check whether argument(SoaDate) exits and is legal")
	}
	if settlement.Amount <= 0 {
		return shim.Error("Please check whether argument(Amount) exits and it must be greater than 0")
	}
	if settlement.DueInterest <= 0 {
		return shim.Error("Please check whether argument(DueInterest) exits and it must be greater than 0")
	}

  // ==== Create settlement compositekey ====
	indexName := "settlement"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(settlement.KingmiId)})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(indexKey, value)    // Save index entry to state.

	// ==== Check if settlement already exists ====
	settlementInfo, err := stub.GetState(indexKey)
	if err != nil {
		return shim.Error("Failed to get settlementInfo: " + err.Error())
	} else if settlementInfo != nil {
		return shim.Error("The settlementInfo already exists")
	}

	// ==== Marshal settlementInfo to JSON ====
  settlementAsBytes, err := json.Marshal(settlement)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Save settlement to state ====
  err = stub.PutState(indexKey, settlementAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke successfully"))
}


// ==========================================================================
//       QuerySettlement - Query all settlements from chaincode state
// ==========================================================================
func (m *KingmiChaincode) QuerySettlement(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  fmt.Printf("args:\n%s\n\n", args)

  settlementQuery := new(SettlementQuery)
  settlementQuery.KingmiId = -1
  settlementQuery.DueDate = -1
  settlementQuery.SoaDate = -1

  err := json.Unmarshal([]byte(args[0]), &settlementQuery)
	if err != nil {
		return shim.Error("Fail to unmarshal settlement queryString " + err.Error())
	}

  if settlementQuery.KingmiId < 0 {
    return shim.Error("Please check whether argument(KingmiId) exits and is legal")
  }
  if settlementQuery.DueDate <= 0 {
    return shim.Error("Please check whether argument(DueDate) exits and is legal")
  }
  if settlementQuery.SoaDate <= 0 {
		return shim.Error("Please check whether argument(SoaDate) exits and is legal")
	}

	queryString := fmt.Sprintf("{\"selector\":" + args[0] +"}")
	queryResults, err := getResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}
