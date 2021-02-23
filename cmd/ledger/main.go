package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric/internal/ledger"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("ledger", "Ledger Utility Tool")

	compare      = app.Command("compare", "Compare two ledgers via their snapshots.")
	comparePath1 = compare.Arg("path1", "File path to first ledger snapshot.").Required().String()
	comparePath2 = compare.Arg("path2", "File path to second ledger snapshot").Required().String()

	troubleshoot = app.Command("troubleshoot", "Identify potentially divergent transactions.")

	args = os.Args[1:]
)

func main() {

	kingpin.Version("0.0.1")

	command, err := app.Parse(args)

	if err != nil {
		kingpin.Fatalf("parsing arguments: %s. Try --help", err)
		return
	}

	switch command {

	case compare.FullCommand():

		result, err := ledger.Compare(*comparePath1, *comparePath2)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result)

	case troubleshoot.FullCommand():

		fmt.Println("Command TBD")

	}

}
