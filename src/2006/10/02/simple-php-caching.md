---
category:
- Geek
creator: admin
date: 2006-10-02
layout: layout.njk
permalink: /2006/10/02/simple-php-caching/
tags:
- PHP
- Programming
title: Simple PHP Caching
type: post
wp_id: "57"
---

Today I was working on a script that generates two drop-down select box with 28472 and 6197 options respectively.  Those are long running items, mostly because it's based on a javascript selection that shows or hides the two select boxes as needed.  That means the javascript engine is trying to cope with appending and removing two giant select boxes into the DOM tree.  I didn't realize this at first, I thought it was just a long query so I made this cache.  Anyway, it should save some time on the PHP side while I try to put together a better solution.

Essentially it works like this:

```
IF (# Of Rows From DB != # Of Rows From Cache)
  Get Rows From DB
ELSE
  PRINT Cache
WRITE Page TO Cache
```

Wow, I totally made up my own ugly version of a modeling language there, I'm so cool.

```php
// Simple Caching Code
$filename = "schoolSwap.js.cache";
$handle = fopen ($filename, 'r');
fgets($handle);
$versionHS = fgets($handle);
$versionCO = fgets($handle);
fclose($handle);
$versionHS = intval(trim($versionHS));
$versionCO = intval(trim($versionCO));
$query = "  SELECT *
		   FROM highschools
		   ORDER BY highschool_name
		";
$resultHS = MSSQL_QUERY($query);
$numRowsHS = MSSQL_NUM_ROWS($resultHS);

$query = "  SELECT *
		   FROM colleges
		   ORDER BY college_name
		";
$resultCO = MSSQL_QUERY($query);
$numRowsCO = MSSQL_NUM_ROWS($resultCO);

if($versionCO == $numRowsCO && $versionHS == $numRowsHS)
{
	// Use the cached version
	$handle = fopen ($filename, 'r');
	print fread($handle, filesize($filename));
	fclose($handle);
	exit();
}
$contents = "<!--
$numRowsHS
$numRowsCO
-->";

/* In this section you append the page content
    with $contents .= statements.
*/

print $contents;

if (is_writable($filename)) {

   if (!$handle = fopen ($filename, 'w+')) {
         echo "<!-- Cannot open file ($filename) -->";
         exit;
   }

   // Write $content to our opened file.
   if (fwrite($handle, $contents) === FALSE) {
       echo "<!-- Cannot write to file ($filename) -->";
       exit;
   }
  
   echo "<!-- Success, wrote cache to file ($filename) -->";
  
   fclose ( $handle ) ;

} else {
   echo "<!-- The file $filename is not writable -->";
}
```
