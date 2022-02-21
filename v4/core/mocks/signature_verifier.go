package mocks

import "github.com/stretchr/testify/mock"

type SignatureVerifier struct {
	mock.Mock
}

func (m *SignatureVerifier) Verify(data, signatureBase64 string) bool {
	args := m.Called(data, signatureBase64)
	return args.Bool(0)
}
