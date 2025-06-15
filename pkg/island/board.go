package island

type TileType uint8

const (
	TileEmpty TileType = iota
	TileLand
	TileSea
	TileBridge
)

type Tile struct {
	Type TileType
}

type Board struct {
	Width     int
	Height    int
	Tiles     []Tile
	UnionFind *UnionFind
	Islands   []int // Indices of land tiles
}

func NewBoard(width, height int) *Board {
	tiles := make([]Tile, width*height)
	for i := range tiles {
		tiles[i] = Tile{Type: TileSea}
	}
	
	return &Board{
		Width:     width,
		Height:    height,
		Tiles:     tiles,
		UnionFind: NewUnionFind(width * height),
		Islands:   []int{},
	}
}

func (b *Board) GetTile(x, y int) *Tile {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return nil
	}
	return &b.Tiles[y*b.Width+x]
}

func (b *Board) SetTile(x, y int, tileType TileType) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	idx := y*b.Width + x
	b.Tiles[idx].Type = tileType
	
	if tileType == TileLand {
		b.Islands = append(b.Islands, idx)
	}
}

func (b *Board) CanBuildBridge(x, y int) bool {
	tile := b.GetTile(x, y)
	if tile == nil || tile.Type != TileSea {
		return false
	}
	
	// Check if adjacent to land or bridge
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	hasConnection := false
	
	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		neighbor := b.GetTile(nx, ny)
		if neighbor != nil && (neighbor.Type == TileLand || neighbor.Type == TileBridge) {
			hasConnection = true
			break
		}
	}
	
	return hasConnection
}

func (b *Board) BuildBridge(x, y int) {
	if !b.CanBuildBridge(x, y) {
		return
	}
	
	b.SetTile(x, y, TileBridge)
	idx := y*b.Width + x
	
	// Connect with adjacent land/bridges
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		neighbor := b.GetTile(nx, ny)
		if neighbor != nil && (neighbor.Type == TileLand || neighbor.Type == TileBridge) {
			nidx := ny*b.Width + nx
			b.UnionFind.Union(idx, nidx)
		}
	}
}

func (b *Board) IsAllConnected() bool {
	if len(b.Islands) <= 1 {
		return true
	}
	
	// Check if all islands are in the same component
	firstIsland := b.Islands[0]
	for i := 1; i < len(b.Islands); i++ {
		if !b.UnionFind.Connected(firstIsland, b.Islands[i]) {
			return false
		}
	}
	
	return true
}

// SetupLevel1 creates a simple level for MVP
func (b *Board) SetupLevel1() {
	// Clear board
	for i := range b.Tiles {
		b.Tiles[i].Type = TileSea
	}
	b.Islands = []int{}
	
	// Create 3 islands
	// Island 1 (top-left)
	b.SetTile(1, 1, TileLand)
	
	// Island 2 (top-right)
	b.SetTile(3, 1, TileLand)
	
	// Island 3 (bottom-center)
	b.SetTile(2, 3, TileLand)
	
	// Reinitialize UnionFind for the new level
	b.UnionFind = NewUnionFind(b.Width * b.Height)
}