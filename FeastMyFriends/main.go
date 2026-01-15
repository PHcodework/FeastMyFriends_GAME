//main method
package main

import(
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
	"strconv"
)

//func RandomizeBackground: randomly generates a background for each round of the game
func RandomizeBackground() string{
	backgrounds := []string{"assets/backgrounds/forest.png","assets/backgrounds/sidewalk.png", "assets/backgrounds/picnic.png", "assets/backgrounds/beach.png"}
	background := backgrounds[rand.Intn(len(backgrounds))] //selects a sprite from a random index in the slice
	return background
}

//func main: main method
func main(){
	rl.InitWindow(1920, 1080, "FINAL")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyQ)	//sets the quit/exit button to Q instead of esc (esc opens the pause menu)

	//initializes the timer used in the game
	initialTimer := float32(50)
	timer := initialTimer

	spawnTimer := float32(0) //initializes the spawn timer (used to time the ants spawning to avoid overpopulation/underpopulation)

	//sets the background of the game
	background := rl.LoadTexture(RandomizeBackground())
	defer rl.UnloadTexture(background)

	//load all screens present in the game (game over, pause, and start)
	gameoverscreen := rl.LoadTexture("assets/textures/gameover.png")
	startscreen := rl.LoadTexture("assets/textures/startscreen.png")
	pausescreen := rl.LoadTexture("assets/textures/pausemenu.png")
	defer rl.UnloadTexture(gameoverscreen)
	defer rl.UnloadTexture(startscreen)
	defer rl.UnloadTexture(pausescreen)

	//load the sprites used for the player, food, and ants
	playerSprite := rl.LoadTexture("assets/sprites/mainant.png")
	antSprite := rl.LoadTexture("assets/sprites/ant.png")
	foodSprite := rl.LoadTexture(RandomizeFood())

	//initializes the player object
	player := Player {
		Pos: rl.NewVector2(500,500),	//set initial position to (500, 500)
		Sprite: playerSprite,			
		Speed: 5,						
		AntsFollowing: 0,
		Direction: 1,
	}

	//initializes the food object
	food := Food {
		Pos: rl.NewVector2(rand.Float32() * (1220) , rand.Float32() * (680)),	//randomly generates initial position
		Sprite: foodSprite,
		Width: 400,
		Height: 400,
		Anim: Animation{	//creates the animation used for the food
			Pos:         rl.NewVector2(0,0),
			SpriteSheet: foodSprite,
			Index:       0,
			MaxIndex:    2,
			Opacity:	 255,
		},
		AntsOnFood: 0,
		Level: 10,	//sets the level to 10 (10 ants are needed to eat the food)
	}
	food.Anim.Pos = food.Pos //match the animation positon to the food's position

	ants := []Ant{} //create an empty slice to hold all ants in the game

	//initialize all background/screen music
	rl.InitAudioDevice()
	bgmusic := rl.LoadMusicStream("assets/audio/bgmusic.mp3")
	gameovermusic := rl.LoadMusicStream("assets/audio/gameovermusic.mp3")
	startmusic := rl.LoadMusicStream("assets/audio/startmusic.mp3")
	rl.PlayMusicStream(bgmusic)
	rl.PlayMusicStream(gameovermusic)
	rl.PlayMusicStream(startmusic)

	//initialize all sound effects
	join := rl.LoadSound("assets/audio/join.mp3")
	rottenaway := rl.LoadSound("assets/audio/rottenaway.mp3")
	win := rl.LoadSound("assets/audio/win.mp3")
	defer rl.UnloadSound(join)
	defer rl.UnloadSound(rottenaway)
	defer rl.UnloadSound(win)

	//set "start" gamestate to True at the beginning of the game
	gameover := false
	paused := false
	start := true

	//initialize score and # of rotten foods
	var score int = 0
	var rotcount int = 0

	//beginning of the game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		//if the user is on the start screen, draw the start screen
		if start == true{
			rl.UpdateMusicStream(startmusic) //play the start screen music

			//draw the start screen background
			rl.ClearBackground(rl.RayWhite)
			rl.DrawTexture(startscreen, 0, 0, rl.White)

			//write control info
			rl.DrawText("Press SPACE\nto start!", 200, 600, 50, rl.Black)
			rl.DrawText("Press Q\nto quit :(", 1220, 600, 50, rl.Black)

			//user input handling:
			if rl.IsKeyPressed(rl.KeySpace){ //if the user presses enter, start the game
				start = false
			}
			if rl.IsKeyPressed(rl.KeyQ){ //if the user presses escape, quit
				rl.CloseWindow()
			}
		}

		//if the game isn't over and the user isn't on the start screen, play the game!
		if !gameover && !start {
			//draw the background of the game
			rl.ClearBackground(rl.RayWhite)
			rl.DrawTexture(background, 0, 0, rl.White)

			//if the user has paused the game, freeze the game and draw the pause window
			if paused{
				rl.DrawTexture(pausescreen, 50, 50, rl.White) //draw the pause window

				//user input handling:
				if rl.IsKeyPressed(rl.KeyQ){ //if the user presses Q, quit the game
					rl.CloseWindow()
				}
				if rl.IsKeyPressed(rl.KeyM){ //if the user presses M, return to the menu
					start = true
				}
				if rl.IsKeyPressed(rl.KeySpace){ //if the user presses Space, resume the game
					paused = false
				}
			} else { //if the user hasn't paused, play the game normally

				//count the timer down, stop it at zero
				if timer > 0{
					timer -= rl.GetFrameTime()
				} else {
					timer = 0
				}

				rl.UpdateMusicStream(bgmusic) //continue playing background music while the game is being played

				food.IsTouching(ants) //determines the ants currently touching the food
				if food.AntsOnFood >= food.Level{ //if the number of ants on the food matches the level, eat the food!
					rl.PlaySound(win)
					score++ //increment score

					//iterates through the slice of existing ants, deleting the ants present on the food (saves memory, makes game more challenging)
					for i := len(ants) - 1; i >= 0; i-- {
						if ants[i].IsOnFood(food){
							ants = append(ants[:i], ants[i+1:]...)
						}
					}
					food.SpawnFood() //spawn new location for the food
					timer = initialTimer //reset timer
				}

				food.Anim.DrawFoodAnimation() //draw the food item

				//increment spawn timer
				spawnTimer += rl.GetFrameTime()
				if spawnTimer >= 3 { //if it's been 3 seconds since any new ants were spawned, spawn more ants!
					SpawnRandom(&ants, antSprite)
					SpawnRandom(&ants, antSprite)
					spawnTimer = 0 //resets spawn timer
				}

				//update the player object (move, drawing, and counting the number of ants following)
				player.Move()
				player.DrawPlayer()
				player.UpdateCount(ants)

				//iterates through all existing ants in the game
				for i := 0; i < len(ants); i++ {
					ants[i].DrawAnt()				//draws the current ant
					ants[i].Follow(player, join)	//checks to see if the ant is following the player, acts accordingly

					//if the ant has not yet followed the player, follow the automated marching method
					if !ants[i].HasFollowed {
						ants[i].March()
					}
					
					//if the ant wanders outside the other end of the screen (i.e. it makes a full voyage across the screen, despawn it to save memory)
					if (ants[i].InitialSpawn == 0 && ants[i].Pos.X > 1900) || (ants[i].InitialSpawn == 1 && ants[i].Pos.X < -30) {
						ants = append(ants[:i], ants[i+1:]...)
						i--
					}
					
				}

				//if the timer hits zero, the food rots away
				if(timer == 0){
					food.Rot(rottenaway)  //change food sprite, make noise
					food.IsRotten = true
					rotcount++			  //increment the rot counter
					timer = initialTimer  //reset the timer
				}

				//if the food is rotten, slowly fade away
				if (food.IsRotten){
					if (food.Anim.Opacity > 0){
						food.Anim.Opacity -= 3	//decrease the opacity bit by bit
					} else{ //if the food is fully faded, randomly spawn a randomly generated food somewhere else in the map
						food.SpawnFood()
					}
				}

				if rl.IsKeyPressed(rl.KeyEscape){ //if the user presses escape, pause the game
					paused = true
				}

				//draw the text elementson screen (score, rotcount, number of ants following) in the upper left side of the screen
				rl.DrawText("Food consumed: " + strconv.Itoa(score), 30, 30, 30, rl.Black)
				rl.DrawText("Food ROTTEN: " + strconv.Itoa(rotcount), 30, 60, 30, rl.Black)
				rl.DrawText("Friends following: " + strconv.Itoa(player.AntsFollowing), 30, 90, 30, rl.Black)

				//draw the # of ants above/below the food items (depending on the food's location)
				if food.Pos.Y < 200{ //if the food's too high on the screen, adjust the text to be under the food
					rl.DrawText("Friends CONSUMING: " + strconv.Itoa(int(food.AntsOnFood)), int32(food.Pos.X + (food.Width / 2) - 200), int32(food.Pos.Y + 325), 30, rl.Black)
				} else{ //if the food isn't too high on the screen, the text is displayed above the food
					rl.DrawText("Friends CONSUMING: " + strconv.Itoa(int(food.AntsOnFood)), int32(food.Pos.X + (food.Width / 2) - 200), int32(food.Pos.Y - 50), 30, rl.Black)
				}

				//display the timer in the upper right side of the screen
				rl.DrawRectangle(1380, 30, 600, 50, rl.Black)
				rl.DrawText("TIMER: " + strconv.Itoa(int(timer)), 1400, 30, 50, rl.White)

				if rotcount >= 3{ //if 3 foods have rotted away, GAME OVER
					gameover = true
				}
			}
		} else if gameover && !start{ //if the player has lost, display the game over screen
			//draw the game over screen
			rl.ClearBackground(rl.RayWhite)
			rl.DrawTexture(gameoverscreen, 0, 0, rl.White)

			rl.UpdateMusicStream(gameovermusic) //play the game over music

			ants = nil //destroy all ants

			//if the user presses r, restart the game and reset the game variables
			if rl.IsKeyPressed(rl.KeyR){
				gameover = false
				rotcount = 0
				score = 0
				timer = initialTimer
				food.SpawnFood()
			}
		}
		rl.EndDrawing()
	}
}