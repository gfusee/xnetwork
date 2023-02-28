package process

import (
	"errors"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-tools-go/elasticreindexer/config"
)

type interval struct {
	start int64
	stop  int64
}

type reindexerMultiWrite struct {
	indicesNoTimestamp   []string
	indicesWithTimestamp []string
	numParallelWrite     int
	blockChainStartTime  int64
	enabled              bool

	reindexerClient ReindexerHandler
}

func NewReindexerMultiWrite(reindexer ReindexerHandler, cfg config.IndicesConfig) (*reindexerMultiWrite, error) {
	if reindexer == nil {
		return nil, errors.New("nil ReindexerHandler")
	}
	if cfg.WithTimestamp.BlockchainStartTime <= 0 {
		return nil, errors.New("blockchainStartTime cannot be less than zero")
	}

	return &reindexerMultiWrite{
		reindexerClient:      reindexer,
		indicesNoTimestamp:   cfg.Indices,
		indicesWithTimestamp: cfg.WithTimestamp.IndicesWithTimestamp,
		numParallelWrite:     cfg.WithTimestamp.NumParallelWrites,
		blockChainStartTime:  cfg.WithTimestamp.BlockchainStartTime,
		enabled:              cfg.WithTimestamp.Enabled,
	}, nil
}

func (rmw *reindexerMultiWrite) ProcessNoTimestamp(overwrite bool, skipMappings bool) error {
	for _, index := range rmw.indicesNoTimestamp {
		if index == "" {
			continue
		}

		err := rmw.reindexerClient.Process(overwrite, skipMappings, index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rmw *reindexerMultiWrite) ProcessWithTimestamp(overwrite bool, skipMappings bool) error {
	if !rmw.enabled {
		return nil
	}

	currentTimestampUnix := time.Now().Unix()
	intervals, err := computeIntervals(rmw.blockChainStartTime, currentTimestampUnix, int64(rmw.numParallelWrite))
	if err != nil {
		return err
	}

	for _, index := range rmw.indicesWithTimestamp {
		if index == "" {
			continue
		}

		err = rmw.reindexBasedOnIntervals(index, intervals, overwrite, skipMappings)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rmw *reindexerMultiWrite) reindexBasedOnIntervals(
	index string,
	intervals []*interval,
	overwrite bool,
	skipMappings bool,
) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(intervals))

	log.Info("starting reindexing", "index", index)

	count := uint64(0)

	for idx, interv := range intervals {
		go func(startTime, stopTime int64, idx int, w *sync.WaitGroup) {
			defer func() {
				log.Info("done", "interval nr", idx)
				w.Done()
			}()

			errIndex := rmw.reindexerClient.ProcessIndexWithTimestamp(index, overwrite, skipMappings, startTime, stopTime, &count)
			if errIndex != nil {
				log.Warn("rmw.processIndexWithTimestamp", "index", index, "error", errIndex.Error())
			}
		}(interv.start, interv.stop, idx, wg)

		time.Sleep(time.Second)
	}

	wg.Wait()

	return nil
}

func computeIntervals(startTime, endTime int64, numIntervals int64) ([]*interval, error) {
	if startTime > endTime {
		return nil, errors.New("blockchain start time is greater than current timestamp")
	}
	if numIntervals < 2 {
		return []*interval{{
			start: startTime,
			stop:  endTime,
		}}, nil
	}

	difference := endTime - startTime

	step := difference / numIntervals

	intervals := make([]*interval, 0)
	for idx := int64(0); idx < numIntervals; idx++ {
		start := startTime + idx*step
		stop := startTime + (idx+1)*step

		if idx == numIntervals-1 {
			stop = endTime
		}

		intervals = append(intervals, &interval{
			start: start,
			stop:  stop,
		})
	}

	return intervals, nil
}
