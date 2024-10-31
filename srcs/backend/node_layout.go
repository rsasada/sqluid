package backend

type nodeHeader struct {
	nodeType 		NodeType //uint8
	isRoot	 		bool
	parentPointer	uint32
}

const (
	nodeTypeSize = 1
	nodeTypeOffset = 0
	isRootSize = 1
	isRootOffset = nodeTypeSize
	parentPointerSize = 4
	parentPointerOffset = isRootSize + isRootOffset
	nodeHeaderSize = nodeTypeSize + isRootSize + parentPointerSize
)


type leafNodeHeader struct {
	common		nodeHeader
	numCells	uint32
}

const (
	numCellsSize = 4
	numCellsOffset = nodeHeaderSize
	leafNodeHeaderSize = numCellsSize + nodeHeaderSize
)


type leafNodeBody struct {

	cell	[]nodeCell
}

type nodeCell struct {
	key		uint32
	value	interface{}	//テーブルのrowsizeに依存する
}

const (
	leafCellKeySize = 4
	leafCellKeyOffset = 0
	leafCellValueOffset = leafCellKeySize + leafCellKeyOffset
	leafNodeSpaceForCells = PageSize - leafNodeHeaderSize
)

func (t *Table) leafCellValueSize() {

	return t.RowSize()
}

func (t *Table) leafCellSize() {

	return t.RowSize + leafCellKeySize
}

func (t *Table) leafNodeMaxCells() {

	return leafNodeSpaceForCells / leafCellSize()
}

func (t *Table)
