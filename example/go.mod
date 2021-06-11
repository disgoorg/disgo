module github.com/DisgoOrg/disgo/example

go 1.16

replace (
	github.com/DisgoOrg/disgo => ../
)

require (
	github.com/DisgoOrg/disgo v0.3.2
	github.com/PaesslerAG/gval v1.1.1
	github.com/sirupsen/logrus v1.8.1
)
