package syncx

import "sync/atomic"

type Ongoing struct{ v uint32 }

func (x *Ongoing) Set() bool { return atomic.CompareAndSwapUint32(&x.v, 0, 1) }
func (x *Ongoing) Reset()    { atomic.StoreUint32(&x.v, 0) }
func (x Ongoing) Bool() bool { return x.v == 1 }
