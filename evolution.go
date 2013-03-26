package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Population struct {
	organisms Organisms
	birthRate float32 // Births per organism (percentage)
	deathRate float32 // Deaths per organism (percentage)
}

type Organisms []Organism

// sort.Interface implementation
func (o Organisms) Len() int {
	return len(o)
}
func (o Organisms) Less(i, j int) bool {
	return o[i].value() < o[j].value()
}
func (o Organisms) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type ReverseSort struct {
	sort.Interface
}

func (r ReverseSort) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func (p Population) value() int {
	sum := 0
	for _, org := range p.organisms {
		sum += org.value()
	}
	return sum
}

func (p Population) best() Organism {
	if len(p.organisms) == 0 {
		return nil
	}

	bestOrg := p.organisms[0]
	for _, org := range p.organisms[1:] {
		if org.value() > bestOrg.value() {
			bestOrg = org
		}
	}

	return bestOrg
}

func (p Population) numParents() int {
	numBirths := int(float32(len(p.organisms)) * p.birthRate)
	numParents := numBirths * 2

	if numParents > len(p.organisms) {
		numParents = len(p.organisms)

		// Ensure numParents is even.
		numParents /= 2
		numParents *= 2
	}

	return numParents
}

func (p Population) numToKill() int {
	numToKill := int(float32(len(p.organisms)) * p.deathRate)
	if numToKill > len(p.organisms) {
		numToKill = len(p.organisms)
	}
	return numToKill
}

// Select parents to breed from the population, given its birth rate.
func (p Population) selectParents() Organisms {
	numParents := p.numParents()
	sort.Sort(ReverseSort{p.organisms})
	mostFit := p.organisms[:numParents]

	// Shuffle
	for i := range mostFit {
		j := rand.Intn(i + 1)
		mostFit[i], mostFit[j] = mostFit[j], mostFit[i]
	}

	return p.organisms[:numParents]
}

func (p *Population) killWeakestOrganisms() int {
	numToKill := p.numToKill()
	sort.Sort(p.organisms)
	p.organisms = p.organisms[numToKill:]
	return numToKill
}

type Organism interface {
	value() int
	crossover(o2 Organism) Organism
	mutate() Organism
}

type IntOrganism []int

func NewIntOrganism(size int) IntOrganism {
	organism := make(IntOrganism, size)
	for i := 0; i < size; i++ {
		organism[i] = rand.Intn(2) // 0 or 1
	}
	return organism
}

// Organism implementation
func (o IntOrganism) value() int {
	sum := 0
	for _, val := range o {
		sum += val
	}
	return sum
}

func (o IntOrganism) crossover(other Organism) Organism {
	o2 := other.(IntOrganism)
	if len(o) != len(o2) {
		panic("Can't crossover IntOrganisms of different length")
	}

	i := rand.Intn(len(o))
	child := make(IntOrganism, 0, len(o))
	child = append(child, o[:i]...)
	child = append(child, o2[i:]...)

	//fmt.Printf("%v[:%d] + %v[%d:] -> %v\n", o, i, o2, i, child)
	return child
}

func (o IntOrganism) mutate() Organism {
	return o
}

func (p Population) evolve(iters int) {
	fmt.Printf("Iteration: 0; Size: %d; Value: %d; Best: %d\n",
		len(p.organisms), p.value(), p.best().value())

	for i := 1; i <= iters; i++ {
		parents := p.selectParents()
		children := make(Organisms, 0, len(parents)/2)

		for j := 0; j < len(parents); j += 2 {
			child := parents[j].crossover(parents[j+1])
			children = append(children, child.mutate())
		}

		numKilled := p.killWeakestOrganisms()

		fmt.Printf("%d children born, %d organisms die\n", len(children),
			numKilled)

		p.organisms = append(p.organisms, children...)

		fmt.Printf("Iteration: %d; Size: %d; Value: %d; Best: %d\n",
			i, len(p.organisms), p.value(), p.best().value())
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	NUM := 100
	organisms := make(Organisms, NUM)
	for i := 0; i < NUM; i++ {
		organisms[i] = NewIntOrganism(100)
	}

	pop := Population{organisms, .4, .4}
	pop.evolve(100)
}
