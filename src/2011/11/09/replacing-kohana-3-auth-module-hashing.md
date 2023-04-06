---
category:
- Geek
creator: admin
date: 2011-11-10
layout: layout.njk
permalink: /2011/11/09/replacing-kohana-3-auth-module-hashing/
tags:
- Hashing
- Kohana
- Passwords
- Snippets
title: Replacing Kohana 3 Auth module hashing
type: post
wp_id: "2006"
---

The password hashing in the Auth module provided with Kohana 3.1 is not very good.  By default it is a simple sha256 hmac with a global salt.

```php
// modules/auth/classes/kohana/auth.php
public function hash($str)
{
  if ( ! $this->_config['hash_key'])
    throw new Kohana_Exception('A valid hash key must be set in your auth config.');

  return hash_hmac($this->_config['hash_method'], $str, $this->_config['hash_key']);
}
```

This [isn't strong](http://security.stackexchange.com/questions/4687/are-salted-sha-256-512-hashes-still-safe-if-the-hashes-and-their-salts-are-expos).  If you loose the hashes and the salt it's just a matter of winding up a [GPU](http://www.golubev.com/blog/?p=35).


So how can we fix this? Well, thanks to Kohana's structure we can easily override the `Auth` class and tweak it.  However, due to Auth's structure, we can't drop the global salt. The hash function has to stand alone, so no passing in salts from the database.


That leaves us with [key stretching](http://en.wikipedia.org/wiki/Key_stretching).

Now, I don't want to deal with a custom key stretching implementation, I'm not a cryptographer.  So, let's find an existing algorithm.


One that pops to mind is [PBKDF2](http://en.wikipedia.org/wiki/PBKDF2). This is a pretty simple algorithm, so it was easy to find and spot check [a PHP implementation](http://www.itnewb.com/tutorial/Encrypting-Passwords-with-PHP-for-Storage-Using-the-RSA-PBKDF2-StandardL)

We just take some [test vectors from RFC 3962](http://tools.ietf.org/html/rfc3962#appendix-B) and run them against the code we found.

```php
<?php
  require_once( 'pbkdf2.php' );

  header( 'Content-Type: text/plain' );

  $tests = array(
    array(
      'rounds' => 1,
      'bits' => 128,
      'expected' => "cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15"
    ),
    array(
      'rounds' => 1,
      'bits' => 256,
      'expected' => "cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15 0a d1 f7 a0 4b b9 f3 a3 33 ec c0 e2 e1 f7 08 37"
    ),
    array(
      'rounds' => 2,
      'bits' => 128,
      'expected' => "01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d"
    ),
    array(
      'rounds' => 2,
      'bits' => 256,
      'expected' => "01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d a0 53 78 b9 32 44 ec 8f 48 a9 9e 61 ad 79 9d 86"
    ),
    array(
      'rounds' => 1200,
      'bits' => 128,
      'expected' => "5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b"
    ),
    array(
      'rounds' => 1200,
      'bits' => 256,
      'expected' => "5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b a7 e5 2d db c5 e5 14 2f 70 8a 31 e2 e6 2b 1e 13"
    ),
  );

  foreach( $tests as $test ) {
    print $test['rounds'] . ' rounds at ' . $test['bits'] . ' bits ' . "\n";
    $start = microtime( TRUE );
    $result = trim( preg_replace( '/(..)/', '\1 ', bin2hex( pbkdf2( 'password', 'ATHENA.MIT.EDUraeburn', $test['rounds'], $test['bits']/8, 'sha1' ) ) ) );
    $diff = microtime( TRUE ) - $start;
    print 'Expected: ' . $test['expected'] . "\n";
    print '     Got: ' . $result . "\n";
    if( $result == $test['expected'] ) {
      print "MATCH\n";
    }
    else { 
      print "NO MATCH\n";
    }
    print 'Took ' . number_format( $diff, 10 ) . "\n\n";
  }
```

Run it, and everything checks out:

```text
1 rounds at 128 bits 
Expected: cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15
     Got: cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15
MATCH
Took 0.0000329018

1 rounds at 256 bits 
Expected: cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15 0a d1 f7 a0 4b b9 f3 a3 33 ec c0 e2 e1 f7 08 37
     Got: cd ed b5 28 1b b2 f8 01 56 5a 11 22 b2 56 35 15 0a d1 f7 a0 4b b9 f3 a3 33 ec c0 e2 e1 f7 08 37
MATCH
Took 0.0000190735

2 rounds at 128 bits 
Expected: 01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d
     Got: 01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d
MATCH
Took 0.0000147820

2 rounds at 256 bits 
Expected: 01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d a0 53 78 b9 32 44 ec 8f 48 a9 9e 61 ad 79 9d 86
     Got: 01 db ee 7f 4a 9e 24 3e 98 8b 62 c7 3c da 93 5d a0 53 78 b9 32 44 ec 8f 48 a9 9e 61 ad 79 9d 86
MATCH
Took 0.0000200272

1200 rounds at 128 bits 
Expected: 5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b
     Got: 5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b
MATCH
Took 0.0019500256

1200 rounds at 256 bits 
Expected: 5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b a7 e5 2d db c5 e5 14 2f 70 8a 31 e2 e6 2b 1e 13
     Got: 5c 08 eb 61 fd f7 1e 4e 4e c3 cf 6b a1 f5 51 2b a7 e5 2d db c5 e5 14 2f 70 8a 31 e2 e6 2b 1e 13
MATCH
Took 0.0144000053
```

So now all that's left is to drop it in, which is pretty simple.  One thing to note is that I wanted this to stay compatible with the default auth config file, so I just extended that a little bit.

```php
<?php
  // application/classes/auth.php

  abstract class Auth extends Kohana_Auth {

    public function hash ( $str ) {
      if ( ! $this->_config['hash_key'] )
        throw new Kohana_Exception( 'A valid hash key must be set in your auth config.' );

      if ( 'pbkdf2' == $this->_config['hash_method'] ) {
        return base64_encode( self::pbkdf2( 
          $str, 
          $this->_config['hash_key'], 
          Arr::get( $this->_config['pbkdf2'], 'rounds', 1000 ),
          Arr::get( $this->_config['pbkdf2'], 'length', 45 ),
          Arr::get( $this->_config['pbkdf2'], 'method', 'sha256' )
        ) );
      }
      else {
        return parent::hash( $str );
      }
    }

    /** PBKDF2 Implementation (described in RFC 2898)
     *
     *  @param string p password
     *  @param string s salt
     *  @param int c iteration count (use 1000 or higher)
     *  @param int kl derived key length
     *  @param string a hash algorithm
     *
     *  @return string derived key
     *
     *  @url http://www.itnewb.com/tutorial/Encrypting-Passwords-with-PHP-for-Storage-Using-the-RSA-PBKDF2-StandardL
    */
    public static function pbkdf2 ( $p, $s, $c, $kl, $a = 'sha256' ) {

        $hl = strlen(hash($a, null, true)); # Hash length
        $kb = ceil($kl / $hl);              # Key blocks to compute
        $dk = '';                           # Derived key

        # Create key
        for ( $block = 1; $block <= $kb; $block ++ ) {

            # Initial hash for this block
            $ib = $b = hash_hmac($a, $s . pack('N', $block), $p, true);

            # Perform block iterations
            for ( $i = 1; $i < $c; $i ++ )

                # XOR each iterate
                $ib ^= ($b = hash_hmac($a, $b, $p, true));

            $dk .= $ib; # Append iterated block
        }

        # Return derived key of correct length
        return substr($dk, 0, $kl);
    }

  }
```

```php
<?php defined('SYSPATH') or die('No direct access allowed.');
  // application/config/auth.php

  return array(
    'driver'       => 'orm',
    'hash_method'  => 'pbkdf2',
    'hash_key'     => 'zomg',
    'lifetime'     => 1209600,
    'session_key'  => 'auth_user',
    'pbkdf2'       => array(
      'method'  => 'sha256',
      'rounds'  => 1000,
      'length'  => 45,
    )
  );
```

One item to note is that I am packing these with `base64_encode`.  This is to fit into the default field type for the ORM driver.  That is also why my length is stunted to 45.  If you really want to go all out, alter your table to use a `TINYBLOB`, up the length to 256 bit and up the rounds.

So that is how I replace weak hashing in K3 with something a bit better.

How do you do it?
