package main

import ()

func main() {

}

func switchStatement() {
	switch coinflip() {
	case "heads":
		heads++
	case "tails":
		tails++
	default:
		fmt.Println("landed on edge!")
	}

	switch {
	case x > 0:
		return +1
	default:
		return 0
	case x < 0:
		return -1
	}
}

// Named Types
// type delcaration makes it possible ot give a name to an existing type
// since struct types are long they are nearly always named

//Pointers
//Go provides pointers or values that contain the address of a variable
// Pointers are explictly visible
// The & operator yields the address and the * operator retrieves the variable that the pointer refers to
// There is no pointer arithmetic

//Methods and interfaces
// method is a function associated with a type
// Interfaces are abstract types that let us treat different concrete types in the same way
