package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Car struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Color string `json:"color"`
	Year  string `json:"year"`
	Owner string `json:"owner"`
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("Error when starting SmartContract", err)
	}
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	cars := []Car{
		Car{ID: "1", Make: "Tesla", Color: "Red", Year: "2015", Owner: "Yigit"},
		Car{ID: "2", Make: "Toyota", Color: "Blue", Year: "2010", Owner: "Kemal"},
		Car{ID: "3", Make: "Honda", Color: "Black", Year: "2000", Owner: "Merve"},
		Car{ID: "4", Make: "Open", Color: "White", Year: "2008", Owner: "Damla"},
		Car{ID: "5", Make: "BMW", Color: "Grey", Year: "1995", Owner: "Begum"},
	}
	i := 0
	for i < len(cars) {
		carAsByte, _ := json.Marshal(cars[i])
		stub.PutState("CAR"+cars[i].ID, carAsByte)
		fmt.Println("CAR"+cars[i].ID, "added, values: ", cars[i])
		i = i + 1
	}
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "createCar" {
		return s.createCar(stub, args)
	} else if function == "queryAllCars" {
		return s.queryAllCars(stub)
	} else if function == "getCarHistory" {
		return s.getCarHistory(stub, args)
	} else if function == "deleteCar" {
		return s.deleteCar(stub, args)
	} else if function == "changeCarOwner" {
		return s.changeCarOwner(stub, args)
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 5 {
		return shim.Error("invalid number of arguments")
	}

	car := Car{ID: args[0], Make: args[1], Color: args[2], Year: args[3], Owner: args[4]}
	carAsBytes, _ := json.Marshal(car)
	stub.PutState("CAR"+args[0], carAsBytes)
	fmt.Println("CAR"+args[0], "added, values: ", car)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCars(stub shim.ChaincodeStubInterface) peer.Response {

	startKey := "CAR0"
	endKey := "CAR99999"

	iterator, err := stub.GetStateByRange(startKey, endKey)

	if err != nil {
		shim.Error(err.Error())
	}

	iterator.Close()
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println("Key: ", queryResponse.Key, " value: ", string(queryResponse.Value))
	}
	return shim.Success(nil)
}

func (s *SmartContract) getCarHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}

	iterator, _ := stub.GetHistoryForKey("CAR" + args[0])
	iterator.Close()

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			shim.Error(err.Error())
		}
		fmt.Println(queryResponse)
	}

	return shim.Success(nil)
}

func (s *SmartContract) deleteCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("invalid number of arguments")
	}

	err := stub.DelState("CAR" + args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SmartContract) changeCarOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("invalid number of arguments")
	}
	carAsBytes, _ := stub.GetState("CAR" + args[0])
	car := Car{}
	json.Unmarshal(carAsBytes, &car)
	car.Owner = args[1]
	carAsBytes, _ = json.Marshal(car)
	stub.PutState("CAR"+args[0], carAsBytes)

	return shim.Success(nil)
}
