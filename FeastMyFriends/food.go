//food.go: defines the properties and behavior of the food objects
package main

import(
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

//struct definition of the Food object
type Food struct {
	Pos rl.Vector2			//current position of the food item
	Sprite rl.Texture2D		//sprite used for the food item
	Width float32			//width of the food item (same for all sprites)
	Height float32			//height of the food item (same for all sprites)
	Anim Animation			//animation (spritesheet) used for the food
	AntsOnFood int			//number of ants currently on the food
	Level int				//level of the food (# of ants required to eat the food)
	IsRotten bool			//status on if the food is rotten or not
}

//func IsTouching: used to determine how many ants are currently on the food item
func (f *Food) IsTouching(ants []Ant){
	f.AntsOnFood = 0

	//iterates through all ants in the map
	for i := 0; i < len(ants); i++{
		a := ants[i]
		if a.Pos.X <= (f.Pos.X + f.Width) && a.Pos.X >= f.Pos.X && a.Pos.Y <= (f.Pos.Y + f.Height) && a.Pos.Y >= f.Pos.Y{
			f.AntsOnFood++ //if an ant is located within the dimensions of the food item, it is on the food
			//^^update the count
		}
	}
}

//func DrawFood: used to draw the food item
func (f *Food) DrawFood(){
	rl.DrawTextureEx(f.Sprite, f.Pos, 0, .8, rl.White) //note: scaled down a bit because the 500 x 500 dimensions were a tad overwhelming
}

//func SpawnFood: used to randomly spawn a (randomly generated) food item somewhere in the map
func (f *Food) SpawnFood(){
	//randomly generate a position for the food to spawn in:
	f.Pos = rl.NewVector2(rand.Float32() * (1000) + 550 , rand.Float32() * (530) + 150)

	//sets the number of ants on food and rotten status to initial values
	f.AntsOnFood = 0
	f.IsRotten = false

	//animation settings reverted/adjusted
	f.Sprite = rl.LoadTexture(RandomizeFood()) //generates a random sprite for the food to be (one of 5 pastries)
	f.Anim.Pos = f.Pos	//draws the sprite in the food's current position
	f.Anim.SpriteSheet = f.Sprite
	f.Anim.Index = 0	//sets the sprite to unrotten/uneaten
	f.Anim.Opacity = 255	//sets the sprite to full opacity
}

//func RandomizeFood: used to randomly generate a sprite for the food
func RandomizeFood() string{
	foodsprites := []string{"concha", "marranito", "smiley", "sprinkle", "donutsheet"}	
	foodsprite := "assets/sprites/" + foodsprites[rand.Intn(len(foodsprites))] + ".png"	//selects a sprite from a random index in the slice
	return foodsprite
}

//func Rot: used to "rot" the food when the timer runs out in main.go
func (f *Food) Rot(rot rl.Sound){
	f.Anim.Index = 1	//set the sprite to the "rotten" sprite
	rl.PlaySound(rot)	//plays the rotting sound
}