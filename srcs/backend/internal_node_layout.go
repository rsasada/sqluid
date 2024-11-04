package backend

import (

)

type internalNodeHeader struct {
	common			nodeHeader
	nodeNumKey		uint32
	rightChild		uint32
}

const (

	internalNodeNumKeysSize		= 4
	internalNodeNumKeysOffset	= nodeHeaderSize
	internalNodeRightChildSize	= 4
	internalNodeRightChildOffset = internalNodeNumKeysSize + internalNodeNumKeysOffset
	internalNodeHeaderSize		= internalNodeNumKeysSize + internalNodeRightChildSize + nodeHeaderSize
)

type internalNodeBody struct {
	cells []internalCell
}

type internalCell struct {
	key		uint32
	child 	uint32
}

const (
	internalNodeKeySize = 4
	internalNodeChildSize = 4
	internalCellSize = internalNodeKeySize + internalNodeChildSize
)
