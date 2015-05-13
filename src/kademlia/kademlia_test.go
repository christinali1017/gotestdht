package kademlia

import (
	"math/rand"
	"net"
	"strconv"
	"testing"
)

func CreateIdForTest(id string) (ret ID) {
	if len(id) > 160 {
		id = id[0:160]
	}
	for i := 0; i < len(id); i++ {
		ret[i] = id[i]
	}
	return 
}


func StringToIpPort(laddr string) (ip net.IP, port uint16, err error) {
	hostString, portString, err := net.SplitHostPort(laddr)
	if err != nil {
		return
	}
	ipStr, err := net.LookupHost(hostString)
	if err != nil {
		return
	}
	for i := 0; i < len(ipStr); i++ {
		ip = net.ParseIP(ipStr[i])
		if ip.To4() != nil {
			break
		}
	}
	portInt, err := strconv.Atoi(portString)
	port = uint16(portInt)
	return
}

func TestPing(t *testing.T) {
	instance1 := NewKademlia(CreateIdForTest(string(1)), "localhost:7890")
	instance2 := NewKademlia(CreateIdForTest(string(2)), "localhost:7891")
	host2, port2, _ := StringToIpPort("localhost:7891")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	contact1, err := instance2.FindContact(instance1.NodeID)
	if err != nil {
		t.Error("Instance 1's contact not found in Instance 2's contact list")
		return
	}
	if contact1.NodeID != instance1.NodeID {
		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
	}
	if contact2.NodeID != instance2.NodeID {
		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
	}
	return
}

func TestFindNode(t *testing.T) {
	instance1 := NewKademlia(CreateIdForTest(string(1)), "localhost:7892")
	instance2 := NewKademlia(CreateIdForTest(string(2)), "localhost:7893")
	instance3 := NewKademlia(CreateIdForTest(string(3)), "localhost:7894")
	host2, port2, _ := StringToIpPort("localhost:7893")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	contact1, err := instance2.FindContact(instance1.NodeID)
	if err != nil {
		t.Error("Instance 1's contact not found in Instance 2's contact list")
		return
	}
	if contact1.NodeID != instance1.NodeID {
		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
	}
	if contact2.NodeID != instance2.NodeID {
		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
	}
	instance3.DoPing(host2, port2)
	instance1ID := instance1.SelfContact.NodeID
	instance2ID := instance2.SelfContact.NodeID
	instance3ID := instance3.SelfContact.NodeID
	contact, err := instance1.FindContact(instance2ID)
	if err != nil {
		t.Error("ERR: Unable to find contact with node ID")
		return
	}
	var res []Contact
	res = instance2.FindClosestContacts(instance3ID, instance1ID)
	resstring := instance2.ContactsToString(res)
	response := instance1.DoFindNode(contact, instance3ID)
	if response != "ok, result is: "+resstring {
		t.Error("Node in Instance2 are stored incorrectly")
	}
	return
}

func TestStore(t *testing.T) {
	instance1 := NewKademlia(CreateIdForTest(string(1)), "localhost:7895")
	instance2 := NewKademlia(CreateIdForTest(string(2)), "localhost:7896")
	host2, port2, _ := StringToIpPort("localhost:7896")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	contact1, err := instance2.FindContact(instance1.NodeID)
	if err != nil {
		t.Error("Instance 1's contact not found in Instance 2's contact list")
		return
	}
	if contact1.NodeID != instance1.NodeID {
		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
	}
	if contact2.NodeID != instance2.NodeID {
		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
	}
	instance1ID := instance1.SelfContact.NodeID
	instance2ID := instance2.SelfContact.NodeID
	contact, err := instance1.FindContact(instance2ID)
	if err != nil {
		t.Error("ERR: Unable to find contact with node ID")
		return
	}
	svalue := strconv.Itoa(int(rand.Intn(256)))
	value := []byte(svalue)
	instance1.DoStore(contact, instance1ID, value)
	response := instance2.LocalFindValue(instance1ID)
	if response != "OK:"+string(value[:]) {
		t.Error("Value in Instance2 are stored incorrectly")
	}
	return
}

func TestFindValue(t *testing.T) {
	instance1 := NewKademlia(CreateIdForTest(string(1)), "localhost:7897")
	instance2 := NewKademlia(CreateIdForTest(string(2)), "localhost:7898")
	instance3 := NewKademlia(CreateIdForTest(string(3)), "localhost:7899")
	host2, port2, _ := StringToIpPort("localhost:7898")
	instance1.DoPing(host2, port2)
	contact2, err := instance1.FindContact(instance2.NodeID)
	if err != nil {
		t.Error("Instance 2's contact not found in Instance 1's contact list")
		return
	}
	contact1, err := instance2.FindContact(instance1.NodeID)
	if err != nil {
		t.Error("Instance 1's contact not found in Instance 2's contact list")
		return
	}
	if contact1.NodeID != instance1.NodeID {
		t.Error("Instance 1 ID incorrectly stored in Instance 2's contact list")
	}
	if contact2.NodeID != instance2.NodeID {
		t.Error("Instance 2 ID incorrectly stored in Instance 1's contact list")
	}
	instance3.DoPing(host2, port2)
	instance1ID := instance1.SelfContact.NodeID
	instance2ID := instance2.SelfContact.NodeID
	instance3ID := instance3.SelfContact.NodeID
	contact, err := instance1.FindContact(instance2ID)
	if err != nil {
		t.Error("ERR: Unable to find contact with node ID")
		return
	}
	svalue := strconv.Itoa(int(rand.Intn(256)))
	value := []byte(svalue)
	instance1.DoStore(contact, instance1ID, value)
	response := instance3.DoFindValue(contact, instance1ID)
	if response != "ok, result is: "+string(value[:]) && response != "No Record" {
		t.Error("Value in Instance2 are stored incorrectly")
	}
	responsenode := instance3.DoFindNode(contact, instance3ID)
	responsevalue := instance3.DoFindValue(contact, instance3ID)
	if responsenode != responsevalue {
		t.Error("Node in Instance2 are stored incorrectly")
	}
	return
}


func TestIterativeFindNode(t *testing.T) {
	
	for i := 0; i < 2000; i++ {

	}
	return
}

