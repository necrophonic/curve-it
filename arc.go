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

	// Ok, that's the set up stuff out of the way. Now to get on with actually
	// calculating some things!

	/* The theory:
		   Given that, in a uniform tube (top and bottom circumference the same), there
		   is no curve. Unrolled, this forms a rectangle.
		   To increase the size of the top circumference, we increase the length of the
		   top edge of the rectangle.
		   As we do so, the two sides bow outwards proportionally (assuming that the lengthing
		   begins at the midpoint and is applied uniformly in both directions). This leaves
		   us with a parallelogram.
		   Attmpting to re-roll this parallelogram (edge flat to edge) does not yield
		   the correct result and gives us a skewed frustrum (i.e. cone with the top cut off).
		   To achieve a smooth roll with flat top and bottom, a curve needs to be applied to
		   the top and bottom edges.
		   So how to determine that curve? If it's too shallow it's no better than the
		   flat edge; if too deep it'll skew the other way.
		   The theory runs that the correct curve is desribed by the circumference of
		   a circle (i.e. the curve is uniform along it's length and, given the right
		   size of circle, will fit perfectly on the circumference). Therefore we need
		   to find the correct size of circle. The theory follows that, given the two
		   curve lengths we want, at the particular distance from each other we require,
		   they will only exist at a single radius (for the inner, and radius+distance for
		   the outer) of the circle.
	     Therefore the method is determine that radius and the angle that defines
	     the sector.
	     How do we do it? Frankly: guess. Then refine.
	     - Make a guess about the radius that describes the lower arc.
	     - Using the calculated angle that gives is the arc at that radius, add the
	       spacer distance and work out what upper arc that would give us.
	     - If it's too big or small (not within a defined margin of acceptance) then
	       refine our guess and try again.
	*/

	radius, theta := calculateRadiusAndTheta(lower, upper, precision, space, guess, 1)

	log.Printf("Upper arc %v within precision: base radius is %v\n", radius, guess)
	log.Printf("Sagitta is [%0.2f]", makeSagitta(theta, guess))
	log.Printf("Theta is [%0.2f]", theta)

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
	// To determine the sagitta we need to know half the length of the chord
	l := radius * math.Sin(theta/2)

	log.Println("Half length is ", l)

	return radius * (1 - math.Cos(theta/2))

	// //
	// log.Println("RADIUS ", radius)
	// log.Println("PART ", 1-radius*math.Cos(theta/2)*-1)
	// //
	// // return radius - (1 - radius*math.Cos(theta/2))
	// return radius - (radius * math.Cos(theta/2))
}

func calculateRadiusAndTheta(l float64, u float64, m float64, space float64, guess float64, it int) (radius float64, theta float64) {
	// Work out the theta on the lower arc given the guess
	theta = makeTheta(l, guess)
	// Now calculate what the upper arc *would* be given the calculated theta
	radius = makeArc(theta, guess+space)

	// Test whether the calculated arc is within the precision we have accepted.
	// If it is then we take it and return out.
	if radius >= u-m && radius <= u+m {
		log.Printf("Calculated in [%v] recursions", it)
		return
	}

	// Recurse!
	if radius > u {
		log.Printf("Too big [%0.2f], guess again with [%v]", radius, guess*2)
		return calculateRadiusAndTheta(l, u, m, space, math.Pow(guess, 2), it+1)
	}
	log.Printf("Too small [%0.2f], guess again with [%v]", radius, guess/2)
	return calculateRadiusAndTheta(l, u, m, space, guess/2, it+1)
}
