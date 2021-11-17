package main

import (
	"fmt"
	"math"
)

type Position struct {
	name string
	x    float64
	y    float64
	z    float64
}

type Edge struct {
	from Position
	to   Position
}

type Face struct {
	p, q, r Position
}

func main() {
	fmt.Println("hello world")

	var golden_ratio float64 = (1 + math.Sqrt(5)) / 2
	fmt.Println(golden_ratio)

	// define a Regular icosahedron
	// has 20 faces, 12 vertices, 30 edges

	// define vertices about the origin
	var vertices = []Position{
		{x: 0, y: 1, z: golden_ratio},
		{x: 0, y: -1, z: golden_ratio},
		{x: 0, y: 1, z: -golden_ratio},
		{x: 0, y: -1, z: -golden_ratio},

		{x: 1, y: golden_ratio, z: 0},
		{x: -1, y: golden_ratio, z: 0},
		{x: 1, y: -golden_ratio, z: 0},
		{x: -1, y: -golden_ratio, z: 0},

		{x: golden_ratio, y: 0, z: 1},
		{x: golden_ratio, y: 0, z: -1},
		{x: -golden_ratio, y: 0, z: 1},
		{x: -golden_ratio, y: 0, z: -1},
	}

	var nameString string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for i := 0; i < len(vertices); i++ {
		vertices[i].name = nameString[i : i+1]
	}

	// fmt.Println(vertices)

	// edges have a length of 2, so if two verticies are 2 away from each other, they define an edge.
	// find edges
	// for i := 0; i < 12; i++ {
	// 	fmt.Println(vertices[i])
	// }

	// test comparisons
	var v1 Position = vertices[3]
	var v1a Position = vertices[3]
	var v2 Position = vertices[2]

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
		var fromVerts []Position
		var toVerts []Position
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
		var commonVerts []Position = setIntersection(fromVerts, toVerts)
		// could be 0, 1, or 2 faces
		for _, p := range commonVerts {
			// define face
			faces = append(faces, Face{e.from, e.to, p})
		}
	}

	// fmt.Println(faces)
	fmt.Println("Found", len(faces), "faces")
	for _, f := range faces {
		fmt.Println(f.p.name, f.q.name, f.r.name)
	}
}

// find points that are in both sets
func setIntersection(fromVerts []Position, toVerts []Position) []Position {
	var ret []Position

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

func distanceMatches(p, q Position, distance float64, closeness float64) bool {
	var calcDistance float64

	var xpart, ypart, zpart float64
	xpart = p.x - q.x
	ypart = p.y - q.y
	zpart = p.z - q.z
	xpart = math.Pow(xpart, 2)
	ypart = math.Pow(ypart, 2)
	zpart = math.Pow(zpart, 2)

	calcDistance = math.Sqrt(xpart + ypart + zpart)

	var difference = math.Abs(distance - calcDistance)
	// fmt.Println("difference", difference)

	return difference < closeness
}

func positionEquals(p, q Position) bool {
	return p.x == q.x && p.y == q.y && p.z == q.z
}
