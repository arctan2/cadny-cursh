package game

type possibleType []coord

func getHorizPossibles(x, y int) possibleType {
	return []coord{
		{x - 3, y},     // left
		{x + 2, y},     // right
		{x - 2, y + 1}, // bottom-left
		{x + 1, y + 1}, // bottom-right
		{x - 2, y - 1}, // top-left
		{x + 1, y - 1}, // top-right
	}
}

func getVerticPossibles(x, y int) possibleType {
	return []coord{
		{x, y - 3}, // top
		{x, y + 2}, // bottom
		{x - 1, y - 2},
		{x + 1, y - 2},
		{x - 1, y + 1},
		{x + 1, y + 1},
	}
}

func getMidPossibles(x, y int) []possibleType {
	return []possibleType{
		{
			{x - 1, y + 1}, // bottom-left
			{x + 1, y + 1}, // bottom-right
		},
		{
			{x - 1, y - 1}, // top-left
			{x + 1, y - 1}, // top-right
		},
		{
			{x - 1, y - 1}, // top-left
			{x - 1, y + 1}, // bottom-left
		},
		{
			{x + 1, y - 1}, // top-right
			{x + 1, y + 1}, // bottom-right
		},
	}
}
