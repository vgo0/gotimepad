package encrypt

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vgo0/gotimepad/db"
	h "github.com/vgo0/gotimepad/handler"
	"github.com/vgo0/gotimepad/models"
)
var Opt = models.OTOptions{}

func Execute(args []string) {
	encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	Opt.DB = encryptCmd.String("d", "gotimepad.gtp", "One time pad database file")
	Opt.Message = encryptCmd.String("m", "message.txt", "Message to encrypt")
	Opt.B64 = encryptCmd.Bool("b", false, "Base64 message content")
	Opt.ToConsole = encryptCmd.Bool("o", false, "Output encrypted message as base64 to stdout")
	Opt.MessageDirect = encryptCmd.String("s", "", "Encrypt message directly from command line")
	encryptCmd.Parse(args)

	if !Opt.NoValidateExists {
		err := db.CheckValidDatabase(*Opt.DB)
		h.CheckError(err)
	}
	
	log.Printf("Attempting to encrypt with OTP %s", *Opt.DB)

	db.Connect(*Opt.DB)

	msg, err := getDataToEncrypt()
	h.CheckError(err)

	encrypted, err := encryptMessageText(&msg)
	h.CheckError(err)

	postProcessText(&encrypted)
	
	if *Opt.ToConsole {
		log.Printf("%s\n", encrypted)
	} else {
		err := writeEncryptedFile(&encrypted)
		h.CheckError(err)
	}
}

/*
Returns byte array of data to encrypt based on options
Can be directly via command line or from a file
*/
func getDataToEncrypt() ([]byte, error) {
	var msg []byte
	if *Opt.MessageDirect == "" {
		var err error
		msg, err = os.ReadFile(*Opt.Message)
		if err != nil {
			return msg, fmt.Errorf("unable to read file to encrypt: %s", err)
		}
	} else {
		msg = []byte(*Opt.MessageDirect)
	}

	return msg, nil
}

/*
Encrypts and provided byte array and returns new byte array
Encryption is done by fetching enough one time pages from the database to encrypt the message

If we run out of pages, or another issue arises, this will return an error
*/
func encryptMessageText(msg *[]byte) ([]byte, error) {
	var encrypted []byte
	var encb byte

	idx := 0

	for {
		var page models.OTPage
		res := db.DB.Where("used = ?", false).First(&page)
		if res.RowsAffected == 0 {
			return encrypted, fmt.Errorf("unable to find an unused page, you might need a new pad")
		} else if res.Error != nil {
			return encrypted, fmt.Errorf("error attempting to retrieve unused page: %s", res.Error)
		}

		// Convert UUID object into bytes and prepend to message
		// This allows us to find the corresponding one time page during the decryption phase
		uuid_bytes, err := page.ID.MarshalBinary()
		if err != nil {
			return encrypted, fmt.Errorf("error embedding page uuid as bytes: %s", err)
		}
		encrypted = append(encrypted, uuid_bytes...)

		for i := 0; i < len(page.Page); i++ {
			if idx >= len(*msg) {
				// We pad and fill excess page space with tab characters
				encb = 0x9
			} else {
				encb = (*msg)[idx]
			}

			encb += page.Page[i]

			encrypted = append(encrypted, encb)
			idx++
		}

		page.Used = true
		db.DB.Save(&page)

		if idx >= len(*msg) {
			break
		}
	}

	return encrypted, nil
}

/*
Attempts to write encrypted data to a file on disk
File is named after the message it is read from (or the default if passed via command line -s)

Returns an error if an issue is encountered
*/
func writeEncryptedFile(content *[]byte) error {
	file_name := *Opt.Message + "." + fmt.Sprint(time.Now().Unix()) + ".gtp"

	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to create encrypted message output file: %s", err)
	}

	defer file.Close()

	if _, err := file.Write(*content); err != nil {
		return fmt.Errorf("error writing encrypted file: %s", err)
	}

	log.Printf("Successfully encrypted message in %s", file_name)
	return nil
}

/*
Used to base64 encrypted text if needed
*/
func postProcessText(encrypted *[]byte) {
	if *Opt.B64 || *Opt.ToConsole {
		var b64_encrypted []byte
		enc_len := base64.StdEncoding.EncodedLen(len(*encrypted))

		b64_encrypted = make([]byte, enc_len)

		base64.StdEncoding.Encode(b64_encrypted, *encrypted)
		*encrypted = b64_encrypted
	}
}