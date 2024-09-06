package kernel

import (
	"bufio"
	"context"
	"debug/elf"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/syspulse/tracker/linux/task"

	"github.com/cilium/ebpf"
	"github.com/hashicorp/go-multierror"
	"github.com/ianlancetaylor/demangle"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -no-global-types -target bpfel -cc clang -cflags -O2 OnCPUProfiling  ../../ebpf/oncpu.c -- -I/usr/include/x86_64-linux-gnu -I/usr/include/asm-generic/

const (
	KernelSymbolFilePath      = "/proc/kallsyms"
	FuncName                  = "do_perf_event"
	MissingSymbol             = "[MISSING]"
	ModuleTypeExec       int8 = iota
	ModuleTypeSo
	ModuleTypePerfMap
	ModuleTypeVDSO
	ModuleTypeUnknown
)

var (
	NotSupportProfilingExe = []string{
		"java", "python", "node", "bash", "ruby", "ssh",
	}
	profilingStatFinderList = []StatFinder{
		NewGoLibrary(),
	}
	ErrNotSupport       = fmt.Errorf("not support")
	mapFileContentRegex = regexp.MustCompile("(?P<StartAddr>[a-f\\d]+)\\-(?P<EndAddr>[a-f\\d]+)\\s(?P<Perm>[^\\s]+)" +
		"\\s(?P<Offset>[a-f\\d]+)\\s[a-f\\d]+\\:[a-f\\d]+\\s\\d+\\s+(?P<Name>[^\\n]+)")
)

type OnCPUProfilingObj struct {
	Counts *ebpf.Map     `ebpf:"counts"`
	Stacks *ebpf.Map     `ebpf:"stacks"`
	Prog   *ebpf.Program `ebpf:"do_perf_event"`
}

type Event struct {
	UserStackID   uint32
	KernelStackID uint32
}

type ModuleRange struct {
	StartAddr, EndAddr, FileOffset uint64
}

type Symbol struct {
	Name     string
	Location uint64
	Size     uint64
}

type Module struct {
	Ranges           []*ModuleRange
	Name             string
	Path             string
	Type             int8
	SoOffset, SoAddr uint64
	Symbols          []*Symbol
}

func (m *Module) contains(addr uint64) (uint64, bool) {
	for _, r := range m.Ranges {
		if addr >= r.StartAddr && addr < r.EndAddr {
			log.Default().Printf("found module %s could hanlde address: %d", m.Name, addr)
			if m.Type == ModuleTypeSo || m.Type == ModuleTypeVDSO {
				offset := addr - r.StartAddr + r.FileOffset
				offset += m.SoAddr - m.SoOffset
				log.Default().Printf("update address %d to offset %d", addr, offset)
				return offset, true
			}
			return addr, true
		}
	}
	return 0, false
}

func (m *Module) findAddr(offset uint64) *Symbol {
	start := 0
	end := len(m.Symbols) - 1
	for start < end {
		mid := start + (end-start)/2
		result := int64(offset) - int64(m.Symbols[mid].Location)

		if result < 0 {
			end = mid
		} else if result > 0 {
			start = mid + 1
		} else {
			return m.Symbols[mid]
		}
	}

	if start >= 1 && m.Symbols[start-1].Location < offset && offset < m.Symbols[start].Location {
		return m.Symbols[start-1]
	}
	log.Default().Printf("could not found the address: %d in module %s", offset, m.Name)

	return nil
}

func processSymbolName(name string) string {
	if name == "" {
		return ""
	}
	// fix process demangle symbol name, such as c++ language symbol
	skip := 0
	if name[0] == '.' || name[0] == '$' {
		skip++
	}
	return demangle.Filter(name[skip:])
}

type Info struct {
	Modules           []*Module
	cacheAddrToSymbol map[uint64]string
}

func (i *Info) FindSymbols(addresses []uint64, defaultSymbol string) []string {
	if len(addresses) == 0 {
		return nil
	}
	result := make([]string, 0)
	for _, addr := range addresses {
		if addr <= 0 {
			continue
		}
		s := i.FindSymbolName(addr)
		if s == "" {
			s = defaultSymbol
		}
		result = append(result, s)
	}
	return result
}

func (i *Info) FindSymbolName(address uint64) string {
	if d := i.cacheAddrToSymbol[address]; d != "" {
		return d
	}
	log.Default().Printf("ready to find the symbol from address: %d", address)
	foundModule := false
	for _, mod := range i.Modules {
		offset, c := mod.contains(address)
		if !c {
			continue
		}
		foundModule = true

		if sym := mod.findAddr(offset); sym != nil {
			name := processSymbolName(sym.Name)
			i.cacheAddrToSymbol[address] = name
			return name
		}
	}
	if !foundModule {
		log.Default().Printf("could not found any module to handle address: %d", address)
	}
	return ""
}

type KernelFinder struct {
	kernelFileExists bool
}

func NewInfo(modules map[string]*Module) *Info {
	ls := make([]*Module, 0)
	for _, m := range modules {
		ls = append(ls, m)
	}
	return &Info{Modules: ls, cacheAddrToSymbol: make(map[uint64]string)}
}

func NewKernelFinder() *KernelFinder {
	stat, _ := os.Stat(KernelSymbolFilePath)
	return &KernelFinder{kernelFileExists: stat != nil}
}

func (k *KernelFinder) Analyze(filepath string) (*Info, error) {
	kernelPath, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(kernelPath)
	symbols := make([]*Symbol, 0)
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), " ")
		atoi, err := strconv.ParseUint(info[0], 16, 64)

		if err != nil {
			return nil, fmt.Errorf("error read addr: %s, %v", info[0], err)
		}
		symbols = append(symbols, &Symbol{
			Name:     info[2],
			Location: atoi,
			Size:     0,
		})
	}

	kernelModule := &Module{
		Name:    "kernel",
		Symbols: symbols,
		// kernel module could handling all symbols
		Ranges: []*ModuleRange{
			{
				StartAddr: 0,
				EndAddr:   math.MaxUint64,
			},
		},
	}

	return NewInfo(map[string]*Module{
		"kernel": kernelModule,
	}), nil
}

