package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"kasperaldrin.com/rerego/pkg/core"
)

/*
func TestStoppingCondition(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	A := core.NewSubscription[any]("A", ctx)
	B := core.NewSubscription[any]("B", ctx)
	C := core.NewSubscription[any]("C", ctx)

	Ain := A.NewSender()
	Cout := C.NewReciever()


}
*/

func TestPuzzleGame(t *testing.T) {
	/**
	* In this puzzle game we have a few agents with specialized
	* skills who have to collaborately solve a puzzle. The game takes place on * a NxM grid with the following possible values:
	*	- 0: Empty space
	*	- #: Wall
	*   - S: Start
	*
	*	- Judge: Will check if the puzzle is solved and if so is the case will * 			send a DONE item through it's output which will freeze all *			agents, print the solved puzzle and exit.
	*
	*	- Strategist: Will chose a strategy to get further to solve the puzzle.
	*
	*   - The Agents
	*		-
	*
	*
	*
	*
	*
	 */

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx.Done()

}

func TestJuicer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating two Subscriptions for query and output.
	query := core.NewSubscription[string]("query", ctx)
	output := core.NewSubscription[string]("output", ctx)

	queryChannel := query.NewSender()       // A channel to add messages to.
	responseChannel := output.NewReciever() // ..recieve messages from.

	juicerFunc := func(input core.DataObjectInterface[string]) core.DataObjectInterface[string] {
		dataStr := input.GetData()
		return core.NewData[string](
			"fruit",
			fmt.Sprintf("%s juice", dataStr),
		)
	}

	generateNode := core.
		NewNodeBase[string, string]("generate").
		WithInput(query).
		WithOutput(output).
		UseFunctionWhenAny(juicerFunc)
	go generateNode.Serve(ctx)

	go func() {
		for _, i := range []string{"apple", "banana", "cherry", "date"} {
			time.Sleep(time.Second)

			fmt.Printf("\nsending %s\n", i)

			queryChannel <- core.NewData[string](
				"fruit", i)

			fmt.Printf("%s sent\n", i)
		}
	}()

	// Listen for output in a separate goroutine
	bottle := "milk"
	go func() {
		fmt.Println("waiting for processed fruit")
		for data := range responseChannel {
			fmt.Printf("%v recieved", data)
			if data.GetData() == "cherry juice" {
				bottle = "cherry juice"
				cancel() // Signal to terminate after first output is received
				return
			}

		}
	}()

	// Wait for all processing to be done
	<-ctx.Done()

	if bottle != "cherry juice" {
		t.Errorf("Expected cherry juice, got %s", bottle)
	}
}

