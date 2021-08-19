# goplist
[![](https://github.com/ZhongXiLu/goplist/workflows/Go/badge.svg)](https://github.com/ZhongXiLu/goplist/actions?query=workflow%3A%22go%22)

Tool that helps convert macOS preferences to bash commands. I initially wanted to save some of my macOS preferences as bash commands so I can use them in my dotfiles and whenever I would have to setup a new macOS environment, I can easily import my macOS preferences.

# How to Install

// TODO

# How to Run

- Simply run `goplist` in any directory and it will record any changes to your preferences. Note that it may take some time for the OS to write your changes to disk. If you want to keep your current preference, simply toggle (with a small delay in between them) and `goplist` will remember your last choice. After you are done, simply press <kbd>ctrl + c</kbd> and `goplist` will write your preferences to file (by default `macos.sh`).
- Run `goplist -h` for help.

# Demo

- Changing system preferences:

<p align="center">
  <img src="https://user-images.githubusercontent.com/25816683/130079872-64b27e34-e488-4450-a96c-5ee960a74e82.gif"/>
</p>

- This also works for some other applications, e.g. iTerm2:

<p align="center">
  <img src="https://user-images.githubusercontent.com/25816683/130079877-e1a48cbf-ea00-45a8-9380-c55d492360d3.gif"/>
</p>

