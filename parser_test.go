package main

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/bugsnag/bugsnag-go/v2"
	e "github.com/bugsnag/bugsnag-go/v2/errors"
)

var createdBy = `panic: hello!

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91
created by github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.App.Index
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:14 +0x3e

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a
`

var normalSplit = `panic: hello!

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a
`

var lastGoroutine = `panic: hello!

goroutine 16 [IO wait]:
net.runtime_pollWait(0x911c30, 0x72, 0x0)
	/0/c/go/src/pkg/runtime/netpoll.goc:146 +0x66
net.(*pollDesc).Wait(0xc2080ba990, 0x72, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:84 +0x46
net.(*pollDesc).WaitRead(0xc2080ba990, 0x0, 0x0)
	/0/c/go/src/pkg/net/fd_poll_runtime.go:89 +0x42
net.(*netFD).accept(0xc2080ba930, 0x58be30, 0x0, 0x9103f0, 0x23)
	/0/c/go/src/pkg/net/fd_unix.go:409 +0x343
net.(*TCPListener).AcceptTCP(0xc20803e168, 0x8, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:234 +0x5d
net.(*TCPListener).Accept(0xc20803e168, 0x0, 0x0, 0x0, 0x0)
	/0/c/go/src/pkg/net/tcpsock_posix.go:244 +0x4b
github.com/revel/revel.Run(0xe6d9)
	/0/go/src/github.com/revel/revel/server.go:113 +0x926
main.main()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/tmp/main.go:109 +0xe1a

goroutine 54 [running]:
runtime.panic(0x35ce40, 0xc208039db0)
	/0/c/go/src/pkg/runtime/panic.c:279 +0xf5
github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers.func·001()
	/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go:13 +0x74
net/http.(*Server).Serve(0xc20806c780, 0x910c88, 0xc20803e168, 0x0, 0x0)
	/0/c/go/src/pkg/net/http/server.go:1698 +0x91
`

var stackOverflow = `fatal error: stack overflow

runtime stack:
runtime.throw(0x10cd82b, 0xe)
	/go/src/runtime/panic.go:1116 +0x72
runtime.newstack()
	/go/src/runtime/stack.go:1060 +0x78d
runtime.morestack()
	/go/src/runtime/asm_amd64.s:449 +0x8f

goroutine 1 [running]:
main.stackExhaustion.func1(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/app/cases.go:42 +0x74 fp=0xc020161be0 sp=0xc020161bd8 pc=0x10a7774
main.stackExhaustion.func1(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/app/cases.go:43 +0x5f fp=0xc020163b30 sp=0xc020161be0 pc=0x10a775f
main.stackExhaustion.func1(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/app/cases.go:43 +0x5f fp=0xc020165a80 sp=0xc020163b30 pc=0x10a775f
main.stackExhaustion.func1(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/app/cases.go:43 +0x5f fp=0xc0201679d0 sp=0xc020165a80 pc=0x10a775f
main.stackExhaustion.func1(0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/app/cases.go:43 +0x5f fp=0xc0201679d0 sp=0xc020165a80 pc=0x10a775f
...additional frames elided...
`

var unexpectedSignal = `fatal error: unexpected signal during runtime execution
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x408db3e]

runtime stack:
runtime.throw({0x40ac264, 0x7fff203462d6})
	/usr/local/Cellar/go/1.17.1/libexec/src/runtime/panic.go:1198 +0x71
runtime.sigpanic()
	/usr/local/Cellar/go/1.17.1/libexec/src/runtime/signal_unix.go:719 +0x396

goroutine 1 [syscall]:
runtime.cgocall(0x408daf0, 0xc000064f38)
	/usr/local/Cellar/go/1.17.1/libexec/src/runtime/cgocall.go:156 +0x5c fp=0xc000064f10 sp=0xc000064ed8 pc=0x400643c
main._Cfunc_testSig(0x0)
	_cgo_gotypes.go:40 +0x45 fp=0xc000064f38 sp=0xc000064f10 pc=0x408da05
main.main()
	/Users/bugsnag/develop/cgo_segfault/main.go:9 +0x25 fp=0xc000064f80 sp=0xc000064f38 pc=0x408da85
runtime.main()
	/usr/local/Cellar/go/1.17.1/libexec/src/runtime/proc.go:255 +0x227 fp=0xc000064fe0 sp=0xc000064f80 pc=0x4035667
runtime.goexit()
	/usr/local/Cellar/go/1.17.1/libexec/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc000064fe8 sp=0xc000064fe0 pc=0x405eba1
`

