package messenger

import "sync"

type (
	Messenger struct {
		chans map[uint64][]chan *Message
		sync.RWMutex
	}

	Message struct {
		Json []byte
		Err  error
	}
)

func New() *Messenger {
	return &Messenger{
		chans: make(map[uint64][]chan *Message),
	}
}

func (m *Messenger) Publish(nr uint64, json []byte, err error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.chans[nr]; !ok {
		return
	}

	mes := &Message{
		Json: json,
		Err:  err,
	}

	for _, c := range m.chans[nr] {
		c <- mes
	}

	delete(m.chans, nr)
}

func (m *Messenger) Subscribe(nr uint64, ch chan *Message) {
	m.Lock()
	defer m.Unlock()

	m.chans[nr] = append(m.chans[nr], ch)
}

func (m *Messenger) CreateEntry(nr uint64) {
	m.chans[nr] = []chan *Message{}
}

func (m *Messenger) HasEntry(nr uint64) bool {
	_, ok := m.chans[nr]

	return ok
}
