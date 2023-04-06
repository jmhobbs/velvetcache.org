---
category:
- Geek
creator: admin
date: 2006-10-23
layout: layout.njk
permalink: /2006/10/23/javascript-password-generator/
tags:
- JavaScript
- Programming
title: Javascript Password Generator
type: post
wp_id: "78"
---

Wrote this for a page at UNO.  Simple, quick password generator.

```javascript
  var passwordElements = new Array(36);
    passwordElements[0]="0";
    passwordElements[1]="1";
    passwordElements[2]="2";
    passwordElements[3]="3";
    passwordElements[4]="4";
    passwordElements[5]="5";
    passwordElements[6]="6";
    passwordElements[7]="7";
    passwordElements[8]="8";
    passwordElements[9]="9";
    passwordElements[10]="A";
    passwordElements[11]="B";
    passwordElements[12]="C";
    passwordElements[13]="D";
    passwordElements[14]="E";
    passwordElements[15]="F";
    passwordElements[16]="G";
    passwordElements[17]="H";
    passwordElements[18]="I";
    passwordElements[19]="J";
    passwordElements[20]="K";
    passwordElements[21]="L";
    passwordElements[22]="M";
    passwordElements[23]="N";
    passwordElements[24]="O";
    passwordElements[25]="P";
    passwordElements[26]="Q";
    passwordElements[27]="R";
    passwordElements[28]="S";
    passwordElements[29]="T";
    passwordElements[30]="U";
    passwordElements[31]="V";
    passwordElements[32]="W";
    passwordElements[33]="X";
    passwordElements[34]="Y";
    passwordElements[35]="Z";

  var RandPassword = "Z";
  var rightnow = new Date();
  function generatePassword(firstInput, secondInput)
  {
    for(x = 0; x < 7; x++)
    {
      RandPassword += passwordElements[Math.floor(Math.random(rightnow.getSeconds())*36)];
    }
    document.getElementById(firstInput).value=RandPassword;
    document.getElementById(secondInput).value=RandPassword;
    RandPassword = "Z";
  }
```

#### Update (10/25/06)

I suppose it would be useful to suggest popping an alert box to let them know what their "random" generated password is.  We actually use the handily insecure method of emailing it to them, but throwing up an alert would work nicely.  Though if you aren't on SSL why bother with trying to be secure anyway? *sigh*
