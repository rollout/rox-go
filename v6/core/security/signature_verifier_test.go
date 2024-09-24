package security

import (
	"testing"

	"github.com/rollout/rox-go/v6/core/client"
	"github.com/rollout/rox-go/v6/core/consts"
	"github.com/stretchr/testify/assert"
)

func TestSignatureVerifierShouldVerifyWithROXSignature(t *testing.T) {
	environment := client.NewSaasEnvironment(consts.ROLLOUT_API)
	var sigVerify = NewSignatureVerifier(environment)
	data := `{"__v":0,"application":"5465dd938ede8bfa5c4a40b9","pending_test_devices":false,"_id":"5465e0848ede8bfa5c4a40c8","tweak":{"sandbox_bundle_url":"http://api.rollout.io/device/5465dd938ede8bfa5c4a40b9/5465deb829f2ed0b638982cb/tweaks_bundles/sandbox","bundle":{"app_version":"5465deb829f2ed0b638982cb","bucket":"production","_id":"5465e0848ede8bfa5c4a40c7","__v":0,"creation_date":"2014-11-14T10:59:16.728Z"},"bundles":[]},"structure":{"md5":{"armv7":"5cda09a29184ae745c9ce24cea5f49c8","i386":"e800a916b1c119418b737ea3b7ffcb44"},"force_upload":0,"upload_url":"http://api.rollout.io/device/5465dd938ede8bfa5c4a40b9/5465deb829f2ed0b638982cb/structures"},"devices":{"mode":{"4E922AEE-0616-4135-98A0-1DDF04910087":"sandbox"}},"application_version":{"id":"5465deb829f2ed0b638982cb","release_version_number":"1.0","rollout_api_version":"1.1.0"},"api_version":"1.0.0","creation_date":"2014-11-14T10:59:16.743Z"}`
	signature := "BsJiQvPn0/fH7EABe/mWEwLeldqxQiccQH0SRmjk4p9u76Z+wbmYXym6YqLbwCPYiciHYl7F7HRE0duOMlx4Rz2HMos8mp6DIwVw4cKfzrcBa+abL56PJa6Be9VB99nwjgagesyvSuTl4nd9u/secgHSTP1dh7xJxcFheK1ouXDcHrznvGDTG/LL+fk0FoqovQV2NWjCQFWAqyXkHp5xZ5YveMPjyaHYtHLfPevSidsKK3Pn5Oi7COrw4GWDI8WcvEt/L4hOtsb0nn/hka0VlVmfa6pUPk0aAL5cxQ0kC82YJ7X0xhZ4RqRcUoxaMMr8gM40I5zHyXE7wLe9NWDMwA=="

	assert.True(t, sigVerify.Verify(data, signature))
}

func TestSignatureVerifierShouldNotVerifyWithROXSignature(t *testing.T) {
	environment := client.NewSaasEnvironment(consts.ROLLOUT_API)
	var sigVerify = NewSignatureVerifier(environment)
	data := `{"__v":0,"application":"5465dd938ede8bfa5c4a40b9","pending_test_devices":false,"_id":"5465e0848ede8bfa5c4a40c8","tweak":{"sandbox_bundle_url":"http://api.rollout.io/device/5465dd938ede8bfa5c4a40b9/5465deb829f2ed0b638982cb/tweaks_bundles/sandbox","bundle":{"app_version":"5465deb829f2ed0b638982cb","bucket":"production","_id":"5465e0848ede8bfa5c4a40c7","__v":0,"creation_date":"2014-11-14T10:59:16.728Z"},"bundles":[]},"structure":{"md5":{"armv7":"5cda09a29184ae745c9ce24cea5f49c8","i386":"e800a916b1c119418b737ea3b7ffcb44"},"force_upload":0,"upload_url":"http://api.rollout.io/device/5465dd938ede8bfa5c4a40b9/5465deb829f2ed0b638982cb/structures"},"devices":{"mode":{"4E922AEE-0616-4135-98A0-1DDF04910087":"sandbox"}},"application_version":{"id":"5465deb829f2ed0b638982cb","release_version_number":"1.0","rollout_api_version":"1.1.0"},"api_version":"1.0.0","creation_date":"2014-11-14T10:59:16.743Z"}`
	signature := "AsJiQvPn0/fH7EABe/mWEwLeldqxQiccQH0SRmjk4p9u76Z+wbmYXym6YqLbwCPYiciHYl7F7HRE0duOMlx4Rz2HMos8mp6DIwVw4cKfzrcBa+abL56PJa6Be9VB99nwjgagesyvSuTl4nd9u/secgHSTP1dh7xJxcFheK1ouXDcHrznvGDTG/LL+fk0FoqovQV2NWjCQFWAqyXkHp5xZ5YveMPjyaHYtHLfPevSidsKK3Pn5Oi7COrw4GWDI8WcvEt/L4hOtsb0nn/hka0VlVmfa6pUPk0aAL5cxQ0kC82YJ7X0xhZ4RqRcUoxaMMr8gM40I5zHyXE7wLe9NWDMwA=="

	assert.False(t, sigVerify.Verify(data, signature))
}

func TestSignatureVerifierShouldAcceptAnythingOnSelfManaged(t *testing.T) {
	selfManagedOptions := client.NewSelfManagedOptions(client.SelfManagedOptionsBuilder{
		ServerURL:    "",
		AnalyticsURL: "",
	})
	environment := client.NewSelfManagedEnvironment(selfManagedOptions)
	var sigVerify = NewSignatureVerifier(environment)

	assert.True(t, sigVerify.Verify("", ""))

	data := `randomdata`
	signature := "AsJiQvPn0/fH7EABe/mWEwLeldqxQiccQH0SRmjk4p9u76Z+wbmYXym6YqLbwCPYiciHYl7F7HRE0duOMlx4Rz2HMos8mp6DIwVw4cKfzrcBa+abL56PJa6Be9VB99nwjgagesyvSuTl4nd9u/secgHSTP1dh7xJxcFheK1ouXDcHrznvGDTG/LL+fk0FoqovQV2NWjCQFWAqyXkHp5xZ5YveMPjyaHYtHLfPevSidsKK3Pn5Oi7COrw4GWDI8WcvEt/L4hOtsb0nn/hka0VlVmfa6pUPk0aAL5cxQ0kC82YJ7X0xhZ4RqRcUoxaMMr8gM40I5zHyXE7wLe9NWDMwA=="

	assert.True(t, sigVerify.Verify(data, signature))
}
