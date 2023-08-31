# mblaze-ui

A [mutt][mutt web]-like TUI for the [mblaze][mblaze github] email client.

## Installation

Install using `go install` as follows

    $ go install github.com/nmeum/mblaze-ui@latest

## Key Bindings

The following key bindings are currently implemented:

* `Esc` / `Ctrl-C`: Exit mblaze-ui
* `Ctrl-L`: Redraw the screen
* `Enter`: View the currently selected email using `mshow(1)`
* `Up` / `Down`: Select the next/previous email
* `PageUp` / `PageDown`: Show the next/previous page of mails
* `s` / `S`: Mark the email as seen/unseen using `mflag(1)`
* `f` / `F`: Mark the email as flagged/unflagged using `mflag(1)`
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

[mutt web]: http://www.mutt.org
[mblaze github]: https://github.com/leahneukirchen/mblaze
