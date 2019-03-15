package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Firstname, Lastname string
}

func main() {
	fmt.Printf("Start:\n")
	var people []Person
	john := Person{
		Firstname: "John",
		Lastname:  "Doe",
	}
	jane := Person{
		Firstname: "Jane",
		Lastname:  "Deaux",
	}
	people = append(people, john)
	people = append(people, jane)

	peopleJson, _ := json.Marshal(people)
	fmt.Println(string(peopleJson))

	people = people[:0]

	json.Unmarshal(peopleJson, &people)
	for _, p := range people {
		fmt.Printf("%s: %s\n", p.Firstname, p.Lastname)
	}
}
