---
category:
- Geek
creator: admin
date: 2006-10-02
layout: layout.njk
permalink: /2006/10/02/php-states/
tags:
- PHP
- Programming
- Snippets
title: PHP States
type: post
wp_id: "60"
---

Just wrote and used this little snippet.  Thought I'd post it for when I need it again.

```php
<?php
$states = array(
     "AL","AK","AZ","AR","CA",
     "CO","CT","DE","DC","FL",
     "GA","HI","ID","IL","IN",
     "IA","KS","KY","LA","ME",
     "MD","MA","MI","MN","MS",
     "MO","MT","NE","NV","NH",
     "NJ","NM","NY","NC","ND",
     "OH","OK","OR","PA","RI",
     "SC","SD","TN","TX","UT",
     "VT","VA","WA","WV","WI",
     "WY");
?>
```

An example usage, as if it isn't obvious.

```php
<?php
print "<select name='states'>";
foreach($states as $state)
{
  print "
    <option "
    .(($row['prospect_state'] == $state) ? "selected=\"selected\" " : "").
    "value=\"$state\">$state</option>";
}
print "</select>";
?>
```
In case it's unfamiliar to you, the `(($row['prospect_state'] == $state) ? "selected=\"selected\" " : "")` is called the [ternary operator](http://us2.php.net/manual/en/language.operators.comparison.php#language.operators.comparison.ternary), also known as "fancy little if/else".

#### Update (10/02/06)

Here's a beefier version.

```php
$states_list = array(
  'AL'=>"Alabama",  
  'AK'=>"Alaska",  
  'AZ'=>"Arizona",  
  'AR'=>"Arkansas",  
  'CA'=>"California",  
  'CO'=>"Colorado",  
  'CT'=>"Connecticut",  
  'DE'=>"Delaware",  
  'DC'=>"District Of Columbia",  
  'FL'=>"Florida",  
  'GA'=>"Georgia",  
  'HI'=>"Hawaii",  
  'ID'=>"Idaho",  
  'IL'=>"Illinois",  
  'IN'=>"Indiana",  
  'IA'=>"Iowa",  
  'KS'=>"Kansas",  
  'KY'=>"Kentucky",  
  'LA'=>"Louisiana",  
  'ME'=>"Maine",  
  'MD'=>"Maryland",  
  'MA'=>"Massachusetts",  
  'MI'=>"Michigan",  
  'MN'=>"Minnesota",  
  'MS'=>"Mississippi",  
  'MO'=>"Missouri",  
  'MT'=>"Montana",
  'NE'=>"Nebraska",
  'NV'=>"Nevada",
  'NH'=>"New Hampshire",
  'NJ'=>"New Jersey",
  'NM'=>"New Mexico",
  'NY'=>"New York",
  'NC'=>"North Carolina",
  'ND'=>"North Dakota",
  'OH'=>"Ohio",  
  'OK'=>"Oklahoma",  
  'OR'=>"Oregon",  
  'PA'=>"Pennsylvania",  
  'RI'=>"Rhode Island",  
  'SC'=>"South Carolina",  
  'SD'=>"South Dakota",
  'TN'=>"Tennessee",  
  'TX'=>"Texas",  
  'UT'=>"Utah",  
  'VT'=>"Vermont",  
  'VA'=>"Virginia",  
  'WA'=>"Washington",  
  'WV'=>"West Virginia",  
  'WI'=>"Wisconsin",  
  'WY'=>"Wyoming"
);
```
