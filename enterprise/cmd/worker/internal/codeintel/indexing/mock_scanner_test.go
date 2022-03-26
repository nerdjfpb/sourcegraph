// Code generated by go-mockgen 1.1.4; DO NOT EDIT.

package indexing

import (
	"sync"

	dbstore "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/dbstore"
	shared "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/shared"
)

// MockPackageReferenceScanner is a mock implementation of the
// PackageReferenceScanner interface (from the package
// github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/dbstore)
// used for unit testing.
type MockPackageReferenceScanner struct {
	// CloseFunc is an instance of a mock function object controlling the
	// behavior of the method Close.
	CloseFunc *PackageReferenceScannerCloseFunc
	// NextFunc is an instance of a mock function object controlling the
	// behavior of the method Next.
	NextFunc *PackageReferenceScannerNextFunc
}

// NewMockPackageReferenceScanner creates a new mock of the
// PackageReferenceScanner interface. All methods return zero values for all
// results, unless overwritten.
func NewMockPackageReferenceScanner() *MockPackageReferenceScanner {
	return &MockPackageReferenceScanner{
		CloseFunc: &PackageReferenceScannerCloseFunc{
			defaultHook: func() error {
				return nil
			},
		},
		NextFunc: &PackageReferenceScannerNextFunc{
			defaultHook: func() (shared.PackageReference, bool, error) {
				return shared.PackageReference{}, false, nil
			},
		},
	}
}

// NewStrictMockPackageReferenceScanner creates a new mock of the
// PackageReferenceScanner interface. All methods panic on invocation,
// unless overwritten.
func NewStrictMockPackageReferenceScanner() *MockPackageReferenceScanner {
	return &MockPackageReferenceScanner{
		CloseFunc: &PackageReferenceScannerCloseFunc{
			defaultHook: func() error {
				panic("unexpected invocation of MockPackageReferenceScanner.Close")
			},
		},
		NextFunc: &PackageReferenceScannerNextFunc{
			defaultHook: func() (shared.PackageReference, bool, error) {
				panic("unexpected invocation of MockPackageReferenceScanner.Next")
			},
		},
	}
}

// NewMockPackageReferenceScannerFrom creates a new mock of the
// MockPackageReferenceScanner interface. All methods delegate to the given
// implementation, unless overwritten.
func NewMockPackageReferenceScannerFrom(i dbstore.PackageReferenceScanner) *MockPackageReferenceScanner {
	return &MockPackageReferenceScanner{
		CloseFunc: &PackageReferenceScannerCloseFunc{
			defaultHook: i.Close,
		},
		NextFunc: &PackageReferenceScannerNextFunc{
			defaultHook: i.Next,
		},
	}
}

// PackageReferenceScannerCloseFunc describes the behavior when the Close
// method of the parent MockPackageReferenceScanner instance is invoked.
type PackageReferenceScannerCloseFunc struct {
	defaultHook func() error
	hooks       []func() error
	history     []PackageReferenceScannerCloseFuncCall
	mutex       sync.Mutex
}

// Close delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockPackageReferenceScanner) Close() error {
	r0 := m.CloseFunc.nextHook()()
	m.CloseFunc.appendCall(PackageReferenceScannerCloseFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Close method of the
// parent MockPackageReferenceScanner instance is invoked and the hook queue
// is empty.
func (f *PackageReferenceScannerCloseFunc) SetDefaultHook(hook func() error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Close method of the parent MockPackageReferenceScanner instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *PackageReferenceScannerCloseFunc) PushHook(hook func() error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PackageReferenceScannerCloseFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func() error {
		return r0
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PackageReferenceScannerCloseFunc) PushReturn(r0 error) {
	f.PushHook(func() error {
		return r0
	})
}

func (f *PackageReferenceScannerCloseFunc) nextHook() func() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PackageReferenceScannerCloseFunc) appendCall(r0 PackageReferenceScannerCloseFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of PackageReferenceScannerCloseFuncCall
// objects describing the invocations of this function.
func (f *PackageReferenceScannerCloseFunc) History() []PackageReferenceScannerCloseFuncCall {
	f.mutex.Lock()
	history := make([]PackageReferenceScannerCloseFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PackageReferenceScannerCloseFuncCall is an object that describes an
// invocation of method Close on an instance of MockPackageReferenceScanner.
type PackageReferenceScannerCloseFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PackageReferenceScannerCloseFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PackageReferenceScannerCloseFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// PackageReferenceScannerNextFunc describes the behavior when the Next
// method of the parent MockPackageReferenceScanner instance is invoked.
type PackageReferenceScannerNextFunc struct {
	defaultHook func() (shared.PackageReference, bool, error)
	hooks       []func() (shared.PackageReference, bool, error)
	history     []PackageReferenceScannerNextFuncCall
	mutex       sync.Mutex
}

// Next delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockPackageReferenceScanner) Next() (shared.PackageReference, bool, error) {
	r0, r1, r2 := m.NextFunc.nextHook()()
	m.NextFunc.appendCall(PackageReferenceScannerNextFuncCall{r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the Next method of the
// parent MockPackageReferenceScanner instance is invoked and the hook queue
// is empty.
func (f *PackageReferenceScannerNextFunc) SetDefaultHook(hook func() (shared.PackageReference, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Next method of the parent MockPackageReferenceScanner instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *PackageReferenceScannerNextFunc) PushHook(hook func() (shared.PackageReference, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *PackageReferenceScannerNextFunc) SetDefaultReturn(r0 shared.PackageReference, r1 bool, r2 error) {
	f.SetDefaultHook(func() (shared.PackageReference, bool, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *PackageReferenceScannerNextFunc) PushReturn(r0 shared.PackageReference, r1 bool, r2 error) {
	f.PushHook(func() (shared.PackageReference, bool, error) {
		return r0, r1, r2
	})
}

func (f *PackageReferenceScannerNextFunc) nextHook() func() (shared.PackageReference, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *PackageReferenceScannerNextFunc) appendCall(r0 PackageReferenceScannerNextFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of PackageReferenceScannerNextFuncCall objects
// describing the invocations of this function.
func (f *PackageReferenceScannerNextFunc) History() []PackageReferenceScannerNextFuncCall {
	f.mutex.Lock()
	history := make([]PackageReferenceScannerNextFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// PackageReferenceScannerNextFuncCall is an object that describes an
// invocation of method Next on an instance of MockPackageReferenceScanner.
type PackageReferenceScannerNextFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 shared.PackageReference
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c PackageReferenceScannerNextFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c PackageReferenceScannerNextFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}
