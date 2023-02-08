package sqlite3

import (
	"context"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Configure SQLite.
var (
	Binary []byte // Binary to load.
	Path   string // Path to load the binary from.

	Interpreter bool
)

var sqlite3 sqlite3Runtime

type sqlite3Runtime struct {
	runtime   wazero.Runtime
	compiled  wazero.CompiledModule
	instances atomic.Uint64
	ctx       context.Context
	err       error
}

func (s *sqlite3Runtime) instantiateModule(ctx context.Context) (api.Module, error) {
	s.ctx = ctx
	s.compileModule()
	if s.err != nil {
		return nil, s.err
	}

	cfg := wazero.NewModuleConfig().
		WithName("sqlite3-" + strconv.FormatUint(s.instances.Add(1), 10))
	return s.runtime.InstantiateModule(ctx, s.compiled, cfg)
}

func (s *sqlite3Runtime) compileModule() {
	if Interpreter {
		s.runtime = wazero.NewRuntimeWithConfig(s.ctx, wazero.NewRuntimeConfigInterpreter())
	} else {
		s.runtime = wazero.NewRuntimeWithConfig(s.ctx, wazero.NewRuntimeConfigCompiler())
	}
	s.err = vfsInstantiate(s.ctx, s.runtime)
	if s.err != nil {
		return
	}

	bin := Binary
	if bin == nil && Path != "" {
		bin, s.err = os.ReadFile(Path)
		if s.err != nil {
			return
		}
	}

	s.compiled, s.err = s.runtime.CompileModule(s.ctx, bin)
}
