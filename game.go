package main

import (
	"fmt"
	"github.com/visualfc/atk/tk"
	"strconv"
)

// Game holds a game state and necessary game variables
type Game struct {
	assetsManager    *AssetsManager
	board            *Board
	gameWindow       *GameWindow
	gameLayout       *tk.GridLayout
	gameBoardButtons [][]*tk.Button
	currentBoardSize int
	gameStatus       GameStatus
}

// GameWindow represents a drawable GUI surface for the game.
type GameWindow struct {
	*tk.Window
}

// GameMode holds an information about available game modes.
type GameMode struct {
	Title     string
	BoardSize int
}

// GameStatus holds current game status/
type GameStatus int

const (
	// GameNotStarted indicates a game is not yet started
	GameNotStarted GameStatus = 0
	// GameInProgress indicates a game in progress.
	GameInProgress GameStatus = 1
	// GameWon indicates a game is finished and a player has won.
	GameWon GameStatus = 2
	// GameLost indicates a game is finished and a player has lost
	GameLost GameStatus = 3
)

// Hardcoded available game modes
var availableGameModes = []GameMode{
	{
		Title:     EasyGameModeTitle,
		BoardSize: 9,
	},
	{
		Title:     MediumGameModeTitle,
		BoardSize: 16,
	},
	{
		Title:     HardGameModeTitle,
		BoardSize: 20,
	},
}

// NewGame returns an instance of a Game object
func NewGame() *Game {
	return &Game{
		NewAssetsManager(),
		nil,
		nil,
		nil,
		nil,
		0,
		GameNotStarted,
	}
}

// Start invokes a game
func (g *Game) Start() {
	err := tk.MainLoop(func() {
		g.gameWindow = &GameWindow{tk.RootWindow()}

		g.assetsManager.LoadAssets()

		g.drawNewGameWindow()

		g.gameWindow.SetTitle(GameTitle)
		g.gameWindow.Center(nil)
		g.gameWindow.ShowNormal()
	})

	if err != nil {
		fmt.Println(err.Error())
	}
}

// drawNewGameWindow responsible for drawing main game screen
func (g *Game) drawNewGameWindow() {
	if g.gameLayout != nil {
		g.gameLayout.Destroy()
	}

	g.gameLayout = tk.NewGridLayout(g.gameWindow)
	currentGridRow := 0

	welcomeLbl := tk.NewLabel(g.gameWindow, WelcomeMessage)
	welcomeLbl.SetFont(tk.NewUserFont("San Francisco", 22))

	g.gameLayout.SetBorderWidth(10)
	g.gameLayout.AddWidget(
		welcomeLbl,
		tk.GridAttrRow(currentGridRow),
		tk.GridAttrColumn(0),
		tk.GridAttrColumnSpan(10),
		tk.GridAttrPadx(100),
		tk.GridAttrPady(30),
	)

	currentGridRow++

	for _, gameMode := range availableGameModes {
		gameMode := gameMode
		gameModeBtn := tk.NewButton(g.gameWindow, gameMode.Title)

		gameModeBtn.OnCommand(func() {
			g.currentBoardSize = gameMode.BoardSize
			g.drawPlayBoard()
		})
		g.gameLayout.AddWidget(
			gameModeBtn,
			tk.GridAttrRow(currentGridRow),
			tk.GridAttrColumn(0),
			tk.GridAttrColumnSpan(10),
			tk.GridAttrPady(5),
		)

		currentGridRow++
	}

	quitBtn := tk.NewButton(g.gameWindow, QuitButtonTitle)

	quitBtn.OnCommand(func() {
		tk.Quit()
	})
	g.gameLayout.AddWidget(
		quitBtn,
		tk.GridAttrRow(currentGridRow),
		tk.GridAttrColumn(0),
		tk.GridAttrColumnSpan(10),
		tk.GridAttrPady(25),
	)
	g.gameWindow.Center(nil)
}

