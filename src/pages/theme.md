---
eleventyExcludeFromCollections: true
noindex: true
permalink: /theme/
title: Theme
---

I've always preferred Monokai/Molokai and it's many derivative color schemes in my editor. At my last redesign, I decided to make my website available in six classic Monokai colors, as well as light and dark mode versions.  The default is blue, but pick whatever combination you like best.

<noscript>

## Oh no!

My theme and color scheme switching requires JavaScript to work.

</noscript>

## Color

<ul>
    <li><a href="#" class="theme-switcher" data-theme="blue" aria-label="Switch to blue color theme">Blue</a></li>
    <li><a href="#" class="theme-switcher" data-theme="pink" aria-label="Switch to pink color theme">Pink</a></li>
    <li><a href="#" class="theme-switcher" data-theme="yellow" aria-label="Switch to yellow color theme">Yellow</a></li>
    <li><a href="#" class="theme-switcher" data-theme="green" aria-label="Switch to green color theme">Green</a></li>
    <li><a href="#" class="theme-switcher" data-theme="orange" aria-label="Switch to orange color theme">Orange</a></li>
    <li><a href="#" class="theme-switcher" data-theme="purple" aria-label="Switch to purple color theme">Purple</a></li>
</ul>

## Mode

<ul>
    <li><a href="#" class="scheme-switcher" data-scheme="light" aria-label="Switch to light scheme">Light</a></li>
    <li><a href="#" class="scheme-switcher" data-scheme="dark" aria-label="Switch to dark scheme">Dark</a></li>
</ul>


<style>
ul {
    margin: 0;
    padding: 0;
}
li {
    margin: 0;
    list-style-type: none;
}
li a {
    display: block;
    padding: 0.5em;
    text-align: center;
    text-decoration: none;
}

a[data-scheme="light"] {
    background-color: white;
    color: black;
}
a[data-scheme="dark"] {
    background-color: black;
    color: white;
}

a[data-theme="blue"] {
    background-color: var(--monokai-blue);
    color: black;
}
a[data-theme="pink"] {
    background-color: var(--monokai-pink);
    color: white;
}
a[data-theme="yellow"] {
    background-color: var(--monokai-yellow);
    color: black;
}
a[data-theme="green"] {
    background-color: var(--monokai-green);
    color: black;
}
a[data-theme="orange"] {
    background-color: var(--monokai-orange);
    color: black;
}
a[data-theme="purple"] {
    background-color: var(--monokai-purple);
    color: black;
}
</style>

<script>
  document.querySelectorAll('a.theme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['theme'] = e.target.dataset['theme'];
    window.localStorage.setItem('theme', e.target.dataset['theme']);
  }));
  document.querySelectorAll('a.scheme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['scheme'] = e.target.dataset['scheme'];
    window.localStorage.setItem('scheme', e.target.dataset['scheme']);
  }));
</script>
