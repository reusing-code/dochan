package persist

import (
	"compress/gzip"
	"encoding/gob"
	"io"
)

type DataMap map[interface{}][]string

func PersistData(data DataMap, w io.Writer) error {
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()

	enc := gob.NewEncoder(gw)
	err = enc.Encode(&data)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveData(r io.Reader) (DataMap, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(gr)
	result := make(DataMap)
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
