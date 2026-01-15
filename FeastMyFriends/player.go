//player.go: defines the properties and behavior of the player
package main

import(
	rl "github.com/gen2brain/raylib-go/raylib"
)

//struct definition of the Player object
type Player struct {
	Pos rl.Vector2			//current position of the player
	Sprite rl.Texture2D		//sprite used for the player
	Speed float32			//speed of the player (adjustable for playtesting purposes)
	AntsFollowing int		//number of ants currently following the player
	Direction int			//direction the player is currently facing (used for sprite flipping)
}

//func Move: used to define how the player moves
func (p *Player) Move() {
	if rl.IsKeyDown(rl.KeyW){ //if W is pressed, move up
		p.Pos.Y -= p.Speed
	}

	if rl.IsKeyDown(rl.KeyA){ //if A is pressed, move to the left
		p.Pos.X -= p.Speed
		p.Direction = 1	//change sprite direction to face left as well
	}

	if rl.IsKeyDown(rl.KeyS){ //if S is pressed, move down
		p.Pos.Y += p.Speed
	}

	if rl.IsKeyDown(rl.KeyD){ //if D is pressed, move to the right
		p.Pos.X += p.Speed
		p.Direction = -1 //change sprite direction to face right as well
	}
	
	//the following lines set the boundaries for the screen so the player can't walk out of the window
	if p.Pos.X < 0{ //prevents the player from walking too far left
		p.Pos.X = 0
	}
	if p.Pos.X > 1880{ //prevents the player from walking too farright
		p.Pos.X = 1880
	}
	if p.Pos.Y < 0{ //prevents the player from walking too far up
		p.Pos.Y = 0
	}
	if p.Pos.Y > 1000{ //prevents the player from walking too far down
		p.Pos.Y = 1000
	}

}

//func DrawPlayer: used to draw the player's sprite in the right position & direction
func (p *Player) DrawPlayer() {
	src := rl.NewRectangle(0, 0, float32(p.Sprite.Width), float32(p.Sprite.Height))

	if p.Direction == -1{	//if the player is facing right, flip the sprite to match that
		src.Width = -src.Width
		src.X = float32(p.Sprite.Width)
	}

	destRect := rl.NewRectangle(p.Pos.X, p.Pos.Y, float32(p.Sprite.Width) * 2, float32(p.Sprite.Height) * 2) //sets the destination to the current player position and sprite dimensions
	//note: sprite dimensions were made larger because it was easier to see

	origin := rl.NewVector2(0, 0)
	rl.DrawTexturePro(p.Sprite, src, destRect, origin, 0, rl.White) //draws the sprite (uses drawtexturepro because it works well with directional flipping)
}

//func UpdateCount: used to update the count of ants currently following the player (implemented for the "friends following" feature)
func (p *Player) UpdateCount(ants []Ant){
	p.AntsFollowing = 0
	//iterates through all ants in the map, detecting if they've joined or not
	for i := 0; i < len(ants); i++{
		a := ants[i]
		if a.Joined{
			p.AntsFollowing += 1 //if the ant is currently joined, add it to the total count of ants following
		}
	}
}