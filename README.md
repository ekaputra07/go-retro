# go-retro

[![Build](https://github.com/ekaputra07/go-retro/actions/workflows/go.yaml/badge.svg)](https://github.com/ekaputra07/go-retro/actions/workflows/go.yaml)
[![Release](https://github.com/ekaputra07/go-retro/actions/workflows/release.yaml/badge.svg)](https://github.com/ekaputra07/go-retro/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ekaputra07/go-retro)](https://goreportcard.com/report/github.com/ekaputra07/go-retro)

Minimalist retro board for happy teams 😉

> In version 2, [NATS](https://nats.io/) is used as communication and persistence layer to enable horizontal scaling (real-time communication works across multiple instances). If you prefer simpler single-instance setup, you can still use the [version 1](https://github.com/ekaputra07/go-retro/releases/tag/1.0.0).

![screenshot](https://github.com/ekaputra07/go-retro/blob/main/screenshot.png)

> **What the heck is retro board?** The Retrospective Board usually used as a tool during sprint restrospective meeting by a team to carry out a "lessons learned" or "how do you feel?" assessment of the sprint. Usually helds at the end of the sprint period.

This project is heavily inspired by https://www.dragondropcards.com, here I'm trying to replicate the functionalities as much as possible while also adding some cool features that I think would be useful (timer, online users, standup, etc).

👉Try it: https://go-retro.fly.dev

### Features
- [x] Create/Update/Delete board columns
- [x] Add/Update/Delete cards
- [x] Move cards to other column
- [x] See number of online users
- [x] A timer to allow users fill-in the board with cards within a specified time limit
- [x] React to a card (thumbs up or emoji?)
- [x] Display user name on who's online list
- [x] Standup feature (shuffle users and display who's turn to speak)
- [x] Persistence layer, powered by NATS KV (expires after 24 hours)
- [ ] Group similar cards

### Development

Therefore during development, you need to have a NATS server running locally.

1. Install dependencies
    ```bash
    # install both backend and frontend dependencies
    make setup
    ```

3. Run the app with Docker compose
    ```bash
    # starts both NATS server and the app
    make compose
    ```

### Docker images

```bash
# version 2
docker pull ekaputra07/goretro:latest
docker pull ekaputra07/goretro:2.x.x

# version 1
docker pull ekaputra07/goretro:1.0.0
```

### Credits

As I'm not a UI guy, I steal and modify the board HTML from https://github.com/mithicher/tasksgram by [@mithicher](https://github.com/mithicher/tasksgram) which perfectly suite my needs for a clean design powered by AlpineJS and Tailwind CSS.
