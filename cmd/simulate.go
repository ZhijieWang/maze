// Copyright © 2018 Zhijie (Bill) Wang <wangzhijiebill@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zhijiewang/maze/common"
	"time"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulate called")
		start := time.Now()
		g := common.CreateWorld(NumRobots, Concurrency)
		for i := 0; i < Iterations; i++ {
			g.Simulate(common.RandMove, common.GraphReWeightByRadiation)
		}
		elapsed := time.Since(start)
		fmt.Printf("Simulation took %s for %v iterations \n", elapsed, Iterations)
	},
}

// Concurrency is the flag for marking execution mode
var Concurrency bool

// Iterations is the flag for how many iterations to run
var Iterations int

// NumRobots is the flag for how many robots to assign on the network
var NumRobots int

func init() {
	rootCmd.AddCommand(simulateCmd)
	simulateCmd.Flags().BoolVarP(&Concurrency, "concurrency", "c", false, "Setting for running the simulation in concurrent mode, which consumes more system resouces for speed up)")
	simulateCmd.Flags().IntVar(&Iterations, "i", 100, "Setting for number of iterations in the simulation")
	simulateCmd.Flags().IntVar(&NumRobots, "n", 3, "Setting for number of robots to spawn on the ground")
}