// drawPlayBoard responsible for drawing gaming board of a given size
func (g *Game) drawPlayBoard() {
	size := g.currentBoardSize
	g.gameLayout.Destroy()
	g.gameLayout = tk.NewGridLayout(g.gameWindow)
	g.gameBoardButtons = make([][]*tk.Button, size)
	g.board = NewBoard(size, size, 0.25)
	g.gameStatus = GameInProgress
	restartBtn := tk.NewButton(g.gameWindow, RestartButtonText)
	mainMenuBtn := tk.NewButton(g.gameWindow, MainMenuButtonText)

	for i := 0; i < size; i++ {
		i := i
		g.gameBoardButtons[i] = make([]*tk.Button, size)

		for j := 0; j < size; j++ {
			j := j
			g.gameBoardButtons[i][j] = tk.NewButton(g.gameWindow, "")
			g.gameBoardButtons[i][j].SetImage(g.assetsManager.GetAsset("default"))

			g.gameBoardButtons[i][j].OnCommand(func() {
				g.unfoldCell(i, j)
			})
			g.gameLayout.AddWidget(
				g.gameBoardButtons[i][j],
				tk.GridAttrRow(i),
				tk.GridAttrColumn(j),
			)
		}
	}

	g.gameLayout.AddWidget(
		restartBtn,
		tk.GridAttrRow(size),
		tk.GridAttrColumn(0),
		tk.GridAttrColumnSpan(size/2),
		tk.GridAttrPadx(50),
		tk.GridAttrPady(10),
	)
	g.gameLayout.AddWidget(
		mainMenuBtn,
		tk.GridAttrRow(size),
		tk.GridAttrColumn(size/2+1),
		tk.GridAttrColumnSpan(size/2),
		tk.GridAttrPadx(50),
		tk.GridAttrPady(10),
	)
	restartBtn.OnCommand(func() {
		g.drawPlayBoard()
	})
	mainMenuBtn.OnCommand(func() {
		g.drawNewGameWindow()
	})
	g.gameWindow.Center(nil)
}

// drawLostMessage responsible for drawing a lost message
func (g *Game) drawLostMessage() {
	lostLabel := tk.NewLabel(g.gameWindow, LostMessage)
	lostLabel.SetFont(tk.NewUserFont("San Francisco", 46))
	lostLabel.SetForground("#FF0000")

	g.gameLayout.AddWidget(
		lostLabel,
		tk.GridAttrRow(g.currentBoardSize+2),
		tk.GridAttrColumn(0),
		tk.GridAttrColumnSpan(g.currentBoardSize),
		tk.GridAttrPadx(50),
		tk.GridAttrPady(10),
	)

	g.gameWindow.Center(nil)
}

// drawWonMessage responsible for drawing a won message
func (g *Game) drawWonMessage() {
	wonLabel := tk.NewLabel(g.gameWindow, WonMessage)
	wonLabel.SetFont(tk.NewUserFont("San Francisco", 46))
	wonLabel.SetForground("#00FF00")

	g.gameLayout.AddWidget(
		wonLabel,
		tk.GridAttrRow(g.currentBoardSize+2),
		tk.GridAttrColumn(0),
		tk.GridAttrColumnSpan(g.currentBoardSize),
		tk.GridAttrPadx(50),
		tk.GridAttrPady(10),
	)

	g.gameWindow.Center(nil)
}

// unfoldCell is invoked everytime a user hits a cell
func (g *Game) unfoldCell(i, j int) {
	if g.gameStatus != GameInProgress {
		return
	}

	isBomb, dangerZone, revealedCells := g.board.UnfoldCell(i, j)

	if isBomb {
		g.gameStatus = GameLost
		g.gameBoardButtons[i][j].SetImage(g.assetsManager.GetAsset("bomb"))
		g.drawLostMessage()
	} else if dangerZone > 0 {
		g.gameBoardButtons[i][j].SetImage(g.assetsManager.GetAsset(strconv.Itoa(dangerZone)))
	} else {
		g.gameBoardButtons[i][j].SetImage(g.assetsManager.GetAsset("empty"))

		for _, cellCoords := range revealedCells {
			g.gameBoardButtons[cellCoords[0]][cellCoords[1]].SetImage(g.assetsManager.GetAsset("empty"))
		}
	}

	if g.board.IsWon() {
		g.gameStatus = GameWon
		g.drawWonMessage()
	}
}
