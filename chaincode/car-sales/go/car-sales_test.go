package main

import (
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
		fmt.Println("State value", name, "was not", value, "as expected")
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

	functionAndArgs = append(functionAndArgs, functionName, "1FM5K7B83FG612729", "3", "0")

	// Init A=123 B=234
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("123"), []byte("B"), []byte("234")})

	checkState(t, stub, "A", "123")
	checkState(t, stub, "B", "234")
}
