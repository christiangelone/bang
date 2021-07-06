package sugar

type IfCase struct {
	value interface{}
}

func (self *IfCase) Eval() interface{} {
	return self.value
}

func (self *IfCase) ElseLazy(value func() interface{}) interface{} {
	if self.value == nil {
		return value()
	}
	return self.value
}

func (self *IfCase) Else(value interface{}) interface{} {
	if self.value == nil {
		return value
	}
	return self.value
}

func If(condition bool) func(v interface{}) *IfCase {
	return func(v interface{}) *IfCase {
		if condition {
			return &IfCase{
				value: v,
			}
		}
		return &IfCase{
			value: nil,
		}
	}
}

func IfLazy(condition bool) func(func() interface{}) *IfCase {
	return func(v func() interface{}) *IfCase {
		if condition {
			return &IfCase{
				value: v(),
			}
		}
		return &IfCase{
			value: nil,
		}
	}
}

func Go(fn func()) {
	go fn()
}
