# mblaze-ui

A minimal TUI for the [mblaze][mblaze github] email client.

![Screenshot of mblaze-ui with a list of email summaries.](https://gist.github.com/nmeum/ddb6ddbe84d9ef5bdabd5a81219c93b2/raw/8a56073afb3b1d3d5019e09f2af43f59c245ace4/mblaze-ui.png)

## About

mblaze-ui is [tcell][tcell github]-based terminal user interface for the
[mblaze][mblaze github] email client. Similar to existing Unix utilities
from `mblaze(7)`, it operates on the current message sequence as set by
`mseq(1)`. For each mail of the sequence, it prints a one line summary
similar to `mscan(1)`. Using the arrow keys, a mail from the sequence
can be selected and manipulated using the key bindings described below.
Conceptually, mblaze-ui is therefore similar to `mless(1)` but offers
a pager-independent interface.

## Installation

This software requires at least Go 1.21.0. If such a Go version
is available on your system, install using `go install` as follows:

    $ go install github.com/nmeum/mblaze-ui@latest

## Usage Example

Set a message sequence using `mseq(1)`, for example

    $ mlist -s ~/mail/INBOX | msort -r -d | mseq -S

will set the message sequence to all unseen emails in your INBOX and
will sort these messages ascending by date. This sequence can then
be viewed and modified in a TUI environment using:

    $ mblaze-ui

In the TUI, a specific message can be selected using the arrow keys.
The state of this message can be modified through the key bindings
described in the following section.

## Key Bindings

The following key bindings are currently implemented:

* `Esc` / `Ctrl-C` / `q`: Exit mblaze-ui
* `Ctrl-L`: Redraw the screen
* `Enter`: View the currently selected email using `mshow(1)`
* `Up` / `Down`: Select the next/previous email
* `PageUp` / `PageDown`: Show the next/previous page of mails
* `f` / `F`: Mark the email as unflagged/flagged using `mflag(1)`
* `s` / `S`: Mark the email as unseen/seen using `mflag(1)`
* `t` / `T`: Mark the mail as untrashed/trashed using `mflag(1)`
* `d`: Delete the currently selected email
* `r`: Compose a reply for the selected email

## License

This program is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by the
Free Software Foundation, either version 3 of the License, or (at your
option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
Public License for more details.

You should have received a copy of the GNU General Public License along
with this program. If not, see <https://www.gnu.org/licenses/>.

[mblaze github]: https://github.com/leahneukirchen/mblaze
[tcell github]: https://github.com/gdamore/tcell
