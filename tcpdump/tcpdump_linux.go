package tcpdump

func (u *UDPWatcherService) Watch(destPort string, duration int) (*domain.Result, error) {
	result := domain.Result{
		Start: time.Now().UTC(),
	}

	// tcpdump -s 96 -nn "udp and port 34194"
	filter := fmt.Sprintf("udp and port %s", destPort)
	cmdParts := []string{u.CommandName, "-s", u.PacketSnapLength, "-nn", filter}

	lineChan := make(chan string)
	defer close(lineChan)

	tupleMap := map[string]*domain.NetTuple{}
	go func() {
		for line := range lineChan {
			tuple, err := u.parseLine(line)
			if err != nil {
				log.Printf("ignoring line: %s", line)
				continue
			}

			key := tuple.Key()
			if tup, ok := tupleMap[key]; ok {
				tup.NumPackets++
				continue
			}

			tupleMap[key] = tuple
		}
	}()

	err := util.ExecuteCmd(cmdParts, duration, lineChan)
	if err != nil {
		return nil, err
	}

	result.End = time.Now().UTC()

	tuples := []*domain.NetTuple{}
	for _, tuple := range tupleMap {
		tuples = append(tuples, tuple)
	}
	result.UDPTuples = tuples

	return &result, nil
}
