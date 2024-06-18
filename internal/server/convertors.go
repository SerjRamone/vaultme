package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SerjRamone/vaultme/internal/models"
	pb "github.com/SerjRamone/vaultme/pkg/vaultme_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertItemFromPB(req *pb.Item) (*models.Item, error) {
	d, err := convertItemDataTypeFromPB(req.Data)
	if err != nil {
		return nil, fmt.Errorf("converting item from protobuff error: %w", err)
	}

	b, err := d.Raw()
	if err != nil {
		return nil, fmt.Errorf("encoding to bytes error: %w", err)
	}

	return &models.Item{
		ID:      req.GetId(),
		UserID:  req.GetUserId(),
		Name:    req.GetName(),
		Type:    req.GetType().String(),
		Version: req.Version,
		Data:    b,
		Meta:    convertMetaFromPB(req.GetMeta()),
	}, nil
}

func convertItemToPB(i *models.Item) (*pb.Item, error) {
	pbItem := pb.Item{
		Id:        i.ID,
		UserId:    i.UserID,
		Name:      i.Name,
		Type:      convertItemDataTypeToPB(i.Type),
		Version:   i.Version,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
		Meta:      convertMetaToPB(i.Meta),
	}

	// converting item data type to pb
	switch i.Type {
	case string(models.CredentialType):
		cred := &models.Credential{}
		if err := json.Unmarshal(i.Data, cred); err != nil {
			return nil, fmt.Errorf("encoding credential data error: %w", err)
		}
		c := &pb.Credential{}
		c.Login = cred.Login
		c.Password = cred.Password
		pbItem.Data = &pb.Item_Credential{Credential: c}
	case string(models.TextType):
		text := &models.Text{}
		if err := json.Unmarshal(i.Data, text); err != nil {
			return nil, fmt.Errorf("encoding text data error: %w", err)
		}
		t := &pb.Text{}
		t.Data = text.Data
		pbItem.Data = &pb.Item_Text{Text: t}
	case string(models.CardType):
		card := &models.Card{}
		if err := json.Unmarshal(i.Data, card); err != nil {
			return nil, fmt.Errorf("encoding card data error: %w", err)
		}
		c := &pb.Card{}
		c.Number = card.Number
		c.Owner = card.Owner
		c.ValidityTo = timestamppb.New(card.ValidityTo)
		pbItem.Data = &pb.Item_Card{Card: c}
	case string(models.RawType):
		raw := &models.File{}
		if err := json.Unmarshal(i.Data, raw); err != nil {
			return nil, fmt.Errorf("encoding raw file data error: %w", err)
		}
		r := &pb.Raw{}
		r.Data = raw.Data
		pbItem.Data = &pb.Item_Raw{Raw: r}
	default:
		return nil, errors.New("invalid proto type")
	}

	return &pbItem, nil
}

func convertDataTypeFromPB(dt pb.DataType) models.ItemType {
	switch dt {
	case *pb.DataType_CREDENTIAL.Enum():
		return models.CredentialType
	case *pb.DataType_TEXT.Enum():
		return models.TextType
	case *pb.DataType_CARD.Enum():
		return models.CardType
	case *pb.DataType_RAW.Enum():
		return models.RawType
	default:
		return "invalid type"
	}
}

func convertItemDataTypeToPB(dType string) pb.DataType {
	switch dType {
	case string(models.CredentialType):
		return pb.DataType_CREDENTIAL
	case string(models.TextType):
		return pb.DataType_TEXT
	case string(models.RawType):
		return pb.DataType_RAW
	case string(models.CardType):
		return pb.DataType_CARD
	default:
		return pb.DataType_UNKNOWN
	}
}

func convertItemDataTypeFromPB(iDataType any) (models.ItemDataType, error) {
	switch d := iDataType.(type) {
	case *pb.Item_Credential:
		return &models.Credential{
			Login:    d.Credential.GetLogin(),
			Password: d.Credential.GetPassword(),
		}, nil
	case *pb.Item_Text:
		return &models.Text{
			Data: d.Text.GetData(),
		}, nil
	case *pb.Item_Card:
		return &models.Card{
			Number:     d.Card.GetNumber(),
			Owner:      d.Card.GetOwner(),
			ValidityTo: d.Card.GetValidityTo().AsTime(),
		}, nil
	case *pb.Item_Raw:
		return &models.File{
			Data: d.Raw.GetData(),
		}, nil
	default:
		return nil, errors.New("invalid item type")
	}
}

func convertMetaFromPB(m []*pb.Meta) []*models.Meta {
	metas := make([]*models.Meta, len(m))
	for i := 0; i < len(m); i++ {
		metas[i] = &models.Meta{
			Tag:  m[i].Tag,
			Text: m[i].Text,
		}
	}
	return metas
}

func convertMetaToPB(m []*models.Meta) []*pb.Meta {
	metas := make([]*pb.Meta, len(m))
	for i := 0; i < len(m); i++ {
		metas[i] = &pb.Meta{
			Tag:  m[i].Tag,
			Text: m[i].Text,
		}
	}
	return metas
}
