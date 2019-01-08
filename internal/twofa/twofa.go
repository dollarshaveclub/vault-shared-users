// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twofa //import github.com/dollarshaveclub/vault-shared-users/internal/twofa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"strings"
	"time"
)

const counterLen = 20

//func (c *Keychain) Code(name string) string {
//	k, ok := c.keys[name]
//	if !ok {
//		log.Fatalf("no such key %q", name)
//	}
//	var code int
//	if k.offset != 0 {
//		n, err := strconv.ParseUint(string(c.data[k.offset:k.offset+counterLen]), 10, 64)
//		if err != nil {
//			log.Fatalf("malformed key counter for %q (%q)", name, c.data[k.offset:k.offset+counterLen])
//		}
//		n++
//		code = hotp(k.raw, n, k.digits)
//		f, err := os.OpenFile(c.file, os.O_RDWR, 0600)
//		if err != nil {
//			log.Fatalf("opening keychain: %v", err)
//		}
//		if _, err := f.WriteAt([]byte(fmt.Sprintf("%0*d", counterLen, n)), int64(k.offset)); err != nil {
//			log.Fatalf("updating keychain: %v", err)
//		}
//		if err := f.Close(); err != nil {
//			log.Fatalf("updating keychain: %v", err)
//		}
//	} else {
//		// Time-based key.
//		code = totp(k.raw, time.Now(), k.digits)
//	}
//	return fmt.Sprintf("%0*d", k.digits, code)
//}

func DecodeKey(key string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(key))
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	sum := h.Sum(nil)
	v := binary.BigEndian.Uint32(sum[sum[len(sum)-1]&0x0F:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}

func Totp(key []byte, t time.Time, digits int) int {
	return hotp(key, uint64(t.UnixNano())/30e9, digits)
}
