package iotago

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/iotaledger/hive.go/serializer"
)

// ReferenceUnlockBlock is an unlock block which references a previous unlock block.
type ReferenceUnlockBlock struct {
	// The other unlock block this reference unlock block references to.
	Reference uint16 `json:"reference"`
}

func (r *ReferenceUnlockBlock) Deserialize(data []byte, deSeriMode serializer.DeSerializationMode) (int, error) {
	if deSeriMode.HasMode(serializer.DeSeriModePerformValidation) {
		if err := serializer.CheckMinByteLength(ReferenceUnlockBlockSize, len(data)); err != nil {
			return 0, fmt.Errorf("invalid reference unlock block bytes: %w", err)
		}
		if err := serializer.CheckTypeByte(data, UnlockBlockReference); err != nil {
			return 0, fmt.Errorf("unable to deserialize reference unlock block: %w", err)
		}
	}
	data = data[serializer.SmallTypeDenotationByteSize:]
	r.Reference = binary.LittleEndian.Uint16(data)
	return ReferenceUnlockBlockSize, nil
}

func (r *ReferenceUnlockBlock) Serialize(deSeriMode serializer.DeSerializationMode) ([]byte, error) {
	var b [ReferenceUnlockBlockSize]byte
	b[0] = UnlockBlockReference
	binary.LittleEndian.PutUint16(b[serializer.SmallTypeDenotationByteSize:], r.Reference)
	return b[:], nil
}

func (r *ReferenceUnlockBlock) MarshalJSON() ([]byte, error) {
	jReferenceUnlockBlock := &jsonReferenceUnlockBlock{}
	jReferenceUnlockBlock.Type = int(UnlockBlockReference)
	jReferenceUnlockBlock.Reference = int(r.Reference)
	return json.Marshal(jReferenceUnlockBlock)
}

func (r *ReferenceUnlockBlock) UnmarshalJSON(bytes []byte) error {
	jReferenceUnlockBlock := &jsonReferenceUnlockBlock{}
	if err := json.Unmarshal(bytes, jReferenceUnlockBlock); err != nil {
		return err
	}
	seri, err := jReferenceUnlockBlock.ToSerializable()
	if err != nil {
		return err
	}
	*r = *seri.(*ReferenceUnlockBlock)
	return nil
}

// jsonUnlockBlockSelector selects the json unlock block object for the given type.
func jsonUnlockBlockSelector(ty int) (JSONSerializable, error) {
	var obj JSONSerializable
	switch byte(ty) {
	case UnlockBlockSignature:
		obj = &jsonSignatureUnlockBlock{}
	case UnlockBlockReference:
		obj = &jsonReferenceUnlockBlock{}
	default:
		return nil, fmt.Errorf("unable to decode unlock block type from JSON: %w", ErrUnknownUnlockBlockType)
	}
	return obj, nil
}

// jsonReferenceUnlockBlock defines the json representation of a ReferenceUnlockBlock.
type jsonReferenceUnlockBlock struct {
	Type      int `json:"type"`
	Reference int `json:"reference"`
}

func (j *jsonReferenceUnlockBlock) ToSerializable() (serializer.Serializable, error) {
	block := &ReferenceUnlockBlock{Reference: uint16(j.Reference)}
	return block, nil
}