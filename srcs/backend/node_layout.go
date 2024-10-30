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

type leafNode struct {
	common		nodeHeader
	numCells	uint32
}

const (
	numCellsSize = 4
	numCellsOffset = nodeHeaderSize
	leafNodeSize = numCellsSize + nodeHeaderSize
)

type leafNodeBody struct {
	
}
