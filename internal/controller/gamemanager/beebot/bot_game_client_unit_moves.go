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

	// move + move and attack
	switch unit.UnitType() {
	case objects.UnitTypeQueen:
		// queen doesn't move
	default:
		calcMoves, calcScores := client.calcUnitMove(y, x, unit)
		moves = append(moves, calcMoves...)
		scores = append(scores, calcScores...)
	}

	// indirect attackers
	if unit.UnitType() == objects.UnitTypeWizard || unit.UnitType() == objects.UnitTypeMortar {
		calcMoves, calcScores := client.calcUnitAttack(y, x, unit)
		moves = append(moves, calcMoves...)
		scores = append(scores, calcScores...)
	}

	// by default, doesn't do anything. Staying is prioritized the least
	moves = append(moves, &message.GameMessage{
		Cmd:    message.CmdUnitStay,
		Sender: client.UserID,
		Data: &message.UnitStayMessageData{
			Y1: y,
			X1: x,
		},
	})
	scores = append(scores, client.scorePosition(y, x))

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

	switch unit.UnitMoveType() {
	case objects.MoveTypeGround:
		ge.FillMoveGround(y, x, unit.UnitMoveRange(), unit.GetOwner(), unit.UnitWeight())
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
				moveScore := client.scorePosition(i, j)
				scores = append(scores, moveScore)

				// move and attack
				if _, ok := cmdwhitelist.UnitMoveAndAttackMap[unit.UnitType()]; ok {
					atkRange := unit.UnitAttackRange()
					for ti := i - atkRange; ti <= i+atkRange; ti++ {
						for tj := j - atkRange; tj <= j+atkRange; tj++ {
							var atkDist int
							if ok, atkDist = ge.ValidateAttack(i, j, ti, tj, unit); !ok {
								continue
							}
							targetUnit := gameLoader.Units[ti][tj]
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
	case objects.MoveTypeBlink:
		for i := y - unit.UnitMoveRange(); i <= y+unit.UnitMoveRange(); i++ {
			for j := x - unit.UnitMoveRange(); j <= x+unit.UnitMoveRange(); j++ {
				if !ge.ValidateMove(y, x, i, j) {
					continue
				}
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
				moveScore := client.scorePosition(i, j)
				scores = append(scores, moveScore)
			}
		}
	}

	return moves, scores
}

// calc for UNIT_ATTACK
func (client *BotGameClient) calcUnitAttack(y, x int, unit objects.Unit) ([]*message.GameMessage, []int) {
	gameLoader := client.Hub.GameLoader
	ge := gameLoader.GridEngine
	var moves []*message.GameMessage
	var scores []int

	if _, ok := cmdwhitelist.UnitAttackMap[unit.UnitType()]; ok {
		moveScore := client.scorePosition(y, x)
		atkRange := unit.UnitAttackRange()
		for i := y - atkRange; i <= y+atkRange; i++ {
			for j := x - atkRange; j <= x+atkRange; j++ {
				var atkDist int
				if ok, atkDist = ge.ValidateAttack(y, x, i, j, unit); !ok {
					continue
				}
				targetUnit := gameLoader.Units[i][j]
				// ok, can attack target unit
				moves = append(moves, &message.GameMessage{
					Cmd:    message.CmdUnitAttack,
					Sender: client.UserID,
					Data: &message.UnitAttackMessageData{
						Y1: y,
						X1: x,
						YT: i,
						XT: j,
					},
				})
				moveAndAttackScore := moveScore + client.scoreCombat(unit, targetUnit, atkDist)
				scores = append(scores, moveAndAttackScore)
			}
		}
	}

	return moves, scores
}

// move scorers

func (client *BotGameClient) scorePosition(y, x int) int {
	score := 0
	// nearest queen
	distQueen := 1000000000
	for _, queenPos := range client.otherQueenPositions {
		distQueen = utils.MinInt(distQueen, utils.HexDistance(y, x, queenPos.Y, queenPos.X))
	}
	score -= distQueen
	return score
}

func (client *BotGameClient) scoreCombat(attacker, defender objects.Unit, dist int) int {
	dmgAtk, dmgDef := combat.SimulateCombat(attacker, defender, dist, false)
	score := 0
	// damage dealt
	score += dmgDef * defender.UnitCost() / defender.UnitMaxHP()
	// damage received
	score -= dmgAtk * attacker.UnitCost() / attacker.UnitMaxHP()
	return score
}
