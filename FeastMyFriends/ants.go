//ants.go: defines the properties and behavior of the ants
package main

import(
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

//struct definition of the Ant object
type Ant struct {
	Pos rl.Vector2			//current position of the ant
	Sprite rl.Texture2D		//sprite used for the ant
	Joined bool				//status on if the ant is currently following the player
	HasFollowed bool		//status on if the ant has ever followed the player
	InitialSpawn int		//direction of the initial spawn of the ant
	Direction int			//current direction of the ant
}

//func DrawAnt: used to draw the ant
//note: identical to the DrawPlayer function
func (a *Ant) DrawAnt() {
	src := rl.NewRectangle(0, 0, float32(a.Sprite.Width), float32(a.Sprite.Height))

	if a.Direction == -1{ //if the ant is facing right, flip the sprite to match that
		src.Width = -src.Width
		src.X = float32(a.Sprite.Width)
	}

	destRect := rl.NewRectangle(a.Pos.X, a.Pos.Y, float32(a.Sprite.Width) * 2, float32(a.Sprite.Height) * 2)//sets the destination to the current ant position and sprite dimensions
	//note: sprite dimensions were made larger because it was easier to see

	origin := rl.NewVector2(0, 0)
	rl.DrawTexturePro(a.Sprite, src, destRect, origin, 0, rl.White) //draws the sprite (uses drawtexturepro because it works well with directional flipping)
}

//func Follow: used to make the ant follow the player when within range
func (a *Ant) Follow(target Player, join rl.Sound){
	//calculate the distance between the ant and the player
	dx := target.Pos.X - a.Pos.X	
	dy := target.Pos.Y - a.Pos.Y
	distance := rl.Vector2Length(rl.NewVector2(dx, dy))

	//whent the ant enters the player's following range, follow
	if distance < 200 {
		if !a.Joined { //if the ant isn't already joined, join
			rl.PlaySound(join) //play the "join" sound (sounds like a high-pitched "hi")
			a.Joined = true	//the ant is now joined
			a.HasFollowed = true //the ant has now followed the player at some point
		}

		// if the ant still isn't close enough to the player, move closer to the player
		if distance > 1 {
			dir := rl.Vector2Normalize(rl.NewVector2(dx, dy))
			a.Pos.X += dir.X * 2
			a.Pos.Y += dir.Y * 2
		}
	} else {
		a.Joined = false //if the player moves out of following range, the ant is no longer joined (isn't following anymore)
	}

	//if the ant has joined, alter its direction based on which way it's moving
	if a.Joined == true{
		if dx <= 0 {
			a.Direction = 1	//if the ant is moving left, face left
		} else {
			a.Direction = -1 //if the ant is moving right, face right
		}
	}
	
}

//func IsOnFood: determines if the ant is currently on the food
func (a *Ant) IsOnFood(f Food) bool {
	foodrect := rl.NewRectangle(f.Pos.X, f.Pos.Y, 400, 400)	//creates a rectangle with the dimensions of the food item

	if rl.CheckCollisionPointRec(a.Pos, foodrect) {	//if the ant collides with the rectangle (food), return true
		return true
	}
	return false //if the ant doesn't touch the food, return false
}

//func SpawnRandom: used to randomly spawn ants on either side of the screen
func SpawnRandom(ants *[]Ant, antSprite rl.Texture2D){
	ycoord := rand.Float32() * 980 + 50	//randomly generates a y-coordinate for the ant

	//randomly generates a side for the ant to spawn on
	var xcoord float32 = 0
	sidechooser := rand.Intn(2)
	if sidechooser == 0 { //if the side was 0, spawn offscreen to the left
		xcoord = 0
	} else{ //if the side was 1, spawn offscreen to the right
		xcoord = 1920
	}

	//create a new ant!
	ant := Ant{
		Sprite: antSprite,
		Pos: rl.NewVector2(xcoord, ycoord),	//sets the spawn location to the randomly generated coordinates above
		InitialSpawn: sidechooser,			//sets the initial spawn to the side of the screen the ant spawned on
		Joined: false,
		HasFollowed: false,
		Direction: 1,						//creates a direction, may flip depending on spawn location
	}

	//directional flips based on spawn location:
	if xcoord == 0{ //if the ant spawned on the left side of the screen (and will move right), face right
		ant.Direction = -1
	} else{ //if the ant spawned on the right side of the screen (and will move left), face left
		ant.Direction = 1
	}
	*ants = append(*ants, ant) //add the ant to the slice of all existing ants
}

//func March: autonomous marching function for ants (only applied before the ants have ever joined to the player)
func (a *Ant) March() {
	if a.InitialSpawn == 0 { //if the ant spawned on the left side of the screen, move to the right
		a.Pos.X += 2
	} else { 				 //if the ant spawned on the right side of the screen, move to the left
		a.Pos.X -= 2
	}
}
