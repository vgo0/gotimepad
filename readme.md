[![Test](https://github.com/vgo0/gotimepad/actions/workflows/test.yml/badge.svg)](https://github.com/vgo0/gotimepad/actions/workflows/test.yml)

# GoTimePad
This is a simple one time pad implementation on Golang.

The premise is to generate a set of "pages" (random bytes) stored in a SQLite database file (the "pad")

Any two users with the same database could decrypt any message encrypted via the same database

If messages can be exchanged out of sync it is likely preferable to generate two pads. This would avoid accidently encrypting messages with the same page twice. 

An example would be to generate `john_pad.gtp` and `jane_pad.gtp`. All messages sent by `john` would be encrypted via `john_pad.gtp`, and `jane` could decrypt them with their copy of `john_pad.gtp`. All messages sent by `jane` would be encrypted via `jane_pad.gtp`, and `john` could decrypt them with their copy of `jane_pad.gtp`.

# Usage
## General
```
./gotimepad
Expected 'init (i)', 'encode (e)', 'decode (d)' command
```
## Create Pad
`-p` Number of pages controls approximately how many messages can be sent with one pad.

`-s` Page size determines the maximum message length per page.

If an encrypted message exceeds a page length, it will consume as many pages as needed to fully encrypt, or error out if it runs out. There is no particular reasoning behind the default values.

```
./gotimepad i -h
Usage of init:
  -d string
        One time pad database file (default "gotimepad.gtp")
  -f    Force create (overwrite existing file)
  -p uint
        Number of pages (default 500)
  -s uint
        Page size (max number of characters per page) (default 2000)
```

```
./gotimepad i -d john_pad.gtp
```

## Encrypt Message
```
./gotimepad e -h
Usage of encrypt:
  -b    Base64 message content
  -d string
        One time pad database file (default "gotimepad.gtp")
  -m string
        Message to encrypt (default "message.txt")
  -o    Output encrypted message as base64 to stdout
  -s string
        Encrypt message directly from command line
```
```
./gotimepad e -d john_pad.gtp -s 'Secret message goes here!'
> Attempting to encrypt with OTP john_pad.gtp
> Successfully encrypted message in message.txt.1563333511.gtp
```
```
./gotimepad e -d john_pad.gtp -s 'Secret message goes here!' -o
> Attempting to encrypt with OTP john_pad.gtp
> AI6kB+V2TSidP8K9y8zcRXe/64AFVho5UDbcMiVM9zl...
```
```
echo 'Secret message goes here!' > secretmessage.txt
./gotimepad e -m secretmessage.txt -d john_pad.gtp -b
> Attempting to encrypt with OTP john_pad.gtp
> Successfully encrypted message in secretmessage.txt.1563333511.gtp
```

## Decrypt
```
./gotimepad d -h
Usage of encrypt:
  -b    Message is Base64 encoded
  -d string
        One time pad database file (default "gotimepad.gtp")
  -m string
        Message to encrypt (default "message.txt")
  -o    Send decrypted message to stdout
  -s string
        Base64'd message to decrypt via comand line
```
```
./gotimepad d -d john_pad.gtp -s 'AI6kB+V2TSidP8K9y8zcRXe/64A...'
> Attempting to decrypt with OTP john_pad.gtp
> Warning, page 008ea407-e576-4d28-9d3f-c2bdcbccdc45 may have already been used!
> Successfully decrypted message in gtp.1563334511.txt

cat gtp.1563334511.txt
Secret message goes here!
```
```
./gotimepad d -d john_pad.gtp -s 'AI6kB+V2TSidP8K9y8zcRXe/64A...' -o
> Attempting to decrypt with OTP john_pad.gtp
> Warning, page 008ea407-e576-4d28-9d3f-c2bdcbccdc45 may have already been used!
> Secret message goes here!
```
```
./gotimepad d -d john_pad.gtp -m secretmessage.txt.1563333511.gtp -o
> Attempting to decrypt with OTP john_pad.gtp
> Warning, page 008ea407-e576-4d28-9d3f-c2bdcbccdc45 may have already been used!
> Secret message goes here!
```
# Encrypted Message Structure
The raw encrypted message will be (16 + page_size) * number_of_pages bytes. Base64 encoded length will vary.

A 1000 character message with a page size of 2000 characters will be 2016 bytes and 1 page.
A 2001 character message with a page size of 2000 cahracters will be 4032 bytes and 2 pages.

The message is structured as [UUID][EncodedMessage], so multi-page messages become [UUID][EncodedMessage][UUID][EncodedMessage]. The UUID is 16 bytes and is included to find the applicable "page" of data within your copy of the database.

# Test
```
go test -v ./..
```

# Build
```
go build
```

# Install
```
go install
```
