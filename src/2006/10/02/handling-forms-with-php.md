---
category:
- Geek
creator: admin
date: 2006-10-03
layout: layout.njk
permalink: /2006/10/02/handling-forms-with-php/
tags:
- PHP
- Programming
- Tutorial
title: Handling Forms With PHP
type: post
wp_id: "53"
---

This is a quick tutorial I wrote about, well, form handling with php.  It's not the end-all ultimate guide, it's just an introduction.


_For this tutorial we'll create a simple login script that will compare values of user input to a preset array to validate the user. Here's the example form we will be using..._

<!--nextpage-->

Most forms on the internet are submitted using either the `GET` method or the `POST` method.  The most obvious difference between these two methods is that using `GET` you will see the values of the form fields displayed in the address of the page.  For example `login.php?user=UserName&pass=PassWord` would retrieve a page called `login.php` and pass it the variables `user` and `pass` with values `UserName` and `PassWord` respectively.


`POST` however works more "beneath the surface" and does all of it's data transmission in a manner completely transparent to the common end user.  Obviously there are advantages to "hiding" this information by using `POST`, but it is important to grasp early on that it is nearly as easy to manipulate `POST` data as `GET` data.  One of the upsides of using `GET` is that a page retrieval using this method can be bookmarked and will likely work, unless the application is written in a way that disallows this.

For this tutorial we'll create a simple login script that will compare values of user input to a preset array to validate the user. Here's the example form we will be using.

```html
<!-- form.html -->
<form action="login.php" method="POST">
Username: <input type="text" name="user" /><br/>
Password: <input type="password" name="pass" /><br/>
<input type="submit" value="Log In" />
</form>
```

This renders in Firefox to look like this:

![Form Render](https://static.velvetcache.org/pages/2006/10/02/handling-forms-with-php/form.png)

A simple, everyday web form, nothing special there.  Note that we chose the `POST` method to give at least token privacy and prevent users from passing their name and password on the address line.  Again its important to emphasize that allthough the password is entered as asterisks on the browser end, this _will_ be sent in clear text unless the connection is made over SSL (https, not http).


Now we need to develop the `login.php` to handle the input from the form.  The first step in this process is learning how to access the values passed by the browser in the `POST`.

PHP 4.1.0 onward has predefined arrays of variables know as "superglobals".  These arrays can be accessed from anywhere in a script or class and contain a large quantity of useful data.  They are:

- `$GLOBALS`
- `$_GET`
- `$_POST`
- `$_SERVER`
- `$_COOKIE`
- `$_FILES`
- `$_ENV`
- `$_SESSION`
- `$_REQUEST`

The ones we will be using are `$_GET`, `$_POST` and `$_REQUEST`.

As you may imagine the `$_GET` array contains all of the variables and values passed by the `GET` method to the script.  Likewise the `$_POST` array contains all the variables passed by the `POST` method to the script.  `$_REQUEST` is a hybrid, containing all the variables passed by the `GET` method, the `POST` method as well as variables from `$_COOKIE`.


So let's start off by simply printing the data that the user entered in the form.  To do this we'll use the `$_POST` superglobal array.

```php
<?php
  // login.php
  print $_POST['user'];
  print "<br/>";
  print $_POST['pass'];
?>
```

There it is, our first form-processing script.  So if we enter `john` as the user name, and `hobbs` as the password, the script should print that back with a break, like so:

```html
john
<br/>
hobbs
```

All right, now we need to check and see if this is a valid user/password combination, and for that we need to set up an associated array of valid combinations in `login.php`

```php
<?php
  // login.php
  $validUsers = array("peter"=>"pan", "john"=>"smith", "robert"=>"scoble");
  print $_POST['user'];
  print "<br/>";
  print $_POST['pass'];
?>
```

Now that we have these valid users and their passwords (a.k.a last names), we can use the `[array_key_exists()](http://us2.php.net/manual/en/function.array-key-exists.php)` function and simple comparators to check if we the submitted information is valid.

```php
<?php
  // login.php
  $validUsers = array("peter"=>"pan", "john"=>"smith", "robert"=>"scoble");
  print $_POST['user'];
  if(array_key_exists($_POST['user'], $validUsers))
  {
    $user = $_POST['user'];
    print "<br/>".$validUsers["$user"];
  }
  else
    print "<br/>No matches found.";
?>
```

If your submitted information matches a username, it prints out the password.  If not it tells you so. Now we are at the point where we can use use the password to make the final decision about whether to admit access or not.  Here's the rest of it.

```php
<?php
  // login.php
  $validUsers = array("peter"=>"pan", "john"=>"smith", "robert"=>"scoble");
  if(array_key_exists($_POST['user'], $validUsers))
  {
    $user = $_POST['user'];
    if($validUsers["$user"] == $_POST['pass'])
    {
      print "This is the protected information.";
    }
    else
    {
      print "Sorry, the password you supplied doesn't match up.";
      exit();
    }
  }
  else
  {
    print "The username ".$validUsers["$user"]." was not found.";
    exit();
  }
?>
```

There you have it, the absolute essentials of form handling and a minor example of password protection.