var result = []e.StackFrame{
	{File: "/0/c/go/src/pkg/runtime/panic.c", LineNumber: 279, Name: "panic", Package: "runtime"},
	{File: "/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go", LineNumber: 13, Name: "func.001", Package: "github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers"},
	{File: "/0/c/go/src/pkg/net/http/server.go", LineNumber: 1698, Name: "(*Server).Serve", Package: "net/http"},
}

var resultCreatedBy = append(result,
	e.StackFrame{File: "/0/go/src/github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers/app.go", LineNumber: 14, Name: "App.Index", Package: "github.com/loopj/bugsnag-example-apps/go/revelapp/app/controllers", ProgramCounter: 0x0})

func TestParsePanic(t *testing.T) {

	todo := map[string]string{
		"createdBy":     createdBy,
		"normalSplit":   normalSplit,
		"lastGoroutine": lastGoroutine,
	}

	for key, val := range todo {
		Err, metadata, err := parsePanic(val)

		if err != nil {
			t.Fatal(err)
		}

		if Err.typeName != "panic" {
			t.Errorf("Wrong type: %s", Err.TypeName())
		}

		if Err.Error() != "hello!" {
			t.Errorf("Wrong message: %s", Err.Error())
		}

		if Err.StackFrames()[0].Func() != nil {
			t.Errorf("Somehow managed to find a func...")
		}

		if metadata != nil {
			t.Errorf("Unexpectedly returned metadata...")
		}

		result := result
		if key == "createdBy" {
			result = resultCreatedBy
		}

		if !reflect.DeepEqual(Err.StackFrames(), result) {
			t.Errorf("Wrong stack for %s: %#v", key, Err.StackFrames())
		}
	}
}

var concurrentMapReadWrite = `fatal error: concurrent map read and map write

goroutine 1 [running]:
runtime.throw(0x10766f5, 0x21)
	/usr/local/Cellar/go/1.15.5/libexec/src/runtime/panic.go:1116 +0x72 fp=0xc00003a6c8 sp=0xc00003a698 pc=0x102d592
runtime.mapaccess1_faststr(0x1066fc0, 0xc000060000, 0x10732e0, 0x1, 0xc000100088)
	/usr/local/Cellar/go/1.15.5/libexec/src/runtime/map_faststr.go:21 +0x465 fp=0xc00003a738 sp=0xc00003a6c8 pc=0x100e9c5
main.concurrentWrite()
	/myapps/go/fatalerror/main.go:14 +0x7a fp=0xc00003a778 sp=0xc00003a738 pc=0x105d83a
main.main()
	/myapps/go/fatalerror/main.go:41 +0x25 fp=0xc00003a788 sp=0xc00003a778 pc=0x105d885
runtime.main()
	/usr/local/Cellar/go/1.15.5/libexec/src/runtime/proc.go:204 +0x209 fp=0xc00003a7e0 sp=0xc00003a788 pc=0x102fd49
runtime.goexit()
	/usr/local/Cellar/go/1.15.5/libexec/src/runtime/asm_amd64.s:1374 +0x1 fp=0xc00003a7e8 sp=0xc00003a7e0 pc=0x105a4a1

goroutine 5 [runnable]:
main.concurrentWrite.func1(0xc000060000)
	/myapps/go/fatalerror/main.go:10 +0x4c
created by main.concurrentWrite
	/myapps/go/fatalerror/main.go:8 +0x4b
`

func TestParseFatalError(t *testing.T) {

	Err, metadata, err := parsePanic(concurrentMapReadWrite)

	if err != nil {
		t.Fatal(err)
	}

	if Err.TypeName() != "fatal error" {
		t.Errorf("Wrong type: %s", Err.TypeName())
	}

	if Err.Error() != "concurrent map read and map write" {
		t.Errorf("Wrong message: '%s'", Err.Error())
	}

	if Err.StackFrames()[0].Func() != nil {
		t.Errorf("Somehow managed to find a func...")
	}

	if metadata != nil {
		t.Errorf("Unexpectedly returned metadata...")
	}

	var result = []e.StackFrame{
		{File: "/usr/local/Cellar/go/1.15.5/libexec/src/runtime/panic.go", LineNumber: 1116, Name: "throw", Package: "runtime"},
		{File: "/usr/local/Cellar/go/1.15.5/libexec/src/runtime/map_faststr.go", LineNumber: 21, Name: "mapaccess1_faststr", Package: "runtime"},
		{File: "/myapps/go/fatalerror/main.go", LineNumber: 14, Name: "concurrentWrite", Package: "main"},
		{File: "/myapps/go/fatalerror/main.go", LineNumber: 41, Name: "main", Package: "main"},
		{File: "/usr/local/Cellar/go/1.15.5/libexec/src/runtime/proc.go", LineNumber: 204, Name: "main", Package: "runtime"},
		{File: "/usr/local/Cellar/go/1.15.5/libexec/src/runtime/asm_amd64.s", LineNumber: 1374, Name: "goexit", Package: "runtime"},
	}

	if !reflect.DeepEqual(Err.StackFrames(), result) {
		t.Errorf("Wrong stack for concurrent write fatal error:")
		for i, frame := range result {
			t.Logf("[%d] %#v", i, frame)
			if len(Err.StackFrames()) > i {
				t.Logf("    %#v", Err.StackFrames()[i])
			}
		}
	}
}

