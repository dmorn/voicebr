[![GoDoc](https://godoc.org/github.com/jecoz/voicebr?status.svg)](https://godoc.org/github.com/jecoz/voicebr)
[![Build Status](https://travis-ci.org/jecoz/voicebr.svg?branch=master)](https://travis-ci.org/jecoz/voicebr)

# voicebr
Broadcasts phone calls to a list of contacts, defined into a static file. `voicebr`
has to be registered to a "voice application" into the nexmo platform, and also a
phone number as to be provided, together with the private key used to sign nexmo's
tokens. All these data can be retrived from nexmo's dashboard.
`voicebr` then spawns a web server that will be contacted from nexmo when a phone
call is made to the registered number. The voice of the caller is then recorded,
saved locally, and reproduced into an outbound call made to each contact managed
by `voicebr`.
