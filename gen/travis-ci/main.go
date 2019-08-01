package main

import (
	"flag"
	"log"
	"os"

	"text/template"
)

func main() {
	t := template.Must(template.New("travis-ci").Parse(travisGo))

	data := jobData{}
	flag.StringVar(&data.Repository, "repository", "github.com/username/repo", "Name of the repository e.g. github.com/hello-world")
	flag.StringVar(&data.Docker, "docker", "github.com/username/repo", "Name of the docker container with registry e.g. github.com/hello-world")
	flag.Parse()

	err := t.Execute(os.Stdout, data)
	if err != nil {
		log.Println("executing template:", err)
	}
}

const usage = `
go run github.com/kastenpotential/mage-utils/gen/travis-ci <repository> <docker>
`

type jobData struct {
	Repository string
	Docker     string
}

const travisGo = `
language: go

go:
    - 1.12.x

sudo: required

services:
    - docker

install:
    - go get -u -d {{.Repository}}
    - cd $GOPATH/src/github.com/magefile/mage
    - go run bootstrap.go
    - cd $TRAVIS_BUILD_DIR

script: mage -v buildDockerContainer

before_deploy:
    - docker login -u "$REGISTRY_USER" -p "$REGISTRY_PASS"

deploy:
    - provider: script
      skip_cleanup: true
      script: mage -v deployLatest
      on:
          branch: master
    - provider: script
      skip_cleanup: true
      script: mage -v deployVersion
      on:
          tags: true

notifications:  
    email:  
        on_success: always  
        on_failure: always
`