func TestRobot(t *testing.T) {

	/*
		A robot system which navigates through a NxM matrix.
		THere are three values in the matrix map, " "(space) denotes an empty space which the robot can occupy, # denotes a wall which the robot cannot occupy, and @ denotes the goal, where the robot needs to go.

		THe position and angle of the robot is marked by <(left), >(right), A(up), and, v(down).

		The robot has the following functions
		- scanning, the robot can scan the environment, in which all surrounding tiles are revealed at a distance of 2 units.
		- rotation, the robot can rotate by 90 degrees clock or counter-clockvise.
		- movement, the robot can move one block at a time in the direction it faces (not backways or sideways)

		- BONUS: if the robot bounces on a wall, game over.

		The vision function will assert that IndexOutOfBounds are walls (all negative indexes, and indexes larger than the map bounds.

		The world is as seen below, but the initial world map for the strategic agent has all tiles hidden except for the robot position. Hidden tiles are denoted by empty strings ("").

	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	real_world := [][]string{
		{"#", "#", "#", "#", "#", "#", "#", "#", "#", "#"},
		{"#", " ", " ", " ", " ", " ", " ", " ", " ", "#"},
		{"#", " ", " ", " ", " ", " ", " ", " ", "@", "#"},
		{"#", " ", " ", " ", " ", " ", " ", "#", "#", "#"},
		{"#", " ", " ", " ", " ", " ", " ", " ", " ", "#"},
		{"#", " ", " ", "#", " ", " ", " ", " ", " ", "#"},
		{"#", "#", "#", "#", " ", "#", " ", " ", " ", "#"},
		{"#", " ", " ", " ", " ", "#", "#", " ", " ", "#"},
		{"#", "A", " ", " ", " ", "#", " ", " ", " ", "#"},
		{"#", "#", "#", "#", "#", "#", "#", "#", "#", "#"},
	}

	rotateClockvise := func(
		world [][]string) [][]string {
		for i := range world {
			for j := range world[i] {
				switch world[i][j] {
				case ">":
					world[i][j] = "v"
				case "v":
					world[i][j] = "<"
				case "<":
					world[i][j] = "A"
				case "A":
					world[i][j] = ">"
				}
			}
		}
		return world
	}

	rotateCounterClockvise := func(
		world [][]string) [][]string {
		for i := range world {
			for j := range world[i] {
				switch world[i][j] {
				case ">":
					world[i][j] = "A"
				case "v":
					world[i][j] = ">"
				case "<":
					world[i][j] = "v"
				case "A":
					world[i][j] = "<"
				}
			}
		}
		return world
	}
	move := func(
		world [][]string) [][]string {
		for i := range world {
			for j := range world[i] {
				switch world[i][j] {
				case ">":
					fmt.Println("moving right, i,j", i, j)
					if len(world[i]) >= j+1 {
						world[i][j] = " "
						world[i][j+1] = ">"
						return world
					}
				case "v":
					if len(world) >= i+1 {
						world[i][j] = " "
						world[i+1][j] = "v"
						return world
					}
				case "<":
					if j-1 >= 0 {
						world[i][j] = " "
						world[i][j-1] = "<"
						return world
					}
				case "A":
					if i-1 >= 0 {
						world[i][j] = " "
						world[i-1][j] = "A"
						return world
					}
				}
			}
		}
		return world
	}

	isIndexOutOfBounds := func(matrix [][]string, row, col int) bool {
		if row < 0 || row >= len(matrix) {
			return true
		}
		if col < 0 || col >= len(matrix[row]) {
			return true
		}
		return false
	}

	scan := func(world [][]string) [][]string {
		// Detect robot position
		for i := range world {
			for j := range world[i] {
				t := world[i][j]
				if t == "A" || t == "v" || t == ">" || t == "<" {
					// Position found at i,j
					indexes_to_fill := [][]int{
						{i + 1, j + 1}, {i + 1, j}, {i + 1, j - 1},
						{i, j + 1}, {i, j - 1},
						{i - 1, j}, {i - 1, j - 1}, {i - 1, j + 1},
					}
					for _, ij := range indexes_to_fill {

						if isIndexOutOfBounds(world, ij[0]-1, ij[1]-1) ||
							isIndexOutOfBounds(world, ij[0]+1, ij[1]+1) {
							world[ij[0]][ij[1]] = "#"
						} else {
							world[ij[0]][ij[1]] = real_world[ij[0]][ij[1]]
						}
					}
					return world
				}
			}
		}
		return world
	}

	/*
		   Logic:

		   There are these types of Stores (Subscribers):
		   - ActionRequest - consumed by actors who do the actions, pushed to by the thinking organs.
		   - ActionResult - consumed by the core, which can either do actions or push to Done once goal is reached.
		   - Goal - consumed by end of program.

			One start the game by pushing the percieved map (everything empty but the position of the actor)

	*/

	type Action int
	const (
		Move Action = iota
		RotationClockwise
		RotationCounterClockwise
		Scan // Scanning the area in front of the robot
	)
	type ActionRequestItem struct {
		Action Action
		World  [][]string
	}

	actionRequest := core.NewSubscription[ActionRequestItem]("areq", ctx)
	actionResult := core.NewSubscription[[][]string]("ares", ctx)
	goal := core.NewSubscription[ActionRequestItem]("goal", ctx)

	gameInitiator := actionResult.NewSender() // send map to start game
	gameGoal := goal.NewReciever()            // recieve to end game

	agentFunction := func(
		input core.DataObjectInterface[ActionRequestItem],
	) core.DataObjectInterface[[][]string] {
		switch input.GetData().Action {
		case Move:
			return core.NewData("move", move(input.GetData().World))
		case RotationClockwise:
			return core.NewData("rotation",
				rotateClockvise(input.GetData().World))
		case RotationCounterClockwise:
			return core.NewData("rotation-counterclockwise",
				rotateCounterClockvise(input.GetData().World))
		case Scan:
			return core.NewData("scan", scan(input.GetData().World))
		default:
			return core.NewData("error", input.GetData().World)
		}
	}

	moveAgent := core.
		NewNodeBase[ActionRequestItem, [][]string]("move_agent").
		WithInput(actionRequest).
		WithOutput(actionResult).
		UseFunctionWhenAny(agentFunction)

	// Position represents a point in the matrix
	type Position struct {
		row, col int
	}

	// directions represent possible movements in the grid
	var directions = []Position{
		{-1, 0}, // up
		{1, 0},  // down
		{0, -1}, // left
		{0, 1},  // right
	}

	var chartoDir = map[string]Position{
		"A": directions[0],
		"v": directions[1],
		"<": directions[2],
		">": directions[3],
	}

	// bfs finds the shortest path from 'A' to '@' using BFS
	bfs := func(grid [][]string, start Position) ([]Position, bool) {
		rows := len(grid)
		cols := len(grid[0])

		// Queue for BFS
		queue := []Position{start}

		// To keep track of visited positions
		visited := make([][]bool, rows)
		for i := range visited {
			visited[i] = make([]bool, cols)
		}
		visited[start.row][start.col] = true

		// To keep track of the path
		parent := make(map[Position]Position)

		for len(queue) > 0 {
			// Dequeue the front position
			current := queue[0]
			queue = queue[1:]

			// If we have reached the target '@'
			if grid[current.row][current.col] == "@" {
				path := []Position{}
				for p := current; p != start; p = parent[p] {
					path = append([]Position{p}, path...)
				}
				path = append([]Position{start}, path...)
				return path, true
			}

			// Explore all possible directions
			for _, dir := range directions {
				newRow := current.row + dir.row
				newCol := current.col + dir.col

				if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols && !visited[newRow][newCol] && grid[newRow][newCol] != "#" {
					visited[newRow][newCol] = true
					queue = append(queue, Position{newRow, newCol})
					parent[Position{newRow, newCol}] = current
				}
			}
		}

		// No path found
		return nil, false
	}

	strategistAgent := core.
		NewNodeBase[[][]string, ActionRequestItem]("strategist_agent").
		WithInput(actionResult).
		WithOutput(actionRequest).
		WithOutput(goal).
		UseFunctionWhenAny(func(input core.DataObjectInterface[[][]string]) core.DataObjectInterface[ActionRequestItem] {
			// get current pos
			world := input.GetData()
			for i := range world {
				for j := range world[i] {
					angle := world[i][j]
					if angle == "A" || angle == "v" || angle == ">" || angle == "<" {
						// Position found at i,j
						start := Position{i, j}
						if path, found := bfs(input.GetData(), start); found {
							// find the next move
							for _, p := range path {
								if p != start {
									// find the direction
									dir := Position{p.row - start.row, p.col - start.col}
									if dir == chartoDir[angle] {
										return core.NewData("move", ActionRequestItem{Move,
											input.GetData()})
									} else {
										return core.NewData("rotation", ActionRequestItem{RotationClockwise, input.GetData()})
									}
								}
							}

						}
					}
				}
			}
			return core.NewData("error", ActionRequestItem{Action: Move, World: input.GetData()})
		})

		/*
			The robot only have access to a smal part of the world at the beginning
		*/

	// Create a new matrix of the same size as world
	robot_world := make([][]string, len(real_world))
	for i := range robot_world {
		robot_world[i] = make([]string, len(real_world[i]))
	}

	// Populate the new matrix based on conditions
	for i := range real_world {
		for j := range real_world[i] {
			if real_world[i][j] == " " || real_world[i][j] == "#" {
				robot_world[i][j] = ""
			} else {
				robot_world[i][j] = real_world[i][j] // Keep the original value
			}
		}
	}

	go strategistAgent.Serve(ctx)
	go moveAgent.Serve(ctx)

	printMatrix := func(matrix [][]string) {
		for _, row := range matrix {
			for _, cell := range row {
				if cell == "" {
					fmt.Printf("0")
				}
				fmt.Printf("%s ", cell)
			}
			fmt.Println()
		}
	}

	// Start the game
	//gameInitiator <- core.NewData("start", robot_world)
	gameInitiator <- core.NewData("start", real_world)
	// Listen for output in a separate goroutine
	goalFlag := true
	go func() {
		fmt.Println("waiting for processed fruit")
		for data := range gameGoal {
			printMatrix(data.GetData().World)
			goalFlag = true
			for _, row := range data.GetData().World {
				for _, cell := range row {
					if cell == "@" {
						goalFlag = false
					}
				}
			}
			if goalFlag {
				cancel() // Signal to terminate after first output is received
				fmt.Printf("Victory!")
				return
			}

		}
	}()

	// Wait for all processing to be done
	<-ctx.Done()

	if !goalFlag {
		t.Errorf("Expected cherry juice, got %s", "milk")
	}

}
