package hashid

type HashData struct {
	ID      int
	Hash    string
	options Option
}

func (h *HashData) IsUser() bool {
	return h.hasUserBit()
}

func (h *HashData) IsMale() bool {
	return h.hasMaleBit()
}

func (h *HashData) hasUserBit() bool {
	i := int(h.options)

	return i-((i>>1)<<1) == 1
}

func (h *HashData) hasCommentatorBit() bool {
	i := int(h.options)

	return (i>>1)-((i>>2)<<1) == 1
}

func (h *HashData) hasMaleBit() bool {
	i := int(h.options)

	return (i>>2)-((i>>3)<<1) == 1
}

func (h *HashData) hasFemaleBit() bool {
	i := int(h.options)

	return (i>>3)-((i>>4)<<1) == 1
}

func (h *HashData) validOptions() bool {
	return ((h.hasUserBit() && !h.hasCommentatorBit()) || (!h.hasUserBit() && h.hasCommentatorBit())) &&
		((h.hasMaleBit() && !h.hasFemaleBit()) || (!h.hasMaleBit() && h.hasFemaleBit()))
}
