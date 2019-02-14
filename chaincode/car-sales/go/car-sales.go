package main

// https://github.com/hyperledger/fabric-sdk-rest/blob/master/tests/input/src/marbles02/marbles_chaincode.go
// https://medium.com/wearetheledger/hyperledger-fabric-couchdb-fantastic-queries-and-where-to-find-them-f8a3aecef767
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Damage struct {
	Description string  `json:"description"`
	RepairPrice float64 `json:"repairPrice"`
}

type Car struct {
	Chassis string   `json:"chassis"`
	Make    string   `json:"make"`
	Model   string   `json:"model"`
	Year    int      `json:"year"`
	Color   string   `json:"color"`
	Owner   int      `json:"ownerID"`
	Price   float64  `json:"price"`
	Damages []Damage `json:"damages"`
}

type Person struct {
	ID             int     `json:"id"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	AccountBalance float64 `json:"accountBalance"`
}

//TODO when repairing dmg, decrese money of owner

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	if function == "initLedger" { //create a new marble
		return s.initLedger(APIstub)
	} else if function == "changeColor" {
		return s.changeColor(APIstub, args)
	} else if function == "getByColor" {
		return s.getByColor(APIstub, args)
	} else if function == "getByColorOwner" {
		return s.getByColor(APIstub, args)
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

func (s *SmartContract) getByColor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	keysIter, err := stub.GetStateByRange("", "")
	cars := make([]Car, 5)
	considerOwner := len(args) == 2
	ownerID := 0
	if considerOwner {
		ownerID, err = strconv.Atoi(args[1])
		return shim.Error("Owner must be int")
	}
	if len(args) > 2 {
		return shim.Error("Incorrect number of arguments. Only chassis number needed.")
	}

	for keysIter.HasNext() {
		carAsBytes, err := keysIter.Next()
		var object interface{}
		car := Car{}
		if err != nil {
			return shim.Error("Error while iterating through world state " + err.Error())
		}

		err = json.Unmarshal(carAsBytes.GetValue(), &object)
		switch object.(type) {
		case Car:
			car = Car(object.(Car))
		default:
			continue
		}

		if car.Color == args[0] || (considerOwner && car.Owner == ownerID) {
			cars = append(cars, car)
		}
	}
	carSliceBytes, err := json.Marshal(cars)
	if err != nil {
		return shim.Error("Error while marshaling car slice " + err.Error())
	}
	return shim.Success(carSliceBytes)
}

func (s *SmartContract) getCarHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Only chassis number needed.")
	}

	it, err := stub.GetHistoryForKey(args[0])
	history := make([]Car, 5)
	if err != nil {
		return shim.Error("Error while getting history" + err.Error())
	}

	for it.HasNext() {
		change, err := it.Next()
		if err != nil {
			return shim.Error("There is an error while iterating through history for key" + args[0])
		}
		carAsBytes := change.GetValue()
		car := Car{}
		json.Unmarshal(carAsBytes, &car)
		history = append(history, car)
	}
	historyAsBytes, _ := json.Marshal(history)
	defer it.Close()
	return shim.Success(historyAsBytes)
}

func (s *SmartContract) changeColor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.GetState(args[0])
	newColor := args[1]
	car := Car{}

	json.Unmarshal(carAsBytes, &car)

	car.Color = newColor

	carAsBytes, _ = json.Marshal(car)
	stub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) noteDamage(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.GetState(args[0])

	car := Car{}

	json.Unmarshal(carAsBytes, &car)

	newDamage := Damage{}
	newDamage.Description = args[1]
	newDamage.RepairPrice, _ = strconv.ParseFloat(args[2], 64)
	dmgCost := 0.0
	for _, dmg := range car.Damages {
		dmgCost += dmg.RepairPrice
	}

	if dmgCost > car.Price {
		err := stub.DelState(args[0])
		if err != nil {
			return shim.Error("Failed to delete car with chassis number " + car.Chassis + " " + err.Error())
		}
	}

	car.Damages = append(car.Damages, newDamage)
	carAsBytes, _ = json.Marshal(car)

	stub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) repairDamage(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.GetState(args[0])
	car := Car{}
	err := json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	car.Damages = make([]Damage, 0)

	carAsBytes, _ = json.Marshal(car)

	stub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeOwner(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}
	carAsBytes, _ := stub.GetState(args[0])
	car := Car{}
	err := json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	payTheDmg := args[2] == "1"

	sellerAsBytes, _ := stub.GetState(string(car.Owner))
	seller := Person{}
	err = json.Unmarshal(sellerAsBytes, &seller)
	if err != nil {
		return shim.Error("Error while unmarshaling seller " + err.Error())
	}
	if len(car.Damages) > 0 && !payTheDmg {
		return shim.Success([]byte("Car has damage and buyer won't take damaged cars"))
	}

	buyerAsBytes, _ := stub.GetState(args[1])
	buyer := Person{}
	err = json.Unmarshal(buyerAsBytes, &buyer)
	if err != nil {
		return shim.Error("Error while unmarshaling buyer " + err.Error())
	}

	if buyer.AccountBalance < car.Price {
		return shim.Error("Buyer doesnt have enough money")
	}

	if payTheDmg {
		sum := 0.0
		for _, dmg := range car.Damages {
			sum += dmg.RepairPrice
		}
		if sum+car.Price > buyer.AccountBalance {
			return shim.Error("Buyer doesn't have enough money to pay the damage.")
		}
	}

	seller.AccountBalance += car.Price
	buyer.AccountBalance -= car.Price
	car.Owner = buyer.ID
	return shim.Success(nil)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	people := []Person{
		Person{ID: 1, FirstName: "Ermentrude", LastName: "Keam", Email: "ekeam0@netscape.com", AccountBalance: 808275.03},
		Person{ID: 2, FirstName: "Sayers", LastName: "Webb-Bowen", Email: "swebbbowen1@paginegialle.it", AccountBalance: 429843.38},
		Person{ID: 3, FirstName: "Arlee", LastName: "Kayley", Email: "akayley2@businessinsider.com", AccountBalance: 606454.07},
		Person{ID: 4, FirstName: "Tyson", LastName: "Chidler", Email: "tchidler3@wix.com", AccountBalance: 196081.08},
		Person{ID: 5, FirstName: "Aggie", LastName: "Garcia", Email: "agarcia4@vinaora.com", AccountBalance: 628858.16},
	}

	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Color: "blue", Owner: 1},
		Car{Make: "Ford", Model: "Mustang", Color: "red", Owner: 2},
		Car{Make: "Hyundai", Model: "Tucson", Color: "green", Owner: 3},
		Car{Make: "Volkswagen", Model: "Passat", Color: "yellow", Owner: 3},
		Car{Make: "Tesla", Model: "S", Color: "black", Owner: 4},
		Car{Make: "Peugeot", Model: "205", Color: "purple", Owner: 5},
		Car{Make: "Chery", Model: "S22L", Color: "white", Owner: 5},
		Car{Make: "Fiat", Model: "Punto", Color: "violet", Owner: 1},
		Car{Make: "Tata", Model: "Nano", Color: "indigo", Owner: 3},
		Car{Make: "Holden", Model: "Barina", Color: "brown", Owner: 4},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		APIstub.PutState(car.Chassis, carAsBytes)
		fmt.Println("Added", cars[i])
	}
	for _, persone := range people {
		personAsBytes, _ := json.Marshal(persone)
		APIstub.PutState(string(persone.ID), personAsBytes)
		fmt.Println("Added", persone)
	}

	return shim.Success(nil)
}
