<h1 align="center">
  Meetup Assistant (wip)
</h1>

<p align="center">
  🤖 Automating Meetup and other events tasks
</p>

## Features

1. Schedule emails to hosts and speakers
2. Schedule twitter to hosts and speakers

## Requirements

1. Notion as the data backend. Feel free to open PRs with more implementations, but even Google Drive would do.
2. Runtime which can persist files, as the state and idempotency storage will rely in an embedded database
3. Mailgun API key for email tasks
4. Twitter API key for twitter tasks
5. Github token for template storage using Github

## Todo

- Implement noop twitter
- Implement real twitter
- Implement and test state machine

## Technologies I want to try

- slog
- wire
- cobra
