package loader

import "gitlab.com/beewar/beewar-be/internal/access/model"

// this function takes in a slice of game users and the intended total slice length, and pads the missing slots with default values.
// this function expects valid data (e.g. no duplicate player_order)
func padGameUsers(gameUsers []*model.GameUser, playerCount int) []*model.GameUser {
	newGameUsers := make([]*model.GameUser, playerCount)
	for _, gu := range gameUsers {
		newGameUsers[gu.PlayerOrder-1] = gu
	}
	for i := range newGameUsers {
		if newGameUsers[i] == nil {
			newGameUsers[i] = &model.GameUser{
				PlayerOrder: uint8(i + 1),
			}
		}
	}
	return newGameUsers
}
