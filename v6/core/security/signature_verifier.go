package security

import (
	"crypto/x509"
	"encoding/base64"

	"github.com/rollout/rox-go/v6/core/model"
)

const (
	roxCertificateBase64 = "MIIDWDCCAkACCQDR039HDUMyzTANBgkqhkiG9w0BAQUFADBuMQswCQYDVQQHEwJjYTETMBEGA1UEChMKcm9sbG91dC5pbzERMA8GA1UECxMIc2VjdXJpdHkxFzAVBgNVBAMTDnd3dy5yb2xsb3V0LmlvMR4wHAYJKoZIhvcNAQkBFg9leWFsQHJvbGxvdXQuaW8wHhcNMTQwODE4MDkzNjAyWhcNMjQwODE1MDkzNjAyWjBuMQswCQYDVQQHEwJjYTETMBEGA1UEChMKcm9sbG91dC5pbzERMA8GA1UECxMIc2VjdXJpdHkxFzAVBgNVBAMTDnd3dy5yb2xsb3V0LmlvMR4wHAYJKoZIhvcNAQkBFg9leWFsQHJvbGxvdXQuaW8wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDq8GMRFLyaQVDEdcHlYm7NnGrAqhLP2E/27W21yTQein7r8FOT/7jJ0PLpcGLw/3zDT5wzIJ3OtFy4HWre2hn7wmt+bI+bbS/9kKrmqkpjAj1+PwnB4lhEad27lolMCuz5purqi209k7q51IMdfq0/Ot7P/Bmp+LBNs2F4jMsPYxZUUYkVTAmPqgnwxuWoJZan/OGNjtj9OGg8eOcOfcyxC4GDR/Yail+kht4I/HHesSXVukqXntsbdgnXKFkX682TuFPc3pd8ly+6N6OSWpbNV8UmEVZygnxWT3vxBT2TWvFexbW52KOFY91wIkjt+IPEMPJBPPDiN9J2nuttvfMpAgMBAAEwDQYJKoZIhvcNAQEFBQADggEBAIXrD6YsIhZa6fYDAR8huP0V3BRwMKjeLGLCXLzvuPaoQGDhn4RJNgz3leNcomIkV/AwneeS9BXgBAcEKjNeLD+nW58RSRnAfxDT5cUtQgIeR6dFmEK05u+8j/cK3VO410xr0taNMbmJfEn07WjfCdcJS3hsGJuVmEUC85KYznbIcafQMGklLYArXYVnR3XKqzxcLohSPX99weujH5wt78Zy3pXxuYCDETwhgcCYCQaZz7mpvtSOub3JQT+Ir5cBSdyI1oPI2dIamUL5+ntTyll/1rbYj83qREw8PKA9Q0KIIgfpggy19TS9zknwOLz44wRdLyT2tFoaiRqHvm6JKaA="
)

type SignatureVerifier interface {
	Verify(data, signatureBase64 string) bool
}

type signatureVerifier struct {
	environment model.Environment
}

func NewSignatureVerifier(environment model.Environment) SignatureVerifier {
	return &signatureVerifier{
		environment: environment,
	}
}

func (sv *signatureVerifier) Verify(data, signatureBase64 string) bool {
	if sv.environment.IsSelfManaged() {
		return true
	}

	certificateBytes, err := base64.StdEncoding.DecodeString(roxCertificateBase64)
	if err != nil {
		return false
	}

	cert, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return false
	}

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false
	}

	err = cert.CheckSignature(x509.SHA256WithRSA, []byte(data), signature)
	return err == nil
}

type disabledSignatureVerifier struct{}

func NewDisabledSignatureVerifier() SignatureVerifier {
	return &disabledSignatureVerifier{}
}

func (dsv *disabledSignatureVerifier) Verify(data, signatureBase64 string) bool {
	return true
}
