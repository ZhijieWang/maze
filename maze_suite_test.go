package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMaze(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Maze Suite")
}
