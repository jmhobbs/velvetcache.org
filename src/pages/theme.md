---
eleventyExcludeFromCollections: true
noindex: true
permalink: /theme/
---

<hr/>
<h1>Theme Colors</h1>
<hr/>

<table class="theme-demo">
  <tr><td class="light"></td><td><code>--main-font-color-light</code></td></tr>
  <tr><td class="dark"></td><td><code>--main-font-color-dark</code></td></tr>
  <tr><td class="weak"></td><td><code>--main-font-color-weak</code></td></tr>
  <tr><td class="strong"></td><td><code>--main-font-color-strong</code></td></tr>
  <tr><td class="theme"></td><td><code>--theme-color</code></td></tr>
  <tr><td class="theme-light"></td><td><code>--theme-color-light</code></td></tr>
  <tr><td class="theme-dark"></td><td><code>--theme-color-dark</code></td></tr>
  <tr><td class="theme-weak"></td><td><code>--theme-color-weak</code></td></tr>
  <tr><td class="theme-strong"></td><td><code>--theme-color-strong</code></td></tr>
</table>

<hr/>
<h1>Example Tags</h1>
<hr/>

<h1>Heading One</h1>
<h2>Heading Two</h2>
<h3>Heading Three</h3>
<h4>Heading Four</h4>

<blockquote>
  Block Quote, Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
</blockquote>

<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>

<pre>

     .--------.
     | Hello. |
     '--. .--'     /`-._
         \|       /_,.._`:-
          \   ,.-'  ,   `-:..-')
             : o ):';      _  {
              `-._ `'__,.-'\`-.)
                  `\\  \,.-'`

</pre>

```go
package main

import "fmt"

func main() {
    messages := make(chan string)

    go func() { messages <- "ping" }()

    msg := <-messages
    fmt.Println(msg)
}
```

<hr/>
<h1>Listing View</h1>
<hr/>

<div class="listing">
  <div>
    <h2><a href="/2023/04/05/moving-from-wordpress-to-11ty/">Moving from WordPress to 11ty</a></h2>
    <time datetime="2023-04-05T00:00:00.000Z">Apr 5, 2023</time>
    <p>After 17 years with WordPress, it's time for something different. Here's how I migrated 500 posts to a static site generator.</p>
  </div>
  <div>
    <h2><a href="/2023/03/26/a-peek-inside-pinentry/">A peek inside pinentry</a></h2>
    <time datetime="2023-03-27T00:00:00.000Z">Mar 27, 2023</time>
    <p>I interact with pinenty daily, but I don't really understand it. This post dives into how it is invoked and can be used outside of GPG for your own projects.</p>
  </div>
</div>

<div class="listing">
  <nav>
    <ol>
      <li>« Previous</li>
      <li><a href="#" rel="next">Next »</a></li>
    </ol>
  </nav>
</div>


<hr/>
<h1>Article</h1>
<hr/>

<article>
  <h1>Sup</h1>
</article>
