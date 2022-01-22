package urcrypt

// #cgo LDFLAGS: -l urcrypt -l aes_siv
// #include <stdlib.h>
// #include "./urcrypt/urcrypt.h"
import (
	"C"
)
import (
	"fmt"
	"math/big"
	"unsafe"

	"github.com/stevelacy/go-ames/noun"
)

func UrcryptEdShar(public, seed [32]byte) []byte {
	out := C.malloc(32)
	public1 := (*C.uint8_t)(C.CBytes(public[:]))
	seed1 := (*C.uint8_t)(C.CBytes(seed[:]))

	C.urcrypt_ed_shar(public1, seed1, (*C.uint8_t)(out))

	out1 := C.GoBytes(out, 32)
	C.free(unsafe.Pointer(out))
	C.free(unsafe.Pointer(public1))
	C.free(unsafe.Pointer(seed1))
	return out1
}

func UrcryptAESSivcEn(message *big.Int, AESSivData [][]byte, key [64]byte) (error, [16]byte, *big.Int) {
	b := noun.BigToLittle(message)
	message1 := (*C.uint8_t)(C.CBytes(b[:]))
	msgLen := len(b)
	msgLenU := (C.ulong)(msgLen)
	iv := C.malloc(16)

	accum := []C.urcrypt_aes_siv_data{}
	for _, v := range AESSivData {
		item := C.urcrypt_aes_siv_data{
			length: (C.ulong)(len(v)),
			bytes:  (*C.uint8_t)(C.CBytes(v)),
		}
		accum = append(accum, item)
	}
	accumLenU := (C.ulong)(len(accum))

	out := C.malloc(msgLenU)
	key1 := (*C.uint8_t)(C.CBytes(key[:]))

	cerr := C.urcrypt_aes_sivc_en(message1, msgLenU, (*C.urcrypt_aes_siv_data)(&accum[0]), accumLenU, key1, (*C.uint8_t)(iv), (*C.uint8_t)(out))
	if cerr != 0 {
		C.free(unsafe.Pointer(message1))
		C.free(unsafe.Pointer(out))
		C.free(unsafe.Pointer(key1))
		C.free(unsafe.Pointer(iv))
		return fmt.Errorf("urcrypt_aes_sivc_en: Failed to encrypt received error code: %d", cerr), [16]byte{}, big.NewInt(0)
	}

	out1 := C.GoBytes(out, (C.int)(msgLen))
	iv1 := C.GoBytes(iv, 16)
	C.free(unsafe.Pointer(message1))
	C.free(unsafe.Pointer(out))
	C.free(unsafe.Pointer(key1))
	C.free(unsafe.Pointer(iv))

	var iv2 [16]byte
	copy(iv2[:], iv1)

	b2 := noun.LittleToBig(out1)
	return nil, iv2, b2
}

func UrcryptAESSivcDe(message *big.Int, AESSivData [][]byte, key [64]byte, iv [16]byte) (error, *big.Int) {
	b := noun.BigToLittle(message)
	message1 := (*C.uint8_t)(C.CBytes(b[:]))
	msgLen := len(b)
	msgLenU := (C.ulong)(msgLen)
	iv1 := (*C.uint8_t)(C.CBytes(iv[:]))

	accum := []C.urcrypt_aes_siv_data{}
	for _, v := range AESSivData {
		item := C.urcrypt_aes_siv_data{
			length: (C.ulong)(len(v)),
			bytes:  (*C.uint8_t)(C.CBytes(v)),
		}
		accum = append(accum, item)
	}
	accumLenU := (C.ulong)(len(accum))

	out := C.malloc(msgLenU)
	key1 := (*C.uint8_t)(C.CBytes(key[:]))

	cerr := C.urcrypt_aes_sivc_de(message1, msgLenU, (*C.urcrypt_aes_siv_data)(&accum[0]), accumLenU, key1, (*C.uint8_t)(iv1), (*C.uint8_t)(out))
	if cerr != 0 {
		C.free(unsafe.Pointer(message1))
		C.free(unsafe.Pointer(out))
		C.free(unsafe.Pointer(key1))
		C.free(unsafe.Pointer(iv1))
		return fmt.Errorf("urcrypt_aes_sivc_de: Failed to encrypt received error code: %d", cerr), big.NewInt(0)
	}
	out1 := C.GoBytes(out, (C.int)(msgLen))
	C.free(unsafe.Pointer(message1))
	C.free(unsafe.Pointer(out))
	C.free(unsafe.Pointer(key1))
	C.free(unsafe.Pointer(iv1))

	b2 := big.NewInt(0)
	b2.SetBytes(out1)
	return nil, b2
}
