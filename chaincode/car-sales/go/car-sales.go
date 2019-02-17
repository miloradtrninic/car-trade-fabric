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
	Owner   string   `json:"ownerID"`
	Price   float64  `json:"price"`
	Damages []Damage `json:"damages"`
}

type Person struct {
	ID             string  `json:"id"`
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
	logger := shim.NewLogger("infoLogger")

	if function == "initLedger" { //create a new marble
		return s.initLedger(APIstub)
	} else if function == "changeColor" {
		return s.changeColor(APIstub, args)
	} else if function == "getByKey" {
		return s.getByKey(APIstub, args)
	} else if function == "getPersonByKey" {
		return s.getPersonByKey(APIstub, logger, args)
	} else if function == "getByColor" {
		return s.getByColor(APIstub, logger, args)
	} else if function == "getByColorOwner" {
		return s.getByColor(APIstub, logger, args)
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

func (s *SmartContract) getByColor(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) sc.Response {
	keysIter, err := stub.GetStateByRange("", "")
	cars := make([]Car, 0)
	considerOwner := len(args) == 2
	var ownerID string
	logger.Info("Started get by color %s", args[0])
	if considerOwner {
		ownerID = args[1]
	}
	if len(args) > 2 {
		return shim.Error("Incorrect number of arguments. Color (and owner).")
	}

	for keysIter.HasNext() {
		logger.Info("Next object check")
		carAsBytes, err := keysIter.Next()
		car := Car{}
		if err != nil {
			return shim.Error("Error while iterating through world state " + err.Error())
		}

		err = json.Unmarshal(carAsBytes.GetValue(), &car)
		if err != nil {
			switch err.(type) {
			case *json.UnsupportedTypeError:
				logger.Info("Object is not a car")
				continue
			default:
				return shim.Error("Error while unmarshaling object" + err.Error())
			}
		}

		if car.Color == args[0] {
			if considerOwner && car.Owner == ownerID {
				cars = append(cars, car)
			} else if !considerOwner {
				cars = append(cars, car)
			}
		}
	}
	carSliceBytes, err := json.Marshal(cars)
	if err != nil {
		return shim.Error("Error while marshaling car slice " + err.Error())
	}
	return shim.Success(carSliceBytes)
}

func (s *SmartContract) getByKey(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Only chassis number needed.")
	}

	carAsBytes, err := stub.GetState(args[0])
	if carAsBytes == nil {
		return shim.Error("No car with chassis = " + args[0])
	}
	if err != nil {
		return shim.Error("Error while getting state" + err.Error())
	}
	return shim.Success(carAsBytes)
}

func (s *SmartContract) getPersonByKey(stub shim.ChaincodeStubInterface, logger *shim.ChaincodeLogger, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Only person id needed.")
	}
	logger.Info("Person ID param = " + args[0])
	personAsBytes, err := stub.GetState(args[0])
	if personAsBytes == nil {
		return shim.Error("No person with ID = " + args[0])
	}
	if err != nil {
		return shim.Error("Error while getting state" + err.Error())
	}
	return shim.Success(personAsBytes)
}

func (s *SmartContract) getCarHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Only chassis number needed.")
	}

	it, err := stub.GetHistoryForKey(args[0])
	history := make([]Car, 0)
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
		err = json.Unmarshal(carAsBytes, &car)
		if err != nil {
			return shim.Error("Error while unmarshaling car " + err.Error())
		}
		history = append(history, car)
	}
	historyAsBytes, err := json.Marshal(history)
	if err != nil {
		return shim.Error("Error while marshaling history " + err.Error())
	}
	defer it.Close()
	return shim.Success(historyAsBytes)
}

