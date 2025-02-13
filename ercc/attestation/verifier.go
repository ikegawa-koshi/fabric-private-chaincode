//go:build WITH_PDO_CRYPTO
// +build WITH_PDO_CRYPTO

/*
Copyright IBM Corp. All Rights Reserved.
Copyright 2020 Intel Corporation

SPDX-License-Identifier: Apache-2.0
*/

package attestation

import "fmt"

// #cgo CFLAGS: -I${SRCDIR}/../../common/crypto
// #cgo LDFLAGS: -L${SRCDIR}/../../common/crypto/_build -L${SRCDIR}/../../common/logging/_build -Wl,--start-group -lupdo-crypto-adapt -lupdo-crypto -Wl,--end-group -lcrypto -lulogging -lstdc++ -lgcov
// #include <stdio.h> /* needed for free */
// #include <stdlib.h>
// #include <string.h>
// #include <stdbool.h>
// #include "attestation-api/evidence/verify-evidence.h"
import "C"

func NewVerifier() VerifierInterface {
	return &VerifierImpl{}
}

type VerifierImpl struct {
}

func (v *VerifierImpl) VerifyEvidence(evidenceBytes []byte, expectedStatementBytes []byte, expectedMrEnclave string) error {
	evidencePtr := C.CBytes(evidenceBytes)
	defer C.free(evidencePtr)
	evidenceLen := len(evidenceBytes)

	expectedStatementPtr := C.CBytes(expectedStatementBytes)
	defer C.free(expectedStatementPtr)
	expectedStatementLen := len(expectedStatementBytes)

	expectedMrEnclaveBytes := []byte(expectedMrEnclave)
	expectedMrEnclavePtr := C.CBytes(expectedMrEnclaveBytes)
	defer C.free(expectedMrEnclavePtr)
	expectedMrEnclaveLen := len(expectedMrEnclaveBytes)

	if !C.verify_evidence(
		(*C.uint8_t)(evidencePtr), C.uint32_t(evidenceLen),
		(*C.uint8_t)(expectedStatementPtr), C.uint32_t(expectedStatementLen),
		(*C.uint8_t)(expectedMrEnclavePtr), C.uint32_t(expectedMrEnclaveLen)) {
		return fmt.Errorf("evidence verification failed")
	}

	return nil
}
