module github.com/DisgoOrg/disgo/_examples/test

go 1.16

replace (
	github.com/DisgoOrg/disgo => ../../
	github.com/DisgoOrg/log => ../../../log
)

require (
	github.com/DisgoOrg/disgo v0.5.6
	github.com/DisgoOrg/log v1.1.0
	github.com/PaesslerAG/gval v1.1.1
)
