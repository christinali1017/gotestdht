package kademlia

import (
	"math/rand"
	"net"
	"strconv"
	"testing"
	// "encoding/json"
	// "strings"
	// "io"
	"sort"
	"fmt"
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
	fmt.Println(".........Begin test find node......")
	numberOfNodes := 300
	numberOfContactsPerNode := 30
	instances := make([]Kademlia, numberOfNodes)
	instancesAddr := make([]string, numberOfNodes)


	//create 100 kademlia instance
	fmt.Println("........Create instances......")
	for i := 0; i < numberOfNodes; i++ {
		port := i + 8000
		address := "localhost:" + strconv.Itoa(port)

		// fmt.Println("port is " + address)
		instancesAddr[i] = address
		instances[i] = *NewKademlia(CreateIdForTest(string(i)), address)
		//instances[i] = *NewKademlia(CreateIdForTest(strconv.Itoa(i)), address)
	}

	fmt.Println(".........ping ......")

	for i := 0; i < numberOfNodes; i++ {
		address := instancesAddr[i]
		host, port, _ := StringToIpPort(address)
		start := i - numberOfContactsPerNode/2;
		end := i + numberOfContactsPerNode/2;
		if i >= numberOfContactsPerNode/2 && i <= numberOfNodes - numberOfContactsPerNode/2 {
			for j := start; j < end; j ++ {
				instances[j].DoPing(host, port)
			}
		} else {
			if (i < numberOfContactsPerNode/2) {
				for j := 0; j < numberOfContactsPerNode; j++ {
					instances[j].DoPing(host, port)
				}
			} else if i > numberOfNodes - numberOfContactsPerNode/2 {
				for j := numberOfNodes - numberOfContactsPerNode; j < numberOfNodes; j++ {
					instances[j].DoPing(host, port)
				}
			}
		}

	}

	fmt.Println("..............Check if contacts are stored......")

	for i := 0; i < numberOfNodes; i++ {
		instance := instances[i]
		start := i - numberOfContactsPerNode/2;
		end := i + numberOfContactsPerNode/2;
		if i >= numberOfContactsPerNode/2 && i <= numberOfNodes - numberOfContactsPerNode/2 {
			for j := start; j < end; j ++ {
				contact, err := instances[j].FindContact(instance.NodeID)
				if err != nil {
					t.Error("Instance" + string(i) + "'s contact not found in Instance" + string(j) + "'s contact list")
					return
				}

				if !contact.NodeID.Equals(instance.NodeID) {
					t.Error("Instance" + string(i) + "'s contact incorrectly stored in Instance" + string(j) + "'s contact list")
				}
			}
		}else {
			if (i < numberOfContactsPerNode/2) {
				for j := 0; j < numberOfContactsPerNode; j++ {
					contact, err := instances[j].FindContact(instance.NodeID)
					if err != nil {
						t.Error("Instance" + string(i) + "'s contact not found in Instance" + string(j) + "'s contact list")
						return
					}

					if !contact.NodeID.Equals(instance.NodeID) {
						t.Error("Instance" + string(i) + "'s contact incorrectly stored in Instance" + string(j) + "'s contact list")
					}
				}
			} else if i > numberOfNodes- numberOfContactsPerNode/2 {
				for j := numberOfNodes - numberOfContactsPerNode; j < numberOfNodes; j++ {
					contact, err := instances[j].FindContact(instance.NodeID)
					if err != nil {
						t.Error("Instance" + string(i) + "'s contact not found in Instance" + string(j) + "'s contact list")
						return
					}

					if !contact.NodeID.Equals(instance.NodeID) {
						t.Error("Instance" + string(i) + "'s contact incorrectly stored in Instance" + string(j) + "'s contact list")
					}
				}
			}
		}
		
	}


	for i := 0; i < numberOfNodes; i++ {
		instance := instances[i]
		start := i - numberOfContactsPerNode/2;
		end := i + numberOfContactsPerNode/2;

		if i >= numberOfContactsPerNode/2 && i <= numberOfNodes - numberOfContactsPerNode/2 {
			for j := start; j < end; j ++ {
				contact, err := instance.FindContact(instances[j].NodeID)
				if err != nil {
					t.Error("Instance" + string(j) + "'s contact not found in Instance" + string(i) + "'s contact list")
					return
				}

				if !contact.NodeID.Equals(instances[j].NodeID) {
					t.Error("Instance" + string(j) + "'s contact incorrectly stored in Instance" + string(i) + "'s contact list")
				}
			}
		}  else {
			if (i < numberOfContactsPerNode/2) {
				for j := 0; j < numberOfContactsPerNode; j++ {
					contact, err := instance.FindContact(instances[j].NodeID)
					if err != nil {
						t.Error("Instance" + string(j) + "'s contact not found in Instance" + string(i) + "'s contact list")
						return
					}

					if !contact.NodeID.Equals(instances[j].NodeID) {
						t.Error("Instance" + string(j) + "'s contact incorrectly stored in Instance" + string(i) + "'s contact list")
					}
				}
			} else {
				for j := numberOfNodes - numberOfContactsPerNode; j < numberOfNodes; j++ {
					contact, err := instance.FindContact(instances[j].NodeID)
					if err != nil {
						t.Error("Instance" + string(j) + "'s contact not found in Instance" + string(i) + "'s contact list")
						return
					}

					if !contact.NodeID.Equals(instances[j].NodeID) {
						t.Error("Instance" + string(j) + "'s contact incorrectly stored in Instance" + string(i) + "'s contact list")
					}
				}
			}
		}
		
	}

	fmt.Println("..............Check Iterative find node function......")

	c, e := instances[2].FindContact(instances[4].NodeID)
	if e != nil {
		fmt.Println("instance 2 didn't have contact of instance 4")
	} else {
		fmt.Println("instance id" + c.NodeID.AsString())
	}

	//check iterative find node 0 find 50
	testerNumber := 0
	testSearchNumber := 50
	resContacts := instances[testerNumber].IterativeFindNode(instances[testSearchNumber].SelfContact.NodeID)

	fmt.Println("..............iterative find node result......")
	// fmt.Println(instances[testerNumber].ContactsToString(resContacts))

	fmt.Println("..............theoretical......")

	//calculate theoretical result
	theoreticalRes := make([]ContactDistance, 0)
	for i := 0; i < numberOfNodes; i++ {
		if i == testSearchNumber {
			continue
		}
		instance := instances[i]
		contact, err := instance.FindContact(instances[testSearchNumber].NodeID)
			if err != nil || !contact.NodeID.Equals(instances[testSearchNumber].NodeID) {
				continue
		}
		if len(theoreticalRes) < MAX_BUCKET_SIZE {
			theoreticalRes = append(theoreticalRes, instances[testerNumber].ContactToDistanceContact(instances[i].SelfContact, instances[testSearchNumber].SelfContact.NodeID))
			sort.Sort(ByDistance(theoreticalRes))
		} else {
			if theoreticalRes[MAX_BUCKET_SIZE-1].Distance.Compare(instances[i].SelfContact.NodeID.Xor(instances[testSearchNumber].SelfContact.NodeID)) == 1 {
				theoreticalRes = theoreticalRes[0: MAX_BUCKET_SIZE - 1]
				theoreticalRes = append(theoreticalRes, instances[testerNumber].ContactToDistanceContact(instances[i].SelfContact, instances[testSearchNumber].SelfContact.NodeID))
				sort.Sort(ByDistance(theoreticalRes))
			}

		}
	}
	fmt.Println("..............theoretical result length ......" + strconv.Itoa(len(theoreticalRes)))

	//convert contactdistance to contact
	theoreticalContact := make([]Contact, 0)
	for i := 0; i < len(theoreticalRes); i++ {
		theoreticalContact = append(theoreticalContact, theoreticalRes[i].SelfContact)
	}
    fmt.Println("..............theoretical result......")
    fmt.Println(instances[0].ContactsToString(theoreticalContact))

	fmt.Println("..............Compare......")
	// fmt.Println("..............find result......" + instances[0].ContactsToString(resContacts))

	//compare result
	for i := 0; i < MAX_BUCKET_SIZE && i < len(resContacts); i++ {
		if !theoreticalRes[i].SelfContact.NodeID.Equals(resContacts[i].NodeID) {
			t.Error("TestIterativeFindNode error, the nodes return are not the closet ones")
		}
	}

	return
}

