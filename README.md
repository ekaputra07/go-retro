# go-retro

A minimalist yet fun retro board for happy teams 😉

![screenshot](https://github.com/ekaputra07/go-retro/blob/main/screenshot.png)

> **What the heck is retro board?** The Retrospective Board usually used as a tool during sprint restrospective meeting by a team to carry out a "lessons learned" or "how do you feel?" assessment of the sprint. Usually helds at the end of the sprint period.

This project is heavily inspired by https://www.dragondropcards.com, here I'm trying to replicate the functionalities as much as possible while also adding some cool features that I think would be useful (timer, online users, etc).

👉Try it: https://go-retro.fly.dev

### Features
- [x] Create/Update/Delete board columns
- [x] Add/Update/Delete cards
- [x] Move cards to other column
- [x] See number of online users
- [x] A timer to allow users fill-in the board with cards within a specified time limit
- [x] React to a card (thumbs up or emoji?)
- [x] Display user name on who's online list
- [ ] Group similar cards
- [ ] Persistence layer, currently all data stored in memory

### Development
The project is Dockerized so you could simply run `docker compose up` and the application should be accessible via `http://localhost:8080`. Or if you're like me, I simply do `go run .` during development.

As I'm not a UI guy, I steal and modify the board HTML from https://github.com/mithicher/tasksgram by [@mithicher](https://github.com/mithicher/tasksgram) which perfectly suite my needs for a clean design powered by AlpineJS and Tailwind CSS.
