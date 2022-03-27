package decrypt

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vgo0/gotimepad/db"
	"github.com/vgo0/gotimepad/models"
	h "github.com/vgo0/gotimepad/handler"

)

var Opt = models.OTOptions{}

func Execute(args []string) {
	decryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	Opt.DB = decryptCmd.String("d", "gotimepad.gtp", "One time pad database file")
	Opt.Message = decryptCmd.String("m", "message.txt", "Message to encrypt")
	Opt.B64 = decryptCmd.Bool("b", false, "Message is Base64 encoded")
	Opt.MessageDirect = decryptCmd.String("s", "", "Base64'd message to decrypt via comand line")
	Opt.ToConsole = decryptCmd.Bool("o", false, "Send decrypted message to stdout")
	decryptCmd.Parse(args)

	if !Opt.NoValidateExists {
		db.CheckValidDatabase(*Opt.DB)
	}

	log.Printf("Attempting to decrypt with OTP %s", *Opt.DB)

	db.Connect(*Opt.DB)

	msg, err := getEncyptedBytes()
	h.CheckError(err)

	decrypted, err := decryptMessageText(&msg)
	h.CheckError(err)

	if *Opt.ToConsole {
		log.Printf("%s\n", decrypted)
	} else {
		err := writeDecryptedFile(&decrypted)
		h.CheckError(err)
	}
}

/*
Generic function to get byte array of encrypted data
Can either be passed directly on command line in base64, or in a file (which may be raw or base64 encoded)

Returns encrypted data or error
*/
func getEncyptedBytes() ([]byte, error) {
	var msg []byte
	var err error

	// Bytes provided as base64 string in command args
	if *Opt.MessageDirect != "" {
		msg, err = base64.StdEncoding.DecodeString(*Opt.MessageDirect)

		if err != nil {
			return msg, fmt.Errorf("unable to base64 decode provided string: %s", err)
		}
	} else {
		var raw []byte
		raw, err = os.ReadFile(*Opt.Message)

		if err != nil {
			return msg, fmt.Errorf("unable to read file to decrypt: %s", err)
		}

		if *Opt.B64 {
			dec_len := base64.StdEncoding.DecodedLen(len(raw))

			msg = make([]byte, dec_len)

			_, err = base64.StdEncoding.Decode(msg, raw)

			if err != nil {
				return msg, fmt.Errorf("unable to base64 decode provided string: %s", err)
			}
		} else {
			msg = raw
		}
	}

	return msg, nil
}

/*
Attempts to decrypt provided data against the specific database

Returns decrypted string or error
*/
func decryptMessageText(msg *[]byte) (string, error) {
	var decrypted []byte
	var uuid_m uuid.UUID
	var decb byte

	idx := 0

	for {
		var page models.OTPage

		// The UUID of the 'page' in the database used for encryption is embedded
		// This has a minimum size of 16 bytes per UUID
		if len(*msg) <= idx+16 {
			return "", fmt.Errorf("malformed message, probably not a valid encrypted message")
		}

		// Create UUID from first 16 bytes
		err := uuid_m.UnmarshalBinary((*msg)[idx:idx+16])
		if err != nil {
			return "", fmt.Errorf("invalid UUID encountered: %s", err)
		}

		idx += 16

		// Attempt to find matching page in our database
		res := db.DB.Where("id = ?", uuid_m).First(&page)
		if res.RowsAffected == 0 {
			return "", fmt.Errorf("unable to find matching UUID page for %s", uuid_m.String())
		} else if res.Error != nil {
			return "", fmt.Errorf("error fetching page for UUID %s: %s", uuid_m.String(), res.Error)
		}

		// Provide warning if page is used, pages should never be used more than once!
		if page.Used {
			log.Printf("Warning, page %s may have already been used!\n", page.ID.String())
		}

		// Decrypt byte by byte based on found page data
		for i := 0; i < len(page.Page); i++ {
			decb = (*msg)[i+idx] - page.Page[i]

			decrypted = append(decrypted, decb)
		}

		page.Used = true
		db.DB.Save(&page)

		// idx holds our place in the overall encrypted content
		// idx is incremented 16 (UUID byte size) + page size per page
		idx += len(page.Page)

		if idx >= len(*msg) {
			break
		}
	}

	// Encryption pads end of message with tab to fill page, strip it out here
	return strings.TrimRight(string(decrypted), "\t"), nil
}

/*
Writes decrypted string output to a file

Returns error if issue encountered
*/
func writeDecryptedFile(text *string) error {
	content := []byte(*text)
	file_name := "gtp." + fmt.Sprint(time.Now().Unix()) + ".txt"

	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to create decrypted message output file: %s", err)
	}

	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("error writing decrypted file: %s", err)
	}

	log.Printf("Successfully decrypted message in %s", file_name)
	return nil
}
