package iotago

import (
	"encoding/json"
	"fmt"

	"github.com/iotaledger/hive.go/serializer"
)

// SignatureUnlockBlock holds a signature which unlocks inputs.
type SignatureUnlockBlock struct {
	// The signature of this unlock block.
	Signature Signature `json:"signature"`
}

func (s *SignatureUnlockBlock) Type() UnlockBlockType {
	return UnlockBlockSignature
}

func (s *SignatureUnlockBlock) Deserialize(data []byte, deSeriMode serializer.DeSerializationMode) (int, error) {
	return serializer.NewDeserializer(data).
		CheckTypePrefix(uint32(UnlockBlockSignature), serializer.TypeDenotationByte, func(err error) error {
			return fmt.Errorf("unable to deserialize signature unlock block: %w", err)
		}).
		ReadObject(&s.Signature, deSeriMode, serializer.TypeDenotationByte, SignatureSelector, func(err error) error {
			return fmt.Errorf("unable to deserialize signature within signature unlock block: %w", err)
		}).
		Done()
}

func (s *SignatureUnlockBlock) Serialize(deSeriMode serializer.DeSerializationMode) ([]byte, error) {
	return serializer.NewSerializer().
		WriteNum(UnlockBlockSignature, func(err error) error {
			return fmt.Errorf("unable to serialize signature unlock block type ID: %w", err)
		}).
		WriteObject(s.Signature, deSeriMode, func(err error) error {
			return fmt.Errorf("unable to serialize signature unlock block signature: %w", err)
		}).
		Serialize()
}

func (s *SignatureUnlockBlock) MarshalJSON() ([]byte, error) {
	jSignatureUnlockBlock := &jsonSignatureUnlockBlock{}
	jSignature, err := s.Signature.MarshalJSON()
	if err != nil {
		return nil, err
	}
	rawMsgJsonSig := json.RawMessage(jSignature)
	jSignatureUnlockBlock.Signature = &rawMsgJsonSig
	jSignatureUnlockBlock.Type = int(UnlockBlockSignature)
	return json.Marshal(jSignatureUnlockBlock)
}

func (s *SignatureUnlockBlock) UnmarshalJSON(bytes []byte) error {
	jSignatureUnlockBlock := &jsonSignatureUnlockBlock{}
	if err := json.Unmarshal(bytes, jSignatureUnlockBlock); err != nil {
		return err
	}
	seri, err := jSignatureUnlockBlock.ToSerializable()
	if err != nil {
		return err
	}
	*s = *seri.(*SignatureUnlockBlock)
	return nil
}

// jsonSignatureUnlockBlock defines the json representation of a SignatureUnlockBlock.
type jsonSignatureUnlockBlock struct {
	Type      int              `json:"type"`
	Signature *json.RawMessage `json:"signature"`
}

func (j *jsonSignatureUnlockBlock) ToSerializable() (serializer.Serializable, error) {
	sig, err := signatureFromJSONRawMsg(j.Signature)
	if err != nil {
		return nil, err
	}

	return &SignatureUnlockBlock{Signature: sig}, nil
}
