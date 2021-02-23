package ledger

import (
	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/privacyenabledstate"
)

// Compare - Compares two ledger snapshots and returns the first divergence
func Compare(snapshotPath1 string, snapshotPath2 string) (result string, err error) {
	result = "Comparing " + snapshotPath1 + " and " + snapshotPath2 + "\n\n"

	// Create snapshot readers to read both snapshots
	snapshotReader1, err := privacyenabledstate.NewSnapshotReader(snapshotPath1, "public_state.data", "public_state.metadata")
	if err != nil {
		return result + "newSnapshotReader error, snapshot 1", err
	}
	snapshotReader2, err := privacyenabledstate.NewSnapshotReader(snapshotPath2, "public_state.data", "public_state.metadata")
	if err != nil {
		return result + "newSnapshotReader error, snapshot 2", err
	}

	// Read each snapshot record and find divergences
	_, snapshotRecord1, err1 := snapshotReader1.Next()
	_, snapshotRecord2, err2 := snapshotReader2.Next()

	for snapshotRecord1 != nil && snapshotRecord2 != nil {
		if err1 != nil {
			return result + "snapshotReader1 error", err1
		}
		if err2 != nil {
			return result + "snapshotReader2 error", err2
		}

		// Determine the difference in records by comparing keys
		res := bytes.Compare(snapshotRecord1.Key, snapshotRecord2.Key)

		if res == 0 {
			if !(proto.Equal(snapshotRecord1, snapshotRecord2)) {
				// Keys are the same but records are different
				result += "Difference in snapshot records\n"
				result += snapshotRecord1.String() + "\n" + snapshotRecord2.String() + "\n\n"
				// Add difference to JSON file
			} else {
				// Matching records
				result += "Match\n"
				result += snapshotRecord1.String() + "\n" + snapshotRecord2.String() + "\n\n"
			}
			_, snapshotRecord1, err1 = snapshotReader1.Next()
			_, snapshotRecord2, err2 = snapshotReader2.Next()
		} else if res == 1 {
			// Snapshot 1 is missing a record
			result += "Snapshot 1 has a missing record\n"
			result += snapshotRecord2.String() + "\n\n"
			_, snapshotRecord2, err2 = snapshotReader2.Next()
		} else if res == -1 {
			// Snapshot 2 is missing a record
			result += "Snapshot 2 has a missing record\n"
			result += snapshotRecord1.String() + "\n\n"
			_, snapshotRecord1, err1 = snapshotReader1.Next()
		}
	}

	// Check for tailing records
	if snapshotRecord1 != nil {
		// Snapshot 2 is missing a record
		for snapshotRecord1 != nil {
			result += "Snapshot 2 has a missing record\n"
			result += snapshotRecord1.String() + "\n\n"
			_, snapshotRecord1, err1 = snapshotReader1.Next()
		}
	} else if snapshotRecord2 != nil {
		// Snapshot 1 is missing a record
		for snapshotRecord2 != nil {
			result += "Snapshot 1 has a missing record\n"
			result += snapshotRecord2.String() + "\n\n"
			_, snapshotRecord2, err2 = snapshotReader2.Next()
		}
	}

	return result, nil
}