type StatFinder interface {
	// IsSupport to stat the executable file for profiling
	IsSupport(filePath string) bool
	// AnalyzeSymbols in the file
	AnalyzeSymbols(filePath string) ([]*Symbol, error)
	// ToModule to init a new module
	ToModule(pid int32, modName, modPath string, moduleRange []*ModuleRange) (*Module, error)
}

func Exists(f string) bool {
	_, e := os.Stat(f)
	return !os.IsNotExist(e)
}

type GoLibrary struct {
}

func NewGoLibrary() *GoLibrary {
	return &GoLibrary{}
}

func (l *GoLibrary) IsSupport(filePath string) bool {
	f, err := elf.Open(filePath)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

func (l *GoLibrary) AnalyzeSymbols(filePath string) ([]*Symbol, error) {
	// read els file
	file, err := elf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// exist symbol data
	symbols, _ := file.Symbols()
	dySyms, _ := file.DynamicSymbols()
	if len(symbols) == 0 && len(dySyms) == 0 {
		return nil, nil
	}
	symbols = append(symbols, dySyms...)

	// adapt symbol struct
	data := make([]*Symbol, len(symbols))
	for i, sym := range symbols {
		data[i] = &Symbol{Name: sym.Name, Location: sym.Value, Size: sym.Size}
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Location < data[j].Location
	})

	return data, nil
}

