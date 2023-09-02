package fsm

import "fmt"

type parsingStateMachine struct {
	pointerIndex  int
	buffer        map[string][]string
	sourcePointer []string
}

func New(sourcePointer *[]string) *parsingStateMachine {
	return &parsingStateMachine{
		pointerIndex:  0,
		buffer:        make(map[string][]string),
		sourcePointer: *sourcePointer,
	}
}

func (self *parsingStateMachine) SwitchState(index int) {
	self.pointerIndex = index
}

func (self *parsingStateMachine) IsValidToProcced() bool {
	return len(self.buffer) > 0
}

func (self *parsingStateMachine) CloseState(index int) error {
	if _, ok := self.buffer[self.sourcePointer[index]]; ok {
		return fmt.Errorf("Argument duplicated [%s]\n", self.sourcePointer[index])
	}

	if index-self.pointerIndex > 1 {
		self.buffer[self.sourcePointer[self.pointerIndex]] = make([]string, index-self.pointerIndex-1)
		copy(
			self.buffer[self.sourcePointer[self.pointerIndex]],
			self.sourcePointer[self.pointerIndex+1:index],
		)
	} else {
		self.buffer[self.sourcePointer[self.pointerIndex]] = make([]string, 0)
	}
	return nil
}

func (self *parsingStateMachine) Finally() *map[string][]string {
	self.buffer[self.sourcePointer[self.pointerIndex]] = make(
		[]string,
		len(self.sourcePointer)-self.pointerIndex-1,
	)
	copy(
		self.buffer[self.sourcePointer[self.pointerIndex]],
		self.sourcePointer[self.pointerIndex+1:len(self.sourcePointer)],
	)
	return &self.buffer
}
