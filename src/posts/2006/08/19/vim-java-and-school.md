---
category:
- Geek
creator: admin
date: 2006-08-20
permalink: /2006/08/19/vim-java-and-school/
tags:
- Linux
- School
title: VIM, Java, and School
type: post
wp_id: "10"
---
I spent a little time setting up VIM today, and thought I'd share my .vimrc.  Most of these configs are from the beastly , but well documented, file at [http://www.amix.dk/vim/vimrc.html](http://www.amix.dk/vim/vimrc.html). Mine's a bit simpler, but it does what I want from it.  Always been a VIM fan.

```
" John's VIM Config
" www.velvetcache.org
" Turn off VI compat
set nocompatible
"Set shell to be bash
set shell=bash
set history=500
" Lets use the mouse!
set mouse=a
" This is the most important line here...
syntax enable
" Set font (GUI)
set gfn=Monospace 10
set encoding=utf-8
hi MatchParen
guifg=#000000
guibg=#D0D090
set ruler
" Highlight Matching parens.
set showmatch
set mat=2
set hlsearch
" Kill all those crummy swap/backup files
set nobackup
set nowb
set noswapfile
set noar
" Make tabs only 2 spaces
set tabstop=2
" Auto-indent / Smart-indent
set ai
set si
" No more \a bells!
set visualbell
```

On other fronts, I start school again on Monday, which in itself partially prompted the VIM editing, since I wanted some highlighting on Java files andwasn't getting it.  This should be an interesting semester, since between my Java class (hehe..), data structures and my job I'll be juggling the syntax for three or four languages at once.  Good times.
