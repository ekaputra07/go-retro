# go-retro

Minimalist, real-time and open-source retro board for happy teams 😉

> **What is retro board?** The Retrospective Board enables you to carry out a 'lessons learned' assessment of a project, or even a phase/sprint within a project. Usually helds at the end of the sprint.

This project is inspired by https://www.dragondropcards.com and here I try to replicate the functionalities as much as possible while also adding some cool features that I think would be useful.

👉Try it: https://go-retro.fly.dev

### Features
- [x] Create/Update/Delete board columns
- [x] Add/Update/Delete cards
- [x] Move cards to other column
- [x] See number of online users
- [ ] Group cards
- [ ] A timer to allow users fill-in the board with cards within a specified time limit
- [ ] Round robin user selection, user are selected randomly after each other during the session
- [ ] Persistence layer, currently all data stored in memory

### Development
The project is Dockerized so you could simply run `docker compose up` and the application should be accessible via `http://localhost:8080`.

As I'm not an UI guy, I'm proudly steal the board design and HTML from https://github.com/mithicher/tasksgram which perfectly suite my needs for a clean design powered by AlpineJS and Tailwind CSS.
