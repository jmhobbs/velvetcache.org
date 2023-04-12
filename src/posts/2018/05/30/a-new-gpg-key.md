---
category:
- Geek
creator: admin
date: 2018-05-30
permalink: /2018/05/30/a-new-gpg-key/
tags:
- GPG
- How To
- Security
title: A New GPG Key
type: post
wp_id: "2779"
summary: >
  After a decade it's time to create a new GPG key. There's no real excuse not to
  use a Yubikey anymore either. This will cover the process I took from start to
  finish, with subkeys, paperkey backups, git signing and SSH with gpg-agent.
---
It's been 12 years since I created my first GPG key and 11 since I've created the one I actually use.  That is far too long, so I decided to create a new pair and deprecate the old.  [In 2013 I started this process](https://pgp.mit.edu/pks/lookup?search=john%40velvetcache.org&op=index), but I didn't follow through and I've since lost access to those keys.  I know where they are, but the machine died so I need to hook up it's HDD and pull the keys out.

Regardless, it is time for new ones, and I did some reading to get a real plan for this.  I would generate a new, strong key offline, with a subkey for each capability. The subkeys would go onto a smart card, in my case a [Yubikey 4](https://www.yubico.com/product/yubikey-4-series/#yubikey-4).  The primary key material would go to offline backup to keep it safe.

## Disclaimer

Nothing in this post is new or novel, but rather collected from many other posts.  I've tried to link to any relevant posts below each section, and I encourage you to read these sources.  Any mistakes I've made I would be glad if you send me an email (GPG encrypted of course ;) to point it out.

## Yubikey Configuration

After I ordered my Yubikey, I had to configure it.  The Yubikey docs expect a fair amount of knowledge before you start, but the steps are pretty simple when you understand it.  Basically, it boils down to:

1. Change the Admin PIN
2. Change the PIN
3. Set a Reset Code
4. Fill in optional metadata

Plug in your card and proceed as follows:

```console
$ gpg --card-edit

Reader ...........: Yubico Yubikey 4 OTP U2F CCID
Application ID ...: D2760001240102010006075857980000
Version ..........: 2.1
Manufacturer .....: Yubico
Serial number ....: 07------
Name of cardholder: [not set]
Language prefs ...: [not set]
Sex ..............: unspecified
URL of public key : [not set]
Login data .......: [not set
Signature PIN ....: not forced
Key attributes ...: rsa4096 rsa4096 rsa2048
Max. PIN lengths .: 127 127 127
PIN retry counter : 3 0 3
Signature counter : 0
Signature key ....: [none]
Encryption key....: [none]
Authentication key: [none]
General key info..: [none]

gpg/card> admin
Admin commands are allowed

gpg/card> passwd
gpg: OpenPGP card no. D2760001240102010006075857980000 detected

1 - change PIN
2 - unblock PIN
3 - change Admin PIN
4 - set the Reset Code
Q - quit

Your selection? 3
PIN changed.

1 - change PIN
2 - unblock PIN
3 - change Admin PIN
4 - set the Reset Code
Q - quit

Your selection? 1
PIN changed.

1 - change PIN
2 - unblock PIN
3 - change Admin PIN
4 - set the Reset Code
Q - quit

Your selection? 4
Reset Code set.

1 - change PIN
2 - unblock PIN
3 - change Admin PIN
4 - set the Reset Code
Q - quit

Your selection? Q


gpg/card> name
Cardholder's surname: Hobbs
Cardholder's given name: John

gpg/card> lang
Language preferences: en

gpg/card> url
URL to retrieve public key: http://static.velvetcache.org/John-Hobbs-Public-Key.asc

gpg/card> login
Login data (account name): john@velvetcache.org

gpg/card>

Reader ...........: Yubico Yubikey 4 OTP U2F CCID
Application ID ...: D2760001240102010006075857980000
Version ..........: 2.1
Manufacturer .....: Yubico
Serial number ....: 07------
Name of cardholder: John Hobbs
Language prefs ...: en
Sex ..............: unspecified
URL of public key : http://static.velvetcache.org/John-Hobbs-Public-Key.asc
Login data .......: john@velvetcache.org
Signature PIN ....: not forced
Key attributes ...: rsa4096 rsa4096 rsa2048
Max. PIN lengths .: 127 127 127
PIN retry counter : 3 3 3
Signature counter : 0
Signature key ....: [none]
Encryption key....: [none]
Authentication key: [none]
General key info..: [none]

gpg/card> quit
```

#### Links

- [https://developers.yubico.com/PGP/Card_edit.html](https://developers.yubico.com/PGP/Card_edit.html)
- [https://developers.yubico.com/yubikey-piv-manager/PIN_and_Management_Key.html](https://developers.yubico.com/yubikey-piv-manager/PIN_and_Management_Key.html)


## Generating Keys

Next, I created my keys. Be sure you set up a clean environment for this, ideally a random directory in `/tmp`, better still on a `ramfs` of an offline, live CD machine.  But that's a bit drastic for my use case.

```console
$ export GNUPGHOME="/tmp/$(pwgen 30 1)/gnupg"
$ echo $GNUPGHOME
/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg
$ mkdir -p "$GNUPGHOME"
$ cd "$GNUPGHOME/.."
$ chmod 0700 gnupg
$ ls -l
total 0
drwx------  2 johnhobbs  wheel  64 May 30 14:23 gnupg
```

You'll want a good base config file in there too.

```console
$ cat <<EOF > gnupg/gpg.conf
# Show long key IDs, not short: https://gwolf.org/node/4070
keyid-format 0xlong

# Display the calculated validity of user IDs during key listings
list-options show-uid-validity
verify-options show-uid-validity

# List keys with their fingerprints
with-fingerprint

# Default preferences used for creating new keys
default-preference-list SHA512 SHA384 SHA256 SHA224 AES256 AES192 AES CAST5 ZLIB BZIP2 ZIP Uncompressed

# Digest used to sign keys
cert-digest-algo SHA512

# Cipher to use for encrypting private keys
s2k-cipher-algo AES256

# Digest to use for mangling passphrases on private keys
s2k-digest-algo SHA512

# Refuse to run if GnuPG cannot get secure memory.
require-secmem
EOF
```

With the directory in place, I can create a primary key, option 4. 4096-bits is as strong as GPG allows right now, and I set it not to expire because I will be keeping offline and it should be ok to revoke manually if needed.

```console
$ gpg --full-gen-key
gpg (GnuPG/MacGPG2) 2.2.3; Copyright (C) 2017 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

gpg: keybox '/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/pubring.kbx' created
Please select what kind of key you want:
   (1) RSA and RSA (default)
   (2) DSA and Elgamal
   (3) DSA (sign only)
   (4) RSA (sign only)
Your selection? 4
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (2048) 4096
Requested keysize is 4096 bits
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 0
Key does not expire at all
Is this correct? (y/N) y

GnuPG needs to construct a user ID to identify your key.

Real name: John Hobbs
Email address: john@velvetcache.org
Comment:
You selected this USER-ID:
    "John Hobbs <john@velvetcache.org>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? O
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
gpg: /tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/trustdb.gpg: trustdb created
gpg: key 0xCF469E79A0A20E10 marked as ultimately trusted
gpg: directory '/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/openpgp-revocs.d' created
gpg: revocation certificate stored as '/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/openpgp-revocs.d/5A4B39AA4C644429718D6EAACF469E79A0A20E10.rev'
public and secret key created and signed.

Note that this key cannot be used for encryption.  You may want to use
the command "--edit-key" to generate a subkey for this purpose.
pub   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                              John Hobbs <john@velvetcache.org>

$ gpg --list-secret-keys
gpg: checking the trustdb
gpg: marginals needed: 3  completes needed: 1  trust model: pgp
gpg: depth: 0  valid:   1  signed:   0  trust: 0-, 0q, 0n, 0m, 0f, 1u
/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/pubring.kbx
-----------------------------------------------------
sec   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                   [ultimate] John Hobbs <john@velvetcache.org>
```

Now it's time to create subkeys.  There are four capabilities that a PGP key can have.

### `C` is for Certify

Your primary key will have the capability of Certification.  Certify is essentially the ability to sign other keys.  A key with Certify can be "parent" to subkeys, create new subkeys, and edit existing ones.  You also need this capability to sign another users public key.

### `S` is for Sign

A key with the Sign capability can sign files and messages, allowing others to verify their integrity.

### `E` is for Encrypt

A key with the Encrypt capability is used for encrypting files. Simple.

### `A` is for Authenticate

An Authentication key is generally used for SSH authentication.

---------------

Generating the subkeys is a bit tedious, but so it goes.

```console
$ gpg --edit-key 0xCF469E79A0A20E10
gpg (GnuPG/MacGPG2) 2.2.3; Copyright (C) 2017 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Secret key is available.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> addkey
Please select what kind of key you want:
   (3) DSA (sign only)
   (4) RSA (sign only)
   (5) Elgamal (encrypt only)
   (6) RSA (encrypt only)
Your selection? 4
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (2048) 4096
Requested keysize is 4096 bits
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 4y
Key expires at Sun May 29 14:16:18 2022 CDT
Is this correct? (y/N) y
Really create? (y/N) y
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> addkey
Please select what kind of key you want:
   (3) DSA (sign only)
   (4) RSA (sign only)
   (5) Elgamal (encrypt only)
   (6) RSA (encrypt only)
Your selection? 6
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (2048) 4096
Requested keysize is 4096 bits
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 4y
Key expires at Sun May 29 14:16:48 2022 CDT
Is this correct? (y/N) y
Really create? (y/N) y
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> save
```

The authentication key requires **E X P E R T  M O D E**. Git gud.

<!-- todo: aside here? -->
> **Note** 2023-03-31
>
> An RSA key is probably not what you want anymore, consider an ECC key, I used ed25519 when renewing this subkey.

```console
$ gpg --expert --edit-key 0xCF469E79A0A20E10
gpg (GnuPG/MacGPG2) 2.2.3; Copyright (C) 2017 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Secret key is available.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> addkey
Please select what kind of key you want:
   (3) DSA (sign only)
   (4) RSA (sign only)
   (5) Elgamal (encrypt only)
   (6) RSA (encrypt only)
   (7) DSA (set your own capabilities)
   (8) RSA (set your own capabilities)
  (10) ECC (sign only)
  (11) ECC (set your own capabilities)
  (12) ECC (encrypt only)
  (13) Existing key
Your selection? 8

Possible actions for a RSA key: Sign Encrypt Authenticate
Current allowed actions: Sign Encrypt

   (S) Toggle the sign capability
   (E) Toggle the encrypt capability
   (A) Toggle the authenticate capability
   (Q) Finished

Your selection? S

Possible actions for a RSA key: Sign Encrypt Authenticate
Current allowed actions: Encrypt

   (S) Toggle the sign capability
   (E) Toggle the encrypt capability
   (A) Toggle the authenticate capability
   (Q) Finished

Your selection? A

Possible actions for a RSA key: Sign Encrypt Authenticate
Current allowed actions: Encrypt Authenticate

   (S) Toggle the sign capability
   (E) Toggle the encrypt capability
   (A) Toggle the authenticate capability
   (Q) Finished

Your selection? E

Possible actions for a RSA key: Sign Encrypt Authenticate
Current allowed actions: Authenticate

   (S) Toggle the sign capability
   (E) Toggle the encrypt capability
   (A) Toggle the authenticate capability
   (Q) Finished

Your selection? Q
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (2048) 4096
Requested keysize is 4096 bits
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 4y
Key expires at Sun May 29 14:20:03 2022 CDT
Is this correct? (y/N) y
Really create? (y/N) y
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> save
```

That's it!  We're in business.


#### Links

  - [https://spin.atomicobject.com/2013/11/24/secure-gpg-keys-guide/](https://spin.atomicobject.com/2013/11/24/secure-gpg-keys-guide/)
  - [https://www.linode.com/docs/security/authentication/gpg-key-for-ssh-authentication/](https://www.linode.com/docs/security/authentication/gpg-key-for-ssh-authentication/)
  - [https://gist.github.com/abeluck/3383449](https://gist.github.com/abeluck/3383449)
  - [https://alexcabal.com/creating-the-perfect-gpg-keypair/](https://alexcabal.com/creating-the-perfect-gpg-keypair/)
  - [https://gist.github.com/graffen/37eaa2332ee7e584bfda](https://gist.github.com/graffen/37eaa2332ee7e584bfda)

## Backups

Before we do anything else, we need to back that thang up.

I'm choosing two methods: backup to a USB key that will live in a fire safe (who has a safety deposit box these days?), and a printed backup in case the USB key fails.  Ideally these two articles would not be co-located.

First we export the keys and move them to the USB stick.  The `export-secret-subkeys` output is less important than the `export-secret-key` output as it doesn't contain a viable certification key, but would be useful as a "middle tier" of backup that wouldn't expose your primary key to risk.

```console
$ # Dump the public key, for giggles.
$ gpg --armor --export 0xCF469E79A0A20E10 > 0xCF469E79A0A20E10.public.asc
$ # This is the all the secret keys together.
$ gpg --armor --export-secret-key 0xCF469E79A0A20E10 > 0xCF469E79A0A20E10.master.asc
$ # This is just the subkeys.
$ gpg --armor --export-secret-subkeys 0xCF469E79A0A20E10 > 0xCF469E79A0A20E10.subkeys.asc
$ ls -l
total 64
-rw-r--r--   1 johnhobbs  wheel  14134 May 30 14:50 0xCF469E79A0A20E10.master.asc
-rw-r--r--   1 johnhobbs  wheel  14134 May 30 14:50 0xCF469E79A0A20E10.public.asc
-rw-r--r--   1 johnhobbs  wheel  12338 May 30 14:50 0xCF469E79A0A20E10.subkeys.asc
drwx------  13 johnhobbs  wheel    416 May 30 14:20 gnupg
```

Now, we could take these ascii armored keys and just print them, but that's a lot of bytes to pray for OCR to recognize.  Instead, we can use a piece of software called Paperkey which strips out everything but the most secret parts of the key and gives you something much shorter to type in.

```console
$ gpg --export-secret-key | paperkey --output 0xCF469E79A0A20E10.master.paper
$ cat 0xCF469E79A0A20E10.master.paper
# Secret portions of key 5A4B39AA4C644429718D6EAACF469E79A0A20E10
# Base16 data extracted Wed May 30 14:53:56 2018
# Created with paperkey 1.5 by David Shaw
#
# File format:
# a) 1 octet:  Version of the paperkey format (currently 0).
# b) 1 octet:  OpenPGP key or subkey version (currently 4)
# c) n octets: Key fingerprint (20 octets for a version 4 key or subkey)
# d) 2 octets: 16-bit big endian length of the following secret data
# e) n octets: Secret data: a partial OpenPGP secret key or subkey packet as
#              specified in RFC 4880, starting with the string-to-key usage
#              octet and continuing until the end of the packet.
# Repeat fields b through e as needed to cover all subkeys.
#
# To recover a secret key without using the paperkey program, use the
# key fingerprint to match an existing public key packet with the
# corresponding secret data from the paper key.  Next, append this secret
# data to the public key packet.  Finally, switch the public key packet tag
# from 6 to 5 (14 to 7 for subkeys).  This will recreate the original secret
# key or secret subkey packet.  Repeat as needed for all public key or subkey
# packets in the public key.  All other packets (user IDs, signatures, etc.)
# may simply be copied from the public key.
#
# Each base16 line ends with a CRC-24 of that line.
# The entire block of data ends with a CRC-24 of the entire block of data.

  1: 00 04 5A 4E 39 AA 4C 64 44 29 71 8D 6E AA CF 46 9E 79 A0 A2 0E 10 745BFD
  ...
  248: 36 96 66 39 EE 0B36C4
  249: D3A56B</pre>

Still not fun to type it all in, but it's better and this is a last ditch sort of thing anyway.

#### Recovery

Backups you don't test aren't backups, they are hopes and dreams.  So let's try recovering from our paperkey output!

<pre>
$ mkdir recovery/
$ # Paperkey wants the public component to be raw.
$ gpg --dearmor 0xCF469E79A0A20E10.public.asc
$ # You can't specify output filename on dearmor so let's move it.
$ mv 0xCF469E79A0A20E10.public.asc.gpg recovery/
$ # Combine the public with the secret to get a GPG compatible keyring.
$ paperkey --pubring recovery/0xCF469E79A0A20E10.public.asc.gpg --secrets 0xCF469E79A0A20E10.master.paper --output recovery/0xCF469E79A0A20E10.master.gpg
$ # Check it out, without importing it.
$ gpg --import --import-options show-only recovery/0xCF469E79A0A20E10.master.gpg
sec   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                              John Hobbs <john@velvetcache.org>
ssb   rsa4096/0xA93E031FD5AB0841 2018-05-30 [S] [expires: 2022-05-29]
ssb   rsa4096/0xC8A284D483920085 2018-05-30 [E] [expires: 2022-05-29]
ssb   rsa4096/0xED7858737F087831 2018-05-30 [A] [expires: 2022-05-29]

gpg: Total number processed: 1
gpg:       secret keys read: 1
```

#### Links

  - [http://www.jabberwocky.com/software/paperkey/](http://www.jabberwocky.com/software/paperkey/)

## The certificate revoke you, secret key!

While not required, we can generate a revocation certificate while we still have the primary key on this machine.

```console
$ gpg --output 0xCF469E79A0A20E10.revocation-certificate.asc --gen-revoke 0xCF469E79A0A20E10

sec  rsa4096/0xCF469E79A0A20E10 2018-05-30 John Hobbs <john@velvetcache.org>

Create a revocation certificate for this key? (y/N) y
Please select the reason for the revocation:
  0 = No reason specified
  1 = Key has been compromised
  2 = Key is superseded
  3 = Key is no longer used
  Q = Cancel
(Probably you want to select 1 here)
Your decision? 0
Enter an optional description; end it with an empty line:
>
Reason for revocation: No reason specified
(No description given)
Is this okay? (y/N) y
ASCII armored output forced.
Revocation certificate created.

Please move it to a medium which you can hide away; if Mallory gets
access to this certificate he can use it to make your key unusable.
It is smart to print this certificate and store it away, just in case
your media become unreadable.  But have some caution:  The print system of
your machine might store the data and make it available to others!
```

Throw that onto your backup drive too while you're at it.

#### Links

  - [https://www.hackdiary.com/2004/01/18/revoking-a-gpg-key/](https://www.hackdiary.com/2004/01/18/revoking-a-gpg-key/)

## Sign!

Ok.  Everything is generated, we have a good backup, we are ready to transition.  To indicate that this key is your new key, you can sign it with your old one, then send it up to the keyservers in the sky (if you're into that)

```console
$ # --local-user lets us specify which key we want to sign with.
$ gpg --local-user 0x2580c0be34eb9490 --sign-key 0xCF469E79A0A20E10
$ gpg --list-sigs 0xCF469E79A0A20E10
pub   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                   [ unknown] John Hobbs <john@velvetcache.org>
sig 3        0xCF469E79A0A20E10 2018-05-30  John Hobbs <john@velvetcache.org>
sig 3        0x2580C0BE34EB9490 2018-05-30  John Hobbs <john@velvetcache.org>
sub   rsa4096/0xC8A284D483920085 2018-05-30 [E] [expires: 2022-05-29]
sig          0xCF469E79A0A20E10 2018-05-30  John Hobbs <john@velvetcache.org>
sub   rsa2048/0xED7858737F087831 2018-05-30 [A] [expires: 2022-05-29]
sig          0xCF469E79A0A20E10 2018-05-30  John Hobbs <john@velvetcache.org>
sub   rsa4096/0xA93E031FD5AB0841 2018-05-30 [S] [expires: 2022-05-29]
sig          0xCF469E79A0A20E10 2018-05-30  John Hobbs <john@velvetcache.org>
$ # Send it all off to the keyservers!
$ gpg --send-keys 0xCF469E79A0A20E10
gpg: sending key 0xCF469E79A0A20E10 to hkps://hkps.pool.sks-keyservers.net
```

#### Links

  - [https://www.apache.org/dev/key-transition.html](https://www.apache.org/dev/key-transition.html)

## To The Smart Card Robin!

![I'm Batman](http://static.velvetcache.org/pages/2018/05/31/a-new-gpg-key/batman.gif)

Moving the keys onto a smart card helps protect them.  They won't exist on your filesystem anymore, only on the card.  That means they can't be read out and stolen by a malicious process, but you can still use them by providing your smart card pin and key password.

Keep in mind, this is a one way trip.  Make sure your backups are really, truly in place. We want to move our Signing, Encryption and Authentication keys onto the card.  The Certification key we will only store offline, as mentioned before.

```console
$ gpg --list-secret-keys
/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/pubring.kbx
-----------------------------------------------------
sec   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                   [ultimate] John Hobbs <john@velvetcache.org>
ssb   rsa4096/0xA93E031FD5AB0841 2018-05-30 [S] [expires: 2022-05-29]
ssb   rsa4096/0xC8A284D483920085 2018-05-30 [E] [expires: 2022-05-29]
ssb   rsa4096/0xED7858737F087831 2018-05-30 [A] [expires: 2022-05-29]
$ gpg --edit-key 0xCF469E79A0A20E10
gpg (GnuPG/MacGPG2) 2.2.3; Copyright (C) 2017 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Secret key is available.

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> key 1

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb* rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> keytocard
Please select where to store the key:
   (1) Signature key
   (3) Authentication key
Your selection? 1

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb* rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> key 2

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb* rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb* rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> key 1

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb* rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> keytocard
Please select where to store the key:
   (2) Encryption key
Your selection? 2

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb* rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb  rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> key 3

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb  rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb* rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb* rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> key 2

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb   rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb* rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> keytocard
Please select where to store the key:
   (3) Authentication key
Your selection? 3

sec  rsa4096/0xCF469E79A0A20E10
     created: 2018-05-30  expires: never       usage: SC
     trust: ultimate      validity: ultimate
ssb   rsa4096/0xA93E031FD5AB0841
     created: 2018-05-30  expires: 2022-05-29  usage: S
ssb  rsa4096/0xC8A284D483920085
     created: 2018-05-30  expires: 2022-05-29  usage: E
ssb* rsa4096/0xED7858737F087831
     created: 2018-05-30  expires: 2022-05-29  usage: A
[ultimate] (1). John Hobbs <john@velvetcache.org>

gpg> save
```

## HOYB

![Hold onto your butts.](http://static.velvetcache.org/pages/2018/05/31/a-new-gpg-key/hoyb.gif)

This is it.  The big moment.  Take out that smart card, secure your backups, and let's delete our primary key material.

```console
$ gpg --delete-secret-key 0xCF469E79A0A20E10
gpg (GnuPG/MacGPG2) 2.2.3; Copyright (C) 2017 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.


sec  rsa4096/0xCF469E79A0A20E10 2018-05-30 John Hobbs <john@velvetcache.org>

Delete this key from the keyring? (y/N) y
This is a secret key! - really delete? (y/N) y
$ gpg --list-secret-keys
```

Now you can have gpg create key stubs for all the keys on your smart card.

```console
$ gpg --card-status
Reader ...........: Yubico Yubikey 4 OTP U2F CCID
Application ID ...: D2760001240102010006075857980000
Version ..........: 2.1
Manufacturer .....: Yubico
Serial number ....: 07------
Name of cardholder: John Hobbs
Language prefs ...: en
...
General key info..: sub  rsa4096/0xCF469E79A0A20E10 2018-05-30 John Hobbs <john@velvetcache.org>
sec#  rsa4096/0xCF469E79A0A20E10  created: 2018-05-30  expires: never
ssb>  rsa4096/0xA93E031FD5AB0841  created: 2018-05-30  expires: 2022-05-29
                                  card-no: 0006 07------
ssb>  rsa2048/0xC8A284D483920085  created: 2018-05-30  expires: 2022-05-29
                                  card-no: 0006 07------
ssb>  rsa4096/0xED7858737F087831  created: 2018-05-30  expires: 2022-05-29
                                  card-no: 0006 07------
</pre>

Now when we list keys, we see that our primary key has a `#` next to it, showing we don't have access to that secret key.  The subkeys have a `&gt;` next to them showing they are stubs for the keys on the smart card.  Success!

<pre>
$ gpg --list-secret-keys
/tmp/mah1zakioboo1Caipa3ORu5ielohga/gnupg/pubring.kbx
-----------------------------------------------------
sec#   rsa4096/0xCF469E79A0A20E10 2018-05-30 [SC]
      Key fingerprint = 5A4B 39AA 4C64 4429 718D  6EAA CF46 9E79 A0A2 0E10
uid                   [ultimate] John Hobbs <john@velvetcache.org>
ssb>  rsa4096/0xA93E031FD5AB0841 2018-05-30 [S] [expires: 2022-05-29]
ssb>  rsa4096/0xC8A284D483920085 2018-05-30 [E] [expires: 2022-05-29]
ssb>  rsa4096/0xED7858737F087831 2018-05-30 [A] [expires: 2022-05-29]
```

## Fin!

That's it.  There is, of course, more to do, like setting up [git signing](https://developers.yubico.com/PGP/Git_signing.html), [SSH access](https://developers.yubico.com/PGP/SSH_authentication/), etc.  But the new keypair is created, and it's on the Yubikey, so that's all for now.

---------

## Update: Git Signing

Turns out git signing is a cinch.  Just throw a couple items into your git config and it's automatic and transparent.

```toml
[user]
  signingKey = 0xF79C72E6EDC70E38
[commit]
  gpgSign = true
[log]
  showSignature = true
[merge]
  verifySignatures = true
```

#### `user.signingKey`

Tells git which key to use for signing, unset it just uses the default key.

#### `commit.gpgSign`

Makes it sign all commits by default, instead of passing `-S` to every `git commit`.

#### `log.showSignature`

By default, git won't show you if a commit is GPG signed.  You can see it with `gpg log --show-signature`, or you can set it as default with this config option.

It makes signed commits much chunkier, so be aware of the reduced screen real estate.

```
commit 6f02c4df4fac400841bf3970c1022c7358298333 (HEAD -> gpg-demo)
gpg: Signature made Wed Jun  6 11:52:40 2018 CDT
gpg:                using RSA key 44DC4F5A950F24A65D3F305801FC8AE9E5070C1D
gpg: Good signature from "John Hobbs <john@velvetcache.org>" [ultimate]
Primary key fingerprint: 5616 12FF A10D 9D7A 7FFB  75F4 F79C 72E6 EDC7 0E38
     Subkey fingerprint: 44DC 4F5A 950F 24A6 5D3F  3058 01FC 8AE9 E507 0C1D
Author: John Hobbs <john@velvetcache.org>
Date:   Wed Jun 6 11:49:33 2018 -0500
```

#### `merge.verifySignatures`

This is the only one I am _not_ setting by default.  If you have it enabled, all merges that include unsigned commits will be rejected.  This really only works if everyone in your organization is signing all their commits.

#### Links

  - [https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work)
  - [https://git-scm.com/docs/git-config](https://git-scm.com/docs/git-config)

---------------------------

## Update: One-Touch Actions

By default, with the smart card in, GPG will happily sign and decrypt things after you enter your PIN the first time, with no further interaction from you.  The Yubikey offers a mode where these actions require a touch on the key to complete, which I like because it makes the action more explicit without requiring me to remove the key between operations.

To enable this, you need a special script, `yubitouch.sh`.  To make it work with my GPG Tools install, I had to hard code the path to `gpg-connect-agent` (`/usr/local/MacGPG2/bin/gpg-connect-agent`) and my admin PIN, since pinentry wasn't working and I didn't want it in my bash history.

```console
$ ./yubitouch.sh sig on
All done!
$ ./yubitouch.sh aut on
All done!
$ ./yubitouch.sh dec on
All done!
```

Now, when GPG needs to sign something, my Yubikey flashes at me until I touch it and give my permission.  Neat.

#### Links

  - [https://developers.yubico.com/PGP/Card_edit.html](https://developers.yubico.com/PGP/Card_edit.html)
  - [https://github.com/a-dma/yubitouch](https://github.com/a-dma/yubitouch)
