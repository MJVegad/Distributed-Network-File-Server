package main

import (
	"fmt"
	"testing"
	"reflect"
)

func ExpectActions(t *testing.T, a []interface{}, b []interface{}) {
	if (len(a) != len(b)) {
		t.Error(fmt.Sprintf("Expected %v elements, found %v elements", len(b), len(a))) // t.Error is visible when running `go test -verbose`
	} else {
		
		if ( len(CompareActions(a, b)) == 0 && len(CompareActions(b, a)) == 0 ) {
		e1 := 1
		e2 := 1
		i := 0
		for ;i<len(a);i++ {
			e1 = 0
			for j:=0;j<len(b);j++ {
				if reflect.DeepEqual(a[i], b[j]) {
					e1 = 1
					break;
				}
			}
			e2 = e2 & e1
			if e2 == 0 {
				break;
			}
		}
		if e2 == 0 {
			t.Fatal(fmt.Sprintf("%v not found in the expected list of actions", a[i]))
		}
		} else {
			t.Fatal("No. of actions mismatched with expected actions list")
		}
		
	}
}

func CompareActions(X, Y []interface{}) []interface{} {
    m := make(map[reflect.Type]int)

    for _, y := range Y {
        m[reflect.TypeOf(y)]++
    }

    var ret []interface{}
    for _, x := range X {
        if m[reflect.TypeOf(x)] > 0 {
            m[reflect.TypeOf(x)]--
            continue
        }
        ret = append(ret, x)
    }

    return ret
}
 //&& a.currentState==b.currentState && a.currentTerm==b.currentTerm && a.majority==b.majority && a.novotes==b.novotes && a.serverId==b.serverId && a.totalvotes==b.totalvotes && a.votedFor==b.votedFor && reflect.DeepEqual(a.log, b.log))
func ExpectStateMachine(t *testing.T, a *StateMachine, b *StateMachine) {
	if (a.commitIndex != b.commitIndex){
		t.Fatal(fmt.Sprintf("Expected %v commitIndex, found %v commitIndex", b.commitIndex, a.commitIndex)) // t.Error is visible when running `go test -verbose`
	} else if (a.currentTerm!=b.currentTerm) {
			t.Fatal(fmt.Sprintf("Expected %v currentTerm, found %v currentTerm", b.currentTerm, a.currentTerm)) // t.Error is visible when running `go test -verbose`
	} else if (a.majority!=b.majority) {
			t.Fatal(fmt.Sprintf("Expected %v majority, found %v majority", b.majority, a.majority)) // t.Error is visible when running `go test -verbose`
	} else if (a.serverId!=b.serverId) {
			t.Fatal(fmt.Sprintf("Expected %v serverId, found %v serverId", b.serverId, a.serverId)) // t.Error is visible when running `go test -verbose`
	} else if (a.novotes!=b.novotes) {
			t.Fatal(fmt.Sprintf("Expected %v novotes, found %v novotes", b.novotes, a.novotes)) // t.Error is visible when running `go test -verbose`
	} else if (a.totalvotes!=b.totalvotes) {
			t.Fatal(fmt.Sprintf("Expected %v totalvotes, found %v totalvotes", b.totalvotes, a.totalvotes)) // t.Error is visible when running `go test -verbose`
	} else if (a.votedFor!=b.votedFor) {
			t.Fatal(fmt.Sprintf("Expected %v votedFor, found %v votedFor", b.votedFor, a.votedFor)) // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.log, b.log)==false) {
			t.Fatal("Log not as expected.") // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.nextIndex, b.nextIndex)==false) {
			t.Fatal("nextIndex not as expected") // t.Error is visible when running `go test -verbose`
	} else if (reflect.DeepEqual(a.matchIndex, b.matchIndex)==false) {
			t.Fatal("matchIndex not as expected") // t.Error is visible when running `go test -verbose`
	} else if (a.currentState!=b.currentState) {
			t.Fatal(fmt.Sprintf("Expected %v state, found %v state", b.currentState, a.currentState)) // t.Error is visible when running `go test -verbose`
	} 
}