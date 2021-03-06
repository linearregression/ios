package consensus

import (
	"github.com/golang/glog"
	"github.com/heidi-ann/ios/msgs"
	"reflect"
)

func doCoordination(view int, startIndex int, endIndex int, entries []msgs.Entry, io *msgs.Io, config Config, prepare bool) bool {
	// PHASE 2: prepare
	if prepare {

		// check that committed is not set
		for i := 0; i < endIndex-startIndex; i++ {
			entries[i].Committed = false
		}

		prepare := msgs.PrepareRequest{config.ID, view, startIndex, endIndex, entries}
		glog.V(1).Info("Starting prepare phase", prepare)
		io.OutgoingBroadcast.Requests.Prepare <- prepare

		// collect responses
		glog.V(1).Info("Waiting for ", config.Quorum.ReplicateSize, " prepare responses")
		for replied := make([]bool, config.N); !config.Quorum.checkReplicationQuorum(replied); {
			msg := <-io.Incoming.Responses.Prepare
			// check msg replies to the msg we just sent
			if reflect.DeepEqual(msg.Request, prepare) {
				glog.V(1).Info("Received ", msg)
				if !msg.Response.Success {
					glog.Warning("Coordinator is stepping down")
					return false
				}
				replied[msg.Response.SenderID] = true
				glog.V(1).Info("Successful response received, waiting for more")
			}
		}
	}

	// PHASE 3: commit
	// set committed so requests will be applied to state machines
	for i := 0; i < endIndex-startIndex; i++ {
		entries[i].Committed = true
	}
	// dispatch commit requests to all
	commit := msgs.CommitRequest{config.ID, startIndex, endIndex, entries}
	glog.V(1).Info("Starting commit phase", commit)
	io.OutgoingBroadcast.Requests.Commit <- commit

	// TODO: handle replies properly
	go func() {
		for replied := make([]bool, config.N); !config.Quorum.checkReplicationQuorum(replied); {
			msg := <-io.Incoming.Responses.Commit
			// check msg replies to the msg we just sent
			if reflect.DeepEqual(msg.Request, commit) {
				glog.V(1).Info("Received ", msg)
				replied[msg.Response.SenderID] = true
			}
		}
	}()

	return true
}

// runCoordinator eturns true if successful
func runCoordinator(state *state, io *msgs.Io, config Config) {

	for {
		glog.V(1).Info("Coordinator is ready to handle request")
		req := <-io.Incoming.Requests.Coordinate
		success := doCoordination(req.View, req.StartIndex, req.EndIndex, req.Entries, io, config, req.Prepare)
		// TODO: check view
		reply := msgs.CoordinateResponse{config.ID, success}
		io.OutgoingUnicast[req.SenderID].Responses.Coordinate <- msgs.Coordinate{req, reply}
		glog.V(1).Info("Coordinator is finished handling request")
		// TOD0: handle failure
	}
}