func (l *GoLibrary) ToModule(pid int32, modName, modPath string, moduleRange []*ModuleRange) (*Module, error) {
	res := &Module{}
	res.Name = modName
	res.Path = modPath
	res.Ranges = moduleRange
	file, err := elf.Open(modPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header := file.FileHeader
	mType := ModuleTypeUnknown
	switch header.Type {
	case elf.ET_EXEC:
		mType = ModuleTypeExec
	case elf.ET_DYN:
		mType = ModuleTypeSo
	}

	if mType == ModuleTypeUnknown {
		if strings.HasSuffix(modPath, ".map") && Exists(modPath) {
			mType = ModuleTypePerfMap
		} else if modName == "[vdso]" {
			mType = ModuleTypeVDSO
		}
	} else if mType == ModuleTypeSo {
		section := file.Section(".text")
		if section == nil {
			return nil, fmt.Errorf("could not found .text section in so file: %s", modName)
		}
		res.SoAddr = section.Addr
		res.SoOffset = section.Offset
	}
	res.Type = mType

	// load all symbols
	res.Symbols, err = l.AnalyzeSymbols(modPath)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type NotSupport struct {
}

func NewNotSupport() *NotSupport {
	return &NotSupport{}
}

func (l *NotSupport) IsSupport(filePath string) bool {
	// handle all not support file
	return true
}

func (l *NotSupport) AnalyzeSymbols(filePath string) ([]*Symbol, error) {
	return nil, ErrNotSupport
}

func (l *NotSupport) ToModule(pid int32, modName, modPath string, moduleRange []*ModuleRange) (*Module, error) {
	return nil, ErrNotSupport
}

type analyzeContext struct {
	pathToFinder map[string]StatFinder
}

func newAnalyzeContext() *analyzeContext {
	return &analyzeContext{
		pathToFinder: make(map[string]StatFinder),
	}
}

func (a *analyzeContext) GetFinder(name string) StatFinder {
	if f := a.pathToFinder[name]; f != nil {
		return f
	}

	// find all finders
	for _, f := range profilingStatFinderList {
		if f.IsSupport(name) {
			a.pathToFinder[name] = f
			return f
		}
	}

	// not support
	n := NewNotSupport()
	a.pathToFinder[name] = n
	return n
}

func isIgnoreModuleName(name string) bool {
	return name != "" &&
		(strings.HasPrefix(name, "//anon") ||
			strings.HasPrefix(name, "/dev/zero") ||
			strings.HasPrefix(name, "/anon_hugepage") ||
			strings.HasPrefix(name, "[stack") ||
			strings.HasPrefix(name, "/SYSV") ||
			strings.HasPrefix(name, "[heap]") ||
			strings.HasPrefix(name, "/memfd:") ||
			strings.HasPrefix(name, "[vdso]") ||
			strings.HasPrefix(name, "[vsyscall]") ||
			strings.HasPrefix(name, "[uprobes]") ||
			strings.HasSuffix(name, ".map"))
}

func parseUInt64InModule(err error, moduleName, key, val string) (uint64, error) {
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseUint(val, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing the %s in maps file: %s", key, moduleName)
	}
	return res, nil
}

func analyzeProfilingInfo(context *analyzeContext, pid int32) (*Info, error) {
	// analyze process mapping
	mapFile, _ := os.Open(fmt.Sprintf("/proc/%d/maps", pid))
	scanner := bufio.NewScanner(mapFile)
	modules := make(map[string]*Module)
	for scanner.Scan() {
		submatch := mapFileContentRegex.FindStringSubmatch(scanner.Text())
		if len(submatch) != 6 {
			continue
		}
		if len(submatch[3]) > 2 && submatch[3][2] != 'x' {
			continue
		}
		moduleName := submatch[5]
		if isIgnoreModuleName(moduleName) {
			continue
		}

		// parsing range
		var err error
		moduleRange := &ModuleRange{}
		moduleRange.StartAddr, err = parseUInt64InModule(err, moduleName, "start address", submatch[1])
		moduleRange.EndAddr, err = parseUInt64InModule(err, moduleName, "end address", submatch[2])
		moduleRange.FileOffset, err = parseUInt64InModule(err, moduleName, "file offset", submatch[4])
		if err != nil {
			return nil, err
		}

		module := modules[moduleName]
		if module != nil {
			module.Ranges = append(module.Ranges, moduleRange)
			continue
		}
		modulePath := fmt.Sprintf("/proc/%d/root%s", pid, moduleName)
		if !Exists(modulePath) {
			fmt.Printf("could not found the module, ignore. name: %s, path: %s", moduleName, modulePath)
			continue
		}

		module, err = context.GetFinder(modulePath).ToModule(pid, moduleName, modulePath, []*ModuleRange{moduleRange})
		if err != nil {
			if err == ErrNotSupport {
				fmt.Printf("not support the module in process(%d): %s, path: %s", pid, moduleName, modulePath)
				continue
			}
			return nil, fmt.Errorf("could not init the module: %s, error: %v", moduleName, err)
		}
		modules[moduleName] = module
	}
	return NewInfo(modules), nil
}

func ProfilingStat(pid int32, exePath string) (*Info, error) {
	stat, err := os.Stat(exePath)
	if err != nil {
		return nil, fmt.Errorf("check file error: %v", err)
	}
	for _, notSupport := range NotSupportProfilingExe {
		if strings.HasPrefix(stat.Name(), notSupport) {
			return nil, fmt.Errorf("not support %s language profiling", notSupport)
		}
	}
	context := newAnalyzeContext()

	// the executable file must have the symbols
	symbols, err := context.GetFinder(exePath).AnalyzeSymbols(exePath)
	if err != nil || len(symbols) == 0 {
		return nil, fmt.Errorf("could not found any symbol in the execute file: %s, error: %v", exePath, err)
	}

	return analyzeProfilingInfo(context, pid)
}

type OnCPUProfilingTask struct {
	Pid                int32
	Exec               string
	Duration           time.Duration
	Interval           time.Duration
	obj                *OnCPUProfilingObjects
	processProfiling   *Info
	kernelProfiling    *Info
	stopChan           chan bool
	perfEventFds       []int
	flushDataNotify    context.CancelFunc
	ShutdownOnce       sync.Once
	stackCounter       map[Event]uint32
	StackNotFoundCache map[uint32]bool
}

func (t *OnCPUProfilingTask) Init() {
	t.stackCounter = make(map[Event]uint32)
	t.stopChan = make(chan bool, 1)

	kernelFinder := NewKernelFinder()
	t.kernelProfiling, _ = kernelFinder.Analyze(KernelSymbolFilePath)

	t.processProfiling, _ = ProfilingStat(t.Pid, t.Exec)

	t.StackNotFoundCache = make(map[uint32]bool)
}

func (t *OnCPUProfilingTask) Run(notify task.RunningSuccessNotify) {
	spec, err := LoadOnCPUProfiling()
	if err != nil {
		panic(err)
	}
	var replacedPid bool
	for i, ins := range spec.Programs[FuncName].Instructions {
		if ins.Reference() == "MONITOR_PID" {
			spec.Programs[FuncName].Instructions[i].Constant = int64(t.Pid)
			spec.Programs[FuncName].Instructions[i].Offset = 0
			replacedPid = true
		}
	}

	if !replacedPid {
		panic(fmt.Errorf("replace the monitor pid failure"))
	}

	var obj OnCPUProfilingObjects

	if err1 := spec.LoadAndAssign(&obj, nil); err1 != nil {
		panic(fmt.Errorf("loading objects: %s", err1))
	}
	defer obj.Close()
	t.obj = &obj

	perfEvents, err := t.openPerfEvent(obj.DoPerfEvent.FD())
	t.perfEventFds = perfEvents
	if err != nil {
		panic(err)
	}

	// notify start success
	notify()
	runtime.SetFinalizer(t, (*OnCPUProfilingTask).Stop)
	<-t.stopChan
	fmt.Println("Done")
}

func (t *OnCPUProfilingTask) openPerfEvent(perfFd int) ([]int, error) {
	eventAttr := &unix.PerfEventAttr{
		Type:   unix.PERF_TYPE_SOFTWARE,
		Config: unix.PERF_COUNT_SW_CPU_CLOCK,
		Bits:   unix.PerfBitFreq,
		Sample: uint64(t.Interval),
		Wakeup: 1,
	}

	fds := make([]int, 0)
	for cpuNum := 0; cpuNum < runtime.NumCPU(); cpuNum++ {
		fd, err := unix.PerfEventOpen(
			eventAttr,
			-1,
			cpuNum,
			-1,
			0,
		)
		if err != nil {
			return fds, err
		}

		// attach ebpf to perf event
		if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_SET_BPF, perfFd); err != nil {
			return fds, err
		}

		// enable perf event
		if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_ENABLE, 0); err != nil {
			return fds, err
		}
		fds = append(fds, fd)
	}

	return fds, nil
}

