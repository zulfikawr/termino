package game

import (
	"math/rand"
	"time"

	"termino/internal/tetromino"
)

type Randomizer struct {
	currentBag []tetromino.Tetromino
	nextBag    []tetromino.Tetromino
	rng        *rand.Rand
}

func NewRandomizer() *Randomizer {
	r := &Randomizer{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	r.currentBag = r.createNewBag()
	r.nextBag = r.createNewBag()
	return r
}

func (r *Randomizer) Next() tetromino.Tetromino {
	// Ensure we always have pieces
	if len(r.currentBag) == 0 {
		r.currentBag = r.nextBag
		r.nextBag = r.createNewBag()
	}

	// Pop from current bag
	piece := r.currentBag[0]
	r.currentBag = r.currentBag[1:]
	return piece
}

func (r *Randomizer) createNewBag() []tetromino.Tetromino {
	// Create new bag of 7 pieces
	pieces := []string{"I", "J", "L", "O", "S", "T", "Z"}
	bag := make([]tetromino.Tetromino, len(pieces))

	// Shuffle pieces
	r.rng.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	for i, name := range pieces {
		bag[i] = tetromino.NewTetromino(name)
	}

	return bag
}
