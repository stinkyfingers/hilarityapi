package storage

import (
	"testing"

	"github.com/stinkyfingers/hilarity/game"
	"github.com/stinkyfingers/hilarity/user"
	"github.com/stretchr/testify/require"
)

func TestSubmitPlay(t *testing.T) {
	cache := NewInMemory()
	g, err := game.NewGame("test-1")
	require.Nil(t, err)
	err = cache.SaveGame(g)
	require.Nil(t, err)
	u := user.User{
		Name: "john",
	}
	g.Join(u)

	gp := game.GamePlay{
		GameName: g.Name,
		User:     u,
		Play: game.Play{
			Answers: [3]string{"chocolate", "vanilla", "strawberry"},
		},
	}
	expected := &game.Game{
		Name:  "test-1",
		Users: map[string]user.User{u.Name: u},
		CurrentRound: game.Round{
			Question: "top 3 ice cream flavors",
			Guesses:  map[string]game.Guess{},
			Plays: map[string]game.Play{
				u.Name: {
					Answers: [3]string{"chocolate", "vanilla", "strawberry"},
				},
			},
		},
	}
	g, err = SubmitPlay(cache, gp)
	require.Nil(t, err)
	require.Equal(t, expected, g)
}