func (obj *OnCPUProfilingObj) Close() error {
	if err := obj.Counts.Close(); err != nil {
		return err
	}
	if err := obj.Stacks.Close(); err != nil {
		return err
	}
	if err := obj.Prog.Close(); err != nil {
		return err
	}
	return nil
}

func (t *OnCPUProfilingTask) closePerfEvent(fd int) error {
	if fd <= 0 {
		return nil
	}
	var result error
	if err := unix.IoctlSetInt(fd, unix.PERF_EVENT_IOC_DISABLE, 0); err != nil {
		result = multierror.Append(result, fmt.Errorf("closing perf event reader: %s", err))
	}
	return result
}

func (t *OnCPUProfilingTask) Stop() error {
	var result error
	t.ShutdownOnce.Do(func() {
		for _, fd := range t.perfEventFds {
			if err := t.closePerfEvent(fd); err != nil {
				result = multierror.Append(result, err)
			}
		}

		// wait for all profiling data been consume finished
		cancel, cancelFunc := context.WithCancel(context.Background())
		t.flushDataNotify = cancelFunc
		select {
		case <-cancel.Done():
		case <-time.After(5 * time.Second):
		}

		if t.obj != nil {
			if err := t.obj.Close(); err != nil {
				result = multierror.Append(result, err)
			}
		}

		close(t.stopChan)
	})
	return result
}