func (s *SmartContract) changeColor(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error while getting state" + err.Error())
	}
	if carAsBytes == nil {
		return shim.Error("No car with chassis = " + args[0])
	}
	newColor := args[1]
	car := Car{}

	err = json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	car.Color = newColor

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return shim.Error("Error while marshaling car " + err.Error())
	}
	stub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) noteDamage(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error while getting state" + err.Error())
	}
	if carAsBytes == nil {
		return shim.Error("No car with chassis = " + args[0])
	}
	car := Car{}

	err = json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	newDamage := Damage{}
	newDamage.Description = args[1]
	newDamage.RepairPrice, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Error while parsing float repair price " + err.Error())
	}
	dmgCost := 0.0
	for _, dmg := range car.Damages {
		dmgCost += dmg.RepairPrice
	}

	if dmgCost > car.Price {
		err = stub.DelState(args[0])
		if err != nil {
			return shim.Error("Error while deleting state" + err.Error())
		}
		if err != nil {
			return shim.Error("Failed to delete car with chassis number " + car.Chassis + " " + err.Error())
		}
	}

	car.Damages = append(car.Damages, newDamage)
	carAsBytes, err = json.Marshal(car)
	if err != nil {
		return shim.Error("Error while marshaling car " + err.Error())
	}
	err = stub.PutState(args[0], carAsBytes)
	if err != nil {
		shim.Error("Error while saving car with number " + car.Chassis + " " + err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) repairDamage(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	carAsBytes, _ := stub.GetState(args[0])
	if carAsBytes == nil {
		return shim.Error("Car with chassis number " + args[0] + " does not exist.")
	}
	car := Car{}
	err := json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	car.Damages = make([]Damage, 0)

	carAsBytes, err = json.Marshal(car)
	if err != nil {
		shim.Error("Error while marshaling car with number " + car.Chassis + " " + err.Error())
	}
	err = stub.PutState(args[0], carAsBytes)
	if err != nil {
		shim.Error("Error while saving car with number " + car.Chassis + " " + err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) changeOwner(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}
	carAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Error while getting car info from world state " + err.Error())
	}
	car := Car{}
	err = json.Unmarshal(carAsBytes, &car)
	if err != nil {
		return shim.Error("Error while unmarshaling car " + err.Error())
	}
	buyerID := args[1]

	if buyerID == car.Owner {
		return shim.Error("Owner same as buyer")
	}
	payTheDmg := args[2] == "1"

	sellerAsBytes, err := stub.GetState(car.Owner)
	if err != nil {
		return shim.Error("Error while getting seller info from world state " + err.Error())
	}
	seller := Person{}
	err = json.Unmarshal(sellerAsBytes, &seller)
	if err != nil {
		return shim.Error("Error while unmarshaling seller " + err.Error())
	}
	if len(car.Damages) > 0 && !payTheDmg {
		return shim.Error("Car has damage and buyer won't take damaged cars")
	}

	buyerAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error("Error while getting buyer info from world state " + err.Error())
	}
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
	carAsBytes, err = json.Marshal(car)
	if err != nil {
		shim.Error("Error while marshaling car with number " + car.Chassis + " " + err.Error())
	}
	err = stub.PutState(car.Chassis, carAsBytes)

	buyerAsBytes, err = json.Marshal(buyer)
	if err != nil {
		shim.Error("Error while marshaling buyer with id " + buyer.ID + " " + err.Error())
	}
	err = stub.PutState(buyer.ID, buyerAsBytes)

	sellerAsBytes, err = json.Marshal(seller)
	if err != nil {
		shim.Error("Error while marshaling seller with id " + seller.ID + " " + err.Error())
	}
	err = stub.PutState(seller.ID, sellerAsBytes)
	return shim.Success(nil)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	people := []Person{
		Person{ID: "1", FirstName: "Ermentrude", LastName: "Keam", Email: "ekeam0@netscape.com", AccountBalance: 808275.03},
		Person{ID: "2", FirstName: "Sayers", LastName: "Webb-Bowen", Email: "swebbbowen1@paginegialle.it", AccountBalance: 429843.38},
		Person{ID: "3", FirstName: "Arlee", LastName: "Kayley", Email: "akayley2@businessinsider.com", AccountBalance: 606454.07},
		Person{ID: "4", FirstName: "Tyson", LastName: "Chidler", Email: "tchidler3@wix.com", AccountBalance: 196081.08},
		Person{ID: "5", FirstName: "Aggie", LastName: "Garcia", Email: "agarcia4@vinaora.com", AccountBalance: 628858.16},
	}

	cars := []Car{
		Car{Chassis: "JH4CL95877C271071", Make: "Toyota", Model: "Prius", Color: "blue", Owner: "1", Price: 3000.00, Year: 2010, Damages: []Damage{}},
		Car{Chassis: "5N1AN0NW4FN197748", Make: "Ford", Model: "Mustang", Color: "red", Owner: "2", Price: 200000.00, Year: 2010, Damages: []Damage{}},
		Car{Chassis: "2C3CDZBT8FH306479", Make: "Hyundai", Model: "Tucson", Color: "green", Owner: "3", Price: 10000.00, Year: 20160, Damages: []Damage{}},
		Car{Chassis: "1FM5K7B83FG612729", Make: "Volkswagen", Model: "Passat", Color: "yellow", Owner: "3", Price: 5200.00, Year: 2008, Damages: []Damage{}},
		Car{Chassis: "1G6AV5S84E0464733", Make: "Tesla", Model: "S", Color: "black", Owner: "4", Price: 22000.00, Year: 2016, Damages: []Damage{}},
		Car{Chassis: "4A31K2DF0BE184156", Make: "Peugeot", Model: "205", Color: "purple", Owner: "5", Price: 800.00, Year: 2000, Damages: []Damage{}},
		Car{Chassis: "WBALZ3C56CD866601", Make: "Chery", Model: "S22L", Color: "white", Owner: "5", Price: 8000.00, Year: 2015, Damages: []Damage{}},
		Car{Chassis: "3VWKX7AJ8AM939252", Make: "Fiat", Model: "Punto", Color: "violet", Owner: "1", Price: 1000.00, Year: 2003, Damages: []Damage{}},
		Car{Chassis: "SCFEFBAC9AG437143", Make: "Tata", Model: "Nano", Color: "indigo", Owner: "3", Price: 500.00, Year: 2005, Damages: []Damage{}},
		Car{Chassis: "2T1KU4EEXDC074983", Make: "Holden", Model: "Barina", Color: "brown", Owner: "4", Price: 3800.00, Year: 2010, Damages: []Damage{}},
	}

	for i, car := range cars {
		carAsBytes, err := json.Marshal(car)
		if err != nil {
			shim.Error("Error while marshaling car with number " + car.Chassis + " " + err.Error())
		}
		err = APIstub.PutState(car.Chassis, carAsBytes)
		if err != nil {
			shim.Error("Error while saving car with number " + car.Chassis + " " + err.Error())
		}
		fmt.Println("Added", cars[i])
	}
	for _, person := range people {
		personAsBytes, err := json.Marshal(person)
		if err != nil {
			shim.Error("Error while marshaling person with number " + person.ID + " " + err.Error())
		}
		err = APIstub.PutState(person.ID, personAsBytes)
		if err != nil {
			shim.Error("Error while saving person with number " + person.ID + " " + err.Error())
		}
		fmt.Println("Added", person)
	}

	return shim.Success(nil)
}
