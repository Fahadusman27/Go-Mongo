package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanAlumni struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NimAlumni     string             `bson:"nim_alumni" json:"nim_alumni"`
	StatusKerja   string             `bson:"status_kerja" json:"status_kerja"`
	JenisIndustri string             `bson:"jenis_industri" json:"jenis_industri"`
	Jabatan       string             `bson:"jabatan" json:"jabatan"`
	Pekerjaan     string             `bson:"pekerjaan" json:"pekerjaan"`
	Gaji          int                `bson:"gaji" json:"gaji"`
	LamaBekerja   int                `bson:"lama_bekerja" json:"lama_bekerja"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Trash struct {
	PekerjaanAlumni `bson:",inline"`
	IsDeleted       time.Time `bson:"is_deleted" json:"is_deleted"`
}
