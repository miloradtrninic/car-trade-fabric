package main


// https://github.com/hyperledger/fabric-sdk-rest/blob/master/tests/input/src/marbles02/marbles_chaincode.go
// https://medium.com/wearetheledger/hyperledger-fabric-couchdb-fantastic-queries-and-where-to-find-them-f8a3aecef767
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Damage struct {
	Description string `json:"description"`
	RepairPrice float64 `json:"repairPrice"`
}

type Car struct {
	Chassis string `json:"chassis"`
	Make string `json:"make"`
	Model string `json:"model"`
	Year int `json:"year"`
	Color string `json:"color"`
	Owner int `json:"ownerID"`
	Price float64 `json:"price"`
	Damages []Damage `json:"damages"`
}

type Person struct {
	Id int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	AccountBalance float64 `json:"accountBalance"`
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "changeColor" {
		return s.changeColor(APIstub, args)
	} else if function == "getByColor" {
		return s.getByColor(APIstub, args)
	} else if function == "getByColorOwner" {
		return s.getByColorOwner(APIstub, args)
	} else if function == "getCarHistory" {
		return s.getCarHistory(APIstub, args)
	} else if function == "changeOwner" {
		return s.changeOwner(APIstub, args)
	} else if function == "noteDamage" {
		return s.noteDamage(APIstub, args)
	} else if function == "repairDamage" {
		return s.repairDamage(APIstub, args)
	} else {
		fmt.Println("I nvoke did not find func: " + function) //error
		return shim.Error("Received unknown function invocation")
	}
}

func (s *SmartContract) changeColor(stub shim.ChaincodeStubInterface, args string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}
	carAsBytes, _ := stub.getState(args[0])
	newColor := args[1]
	car := Car{}

	json.Unmrashal(carAsBytes, &car)

	car.Color = newColor

	carAsBytes, _ = json.Marshal(car)
	stub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) noteDamage(stub shim.ChaincodeStubInterface, args string) sc.Resoponse {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.getState(args[0])

	car := Car{}

	json.Unmarshal(carAsBytes, &car)

	newDamage := Damage{}
	newDamage.Description =  args[1]
	newDamage.RepairPrice = args[2]
	dmgCost := 0
	for dmg := range car.Damages {
		dmgCost += dmg
	}

	if (dmgCost > car.Price) {
		err := stub.DelState(args[0])
		if err != nil {
			return shim.Error("Failed to delete car with chassis number " + car.Chassis + " "  + err.Error())
		}
	}

	car.Damages.append(car)
	carAsBytes, _ = json.Marshal(car)

	stub.PutState(args[0], carAsBytes)
}

func (s *SmartContract) repairDamage(stub shim.ChaincodeStubInterface, args string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.getState(args[0])
	car := Car{}
	json.Unmarshal(carAsBytes, &car)

	car.Damages = make([]Damage, 0)

	carAsBytes, _ = json.Marshal(car)

	stub.PutState(args[0], carAsBytes)
}

func (s *SmartContract) getByColor(args string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorect number of arguments.")
	}
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(cars) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(cars[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", cars[i])
		i = i + 1
	}

	return shim.Success(nil)
}