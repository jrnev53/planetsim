package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/golang/geo/r3"
)

type Edge struct {
	from r3.Vector
	to   r3.Vector
}

type Face struct {
	p, q, r r3.Vector
}

func main() {
	fmt.Println("hello world")

	var golden_ratio float64 = (1 + math.Sqrt(5)) / 2
	fmt.Println(golden_ratio)

	// define a Regular icosahedron
	// has 20 faces, 12 vertices, 30 edges

	// define vertices about the origin
	var vertices = []r3.Vector{
		{X: 0, Y: 1, Z: golden_ratio},
		{X: 0, Y: -1, Z: golden_ratio},
		{X: 0, Y: 1, Z: -golden_ratio},
		{X: 0, Y: -1, Z: -golden_ratio},

		{X: 1, Y: golden_ratio, Z: 0},
		{X: -1, Y: golden_ratio, Z: 0},
		{X: 1, Y: -golden_ratio, Z: 0},
		{X: -1, Y: -golden_ratio, Z: 0},

		{X: golden_ratio, Y: 0, Z: 1},
		{X: golden_ratio, Y: 0, Z: -1},
		{X: -golden_ratio, Y: 0, Z: 1},
		{X: -golden_ratio, Y: 0, Z: -1},
	}

	// var nameString string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// for i := 0; i < len(vertices); i++ {
	// 	vertices[i].name = nameString[i : i+1]
	// }

	// fmt.Println(vertices)

	// edges have a length of 2, so if two verticies are 2 away from each other, they define an edge.
	// find edges
	// for i := 0; i < 12; i++ {
	// 	fmt.Println(vertices[i])
	// }

	// test comparisons
	var v1 r3.Vector = vertices[3]
	var v1a r3.Vector = vertices[3]
	var v2 r3.Vector = vertices[2]

	fmt.Println("v1 equals v1a?", v1 == v1a)
	fmt.Println("v1 equals v2?", v1 == v2)

	var edges []Edge

	for i, p := range vertices {
		fmt.Println(p)

		for _, q := range vertices[i+1:] {
			fmt.Println(" ", q)
			if distanceMatches(p, q, 2, 0.01) {
				edges = append(edges, Edge{p, q})
			}
		}
	}

	fmt.Println(len(edges), "Edges Found")
	for _, e := range edges {
		fmt.Println(e)
	}

	// find 20 faces
	// we want faces to define triangle cones of the planet and atmosphere
	// each edge is involved in two faces.
	var faces []Face

	for i, e := range edges {
		// finding edges that share a vertex with this edge
		fmt.Println("Finding edges linked to", e)
		// for the edge, find all the other edges that include one of the ends of the edge
		var fromVerts []r3.Vector
		var toVerts []r3.Vector
		for _, f := range edges[i+1:] {
			// do these equilvances work?
			if e.from == f.from {
				fromVerts = append(fromVerts, f.to)
			}
			if e.from == f.to {
				fromVerts = append(fromVerts, f.from)
			}

			if e.to == f.from {
				toVerts = append(toVerts, f.to)
			}
			if e.to == f.to {
				toVerts = append(toVerts, f.from)
			}
		}
		fmt.Println("fromVerts size", len(fromVerts), "toVerts size", len(toVerts))

		// find the two vertices that shows up in both lists
		var commonVerts []r3.Vector = setIntersection(fromVerts, toVerts)
		// could be 0, 1, or 2 faces
		for _, p := range commonVerts {
			// define face
			faces = append(faces, Face{e.from, e.to, p})
		}
	}

	// fmt.Println(faces)
	fmt.Println("Found", len(faces), "faces")
	for _, f := range faces {
		fmt.Println(f)

		// verify the length between the vertices in the face
		distPQ := f.p.Distance(f.q)
		distQR := f.q.Distance(f.r)
		distRP := f.r.Distance(f.p)
		fmt.Println(distPQ, distQR, distRP)
		if distPQ+distQR+distRP != 6 {
			panic("something isn't adding up")
		}
	}

	// // fmt.Println(faces)
	// fmt.Println("Found", len(faces), "faces")
	// for _, f := range faces {
	// 	fmt.Println(f.p.name, f.q.name, f.r.name)
	// }

	// find the angles of each face
	fmt.Println("Found", len(faces), "faces")
	for _, f := range faces {
		var angleP string = getAngle(f.p)
		var angleQ string = getAngle(f.q)
		var angleR string = getAngle(f.r)
		fmt.Println(angleP, angleQ, angleR)
	}

}

func getAngle(position r3.Vector) string {
	var angleXY float64 = math.Atan2(position.Y, position.X) / (math.Pi / 180)
	var angleYZ float64 = math.Atan2(position.Z, position.Y) / (math.Pi / 180)
	var angleZX float64 = math.Atan2(position.X, position.Z) / (math.Pi / 180)
	builder := strings.Builder{}
	builder.WriteString("{")
	builder.WriteString(fmt.Sprintf("%0.2f", angleXY))
	builder.WriteString(",")
	builder.WriteString(fmt.Sprintf("%0.2f", angleYZ))
	builder.WriteString(",")
	builder.WriteString(fmt.Sprintf("%0.2f", angleZX))
	builder.WriteString("}")
	return builder.String()
}

// find points that are in both sets
func setIntersection(fromVerts []r3.Vector, toVerts []r3.Vector) []r3.Vector {
	var ret []r3.Vector

	fmt.Println("setIntersection with", fromVerts, "and", toVerts)
	for _, fv := range fromVerts {
		for _, tv := range toVerts {
			// if fv == tv
			if positionEquals(fv, tv) {
				ret = append(ret, fv)
			}
		}
	}

	fmt.Println("return set:", ret)
	fmt.Println("Set has", len(ret), "members")

	return ret
}

func distanceMatches(p, q r3.Vector, distance float64, closeness float64) bool {
	var calcDistance float64

	var xpart, ypart, zpart float64
	xpart = p.X - q.X
	ypart = p.Y - q.Y
	zpart = p.Z - q.Z
	xpart = math.Pow(xpart, 2)
	ypart = math.Pow(ypart, 2)
	zpart = math.Pow(zpart, 2)

	calcDistance = math.Sqrt(xpart + ypart + zpart)

	var difference = math.Abs(distance - calcDistance)
	// fmt.Println("difference", difference)

	return difference < closeness
}

func positionEquals(p, q r3.Vector) bool {
	return p.X == q.X && p.Y == q.Y && p.Z == q.Z
}
