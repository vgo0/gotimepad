package decrypt

import (
	"reflect"
	"testing"

	"github.com/vgo0/gotimepad/db"
	"github.com/vgo0/gotimepad/models"
)

func Test_getEncyptedBytes(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
		options models.OTOptions
	}{
		{
			"Test Direct Base64 Error", 
			[]byte{}, 
			true, 
			models.OTOptions{
				MessageDirect: createString("NOTBASE64!"),
			},
		},
		{
			"Test Direct Base64 Success", 
			[]byte{0x68,0x65,0x6c,0x6c,0x6f,0x20,0x74,0x68,0x69,0x73,0x20,0x69,0x73,0x20,0x61,0x20,0x74,0x65,0x73,0x74}, 
			false, 
			models.OTOptions{
				MessageDirect: createString("aGVsbG8gdGhpcyBpcyBhIHRlc3Q="),
			},
		},
		{
			"Test File Does Not Exist Error", 
			[]byte{}, 
			true, 
			models.OTOptions{
				Message: createString("../does/not/exist"),
				MessageDirect: createString(""),
			},
		},
		{
			"Test File Base64 Error", 
			[]byte{}, 
			true, 
			models.OTOptions{
				Message: createString("../tests/raw_input"),
				B64: createBool(true),
				MessageDirect: createString(""),
			},
		},
		{
			"Test File Base64 Success", 
			[]byte{0x68,0x65,0x6c,0x6c,0x6f,0x20,0x74,0x68,0x69,0x73,0x20,0x69,0x73,0x20,0x61,0x20,0x74,0x65,0x73,0x74, 0x0}, 
			false, 
			models.OTOptions{
				Message: createString("../tests/base64_input"),
				B64: createBool(true),
				MessageDirect: createString(""),
			},
		},
		{
			"Test Raw File Success", 
			[]byte{0x68,0x65,0x6c,0x6c,0x6f,0x20,0x74,0x68,0x69,0x73,0x20,0x69,0x73,0x20,0x61,0x20,0x74,0x65,0x73,0x74}, 
			false, 
			models.OTOptions{
				Message: createString("../tests/raw_input"),
				B64: createBool(false),
				MessageDirect: createString(""),
			},
		},
	}
	for _, tt := range tests {
		var cfg models.OTOptions

		t.Run(tt.name, func(t *testing.T) {
			cfg = Opt
			Opt = tt.options

			got, err := getEncyptedBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("getEncyptedBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getEncyptedBytes() = %v, want %v", got, tt.want)
				}
			}
			

			Opt = cfg
		})
	}
}

func Test_decryptMessageText(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    string
		wantErr bool
		options models.OTOptions
		resetUsed bool
	}{
		{
			"Test Max Length Success",
			[]byte{0x00,0x86,0xba,0x7e,0xe6,0x7d,0x46,0x22,0xa4,0xcf,0xc5,0x93,0x0c,0x4c,0x08,0x3c,0xd7,0x08,0x82,0x4d,0x7d,0xd6,0x5d,0x94,0x56,0xb3,0x3f,0x7e,0xab,0xd9,0x65,0x6e,0xc1,0xcf,0xc3,0xc7},
			"hello this is a test",
			false,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			true,
		},
		{
			"Test Padding Success",
			[]byte{0x00,0x86,0xba,0x7e,0xe6,0x7d,0x46,0x22,0xa4,0xcf,0xc5,0x93,0x0c,0x4c,0x08,0x3c,0xd7,0x08,0x82,0x4d,0x7d,0xd6,0x5d,0x94,0x56,0xb3,0x3f,0x7e,0xab,0xd9,0x65,0x6e,0x56,0x73,0x59,0x5c},
			"hello this is a ",
			false,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			true,
		},
		{
			"Test Multi Page Success",
			[]byte{0x00,0x86,0xba,0x7e,0xe6,0x7d,0x46,0x22,0xa4,0xcf,0xc5,0x93,0x0c,0x4c,0x08,0x3c,0xd7,0x08,0x82,0x4d,0x7d,0xd6,0x5d,0x94,0x56,0xb3,0x3f,0x7e,0xab,0xd9,0x65,0x6e,0xc1,0xcf,0xc3,0xc7,0x00,0x8f,0x3a,0x62,0x01,0xd8,0x44,0x29,0x8c,0x2d,0x54,0x6a,0xc3,0xbc,0x70,0xec,0x10,0x5a,0xf0,0xd1,0xf5,0xc0,0xac,0xde,0xc4,0x2c,0x79,0x79,0x56,0xd0,0x89,0xfc,0x2a,0xd1,0x2b,0x53},
			"hello this is a testhello this is a test",
			false,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			true,
		},
		{
			"Test Missing UUID Error",
			[]byte{0x01,0x86,0xba,0x7e,0xe6,0x7d,0x46,0x22,0xa4,0xcf,0xc5,0x93,0x0c,0x4c,0x08,0x3c,0xd7,0x08,0x82,0x4d,0x7d,0xd6,0x5d,0x94,0x56,0xb3,0x3f,0x7e,0xab,0xd9,0x65,0x6e,0xc1,0xcf,0xc3,0xc7,0x00,0x8f,0x3a,0x62,0x01,0xd8,0x44,0x29,0x8c,0x2d,0x54,0x6a,0xc3,0xbc,0x70,0xec,0x10,0x5a,0xf0,0xd1,0xf5,0xc0,0xac,0xde,0xc4,0x2c,0x79,0x79,0x56,0xd0,0x89,0xfc,0x2a,0xd1,0x2b,0x53},
			"hello this is a testhello this is a test",
			true,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			false,
		},
		{
			"Test Invalid Length Error",
			[]byte{0x00,0x86,0xba,0x7e,0xe6,0x7d,0x46,0x22,0xa4,0xcf,0xc5,0x93,0x0c,0x4c,0x08,0x3c,0xd7,0x08,0x82,0x4d,0x7d,0xd6,0x5d,0x94,0x56,0xb3,0x3f,0x7e,0xab,0xd9,0x65,0x6e,0xc1,0xcf,0xc3,0xc7,0x00,0x8f,0x3a,0x62,0x01,0xd8,0x44,0x29,0x8c,0x2d,0x54,0x6a,0xc3,0xbc,0x70,0xec,0x10,0x5a,0xf0,0xd1,0xf5,0xc0,0xac,0xde,0xc4,0x2c,0x79,0x79,0x56,0xd0,0x89,0xfc,0x2a,0xd1,0x2b,0x53,0x01,0x02},
			"hello this is a testhello this is a test",
			true,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			true,
		},
		{
			"Test Malformed Error",
			[]byte{0x00,0x86,0xba,0x7e,0xe6},
			"hello this is a testhello this is a test",
			true,
			models.OTOptions{
				DB: createString("../tests/test_decrypt.db"),
			},
			true,
		},
	}
	for _, tt := range tests {
		var cfg models.OTOptions

		t.Run(tt.name, func(t *testing.T) {
			cfg = Opt
			Opt = tt.options
			db.Connect(*Opt.DB)

			if tt.resetUsed {
				// Get back to predictable page set for encryption
				db.DB.Exec("UPDATE ot_pages SET used = false")
			}

			got, err := decryptMessageText(&tt.in)

			if (err != nil) != tt.wantErr {
				t.Errorf("encryptMessageText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("encryptMessageText() = %v, want %v", got, tt.want)
				}
			}

			Opt = cfg
		})
	}
}

func createString(in string) *string {
    return &in
}

func createBool(in bool) *bool {
    return &in
}