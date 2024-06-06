//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return nil
	}
	maxOut := guests[0].CheckOutDate
	for _, guest := range guests {
		maxOut = max(maxOut, guest.CheckOutDate)
	}
	info := make([]int, maxOut+1)
	for _, guest := range guests {
		info[guest.CheckInDate]++
		info[guest.CheckOutDate]--
	}
	res := make([]Load, 0)
	count := 0
	for data, cnt := range info {
		count += cnt
		if cnt == 0 {
			continue
		}
		res = append(res, Load{StartDate: data, GuestCount: count})
	}
	return res
}
