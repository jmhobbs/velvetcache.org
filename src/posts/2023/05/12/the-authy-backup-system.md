---
tags:
- Security
- neovim
date: 2023-05-12
title: The Authy Backup System
type: post
summary: >
  A look at how the Authy authenticator app works, in the pursuit of securely backing up my MFA tokens offline.
---

I've been an [Authy](https://authy.com/) user for a long time, about a decade as of writing this post. It's always been a nice alternative to the default of the time, [Google Authenticator](https://googleauthenticator.net/).  The backup feature was the best, as it meant a phone dying wouldn't cause a world of suffering.

Over time it has bugged me more and more that I had no easy way out of Authy.  Yes, there were backups, but I did not control them.  I couldn't port out my TOTP codes to another service without detaching and reattaching MFA on basically every app.  I didn't save my setup QR codes for all of these, who does that?

With a new device to set up, I finally got the itch to actually figure it out, or try to.

# Backups!

Authy has a desktop app, which is an Electron app.  As such, it's pretty easy to get into and mess around with. I found this excellent gist, which works great [gboudreau/94bb0c1](https://gist.github.com/gboudreau/94bb0c11a6209c82418d01a59d958c93)

The problem here is that it's a very manual process, and I'd like to automate it if I can.

# Exploring The API

Poking around in the code, it was clear that it was not using any exotic technologies.  This was a pretty standard Electron app talking to an API over HTTP and doing fairly basic, boring crypto.  Which is good! I want my security tooling to be boring, boring is usually simpler to understand and safer.

With that knowledge in mind, I started up [Charles](https://charlesproxy.com), added `*.authy.com` to the SSL proxy list and started a fresh copy of the Authy desktop app.  Luckily for us, there is no TLS cert pinning or verification going on in the desktop app, so we can slip our trusted proxy cert in with no complications.

# Setup & Registration

The first bit of the flow is device registration, which breaks down into a handful of sub-steps.

Unless specified, all of the following API traffic went to the host `https://api.authy.com`, so I will leave that part out going forward.

## Lookup Account By Phone

`/json/users/1-5558675309/status?locale=en-US`

After entering your phone number, the app makes a `GET` to this URL, swapping in your country code and phone number. 

There are some interesting HTTP headers on this request and those that follow, including a `x-authy-api-key`.  I will [cover these headers later on](#authy-auth), for now we can ignore them.

The response is fairly short and sweet;

```json
{
	"force_ott": false,
	"primary_email_verified": false,
	"message": "active",
	"devices_count": 1,
	"authy_id": 123456,
	"success": true
}
```

Interestingly, this _should_ allow you to enumerate users by phone number.  There doesn't seem to be anything here that would prevent that, although I would hope there is some aggressive rate limiting on this endpoint.

## Start Registration

`/json/users/123456/devices/registration/start?locale=en-US`

The next step is to start registration of this device.  This is a form `POST`, where `123456` is the `authy_id` value we got in the previous step.  This is clearly the user identifier.

This form has five fields:

```
via	push
signature	223f1a5b76ea771e79de335c31cb82ec5fa396848be69a07de3e05ecd13856a1
device_app	authy
device_name	Authy Desktop on cornuta.home
api_key	37b312a3d682b823c439522e1fd31c82
```

Two of those fields are interesting, `api_key` is the same value we saw in the `x-authy-api-key` header on our first request, and indeed it is still present as a header in this request.  Not sure why they are including it twice, but there it is.

The `signature` field is a little more opaque.  It looks like 256 bits, hex encoded.  Given the name, my gut instinct was that this was the output of a hash like a SHA-256. After trying to sort it out, I hopped into the code to see if there were any clues.

It turns out, this value is 256 bits of random, generated in `RegistrationController` and is used as a sort of identifier for the device registration request polling.

```javascript
// This code is found in `js/app.js` from the `app.asar` in the macOS version 2.2.3

// Here's the URL being requested. The second parameter, `t`, becomes `signature`
// Line 25401
createNewDeviceRequest(e, t, n, r, o) {
  return i.post(
    `/json/users/${e}/devices/registration/start`,
    {
      via: n,
      signature: t,
      api_key: s.API_KEY,
      device_app: s.getFlavor(),
      device_name: a.get().getDeviceName(),
    },
    function (e) {
      return r(
        e.message,
        e.request_id,
        e.approval_pin,
        e.provider
      );
    },
    o
  );
}

// Here's the invocation, it's using `this.signature` for that parameter
// Line 15316
this.regApi.createNewDeviceRequest(
  this.userId,
  this.signature,
  e,
  function (e, t, n, r) {
    return o(e, t, n, r);
  },
  t
)

// Here is where `this.signature` is set, using `generateSalt`
// Line 15237
setUserStatusData(e, t, n) {
  return (
    f(this, RegistrationController),
    (this.userId = e),
    (this.countryCode = t),
    (this.cellphone = n),
    (this.signature = h.generateSalt())
  );
}

// Finally, here is `generateSalt`, which gets 256 random bits
// Line 20677
static generateSalt(e = 256) {
  e = s.random.getBytesSync(e / 8);
  return s.util.createBuffer(e).toHex();
}
```

This request has a short response as well, nothing outstanding.

```json
{
	"message": "A request was sent to your other devices.",
	"request_id": "642346129bcdb264004ee265",
	"approval_pin": 18,
	"provider": "push",
	"success": true
}
```

## Poll For Status

`/json/users/123456/devices/registration/642346129bcdb264004ee265/status`

At this point, a push request has been sent out to existing devices with the Authy app installed to prompt them to approve this, new device.

In the meantime, the app begins polling this endpoint for the registration request on about a two second interval, waiting for it to be approved.  It's a `GET` request but I've lopped off the query parameters because they are quite long:

```
locale=en-US
signature=223f1a5b76ea771e79de335c31cb82ec5fa396848be69a07de3e05ecd13856a1
api_key=37b312a3d682b823c439522e1fd31c82
```

Nothing new in there, just the same API key (still in the headers too) and the random value from the previous step. The response is simple:

```json
{
	"message": {
		"request_status": "Request Status."
	},
	"status": "pending",
	"success": true
}
```

Once you approve on another device, you'll get this response instead

```json
{
	"message": {
		"request_status": "Request Status."
	},
	"status": "accepted",
	"pin": "543210",
	"success": true
}
```

## <a name="registration-complete"></a>Registration Complete

`/json/users/123456/devices/registration/complete`

The app now marks the registration as complete with a form `POST`.  The form fields include the `pin` value from the registration request status result, as well as a UUID.

When inspecting the Authy desktop app live, this UUID is stored in as `browser.unique.ids`.  It appears to be generated by the app to identify our device specifically.

```
pin	543210
uuid	9e05524bb01cd30262211073cd313bab
device_app	authy
device_name	Authy Desktop on cornuta.home
api_key	37b312a3d682b823c439522e1fd31c82
```

In response we get an ID for our device from Authy, as well as a couple more interesting bits.

```json
{
	"device": {
		"id": 804509754,
		"secret_seed": "5ae00da039e2d38c7659d68a94545077",
		"api_key": "665b1e3851eefefa3fb878654292f165",
		"reinstall": false
	},
	"authy_id": 123456
}
```

The `api_key` value feels like a device specific API token, based on naming and the fact it's the same size and character set (hex) as the value that we've been passing around in `x-authy-api-key`.

The `secret_seed` value is 32 bytes of hex, without much context.  The name makes it appear to be TOTP related, based on values we see in the backup script in the gist. `secretSeed` here decodes to become the TOTP secret value.

```javascript
var secretSeed = i.secretSeed;
if (typeof secretSeed == 'undefined') {
   secretSeed = i.encryptedSeed;
}
var secret = (i.markedForDeletion === false ? i.decryptedSeed : hex_to_b32(secretSeed));
```

In a few requests from now we will see it's purpose.

At this point, the app is registered as a device, we're most of the way in!

# <a name="the-device-token"></a> The Device Token

`/json/devices/804509754/soft_tokens/804509754/check`

There are a handful of asset related calls in between the previous request and this one.  I'm leaving those out as there isn't anything of interest in them, just references to icons for use in the UI.

The next meaningful request is this `GET`, again with some lengthy parameters:

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
sha=aadb416b48116a232774b412abdc822bd2d63009cca973d11a4ddff1f5702632
```

Oddly, we're still sending the same `api_key` as before, rather than the value we got under `device.api_key` when completing registration.

<aside>
  <p>
    In the context of Authy, there are <a href="https://support.authy.com/hc/en-us/articles/4407123462299-What-is-a-Token-">two types of token</a>, "Authy" and "authenticator".  "Authenticator tokens" are a traditonal TOTP token with no special Authy integration. The "Authy" tokens are configured a bit differently from regular TOTP tokens, and are unique on each registered device.  If you open Authy on your phone, and Authy on your computer, and go the same account, you will get different tokens for the same time.  This allows Authy to tie MFA access to a specific device.
  </p>
</aside>

Confusingly the term "soft token" is also used in the API, which I have to assume is a "software token", as opposed to a hardware token like using a Yubikey.  Generally, "soft token" is referring to an "authenticator token", while "apps" are "Authy tokens".

The structure of the path and presence of the `soft_token` in the URL here hints that we're dealing with a TOTP token which is specific to this device.

With an educated guess we find that it turns out that the `sha` value here is the SHA-256 of the `secret_seed` from the registration complete step.

The response is minimal, nothing to discuss here.

```json
{
	"message": "Token is correct.",
	"success": true
}
```

## <a name="get-device-info"></a>Get Device/Account Info

`/json/users/123456/devices/804509754`

From here the app fetches some account information, I assume to put on the settings page. Interestingly, when it requests to  it includes the following query parameters:

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

The parameters `otp1`, `otp2` and `otp3` are from the device soft token for the current period, and the next two periods.  In this case, the request was at `Tue Mar 28 2023 19:55:12 GMT+0000`, which is `1680033312`.

Authy tokens are 7 characters, with a 10 second window (see [the gist](https://gist.github.com/gboudreau/94bb0c11a6209c82418d01a59d958c93)), so we can check our assumption with a little Go.

```golang
package main

import (
  "encoding/base32"
  "encoding/hex"
  "fmt"
  "log"
  "strings"

  "github.com/xlzd/gotp"
)

func main() {
  secretSeed := "5ae00da039e2d38c7659d68a94545077"

  dec, err := hex.DecodeString(secretSeed)
  if err != nil {
    panic(err)
  }

  // otpauth URI's generally use base32 encoded secrets, sans padding
  // https://github.com/google/google-authenticator/wiki/Key-Uri-Format#secret

  otpSeed := strings.TrimRight(base32.StdEncoding.EncodeToString(dec), "=")

  totp := gotp.NewTOTP(otpSeed, 7, 10, nil)
  fmt.Println(totp.At(1680033312))
  fmt.Println(totp.At(1680033322))
  fmt.Println(totp.At(1680033332))
}
```

When run, we get the same values as the parameters:

```console
$ go run .
9735470
8142428
6063286
```

My assumption here is that Authy uses that device token as part of request authentication, since the `api_key` on this is still the static key from before we ever connected anything. Using the TOTP does mean a replay attack would have to happen inside that 10 second window.

The call returns the account details.

```json
{
	"email": "john@velvetcache.org",
	"cellphone": "555-867-5309",
	"country_code": 1,
	"multidevice_enabled": true,
	"multidevices_enabled": true,
	"primary_email_verified": false,
	"success": true
}
```

## Devices List

`/json/users/123456/devices`

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

This yields a list of all connected devices, in great detail.  Much, much more detail than would be needed in this UI.  I've trimmed the response JSON down to just the Authy desktop device we are currently using, but it is representative of the other entries, and in some cases this is a **smaller** listing.

```json
{
	"message": "Devices List",
	"devices": [{
		"_id": "6423461fa9a21167d008d8a8",
		"_migrations": {},
		"access_token": null,
		"account_id": null,
		"api_key": "665b1e3851eefefa3fb878654292f165",
		"authy_ver": null,
		"aws_arn": null,
		"city": null,
		"country": null,
		"created_at": "2023-03-28T19:55:11Z",
		"device_app": "authy",
		"device_type": "desktop",
		"enabled_unlock_methods": [],
		"ip": null,
		"key_rotation_nonce": null,
		"keys_rotated_at": "2023-03-28T19:55:11Z",
		"last_otp_used": null,
		"last_otp_used_at": null,
		"last_sync_at": null,
		"last_unlock_date": null,
		"last_unlock_method_used": "unknown",
		"master_token_id": 804509754,
		"name": "Authy Desktop on cornuta.home",
		"needs_health_check": true,
		"needs_key_rotation": false,
		"push_token": null,
		"rate_limiters": {},
		"region": null,
		"registered": true,
		"registration_city": "Omaha",
		"registration_country": "United States",
		"registration_device_id": 804676179,
		"registration_ip": "127.0.0.1",
		"registration_method": "push",
		"registration_region": "Nebraska",
		"reinstall": false,
		"remote_user_id": null,
		"risky_device": false,
		"role": "authy",
		"sdk_app_id": null,
		"sdk_only": false,
		"soft_token_ids": [
			"6423461fa9a21167d008d8aa",
			"6423461fa9a21167d008d8ab",
			"6423461fa9a21167d008d8ac",
			"6423461fa9a21167d008d8ad",
			"6423461fa9a21167d008d8ae",
			"6423461fa9a21167d008d8a9"
		],
		"ssl_cypher": null,
		"ssl_protocol": null,
		"twilio_notify_registered": false,
		"twilio_notify_sid": null,
		"updated_at": "2023-03-28T19:55:11Z",
		"used_times": 0,
		"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) AuthyDesktop/2.2.3 Chrome/96.0.4664.110 Electron/16.0.8 Safari/537.36",
		"user_id": "52a057d0f92ea1708b003e8e",
		"uuid": "9e05524bb01cd30262211073cd313bab"
	}],
  "success": true
}
```

I mean, come on, that's a lot.

The `keys_rotated_at` indicates that _some_ set of keys on this could be rotated, so that's promising. The `soft_token_ids` field seems to indicate these six tokens are important enough to make this listing, and are tied directly to this device.  Unfortunately, those ID's don't show up in any other API request I captured.

Most of the rest of this is fluff as far as we are concerned.

## <a name="auth-sync"></a>Auth Sync

`/json/devices/804509754/auth_sync`

The next API call we see is one of the first things the Authy desktop app does when starting up, every time.  Just based on the URL, it seems this is a check in to ensure that the device is in good standing and has not been removed from the account by another device.

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
ga_ver=
authy_ver=2.2.3
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

The response from this request is fairly opaque, we don't have much context on what these fields mean yet.

```json
{
	"moving_factor": "168003331",
	"needs_health_check": true,
	"sync_ga": true,
	"authy_token": {
		"hidden": true,
		"valid": true
	},
	"update": false,
	"sync_password": false,
	"enroll_backup_key": false,
	"success": true
}
```

One oddball we can guess at here is `moving_factor`.  This looks a lot like a unix timestamp, and it appears to be a truncated one. If we add a `0` tog the end we get `1680033310`.  This is `Tue Mar 28 2023 19:55:10 GMT+0000`.  The request itself ran at `Tue Mar 28 2023 19:55:15 GMT+0000`.  That said, I'm still not sure _why_ it is there.

The `needs_health_check` field only appeared on the initial app setup, not on previous opens.

## <a name="get-the-tokens"></a> Get The Tokens

`/json/users/123456/authenticator_tokens`

Jackpot! This is the one we've been waiting for.  The app makes a `GET` and all the authenticator tokens come back.

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
apps=
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

This is a truncated response, because you don't want to scroll though thirty entries, and I don't want to have to redact them.

```json
{
	"message": "success",
	"authenticator_tokens": [{
		"account_type": "authenticator",
		"digits": 6,
		"encrypted_seed": "b4cu5FJSJqv/YfDs2nVvPNeaRG1PyvPbNyZLMau48G4=",
		"issuer": null,
		"key_derivation_iterations": null,
		"logo": null,
		"name": "www.velvetcache.org:VelvetCache Blog",
		"original_name": "www.velvetcache.org:admin",
		"password_timestamp": 1519340474,
		"salt": "MIIm9CTN13yqHfxjr4OSQxEDJlyDzVEm",
		"unique_id": "1518713548"
	}],
	"deleted": [],
	"success": true
}
```

We will come back to this _very_ soon, we just have to make it through a few more API calls.  If you can't take the suspense, you can [jump ahead to decrypting](#decrypting)

On subsequent calls to this endpoint, the desktop app provides a list of `apps` in the parameters as a comma separated list.  These are the `unique_id` values from the previous listing.  Again, I've truncated here, it's a lot.

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
apps=1477942954%2C1518643238%2C1518713548%2C1518729515
otp1=8435018
otp2=2732168
otp3=4933747
device_id=804509754
```

Given this list, the API can return just what has changed rather than the full set.

```
{
	"message": "success",
	"authenticator_tokens": [],
	"deleted": [],
	"success": true
}
```

## <a name="get-the-rsa-key"></a> Get The...RSA key?

`/json/devices/804509754/rsa_key`

This is one I can't explain.  We are sent a 2048 bit RSA private key which, based on the URL, is tied to our device.

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

This RSA key does not have a passphrase.

```
{
	"message": "success",
	"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIxiLq2vqzVE....\n-----END RSA PRIVATE KEY-----\n"
	"success": true
}
```

I'm not sure of the purpose of this key yet, I've not seen any evidence of it's use.

## Device Apps Sync

`/json/users/123456/devices/804509754/apps/sync`

Next up is a `POST` to sync the device specific apps.  These are Authy tokens, not authenticator tokens.  Why a `POST` is used here is not clear, we don't send any new or particularly meaningful information to the server.

```
api_key	37b312a3d682b823c439522e1fd31c82
locale	en-US
last_unlock_method_used	none
last_unlock_date	0
enabled_unlock_methods[]	none
otp1	9735470
otp2	8142428
otp3	6063286
device_id	804509754
```

This brings back a list of our Authy apps.

```json
{
	"message": "App Sync.",
	"apps": [{
		"_id": "5084abd1f91ea1ab41000023a",
		"name": "Cloudflare",
		"serial_id": 43,
		"version": 12,
		"assets_group": "5084abd1f92ea1ab4100002d",
		"background_color": null,
		"circle_background": null,
		"circle_color": null,
		"custom_assets": true,
		"generating_assets": false,
		"labels_color": null,
		"labels_shadow_color": null,
		"timer_color": "#FF8600",
		"token_color": null,
		"authy_id": 804457900,
		"secret_seed": "cdd66a78b1722cb2a3b657fc21f2edd9"
		"digits": 7,
		"member_since": 1564610102,
		"transactional_otp": false
	}],
	"deleted": [],
	"success": true
}
```

The obvious bit here is the `secret_seed`, which we can assume is our TOTP token value.  Keeping in mind that every Authy token is different on every device, this should be a fresh token issued for this device.  As such, there is no decryption to be dealt with for these tokens, they're sent in clear text, hex encoded.

```go
package main

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	dec, err := hex.DecodeString("cdd66a78b1722cb2a3b657fc21f2edd9")
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.TrimRight(base32.StdEncoding.EncodeToString(dec), "="))
}
```

```console
$ go run main.go
ZXLGU6FROIWLFI5WK76CD4XN3E
```

Much like [`authenticator_tokens`](#get-the-tokens), this request is called at every startup of the app.  The method of tracking state is slightly different here though, as instead of an `apps` parameter with a comma separated list of ID's, we get individual form fields for each app, like this:

```
api_key	37b312a3d682b823c439522e1fd31c82
locale	en-US
vs5084abd1f91ea1ab4100003a	12
vs510cd201f91ea171d600027d	5
vs52cdc0e79d19c905e9005a1d	32
vs567b20176110703121000152	23
last_unlock_method_used	none
last_unlock_date	0
enabled_unlock_methods[]	none
otp1	8435018
otp2	2732168
otp3	4933747
device_id	804509754
```

The `vs<hex>` fields here are clearly tied to the apps themselves, with the `<hex>` part being the `_id` from the listing, and the numeric value being the `version` field.

```json
{
	"message": "App Sync.",
	"apps": [],
	"deleted": [],
	"success": true
}
```

## RSA Key Check

`/json/devices/804509754/rsa_key_check`

The app then checks on the status of the RSA private key it fetched before.  Again, `sha` here is a SHA-256 sum of the key like we saw with [the device token check](#the-device-token)

```
api_key=37b312a3d682b823c439522e1fd31c82
locale=en-US
sha=4970f487129ab4967ec57287d96c34ffc3542b4c529e07c5fcbab3d2f2ce1914
otp1=9735470
otp2=8142428
otp3=6063286
device_id=804509754
```

The response is minimal.

```
{
	"message": "RSA key is valid.",
	"success": true
}
```

## Create App Keys

`/json/devices/804509754/public_key`

This is an interesting bit.  The app makes a `POST` for each of the Authy apps, uploading an RSA public key.

```
api_key	37b312a3d682b823c439522e1fd31c82
locale	en-US
key	-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9mfKQEWb9109RNvFboS/
/mF+QEDqdJ1edfR2whKTQWVd7DZG/bOiN7wZjas+HBa9JH2hM243Pq8kBGxLYfWq
cCsCWPinnH2mnBEhwxCviQxnaZVdtXXBN9XQWU6xpANxHV1hjRTLos9F/LxBaUlA
ypM5Wzvwl4ZHelzLlXmULVoaClO2yx6q8tVx8jtM5teqAo59Ux9OuMvPufm30Mzo
Rgfs4v3wrWbZe/XhOzh1XwGqP6ikeOXhS+e1jm2RGaSQS01ByPQk/clZrIf+ZWQz
XTGHRc2FpbT7W3ieH3Rz7rlYyTkuY8HHvlvV3+fwg6zYCeBGqqV+P1sLw/QpyoR/
DwIDAQAB
-----END PUBLIC KEY-----

app_id	510cd201f91ea171d600027d
otp1	9735470
otp2	8142428
otp3	6063286
device_id	804509754
```

I'm unsure what the purpose of these keys are.  They are different for each app, and bear no apparent relationship to the [private key we received earlier](#get-the-rsa-key).

```
{
	"message": "Public key created for device 804509754 and app 43",
	"success": true
}
```

The `43` in the response message is the `serial_id` from the app sync listing.

# Unlock & Update Tokens

`/json/users/123456/authenticator_tokens/update`

After unlocking with our backups password, we see (some) of tokens updated with `POST` requests to this endpoint.

The data sent almost entirely matches the data for each token from the [auth sync](#auth-sync) where we pulled all our authenticator tokens down.

There are some fields from the auth sync that are not sent in the request, but those that are sent all match, with the exception of `account_type`.

It is `authenticator` in the auth sync, and `google` in the `POST` here.  I assume that there was some logic changes around how they display, and possibly handle, tokens by type, so this is a kind of client driven migration to the new values.

```
api_key	37b312a3d682b823c439522e1fd31c82
locale	en-US
name	 john@redacted.com
original_name	Google:john@redacted.com
digits	6
token_id	1607959200
account_type	google
encrypted_seed	PIzpGFA0jeAefTenhvbxaYeD7zFqQKKC63GokXadrmncUCZxKpoSgGg9pMJb/jAz
salt	Tha5DgMhigkbPI2VIvOTPBVfvHFVNFRx
password_timestamp	1519340474
otp1	0998251
otp2	0000283
otp3	4846204
device_id	804509754
```

```json
{
	"message": "Token saved successfully",
	"success": true
}
```

# <a name="authy-auth"></a>Authy API Auth

The authentication scheme for the Authy API is pretty basic, and also a bit inconsistent.

From the first request we see an auth header, `x-authy-api-key`.  This value remains fixed, it doesn't appear to update, even after we've received what appears to be a device specific API key when [completing device registration](#registration-complete).  The value is often in this header, but is also sometimes in the `POST` fields, or even in the query parameters.  Sometimes it is in multiple locations.

Once the device is registered the main authentication method appears to be sending three OTP values tied to the device token along with the device ID.  We saw this first when [getting the device and account information](#get-device-info) and I explored it a bit there.  It's a bit of a novel technique, but it has it's advantages.

# Other Notable Headers

Aside from auth there are some headers which show up in basically every request.

The `x-authy-request-id` header contains a UUID and, despite the name, does not change on every request.  Instead, it appears to be more of a session id. I only observed it changing when moving from an registration to being in an authorized context, as well as on each app open.

The `x-authy-private-ip` header appears to contain local IP's which reference my device.  Buried in a wave of link local IPv6 addresses was my IPv4.  Not sure why Authy needs these, perhaps they have a feature (or planned a feature) to do a same network peer-to-peer type of syncing.

```
127.0.0.1,::1,fe80::1,fe80::acf7:5dff:fe6b:fed9,fe80::acf7:5dff:fe6b:fedb,fe80::acf7:5dff:fe6b:feda,fe80::1461:35c9:4d80:e3d2,192.168.1.147,fe80::e8d4:46ff:fee0:6cca,fe80::e8d4:46ff:fee0:6cca,fe80::1657:cb56:c6af:a314,fe80::a016:3212:6b67:9f72,fe80::ce81:b1c:bd2c:69e
```

The `x-user-agent` header is a bit weird, since it's value, `AuthyDesktop 2.2.3`, is embedded (slightly differently) in the standardized `user-agent` value, `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) AuthyDesktop/2.2.3 Chrome/96.0.4664.110 Electron/16.0.8 Safari/537.36`.  Similarly, `x-authy-device-app` with a value of `authy` seems superfluous.

---

# <a name="decrypting"></a> Decrypting

On to the exciting bit!  We have registered our device, we have captured the encrypted tokens in a network request, and now we want to break them free. We know from some of the JavaScript that it is AES encrypted, but not much else.

```javascript
// This code is found in `js/app.js` from the `app.asar` in the macOS version 2.2.3

// Line 16375
(t = u.decryptAES(r.salt, e.password, r.encryptedSeed)),
```

Luckily, Authy told us in 2018 _exactly_ how this works:  [How Authy 2FA Backups Work](https://authy.com/blog/how-the-authy-two-factor-backups-work/)

In short:

- Your backup password is stretched with PBKDF2, 1000 rounds by default
- This derived key is used with AES-256 in CBC mode
- A different IV used is per account

Perfect! All very boring (if a bit aged) crypto.

Looking at an account, some obvious related fields jump out.

```json
	{
		"account_type": "authenticator",
		"digits": 6,
		"encrypted_seed": "7MXnPTX0BGMhh+s3kKNhjPt4V98nKEVA24/+OiUn3e0=",
		"issuer": null,
		"key_derivation_iterations": null,
		"logo": null,
		"name": "www.velvetcache.org:VelvetCache Blog",
		"original_name": "www.velvetcache.org:admin",
		"password_timestamp": 1519340474,
		"salt": "MIIm9CTN13yqHfxjr4OSQxEDJlyDzVEm",
		"unique_id": "1518713548"
	}
```

The `encrypted_seed` is our cipher text, `key_derivation_iterations` _would_ be a number if that had been increased, but we'll assume the default of 1000.  `salt` is probably our PBKDF2 salt.  But what about the initialization vector?

I went back through the tokens, and none of the 33 had any key that indicated an IV value.  It should be 16 bytes, probably encoded to hex as Authy seems to love encoding to hex.  It didn't seem to be prepended onto the `encrypted_seed` values. I dug back through all the responses I had captured on the API, hoping I'd missed something.  Was this delivered through a side channel for some reason?  Did it elude my capture with Charles?

I was stumped, so back into the obfuscated code.  I found `decryptAESWithKey` which referenced `n.IV`

```javascript
// This code is found in `js/app.js` from the `app.asar` in the macOS version 2.2.3

// Line 20641
static decryptAESWithKey(e, t) {
  (e = s.util.createBuffer(e)),
    (e = s.aes.createDecryptionCipher(e, "CBC"));
  return (
    e.start(s.util.createBuffer(n.IV)),
    (t = s.util.createBuffer(s.util.decode64(t))),
    e.update(t),
    e.finish() ? e.output.data : null
  );
}

// Line 20726
return (
  (n.PBKDF2_PARAMS = { keySize: 8, iterations: 1e3 }),
  (n.IV = s.util.decodeUtf8("\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0")),
  n
);
```

It looks like, the default, and perhaps the _constant_ value of the initialization vector was all null bytes.  That would be unfortunate, but probably not fatal.  The seeds of these MFA tokens are static, they should not need to be re-encrypted unless the key is being changed.

However, AES-CBC with the same key and IV will leak shared prefixes of the plaintext.  A demonstration:

```go
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func main() {
	key, _ := base64.StdEncoding.DecodeString("5azCEYjUOir+fxfMAC3GH3NuSJoKk8Qurbp7apBXfxA=")

	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	samples := [][]byte{
		PKCS5Padding([]byte("Same Prefix Different Endings"), block.BlockSize()),
		PKCS5Padding([]byte("Same Prefix Different Finish"), block.BlockSize()),
	}

	for _, sample := range samples {
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(sample, sample)
		encoded := hex.EncodeToString(sample)
		fmt.Println(encoded[:32], encoded[32:])

	}
}
```

```console
$ go run .
0bb59fab5093e95a49a7abe0a55d5b6b 2e72dea8623bd56e15c07019643d8c26
0bb59fab5093e95a49a7abe0a55d5b6b 7d7c9ad3dc1331de6930024655c73df9
```

That means that if you had two TOTP seeds that had a common prefix, you could see it in their encrypted versions.  However, since each app has it's own salt for PBKDF2, it's essentially impossible for them to have the same key, which is as good as having a different IV.

It's not a real problem, but it does feel uncomfortable somehow, especially since they claim to be doing this differently in that blog post.  Regardless, with this IV in hand we can iterate over our API response tokens and generate TOTP for them at will.

# Wrapping It Up

With all this in hand, we can write our own Authy client, which we can then automate with a cron job, etc.

I've done that in [authy-cli](https://github.com/jmhobbs/authy-cli), which can register, sync, list and export your tokens.

![authy-cli being used to register a new Authy device](https://static.velvetcache.org/pages/2023/05/12/the-authy-backup-system/register.png)

