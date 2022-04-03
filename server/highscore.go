package server

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const highscoreFileName = "highscore.json"

type Score struct {
	Name  string
	Value int
}

type Highscore struct {
	mutex  sync.Mutex
	scores []Score
}

func (h *Highscore) init() {
	if h.scores != nil {
		return
	}

	if _, err := os.Stat(highscoreFileName); os.IsNotExist(err) {
		f, err := os.Create(highscoreFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = f.Write([]byte("[]"))
		if err != nil {
			fmt.Println(err)
			return
		}
		f.Close()
	} else if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.ReadFile(highscoreFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(file, &h.scores)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Highscore) persist() {
	f, err := os.OpenFile(highscoreFileName, os.O_TRUNC|os.O_WRONLY, 0)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.NewEncoder(f).Encode(h.scores)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Highscore) Get() []Score {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.init()

	var highscoreCopy = make([]Score, len(h.scores))
	copy(highscoreCopy, h.scores)
	return highscoreCopy
}

func (h *Highscore) Insert(newScore Score) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.init()

	if len(h.scores) == 0 {
		h.scores = append(h.scores, newScore)
		h.persist()
		return
	}

	for i, score := range h.scores {
		if score.Value <= newScore.Value {
			h.scores = append(h.scores[:i+1], h.scores[i:]...)
			h.scores[i] = newScore
			break
		}
	}

	h.persist()
}
