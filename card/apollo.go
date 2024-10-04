package card

import (
	"encoding/binary"
	"fmt"

	"github.com/ubavic/bas-celik/document"
)

// Apollo is the type of the first smart ID cards.
// Apollo cards are not manufactured anymore, and this code could be removed in the future.
type Apollo struct {
	atr           Atr
	smartCard     Card
	documentFile  []byte
	personalFile  []byte
	residenceFile []byte
	photoFile     []byte
}

var APOLLO_ATR = Atr([]byte{
	0x3B, 0xB9, 0x18, 0x00, 0x81, 0x31, 0xFE, 0x9E, 0x80,
	0x73, 0xFF, 0x61, 0x40, 0x83, 0x00, 0x00, 0x00, 0xDF,
})

func readApolloCard(card Apollo) (*document.IdDocument, error) {
	doc := document.IdDocument{}

	rsp, err := card.readFile(ID_DOCUMENT_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading document file: %w", err)
	}

	card.documentFile = rsp

	err = parseIdDocumentFile(card.documentFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing document file: %w", err)
	}

	rsp, err = card.readFile(ID_PERSONAL_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading personal file: %w", err)
	}

	card.personalFile = rsp

	err = parseIdPersonalFile(card.personalFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing personal file: %w", err)
	}

	rsp, err = card.readFile(ID_RESIDENCE_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading residence file: %w", err)
	}

	card.residenceFile = rsp

	err = parseIdResidenceFile(card.residenceFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing residence file: %w", err)
	}

	rsp, err = card.readFile(ID_PHOTO_FILE_LOC)
	if err != nil {
		return nil, fmt.Errorf("reading photo file: %w", err)
	}

	card.photoFile = trim4b(rsp)

	err = parseAndAssignIdPhotoFile(card.photoFile, &doc)
	if err != nil {
		return nil, fmt.Errorf("parsing photo file: %w", err)
	}

	return &doc, nil
}

func (card Apollo) readFile(name []byte) ([]byte, error) {
	output := make([]byte, 0)

	_, err := card.selectFile(name, 4)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	data, err := read(card.smartCard, 0, 6)
	if err != nil {
		return nil, fmt.Errorf("reading file header: %w", err)
	}

	if len(data) < 5 {
		return nil, fmt.Errorf("file too short")
	}
	length := uint(binary.LittleEndian.Uint16(data[4:]))
	offset := uint(6)

	for length > 0 {
		data, err := read(card.smartCard, offset, length)
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	return output, nil
}

func (card Apollo) selectFile(name []byte, ne uint) ([]byte, error) {
	apu := buildAPDU(0x00, 0xA4, 0x08, 0x00, name, ne)
	rsp, err := card.smartCard.Transmit(apu)
	if err != nil {
		return nil, fmt.Errorf("selecting file: %w", err)
	}

	return rsp, nil
}

func (card Apollo) Atr() Atr {
	return card.atr
}

func (card Apollo) initCard() error {
	return nil
}
