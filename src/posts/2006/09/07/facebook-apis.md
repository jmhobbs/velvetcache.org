---
category:
- Geek
creator: admin
date: 2006-09-07
permalink: /2006/09/07/facebook-apis/
tags:
- Open Source
- Programming
- School
- XML
title: Facebook API's
type: post
wp_id: "22"
---
This article is in response to a Facebook posting about Facebook "selling" personal information through their developers section.  I originally tried to post this on facebook, but it has restrictive message limits, so here it is in it's entirety.  (Pardon the spelling and such, I wrote this is in a rush.)

This is the posting I found on the Facebook site, verbatim:

> IMPORTANT! EVERYONE MUST READ THIS-
> With all the new changes, Facebook has adopted a "Facebook Development Platform." Basically, it allows Facebook Inc. to sell any information on you to anyone. This includes your picture, hometown, current location, interests, political views, musical preferences, relationship status, etc. Pretty much anything that you enter on Facebook is sold. The best part is that you are AUTOMATICALLY ENROLLED! They didn't even tell us! I'm pretty sure Facebook thought they could get away with nobody noticing it since everyone is so overloaded with this new news feed/mini feed junk. If you don't believe me, check it out for yourself in the Facebook Terms of Service. It's black and white. So screw them, to remove yourself from the "Facebook Development Platform" follow these instructions:
>
> 1. Log in to facebook.
> 2. Click "My Privacy" on the left edge of the window.
> 3. Under the network, "Everyone" click "edit settings"
> 4. Scroll to the bottom of the page to the heading "Facebook Development Platform" and uncheck the statement that says "My information may be used according to the restricted Terms of Service."
> 5. Click Save.
> 7. You have official thwarted facebook from whoring out your personal info to the highest bidder! Spread the news! Tell your friends to remove their name from the selling block. This is utterly disgusting!

First of all, they aren't selling the information, they are exposing web services.  What that means is that they've realeased API's to let developers write applications to extend Facebook.  The developers can't access information about you unless you log in to the application they developed.  Additionally all information passess through http://api.facebook.com/,  which means the developers don't have direct acess to your passwords, so they aren't going to steal your account information that way.

Anyone seeking to develop an application has to apply for an API key, which doesn't cost anything, but can uniquely identify, and allow for the immediate disabling of any application written to steal data.

From the Facebook API FAQ:

> Any content delivered to the outside application can be safely displayed to the user. However, in general, content delievered when using a session key should only be stored until that session key expires, or twelve hours, whichever comes later. The exceptions to this are user ids and affiliations information. This is detailed further in the accompanying documentation.

The opening of these API's and web services are a contribution to the Facebook society, and something they're doing for free.  Open source programming is a great benefit to end users, and these API's and those of other sites allow the creation of rich, dynamic and integrated web applications.

If people take the time to read and understand whats going on with the system they can see that there is noting to "be afraid" of, and that their information is just as secure as ever.  Which, I'll admit, isn't saying much.  Social networking systems such as Facebook are hugely complex systems, and Facebooks should be applauded on their content control systems.  They've added a at least a mediocum of security to the information on their system with a complex network of "friend" relationships.  Compared to other sites, Myspace for example, the security of your information inside of Facebook is commendable.
That leads to the final argument here.  Even though Facebook is releasing these API's (which are a good thing) it's all still in your control.  As the post that first sparked this response says, you can opt out.  And barring that Facebook, like other web services such as Flickr, Myspace, etc. is an opt-in.  You signed up, if you don't like it that much, unsubscribe.

I apologize if this was even slightly emotionally charged, I just don't want to see a good thing ruined because of illogical fears.

Here are some links to help understand whats actually going on.

- [http://developers.facebook.com/faq.php](http://developers.facebook.com/faq.php)
- [http://en.wikipedia.org/wiki/Application_programming_interface](http://en.wikipedia.org/wiki/Application_programming_interface)
- [http://en.wikipedia.org/wiki/Web_services](http://en.wikipedia.org/wiki/Web_services)

And here are some example applications built on the API.

- [http://www.blabbook.com/](http://www.blabbook.com/)
- [http://matchrevolution.com/](http://matchrevolution.com/)

Please read and consider these.  If you find a posting on Facebook that goes against the new web services, post a link to this site!  We could potentially lose a great number of wonderfull social services provided by clever programmers and the Facebook API's!
