package beebot

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/cmdwhitelist"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/combat"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/message"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

func (client *BotGameClient) doNextUnitMove(y, x int, unit objects.Unit) *message.GameMessage {
	var moves []*message.GameMessage
	var scores []int
	// by default, doesn't do anything
	moves = append(moves, &message.GameMessage{
		Cmd:    message.CmdUnitStay,
		Sender: client.UserID,
		Data: &message.UnitStayMessageData{
			Y1: y,
			X1: x,
		},
	})
	scores = append(scores, client.scoreNearestQueen(y, x))

	switch unit.GetUnitType() {
	case objects.UnitTypeInfantry:
		calcMoves, calcScores := client.calcUnitMove(y, x, unit)
		moves = append(moves, calcMoves...)
		scores = append(scores, calcScores...)
	}

	var bestMove *message.GameMessage
	var bestScore = -1000000000
	for i, move := range moves {
		// rate the unit move
		if scores[i] > bestScore {
			bestScore = scores[i]
			bestMove = move
		}
	}

	if bestMove == nil {
		panic("beebot panic: no move selected")
	}
	return bestMove
}

// calc functions to brute-force possible moves

// calc for UNIT_MOVE and possibly UNIT_MOVE_ATTACK
func (client *BotGameClient) calcUnitMove(y, x int, unit objects.Unit) ([]*message.GameMessage, []int) {
	gameLoader := client.Hub.GameLoader
	ge := gameLoader.GridEngine
	var moves []*message.GameMessage
	var scores []int

	switch unit.GetMoveType() {
	case objects.MoveTypeGround:
		ge.FillMoveGround(y, x, unit.GetMoveRange(), unit.GetUnitOwner(), unit.GetWeight())
		for i := 0; i < gameLoader.Height; i++ {
			for j := 0; j < gameLoader.Width; j++ {
				if ge.Dist[i][j] <= 0 {
					continue
				}
				if gameLoader.Units[i][j] != nil {
					continue
				}
				// valid move here
				// move only
				moves = append(moves, &message.GameMessage{
					Cmd:    message.CmdUnitMove,
					Sender: client.UserID,
					Data: &message.UnitMoveMessageData{
						Y1: y,
						X1: x,
						Y2: i,
						X2: j,
					},
				})
				moveScore := client.scoreNearestQueen(i, j)
				scores = append(scores, moveScore)

				// move and attack
				if _, ok := cmdwhitelist.UnitMoveAndAttackMap[unit.GetUnitType()]; ok {
					atkRange := unit.GetAttackRange()
					for ti := i - atkRange; ti <= i+atkRange; ti++ {
						for tj := j - atkRange; tj <= j+atkRange; tj++ {
							if ti < 0 || ti >= gameLoader.Height || tj < 0 || tj >= gameLoader.Width {
								continue
							}
							atkDist := utils.HexDistance(i, j, ti, tj)
							if atkDist > atkRange {
								continue
							}
							targetUnit := gameLoader.Units[ti][tj]
							if targetUnit == nil {
								continue
							}
							if targetUnit.GetUnitOwner() == client.PlayerOrder {
								continue // cannot attack friend
							}
							// ok, can attack target unit
							moves = append(moves, &message.GameMessage{
								Cmd:    message.CmdUnitMoveAndAttack,
								Sender: client.UserID,
								Data: &message.UnitMoveAndAttackMessageData{
									Y1: y,
									X1: x,
									Y2: i,
									X2: j,
									YT: ti,
									XT: tj,
								},
							})
							moveAndAttackScore := moveScore + client.scoreCombat(unit, targetUnit, atkDist)
							scores = append(scores, moveAndAttackScore)
						}
					}
				}
			}
		}
		ge.FillMoveGroundReset(y, x)
	}

	return moves, scores
}

// move scorers

func (client *BotGameClient) scoreNearestQueen(y, x int) int {
	distQueen := 1000000000
	for _, queenPos := range client.otherQueenPositions {
		distQueen = utils.MinInt(distQueen, utils.HexDistance(y, x, queenPos.Y, queenPos.X))
	}
	return -distQueen
}

func (client *BotGameClient) scoreCombat(attacker, defender objects.Unit, dist int) int {
	dmgAtk, dmgDef := combat.SimulateNormalCombat(attacker, defender, dist)
	score := 0
	// damage dealt
	score += dmgDef * defender.GetCost() / defender.GetMaxHP()
	// damage received
	score -= dmgAtk * attacker.GetCost() / attacker.GetMaxHP()
	return score
}
