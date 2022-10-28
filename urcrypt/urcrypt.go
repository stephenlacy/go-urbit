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

	"github.com/stevelacy/go-urbit/noun"
)

func UrcryptEdShar(public, seed [32]byte) []byte {
	out := C.malloc(32)
	public1 := (*C.uint8_t)(C.CBytes(public[:]))
	seed1 := (*C.uint8_t)(C.CBytes(seed[:]))

	defer C.free(unsafe.Pointer(out))
	defer C.free(unsafe.Pointer(public1))
	defer C.free(unsafe.Pointer(seed1))

	C.urcrypt_ed_shar(public1, seed1, (*C.uint8_t)(out))

	out1 := C.GoBytes(out, 32)

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

	defer C.free(unsafe.Pointer(message1))
	defer C.free(unsafe.Pointer(out))
	defer C.free(unsafe.Pointer(key1))
	defer C.free(unsafe.Pointer(iv))

	cerr := C.urcrypt_aes_sivc_en(message1, msgLenU, (*C.urcrypt_aes_siv_data)(&accum[0]), accumLenU, key1, (*C.uint8_t)(iv), (*C.uint8_t)(out))
	if cerr != 0 {
		return fmt.Errorf("urcrypt_aes_sivc_en: Failed to encrypt received error code: %d\n", cerr), [16]byte{}, big.NewInt(0)
	}

	out1 := C.GoBytes(out, (C.int)(msgLen))
	iv1 := C.GoBytes(iv, 16)

	var iv2 [16]byte
	copy(iv2[:], iv1)

	b2 := noun.LittleToBig(out1)

	return nil, iv2, b2
}

func UrcryptAESSivcDe(message *big.Int, AESSivData [][]byte, key [64]byte, iv [16]byte) (*big.Int, error) {
	b := noun.BigToLittle(message)
	message1 := (*C.uint8_t)(C.CBytes(b[:]))
	msgLen := len(b)
	msgLenU := (C.ulong)(msgLen)
	iv1 := (*C.uint8_t)(C.CBytes(iv[:]))

	data := []C.urcrypt_aes_siv_data{}
	for _, v := range AESSivData {
		item := C.urcrypt_aes_siv_data{
			length: (C.ulong)(len(v)),
			bytes:  (*C.uint8_t)(C.CBytes(v)),
		}
		data = append(data, item)
	}
	accumLenU := (C.ulong)(len(data))

	out := C.malloc(msgLenU)
	key1 := (*C.uint8_t)(C.CBytes(key[:]))

	defer C.free(unsafe.Pointer(message1))
	defer C.free(unsafe.Pointer(out))
	defer C.free(unsafe.Pointer(key1))
	defer C.free(unsafe.Pointer(iv1))

	cerr := C.urcrypt_aes_sivc_de(message1, msgLenU, (*C.urcrypt_aes_siv_data)(&data[0]), accumLenU, key1, (*C.uint8_t)(iv1), (*C.uint8_t)(out))
	if cerr != 0 {
		return big.NewInt(0), fmt.Errorf("urcrypt_aes_sivc_de: Failed to decrypt received error code: %d\n", cerr)
	}
	out1 := C.GoBytes(out, (C.int)(msgLen))

	b2 := big.NewInt(0)
	b2.SetBytes(out1)

	b3 := noun.LittleToBig(out1)

	return b3, nil
}
