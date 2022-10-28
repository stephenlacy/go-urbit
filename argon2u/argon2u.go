package argon2u

// #cgo LDFLAGS: -L${SRCDIR}/argon2 -l argon2
// #include <stdlib.h>
// #include "./argon2/include/argon2.h"
import (
	"C"
)

import (
	"encoding/hex"
	"fmt"
	"unsafe"
)

const Argon2u = 10
const ArgonVersion = 0x13

type HashOptions struct {
	Pass        []byte
	Salt        []byte
	Type        int
	HashLen     uint64
	Parallelism int
	Mem         int
	Time        int
}
type ArgonResponse struct {
	Hash    []byte
	Encoded string
	HashHex string
}

func Hash(opts HashOptions) (ArgonResponse, error) {
	t_cost := (C.uint32_t)(opts.Time)
	m_cost := (C.uint32_t)(opts.Mem)
	parallelism := (C.uint32_t)(opts.Parallelism)
	pwd := C.CBytes(opts.Pass[:])
	pwdlen := (C.uint64_t)(C.uint64_t(len(opts.Pass)))
	salt := C.CBytes(opts.Salt[:])
	saltlen := (C.uint32_t)(len(opts.Salt))
	hash := C.malloc(10240)
	hashlen := (C.uint32_t)(C.uint64_t(opts.HashLen))
	encoded := C.CString("")
	encodedlen := (C.uint32_t)(10240)
	argon2Type := (C.argon2_type)(Argon2u)
	version := (C.uint32_t)(ArgonVersion)

	defer C.free(unsafe.Pointer(pwd))
	defer C.free(salt)
	defer C.free(hash)
	defer C.free(unsafe.Pointer(encoded))

	res := C.argon2_hash(
		t_cost,
		m_cost,
		parallelism,
		pwd,
		C.ulong(pwdlen),
		salt,
		C.ulong(saltlen),
		hash,
		C.ulong(hashlen),
		encoded,
		C.ulong(encodedlen),
		argon2Type,
		version,
	)

	if res != 0 {
		return ArgonResponse{}, fmt.Errorf("argon2 error code: %d\n", res)
	}

	enc := C.GoBytes(unsafe.Pointer(encoded), (C.int)(encodedlen))
	out := C.GoBytes(hash, (C.int)(hashlen))
	response := ArgonResponse{
		Hash:    out,
		Encoded: string(enc),
		HashHex: hex.EncodeToString(out),
	}

	return response, nil
}
