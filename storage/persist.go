package storage

import (
	"github.com/golang/glog"
	"github.com/heidi-ann/ios/msgs"
  "strconv"
  "os"
)

type FileStorage struct {
  viewFile wal
  logFile  wal
  snapFile fileWriter
}

func MakeFileStorage(diskPath string, persistenceMode string) *FileStorage {
  // create disk path if needs be
  if _, err := os.Stat(diskPath); os.IsNotExist(err) {
    err = os.MkdirAll(diskPath,0777)
    if err != nil {
      glog.Fatal(err)
    }
  }

  logFilename := diskPath + "/log.temp"
  dataFilename := diskPath + "/view.temp"
  snapFilename := diskPath + "/snapshot.temp"

  viewFile := openWriteAheadFile(dataFilename, persistenceMode)
  logFile := openWriteAheadFile(logFilename, persistenceMode)
  snapFile := openWriter(snapFilename)
	s := FileStorage{viewFile, logFile, snapFile}
	return &s
}

func (fs *FileStorage) PersistView(view int) {
  glog.Info("Updating view to ", view, " in persistent storage")
  fs.viewFile.writeAhead([]byte(strconv.Itoa(view)))
}

func (fs *FileStorage) PersistLogUpdate(log msgs.LogUpdate) {
  glog.V(1).Info("Updating log with ", log, " in persistent storage")
  b, err := msgs.Marshal(log)
  if err != nil {
    glog.Fatal(err)
  }
  // write to persistent storage
  fs.logFile.writeAhead(b)
}

func (fs *FileStorage) PersistSnapshot(snap msgs.Snapshot) {
  glog.Info("Saving request cache and state machine snapshot upto index", snap.Index,
    " of size ", len(snap.Bytes))
  fs.snapFile.write([]byte(strconv.Itoa(snap.Index)))
  fs.snapFile.write(snap.Bytes)
}
