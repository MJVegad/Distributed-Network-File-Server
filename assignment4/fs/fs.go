package fs

import (
	_ "fmt"
	"sync"
	"time"
)

type FileInfo struct {
	filename   string
	contents   []byte
	version    int
	absexptime time.Time
	timer      *time.Timer
}

type FS struct {
	sync.RWMutex
	Dir map[string]*FileInfo
}

//var fs = &FS{dir: make(map[string]*FileInfo, 1000)}
//var gversion = 0 // global version

func (fi *FileInfo) cancelTimer() {
	if fi.timer != nil {
		fi.timer.Stop()
		fi.timer = nil
	}
}

func ProcessMsg(msg *Msg,ff *FS,gversion *int) *Msg {
	switch msg.Kind {
	case 'r':
		return processRead(msg,ff)
	case 'w':
		return processWrite(msg,ff,gversion)
	case 'c':
		return processCas(msg,ff,gversion)
	case 'd':
		return processDelete(msg,ff)
	}

	// Default: Internal error. Shouldn't come here since
	// the msg should have been validated earlier.
	return &Msg{Kind: 'I'}
}

func processRead(msg *Msg,  ff *FS) *Msg {
	ff.RLock()
	defer ff.RUnlock()
	if fi := ff.Dir[msg.Filename]; fi != nil {
		remainingTime := 0
		if fi.timer != nil {
			remainingTime := int(fi.absexptime.Sub(time.Now()))
			if remainingTime < 0 {
				remainingTime = 0
			}
		}
		return &Msg{
			Kind:     'C',
			Filename: fi.filename,
			Contents: fi.contents,
			Numbytes: len(fi.contents),
			Exptime:  remainingTime,
			Version:  fi.version,
		}
	} else {
		return &Msg{Kind: 'F'} // file not found
	}
}

func internalWrite(msg *Msg,ff *FS, gversion *int) *Msg {
	fi := ff.Dir[msg.Filename]
	if fi != nil {
		fi.cancelTimer()
	} else {
		fi = &FileInfo{}
	}

	*gversion += 1
	fi.filename = msg.Filename
	fi.contents = msg.Contents
	fi.version = *gversion

	var absexptime time.Time
	if msg.Exptime > 0 {
		dur := time.Duration(msg.Exptime) * time.Second
		absexptime = time.Now().Add(dur)
		timerFunc := func(name string, ver int) func() {
			return func() {
				processDelete(&Msg{Kind: 'D',
					Filename: name,
					Version:  ver},ff)
			}
		}(msg.Filename, *gversion)

		fi.timer = time.AfterFunc(dur, timerFunc)
	}
	fi.absexptime = absexptime
	ff.Dir[msg.Filename] = fi

	return ok(*gversion)
}

func processWrite(msg *Msg, ff *FS, gversion *int) *Msg {
	ff.Lock()
	defer ff.Unlock()
	return internalWrite(msg,ff, gversion)
}

func processCas(msg *Msg, ff *FS, gversion *int) *Msg {
	ff.Lock()
	defer ff.Unlock()

	if fi := ff.Dir[msg.Filename]; fi != nil {
		if msg.Version != fi.version {
			return &Msg{Kind: 'V', Version: fi.version}
		}
	}
	return internalWrite(msg,ff,gversion)
}

func processDelete(msg *Msg, ff *FS) *Msg {
	ff.Lock()
	defer ff.Unlock()
	fi := ff.Dir[msg.Filename]
	if fi != nil {
		if msg.Version > 0 && fi.version != msg.Version {
			// non-zero msg.Version indicates a delete due to an expired timer
			return nil // nothing to do
		}
		fi.cancelTimer()
		delete(ff.Dir, msg.Filename)
		return ok(0)
	} else {
		return &Msg{Kind: 'F'} // file not found
	}

}

func ok(version int) *Msg {
	return &Msg{Kind: 'O', Version: version}
}
