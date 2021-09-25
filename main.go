package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type dice []int

func (d dice) randomizeDice() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(d); i++ {
		d[i] = rand.Intn(7-1) + 1
	}
}

type Player struct {
	order                int
	dice                 dice
	points               int
	playing              bool
	one_from_prev_player int
}

type Scoreboard struct {
	player int
	points int
	dice   dice
}

func init_players(num_player int, num_dice int) []Player {

	var players []Player

	for i := 0; i < num_player; i++ {
		var dice = make([]int, num_dice)
		var player = Player{
			order:                i,
			dice:                 dice,
			points:               0,
			playing:              true,
			one_from_prev_player: 0,
		}
		players = append(players, player)
	}

	return players
}

func init_scoreboards(num_player int) []Scoreboard {

	var scoreboards []Scoreboard

	for i := 0; i < num_player; i++ {
		scoreboard := Scoreboard{
			player: i,
			dice:   []int{},
			points: 0,
		}
		scoreboards = append(scoreboards, scoreboard)
	}

	return scoreboards
}

func evaluate(players *[]Player, scoreboards *[]Scoreboard) {

	for i := 0; i < len(*players); i++ {
		if len((*players)[i].dice) == 0 {
			// if (*players)[i].points > result["points"] {
			// 	result["player"] = (*players)[i].order
			// 	result["points"] = (*players)[i].points
			// }

			for j := 0; j < len(*scoreboards); j++ {
				if (*scoreboards)[j].player == (*players)[i].order {
					(*scoreboards)[j].dice = (*players)[i].dice[:0]
					(*scoreboards)[j].points = (*players)[i].points
				}
			}

			(*players) = append((*players)[:i], (*players)[i+1:]...)
			i--
			continue
		}

		if (*players)[i].playing {
			for j := 0; j < len((*players)[i].dice); j++ {
				if (*players)[i].dice[j] == 6 {
					(*players)[i].points += 1
					(*players)[i].dice = append((*players)[i].dice[:j], (*players)[i].dice[j+1:]...)
					j--
				} else if (*players)[i].dice[j] == 1 {
					next := i + 1
					if next >= len((*players)) {
						next = next - len((*players))
					}
					(*players)[next].one_from_prev_player += 1

					(*players)[i].dice = append((*players)[i].dice[:j], (*players)[i].dice[j+1:]...)
					j--
				}
			}
		}
	}

	for _, player := range *players {
		for i := 0; i < player.one_from_prev_player; i++ {
			player.dice = append(player.dice, 1)
		}

		for j := 0; j < len(*scoreboards); j++ {
			if (*scoreboards)[j].player == player.order {
				(*scoreboards)[j].dice = player.dice[:]
				(*scoreboards)[j].points = player.points
			}
		}

		player.one_from_prev_player = 0
	}
}

func main() {
	var num_players int
	var num_dice int

	fmt.Println("Enter number of players: ")
	fmt.Scanln(&num_players)
	fmt.Println("Enter number of dice: ")
	fmt.Scanln(&num_dice)

	players := init_players(num_players, num_dice)
	scoreboards := init_scoreboards(num_players)

	i := 0
	for len(players) > 1 {
		fmt.Println("========================================== Turn ", i+1)
		for _, player := range players {
			player.dice.randomizeDice()
		}

		fmt.Printf("Lempar dadu %d\n", i+1)
		for _, player := range players {
			var dice []string

			for _, die := range player.dice {
				dice = append(dice, strconv.Itoa(die))
			}

			joined_dice := strings.Join(dice, ", ")

			fmt.Printf("Player %d, Dice: %s, Points: %d\n", player.order, joined_dice, player.points)
		}

		evaluate(&players, &scoreboards)

		fmt.Println("Setelah evaluasi")

		for _, scoreboard := range scoreboards {
			var dice []string

			for _, die := range scoreboard.dice {
				dice = append(dice, strconv.Itoa(die))
			}

			joined_dice := strings.Join(dice, ", ")

			fmt.Printf("Player %d, Dice: %s, Points: %d\n", scoreboard.player, joined_dice, scoreboard.points)
		}
		i++

	}

	highest := 0
	winner := 0

	for _, scoreboard := range scoreboards {
		if scoreboard.points > highest && scoreboard.player != players[0].order {
			highest = scoreboard.points
			winner = scoreboard.player
		}
	}

	fmt.Printf("Player %d is the winner with %d points", winner, highest)
}
