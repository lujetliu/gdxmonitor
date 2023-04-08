package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// BloomBitLength represents the number of bits used in a header log bloom.
	BloomBitLength = 8 * BloomByteLength
)

type Bloom [BloomByteLength]byte

type BlockNonce [8]byte

type Block struct {
	gorm.Model
	Number      string
	ParentHash  string
	UncleHash   string
	Coinbase    string
	Root        string
	TxHash      string
	ReceiptHash string
	GasLimit    uint64
	GasUsed     uint64
	Time        uint64
}

type Page struct {
	PageNum  uint `form:"page_num" json:"page_num"`
	PageSize uint `form:"page_size" json:"page_size"`
	Offset   uint `json:"-"`
}

func (p *Page) Default() {
	if p.PageNum == 0 {
		p.PageNum = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}
	p.Offset = (p.PageNum - 1) * p.PageSize
}

type Where map[string]interface{}

func (b *Block) Create(db *gorm.DB) error {
	return db.Create(b).Error
}

func (b *Block) Get(db *gorm.DB) error {
	return db.Where("id = ? ", b.ID).First(b).Error
}

func (e *Block) Update(db *gorm.DB) error {
	return db.Save(e).Error
}

func (e *Block) Deleted(db *gorm.DB) error {
	return db.Delete(e).Error
}

func (e *Block) Query(db *gorm.DB, whs Where, page Page) ([]Block, int64, error) {
	var total int64
	var es []Block
	query := db.Model(&Block{})
	for field, value := range whs {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}
	query.Count(&total)

	res := query.Offset(int(page.Offset)).Limit(int(page.PageSize)).
		Order("created_at DESC").Find(&es)
	if res.Error != nil {
		return es, total, res.Error
	}
	return es, total, nil
}
