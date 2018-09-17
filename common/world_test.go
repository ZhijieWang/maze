package common

import (
	"testing"
)
func TestCanModifyRobot(t *testing.T){
	g:=CreateWorld(3, false)
	r := g.GetRobots()
	r[0].location=r[1].location
	if r[0].location != g.GetRobots()[0].location{
		t.Errorf("Expect the robots returned to be modifiable\n")
	}
}
