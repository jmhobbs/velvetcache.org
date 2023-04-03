---
category:
- Geek
creator: admin
date: 2015-12-09
layout: layout.njk
permalink: /2015/12/08/using-lets-encrypt-with-dreamhost/
tags:
- DreamHost
- guide
- ssl
- tls
title: Using Let's Encrypt With Dreamhost
type: post
wp_id: "2666"
summary: Manually issuing and installing a Let's Encrypt certificate on Dreamhost
---

#### Update

2016-01-28

As pointed out in the comments, [Dreamhost now supports Let's Encrypt](https://www.dreamhost.com/blog/2015/12/03/lets-encrypt-and-dreamhost/) in the panel.  No more workaround needed!

-------

Let's Encrypt has [entered public beta](https://letsencrypt.org/2015/12/03/entering-public-beta.html), which means I should probably play with it!

This website is hosted on Dreamhost, which has a round about way of installing SSL certs, but it's not too bad.

First, you have to go to your Dreamhost panel, then "Secure Hosting" and select "Add Secure Hosting"

![Add Secure Hosting](http://static.velvetcache.org/pages/2015/12/08/add-secure-hosting.png)

From here, you pick your domain you want to secure.  It's a little bit wonky, in that it doesn't show **www.** domains as subdomains in this list, so if you use that, you'll need to just select the parent domain.

Doing this will issue you a self-signed certificate which will throw up scary browser warnings.  We will fix that next.

I chose to run Let's Encrypt on my laptop, so I followed the [user guide](http://letsencrypt.readthedocs.org/en/latest/using.html) to get things installed.  Basically just a git clone.

Next you have to begin the request process.

```console
$ ./letsencrypt-auto certonly --manual --debug
```

- `certonly` states that we only want a certificate generated, not installed.
- `--manual` means that we are going to manually authenticate it.
- `--debug` is used with the OS X version because it is experimental.

This will probably download some junk with homebrew, then it's going to ask you some questions, the greatest of which is what domain you want to use.

With this in hand, it will generate an authentication string that you need to put into a file on the server.

```console
[claw]$ cd velvetcache.org
[claw]$ mkdir -p .well-known/acme-challenge/
[claw]$ echo -n "EVgSHY-sQeMAy4TTx_-jjrx-mR3Dmr4M5Byt9vBKcLE.9wpTWpx1Ghg8yXEMASBfWbfU-fGgjG6D-ixF4ip3cDU" &gt; .well-known/acme-challenge/EVgSHY-sQeMAy4TTx_-jjrx-mR3Dmr4M5Byt9vBKcLE
```

Once you do that, it spits out your certificate into `/etc/letsencrypt/live/[domain]`

```console
/etc/letsencrypt/live/www.velvetcache.org $ ls
cert.pem	chain.pem	fullchain.pem	privkey.pem
```

Back on the Dreamhost panel, you'll want to click on "Edit" for the domain we are securing, then select "Manual Configuration".

![Edit Secure Hosting](http://static.velvetcache.org/pages/2015/12/08/edit-secure-hosting.png)

You can clear the CSR field and then into the "Certificate" field, enter the content from `cert.pem`.

Into "Intermediate Certificate" I placed the contents of `chain.pem`

Lastly, we have to change the format of the private key file to one Dreamhost understands.

```console
/etc/letsencrypt/live/www.velvetcache.org $ openssl rsa -in privkey.pem -out privkey.key
writing RSA key
```

Then we paste `privkey.key` into the Dreamhost interface for "Private Key", save and wait for our new certificate to get installed.

![Editing Certificates](http://static.velvetcache.org/pages/2015/12/08/edited-keys.png)

It's magic!

![Add Secure Hosting](http://static.velvetcache.org/pages/2015/12/08/success.png)

Now I just have to fix all my asset URLs too...

