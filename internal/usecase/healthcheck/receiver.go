package healthcheck

import "log/slog"

func receiveEvents(num int, urlInput <-chan string, tester Pinger, storage Storage) {
	for range num {
		go func() {
			for url := range urlInput {
				res, err := tester.Ping(url)
				if err != nil {
					slog.Error("error pinging url %s: %s", url, err)
					continue
				}

				storage.StoreStatus(url, res)
			}
		}()
	}
}
