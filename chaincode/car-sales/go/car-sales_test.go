package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var testLog = shim.NewLogger("mycc_test")

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "from world state is", string(bytes))
		fmt.Println("State value", name, "expected is", value)
		fmt.Println("State value as expected")
		t.FailNow()
	}
}
func parseStringSliceToByteSlice(args []string) [][]byte {
	var argsByte [][]byte
	for _, str := range args {
		argsByte = append(argsByte, []byte(str))
	}
	return argsByte
}

func checkInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := parseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status != shim.OK {
		testLog.Info("Invoke", functionAndArgs, "failed", string(res.Message))
		t.FailNow()
	} else {
		testLog.Info("Invoke", functionAndArgs, "successful", string(res.Message))
	}
}

func TestCarSales_Init(t *testing.T) {

	scc := new(SmartContract)
	stub := shim.NewMockStub("carsales", scc)
	var functionAndArgs []string
	functionName := "changeOwner"
	functionAndArgs = append(functionAndArgs, functionName, "1FM5K7B83FG612729", "4", "0")

	checkInvoke(t, stub, []string{"initLedger"})

	checkInvoke(t, stub, functionAndArgs)
	car := Car{Chassis: "1FM5K7B83FG612729", Make: "Volkswagen", Model: "Passat", Color: "yellow", Owner: "4", Price: 5200.00, Year: 2008, Damages: []Damage{}}
	seller := Person{ID: "3", FirstName: "Arlee", LastName: "Kayley", Email: "akayley2@businessinsider.com", AccountBalance: 611654.07}
	buyer := Person{ID: "4", FirstName: "Tyson", LastName: "Chidler", Email: "tchidler3@wix.com", AccountBalance: 190881.08}
	carAsBytes, _ := json.Marshal(car)
	sellerAsBytes, _ := json.Marshal(seller)
	buyerAsBytes, _ := json.Marshal(buyer)

	checkState(t, stub, "1FM5K7B83FG612729", string(carAsBytes))
	checkState(t, stub, "3", string(sellerAsBytes))
	checkState(t, stub, "4", string(buyerAsBytes))
}
