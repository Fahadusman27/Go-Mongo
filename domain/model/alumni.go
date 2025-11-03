package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	NIM        string             `bson:"nim" json:"nim"`
	Nama       string             `bson:"nama" json:"nama"`
	Angkatan   *int               `bson:"angkatan" json:"angkatan"`
	TahunLulus *int               `bson:"tahun_lulus" json:"tahun_lulus"`
	IDFakultas *int               `bson:"id_fakultas" json:"id_fakultas"`
	IDProdi    *int               `bson:"id_prodi" json:"id_prodi"`
	IDSumber   *int               `bson:"id_sumber" json:"id_sumber"`
	Sumber     *string            `bson:"sumber" json:"sumber"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
