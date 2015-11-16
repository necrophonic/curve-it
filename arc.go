package main

import (
	"flag"
	"log"
	"math"
	"os"
)

func main() {
	var lower, upper, space, guess, precision float64

	flag.Float64Var(&lower, "lower", -1, "lower arc")
	flag.Float64Var(&upper, "upper", -1, "upper arc")
	flag.Float64Var(&space, "space", -1, "space between the arcs (height)")
	flag.Float64Var(&guess, "guess", 10, "initial guess at core radius")
	flag.Float64Var(&precision, "precision", 0.01, "precision within which to target search")
	flag.Parse()

	if upper < lower {
		log.Fatalf("Upper arc [%v] must be longer than lower [%v]", upper, lower)
	}
	if upper <= 0 || lower <= 0 {
		log.Fatal("Upper and lower arcs must be positive, non-zero values")
	}
	if guess <= 2 {
		log.Fatalf("Guess [%v] must be larger than 2", guess)
	}

	log.Println("Upper: ", upper, "  Lower: ", lower, "  Space: ", space)
	log.Println("Initial guess is", guess)
	log.Println("Precision is", precision)

	calculate(lower, upper, precision, space, guess, 1)
	os.Exit(0)
}

// Theta is the angle at the circle origin that defines the sector
func makeTheta(arc float64, radius float64) float64 {
	return arc / (radius * (math.Pi / 180))
}

// The arc is the length of an arc described by the radius over a given theta
func makeArc(theta float64, radius float64) float64 {
	return theta * radius * (math.Pi / 180)
}

// The sagitta is the tangential distance between the center of the chord
// of the sector and the outline of the circle
func makeSagitta(theta float64, radius float64) float64 {
	// return radius * (1 - math.Cos(theta/2))

	//
	log.Println("RADIUS ", radius)
	log.Println("PART ", 1-radius*math.Cos(theta/2)*-1)
	//
	// return radius - (1 - radius*math.Cos(theta/2))
	return radius - (radius * math.Cos(theta/2))
}

func calculate(l float64, u float64, m float64, space float64, guess float64, it int) {
	// Work out the theta on the lower arc given the guess
	theta := makeTheta(l, guess)
	// Now calculate what the upper arc *would* be given the calculated theta
	calculatedUpperArc := makeArc(theta, guess+space)

	// Test whether the calculated arc is within the precision we have accepted.
	// If it is then we take it and return out.
	if calculatedUpperArc >= u-m && calculatedUpperArc <= u+m {
		log.Printf("Upper arc %v within precision: base radius is %v\n", calculatedUpperArc, guess)
		log.Printf("Sagitta is [%0.2f]", makeSagitta(theta, guess))
		log.Printf("Theta is [%0.2f]", theta)
		log.Printf("Calculated in [%v] recursions", it)
		return
	}

	// Recurse!
	if calculatedUpperArc > u {
		log.Printf("Too big [%0.2f], guess again with [%v]", calculatedUpperArc, guess*2)
		calculate(l, u, m, space, math.Pow(guess, 2), it+1)
	} else {
		log.Printf("Too small [%0.2f], guess again with [%v]", calculatedUpperArc, guess/2)
		calculate(l, u, m, space, guess/2, it+1)
	}

}
