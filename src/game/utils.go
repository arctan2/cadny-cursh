package game

type possibleType []coord

func getHorizPossibles(x, y int) possibleType {
	return []coord{
		{x - 4, y},     // left
		{x + 1, y},     // right
		{x - 3, y + 1}, // bottom-left
		{x, y + 1},     // bottom-right
		{x - 3, y - 1}, // top-left
		{x, y - 1},     // top-right
	}
}

func getVerticPossibles(x, y int) possibleType {
	return []coord{
		{x, y - 4}, // top
		{x, y + 1}, // bottom
		{x - 1, y}, // left
		{x - 1, y - 1},
		{x - 1, y - 2},
		{x + 1, y}, // right
		{x + 1, y - 1},
		{x + 1, y - 2},
	}
}
