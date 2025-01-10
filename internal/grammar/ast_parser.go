package grammar

type BlockType int

const (
	BlockTypeText = iota
	BlockTypeBold
	BlockTypeCode
	BlockTypeBigCode
)

type Block struct {
	Type BlockType
	Str  string
}

func (a *AshMd) ParseAST(originalStr string) []Block {
	blocks := make([]Block, 0)
	for _, token := range a.Tokens() {
		str := string([]rune(originalStr)[token.begin:token.end])
		block := Block{
			Str: str,
		}
		switch token.pegRule {
		case ruleregularText:
			block.Type = BlockTypeText
		case ruleboldText:
			block.Type = BlockTypeBold
		case rulecodeText:
			block.Type = BlockTypeCode
		case rulebigCodeText:
			block.Type = BlockTypeBigCode
		default:
			continue
		}
		blocks = append(blocks, block)
	}
	return blocks
}