func TestParseStackOverflow(t *testing.T) {
	Err, metadata, err := parsePanic(stackOverflow)

	if err != nil {
		t.Fatal(err)
	}

	if Err.TypeName() != "fatal error" {
		t.Errorf("Wrong type: %s", Err.TypeName())
	}

	if Err.Error() != "stack overflow" {
		t.Errorf("Wrong message: '%s'", Err.Error())
	}

	if Err.StackFrames()[0].Func() != nil {
		t.Errorf("Somehow managed to find a func...")
	}

	if metadata != nil {
		t.Errorf("Unexpectedly returned metadata...")
	}

	var result = []e.StackFrame{
		{File: "/go/src/app/cases.go", LineNumber: 42, Name: "stackExhaustion.func1", Package: "main"},
		{File: "/go/src/app/cases.go", LineNumber: 43, Name: "stackExhaustion.func1", Package: "main"},
		{File: "/go/src/app/cases.go", LineNumber: 43, Name: "stackExhaustion.func1", Package: "main"},
		{File: "/go/src/app/cases.go", LineNumber: 43, Name: "stackExhaustion.func1", Package: "main"},
		{File: "/go/src/app/cases.go", LineNumber: 43, Name: "stackExhaustion.func1", Package: "main"},
	}

	if !reflect.DeepEqual(Err.StackFrames(), result) {
		t.Errorf("Wrong stack:")
		for i, frame := range result {
			t.Logf("[%d] %#v", i, frame)
			if len(Err.StackFrames()) > i {
				t.Logf("    %#v", Err.StackFrames()[i])
			}
		}
	}
}

func TestParseSignal(t *testing.T) {
	Err, metadata, err := parsePanic(unexpectedSignal)

	if err != nil {
		t.Fatal(err)
	}

	if Err.TypeName() != "SIGSEGV" {
		t.Errorf("Wrong type: %s", Err.TypeName())
	}

	if Err.Error() != "segmentation violation" {
		t.Errorf("Wrong message: '%s'", Err.Error())
	}

	if Err.StackFrames()[0].Func() != nil {
		t.Errorf("Somehow managed to find a func...")
	}

	if metadata == nil {
		t.Errorf("Missing Signal metadata...")
	} else {
		expectedMetadata := &bugsnag.MetaData{
			"signal": {
				"code":        "0x1",
				"addr":        "0x0",
				"pc":          "0x408db3e",
			},
		}

		if !reflect.DeepEqual(metadata, expectedMetadata) {
			t.Errorf("Wrong metadata:")
			for k, sigDetail := range *expectedMetadata {
				t.Logf("[%s] %#v", k, sigDetail)
				t.Logf("%-"+strconv.Itoa(len(k))+"s   %#v", "", (*metadata)[k])
			}
		}
	}

	var result = []e.StackFrame{
		{File: "/usr/local/Cellar/go/1.17.1/libexec/src/runtime/cgocall.go", LineNumber: 156, Name: "cgocall", Package: "runtime"},
		{File: "_cgo_gotypes.go", LineNumber: 40, Name: "_Cfunc_testSig", Package: "<generated>.main"},
		{File: "/Users/bugsnag/develop/cgo_segfault/main.go", LineNumber: 9, Name: "main", Package: "main"},
		{File: "/usr/local/Cellar/go/1.17.1/libexec/src/runtime/proc.go", LineNumber: 255, Name: "main", Package: "runtime"},
		{File: "/usr/local/Cellar/go/1.17.1/libexec/src/runtime/asm_amd64.s", LineNumber: 1581, Name: "goexit", Package: "runtime"},
	}

	if !reflect.DeepEqual(Err.StackFrames(), result) {
		t.Errorf("Wrong stack:")
		for i, frame := range result {
			t.Logf("[%d] %#v", i, frame)
			if len(Err.StackFrames()) > i {
				t.Logf("    %#v", Err.StackFrames()[i])
			}
		}
	}
}