func (t *OnCPUProfilingTask) GenerateProfilingData(profilingInfo *Info, stackID uint32, stackMap *ebpf.Map,
	stackType int32, symbolArray []uint64) *task.EBPFProfilingStackMetadata {
	if profilingInfo == nil || stackID <= 0 || stackID == math.MaxUint32 {
		return nil
	}
	if err := stackMap.Lookup(stackID, symbolArray); err != nil {
		if t.StackNotFoundCache[stackID] {
			return nil
		}
		t.StackNotFoundCache[stackID] = true
		log.Default().Printf("error to lookup %v stack: %d, error: %v", stackType, stackID, err)
		return nil
	}
	symbols := profilingInfo.FindSymbols(symbolArray, MissingSymbol)
	if len(symbols) == 0 {
		return nil
	}
	return &task.EBPFProfilingStackMetadata{
		StackType:    stackType,
		StackId:      int32(stackID),
		StackSymbols: symbols,
	}
}

func (t *OnCPUProfilingTask) FlushData() ([]*task.EBPFProfilingData, error) {
	if t.obj == nil {
		return nil, nil
	}
	var stack Event
	var counter uint32
	iterate := t.obj.Counts.Iterate()
	stacks := t.obj.Stacks
	result := make([]*task.EBPFProfilingData, 0)
	stackSymbols := make([]uint64, 100)
	for iterate.Next(&stack, &counter) {
		metadatas := make([]*task.EBPFProfilingStackMetadata, 0)
		// kernel stack
		if d := t.GenerateProfilingData(t.kernelProfiling, stack.KernelStackID, stacks,
			task.EBPFProfilingStackType_PROCESS_KERNEL_SPACE, stackSymbols); d != nil {
			metadatas = append(metadatas, d)
		}

		// user stack
		if d := t.GenerateProfilingData(t.processProfiling, stack.UserStackID, stacks,
			task.EBPFProfilingStackType_PROCESS_USER_SPACE, stackSymbols); d != nil {
			metadatas = append(metadatas, d)
		}

		if len(metadatas) == 0 {
			continue
		}

		// update the counters in memory
		dumpCount := int32(counter)
		existCounter := t.stackCounter[stack]
		if existCounter > 0 {
			dumpCount -= int32(existCounter)
		}
		t.stackCounter[stack] = counter
		if dumpCount <= 0 {
			continue
		}

		result = append(result, &task.EBPFProfilingData{
			Profiling: &task.EBPFProfilingData_OnCPU{
				OnCPU: &task.EBPFOnCPUProfiling{
					Stacks:    metadatas,
					DumpCount: dumpCount,
				},
			},
		})
	}

	// close the flush data notify if exists
	if t.flushDataNotify != nil {
		t.flushDataNotify()
	}

	return result, nil
}

func DoProfiling(pid int32, cmdPath string, duration time.Duration, successNotify task.RunningSuccessNotify, onFinish func(data []*task.EBPFProfilingData)) {
	interval := 1 * time.Microsecond
	isRunning := make(chan bool, 1)

	task := OnCPUProfilingTask{Pid: pid, Exec: cmdPath, Duration: duration, Interval: interval}

	defer task.Stop()

	task.Init()

	go task.Run(func() {
		isRunning <- true
		successNotify()
	})

	<-isRunning

	timer := time.NewTimer(duration)
	<-timer.C
	dataLst, _ := task.FlushData()
	if len(dataLst) > 0 {
		b, _ := json.Marshal(dataLst)
		fmt.Println(string(b))
	}
	onFinish(dataLst)
}

func CreateProfilingTask(jobInfo task.Job, callback task.RunningSuccessNotify, onFinish func(data []*task.EBPFProfilingData)) {
	pid := int32(jobInfo.Pid)
	duration := int32(jobInfo.Duration)
	process, _ := process.NewProcess(pid)
	cmdPath, _ := process.Exe()
	go func() {
		DoProfiling(pid, cmdPath, time.Duration(duration)*time.Minute, callback, onFinish)
	}()
}
