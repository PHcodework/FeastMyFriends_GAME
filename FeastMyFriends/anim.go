//anim.go: defines the properties and behavior of the animations used for the food
//note: could have been implemented in the food.go file but was separated for tidiness
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

//struct definition of the Animation object
type Animation struct {
	Pos         rl.Vector2		//current position of the animation
	SpriteSheet rl.Texture2D	//spritesheet used for the animation
	Index       int32			//current sprite in the sheet being used
	MaxIndex    int32			//total number of sprites in the sheet
	Opacity 	uint8			//opacity of the sprite being drawn
}

//func DrawFoodAnimation: used to draw the animation of the food
func (a Animation) DrawFoodAnimation() {

	//sets the width and height of the current sprite (same for all sprites used in the animation)
	frameWidth := float32(a.SpriteSheet.Width / (a.MaxIndex + 1))
	frameHeight := float32(a.SpriteSheet.Height)

	sourceRect := rl.NewRectangle(frameWidth*float32(a.Index), 0, frameWidth, frameHeight)

	destRect := rl.NewRectangle(a.Pos.X, a.Pos.Y, frameWidth * .5, frameHeight * .5) //sets the destination to the current player position and sprite dimensions
	opacity := rl.Color{255, 255, 255, a.Opacity}	//sets the opacity of the sprite to the current opacity in the game (opacity is altered when the food rots)
	rl.DrawTexturePro(a.SpriteSheet, sourceRect, destRect, rl.Vector2Zero(), 0, opacity) //draws the sprite (uses drawtexturepro because it works well with directional flipping)
}
